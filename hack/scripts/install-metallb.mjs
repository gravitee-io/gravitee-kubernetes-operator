/**
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { LOG, PROJECT_DIR, time, toggleVerbosity } from "./lib/index.mjs";

toggleVerbosity(argv.verbose);

const METALLB_VERSION = $.env.METALLB_VERSION || "v0.14.8";
const METALLB_MANIFEST_URL = `https://raw.githubusercontent.com/metallb/metallb/${METALLB_VERSION}/config/manifests/metallb-native.yaml`;
const METALLB_CONFIG_PATH = path.join(
  PROJECT_DIR,
  "hack",
  "kind",
  "metallb",
  "metallb-config.yaml",
);

// Default fallback IP range (current hardcoded value)
const DEFAULT_IP_RANGE = "172.18.255.200-172.18.255.250";

async function installMetalLB() {
  LOG.blue(`Installing MetalLB ${METALLB_VERSION}...`);
  await $`kubectl apply -f ${METALLB_MANIFEST_URL}`;
}

async function waitForCRDs() {
  LOG.blue("Waiting for MetalLB CRDs to be available...");
  try {
    await $`kubectl wait --for condition=established --timeout=60s crd/ipaddresspools.metallb.io`;
  } catch {
    // CRD might already exist, continue
  }
  try {
    await $`kubectl wait --for condition=established --timeout=60s crd/l2advertisements.metallb.io`;
  } catch {
    // CRD might already exist, continue
  }
}

async function waitForPods() {
  LOG.blue("Waiting for MetalLB pods to be ready...");
  await $`kubectl wait --namespace metallb-system --for=condition=ready pod --selector=app=metallb --timeout=90s`;
}

/**
 * Detects the Docker network subnet for the Kind cluster
 * @returns {Promise<string|null>} The subnet in CIDR notation (e.g., "172.18.0.0/16") or null if detection fails
 */
async function detectKindNetworkSubnet() {
  try {
    // Get the current kubectl context to determine cluster name
    const contextOutput = await $`kubectl config current-context`.quiet();
    const context = contextOutput.stdout.trim();

    // Extract cluster name from context (format: kind-<cluster-name>)
    let clusterName = "gravitee"; // default
    if (context.startsWith("kind-")) {
      clusterName = context.replace("kind-", "");
    }

    // Kind creates a network named "kind" by default, but we'll try both
    const networkNames = ["kind", `kind-${clusterName}`];

    for (const networkName of networkNames) {
      try {
        const networkInfo =
          await $`docker network inspect ${networkName} --format '{{json .}}'`.quiet();
        const network = JSON.parse(networkInfo.stdout.trim());

        // Check IPAM config for subnet
        if (
          network.IPAM &&
          network.IPAM.Config &&
          network.IPAM.Config.length > 0
        ) {
          const subnet = network.IPAM.Config[0].Subnet;
          if (subnet) {
            LOG.blue(`Detected Kind network subnet: ${subnet}`);
            return subnet;
          }
        }
      } catch (e) {
        // Network doesn't exist, try next one
        continue;
      }
    }

    // Fallback: try to find any network with "kind" in the name
    try {
      const allNetworks =
        await $`docker network ls --format '{{.Name}}'`.quiet();
      const networkList = allNetworks.stdout.trim().split("\n");
      const kindNetwork = networkList.find((name) => name.includes("kind"));

      if (kindNetwork) {
        const networkInfo =
          await $`docker network inspect ${kindNetwork} --format '{{json .}}'`.quiet();
        const network = JSON.parse(networkInfo.stdout.trim());

        if (
          network.IPAM &&
          network.IPAM.Config &&
          network.IPAM.Config.length > 0
        ) {
          const subnet = network.IPAM.Config[0].Subnet;
          if (subnet) {
            LOG.blue(`Detected Kind network subnet: ${subnet}`);
            return subnet;
          }
        }
      }
    } catch (e) {
      // Ignore errors in fallback
    }

    return null;
  } catch (error) {
    LOG.yellow(
      `Warning: Could not detect Kind network subnet: ${error.message}`,
    );
    return null;
  }
}

/**
 * Calculates a safe IP range from a subnet CIDR
 * Uses the last octet range .200-.250 to avoid conflicts
 * @param {string} subnet - Subnet in CIDR notation (e.g., "172.18.0.0/16")
 * @returns {string} IP range in format "X.X.X.200-X.X.X.250"
 */
function calculateIPRange(subnet) {
  // Parse CIDR notation (e.g., "172.18.0.0/16")
  const [baseIP, prefixLength] = subnet.split("/");
  const prefixLen = parseInt(prefixLength, 10);

  // Split IP into octets
  const octets = baseIP.split(".").map(Number);

  // For /16 networks (most common for Docker), use first 2 octets
  // For /24 networks, use first 3 octets
  // For other networks, we'll use a conservative approach

  if (prefixLen <= 16) {
    // /16 or larger: use first 2 octets, set 3rd to 255, 4th to 200-250
    return `${octets[0]}.${octets[1]}.255.200-${octets[0]}.${octets[1]}.255.250`;
  } else if (prefixLen <= 24) {
    // /24: use first 3 octets, 4th to 200-250
    return `${octets[0]}.${octets[1]}.${octets[2]}.200-${octets[0]}.${octets[1]}.${octets[2]}.250`;
  } else {
    // Smaller subnets: use the base IP + .200-.250 pattern on last octet
    // This is a conservative fallback
    const base = octets.slice(0, 3).join(".");
    return `${base}.200-${base}.250`;
  }
}

/**
 * Generates MetalLB configuration YAML
 * @param {string} ipRange - IP address range (e.g., "172.18.255.200-172.18.255.250")
 * @returns {string} YAML configuration
 */
function generateMetalLBConfig(ipRange) {
  return `# Copyright (C) 2015 The Gravitee team (http://gravitee.io)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: default
  namespace: metallb-system
spec:
  addresses:
  - ${ipRange}

---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: default
  namespace: metallb-system
spec:
  ipAddressPools:
  - default
`;
}

async function configureMetalLB() {
  LOG.blue("Configuring MetalLB IP address pool...");

  // Try to detect the network subnet dynamically
  const subnet = await detectKindNetworkSubnet();
  let ipRange = DEFAULT_IP_RANGE;
  let usingFallback = true;

  if (subnet) {
    try {
      ipRange = calculateIPRange(subnet);
      usingFallback = false;
      LOG.blue(`Using detected IP range: ${ipRange}`);
    } catch (error) {
      LOG.yellow(
        `Warning: Failed to calculate IP range from subnet ${subnet}, using fallback: ${error.message}`,
      );
    }
  } else {
    LOG.yellow(
      `Warning: Could not detect Kind network subnet, using fallback range: ${ipRange}`,
    );
    LOG.yellow(
      `If you encounter IP conflicts, ensure Docker uses the default 172.18.0.0/16 network`,
    );
  }

  // Generate and apply the configuration
  const configYaml = generateMetalLBConfig(ipRange);

  const tmpFile = path.join(os.tmpdir(), `metallb-config-${Date.now()}.yaml`);
  await fs.writeFile(tmpFile, configYaml);
  try {
    await $`kubectl apply -f ${tmpFile}`;
  } finally {
    await fs.unlink(tmpFile).catch(() => {});
  }

  if (usingFallback) {
    LOG.yellow(
      `⚠️  Using fallback IP range. Consider verifying Docker network configuration.`,
    );
  } else {
    LOG.green(`✓ MetalLB configured with IP range: ${ipRange}`);
  }
}

function isMacOS() {
  try {
    const platform = process.platform;
    return platform === "darwin";
  } catch {
    return false;
  }
}

LOG.blue(`
  ☸ Installing MetalLB
`);

if (isMacOS()) {
  LOG.yellow(`
  ⚠️  macOS DETECTED
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  
  MetalLB IP addresses will NOT be accessible from your macOS host because Docker
  Desktop runs in a VM. The IPs assigned by MetalLB exist only within the Docker
  network and cannot be reached from your Mac.
  
  For macOS, we recommend using cloud-provider-kind instead:
  
    make cloud-lb
  
  This will map LoadBalancer services directly to localhost ports, making them
  accessible from your Mac.
  
  MetalLB will still be installed, but you'll need to use port-forwarding or
  access services from within the cluster to reach them.
  
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  `);

  await new Promise((resolve) => setTimeout(resolve, 2000));
}

await time(installMetalLB);

await time(waitForCRDs);

await time(waitForPods);

await time(configureMetalLB);

LOG.green(`
  ✅ MetalLB installed and configured successfully
`);
