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

import {
  LOG,
  PROJECT_DIR,
  isEmptyString,
  toggleVerbosity,
  time,
  setNoQuoteEscape,
  setQuoteEscape,
} from "./lib/index.mjs";

const VERSION = argv.version;
const IMG = argv.img;
const KEEP_CLUSTER = argv["keep-cluster"] === true;

toggleVerbosity(argv.verbose);

if (isEmptyString(VERSION)) {
  LOG.red("--version is required");
  process.exit(1);
}
if (isEmptyString(IMG)) {
  LOG.red("--img is required");
  process.exit(1);
}

const CLUSTER_NAME = "gko-olm-test";
const KIND_CONFIG = path.join(PROJECT_DIR, "hack", "kind", "kind.olm.yaml");
const BIN_DIR = path.join(PROJECT_DIR, "bin");
const TEST_NS = "gko-olm-test";
const OLM_NS = "olm";
const OLM_VERSION = "v0.28.0";
const OPM_VERSION = "v1.26.0";
const EXPECTED_CRD_COUNT = 11;

const REGISTRY_NAME = "kind-registry";
const REGISTRY_PORT = 5050;

const BUNDLE_IMG = "gko-bundle:test";
const BUNDLE_IMG_REGISTRY = `localhost:${REGISTRY_PORT}/gko-bundle:test`;
const INDEX_IMG_LOCAL = "gko-index:test";
const INDEX_IMG = `localhost:${REGISTRY_PORT}/gko-index:test`;
const OPERATOR_IMG = `${IMG}:${VERSION}`;

const OLM_BASE_URL = `https://github.com/operator-framework/operator-lifecycle-manager/releases/download/${OLM_VERSION}`;

LOG.magenta(`
  OLM Bundle Test
    Version:      ${VERSION}
    Operator:     ${OPERATOR_IMG}
    Cluster:      ${CLUSTER_NAME}
    Keep cluster: ${KEEP_CLUSTER}
`);

function osPlatform() {
  const p = process.platform;
  if (p === "darwin") return "darwin";
  if (p === "win32") return "windows";
  return "linux";
}

function osArch() {
  const a = process.arch;
  if (a === "arm64") return "arm64";
  return "amd64";
}

async function ensureTool(name) {
  try {
    await $`which ${name}`.quiet();
  } catch {
    throw new Error(`Required tool "${name}" not found in PATH`);
  }
}

async function ensureOpm() {
  const opmPath = path.join(BIN_DIR, "opm");
  try {
    await fs.access(opmPath, fs.constants.X_OK);
    LOG.green(`  opm already present at ${opmPath}`);
    return opmPath;
  } catch {
    // need to download
  }

  const platform = osPlatform();
  const arch = osArch();
  const url = `https://github.com/operator-framework/operator-registry/releases/download/${OPM_VERSION}/${platform}-${arch}-opm`;

  LOG.blue(`  Downloading opm ${OPM_VERSION} for ${platform}/${arch} ...`);
  await fs.ensureDir(BIN_DIR);
  await $`curl -fsSL -o ${opmPath} ${url}`;
  await $`chmod +x ${opmPath}`;
  LOG.green(`  opm installed at ${opmPath}`);
  return opmPath;
}

async function startRegistry() {
  try {
    const result =
      await $`docker inspect ${REGISTRY_NAME} --format={{.State.Running}}`.quiet();
    if (result.stdout.trim() === "true") {
      LOG.green(`  Registry "${REGISTRY_NAME}" already running`);
      return;
    }
  } catch { /* not running */ }

  await $`docker rm -f ${REGISTRY_NAME}`.nothrow().quiet();
  await $`docker run -d --restart=always -p 127.0.0.1:${REGISTRY_PORT}:5000 --name ${REGISTRY_NAME} registry:2`;
  LOG.green(`  Registry started on localhost:${REGISTRY_PORT}`);
}

async function connectRegistryToKind() {
  await $`docker network connect kind ${REGISTRY_NAME}`.nothrow().quiet();

  const registryDir = `/etc/containerd/certs.d/localhost:${REGISTRY_PORT}`;
  const node = `${CLUSTER_NAME}-control-plane`;
  await $`docker exec ${node} mkdir -p ${registryDir}`;

  const hostsToml = [
    `[host."http://${REGISTRY_NAME}:5000"]`,
    `  capabilities = ["pull", "resolve"]`,
    "",
  ].join("\n");
  const tmpFile = path.join(os.tmpdir(), "hosts.toml");
  await fs.writeFile(tmpFile, hostsToml);
  await $`docker cp ${tmpFile} ${node}:${registryDir}/hosts.toml`;
  await fs.unlink(tmpFile).catch(() => {});

  LOG.green(`  Registry connected to Kind network and configured on node`);
}

async function stopRegistry() {
  await $`docker rm -f ${REGISTRY_NAME}`.nothrow().quiet();
  LOG.green(`  Registry removed`);
}

async function deleteClusterIfExists() {
  try {
    const result = await $`kind get clusters`.quiet();
    const clusters = result.stdout.trim().split("\n");
    if (clusters.includes(CLUSTER_NAME)) {
      LOG.yellow(`  Deleting pre-existing cluster "${CLUSTER_NAME}" ...`);
      await $`kind delete cluster --name ${CLUSTER_NAME}`;
    }
  } catch {
    // kind not running or no clusters — fine
  }
}

async function createCluster() {
  LOG.blue("  Creating Kind cluster ...");
  setNoQuoteEscape();
  await $`kind create cluster --config ${KIND_CONFIG} --wait 60s`;
  setQuoteEscape();
  LOG.green(`  Cluster "${CLUSTER_NAME}" created`);
}

async function installOLM() {
  LOG.blue(`  Installing OLM ${OLM_VERSION} ...`);

  await $`kubectl apply --server-side --force-conflicts -f ${OLM_BASE_URL}/crds.yaml`;
  await $`kubectl wait --for=condition=Established -f ${OLM_BASE_URL}/crds.yaml --timeout=60s`;
  await $`kubectl apply --server-side --force-conflicts -f ${OLM_BASE_URL}/olm.yaml`;

  LOG.blue("  Waiting for OLM pods ...");
  await $`kubectl rollout status deployment/olm-operator -n ${OLM_NS} --timeout=120s`;
  await $`kubectl rollout status deployment/catalog-operator -n ${OLM_NS} --timeout=120s`;
  LOG.green("  OLM installed");
}

async function buildBundleImage() {
  LOG.blue("  Building bundle image ...");
  await $`docker build -f bundle.Dockerfile -t ${BUNDLE_IMG} .`;
  LOG.blue("  Pushing bundle image to local registry ...");
  await $`docker tag ${BUNDLE_IMG} ${BUNDLE_IMG_REGISTRY}`;
  await $`docker push ${BUNDLE_IMG_REGISTRY}`;
  LOG.green(`  Built and pushed ${BUNDLE_IMG_REGISTRY}`);
}

async function buildOperatorImage() {
  LOG.blue("  Building operator image ...");
  await $`docker build -t ${OPERATOR_IMG} .`;
  LOG.green(`  Built ${OPERATOR_IMG}`);
}

async function buildCatalogIndex(opmPath) {
  LOG.blue("  Building catalog index with opm ...");
  await $`${opmPath} index add --bundles ${BUNDLE_IMG_REGISTRY} --tag ${INDEX_IMG_LOCAL} --container-tool docker`;
  LOG.blue("  Pushing index image to local registry ...");
  await $`docker tag ${INDEX_IMG_LOCAL} ${INDEX_IMG}`;
  await $`docker push ${INDEX_IMG}`;
  LOG.green(`  Built and pushed ${INDEX_IMG}`);
}

async function loadImagesIntoKind() {
  LOG.blue("  Loading operator image into Kind ...");
  await $`kind load docker-image ${OPERATOR_IMG} --name ${CLUSTER_NAME}`;
  LOG.green("  Operator image loaded");
}

async function applyOLMResources() {
  LOG.blue("  Creating namespace and OLM resources ...");

  await $`kubectl create ns ${TEST_NS}`;

  const catalogSource = YAML.stringify({
    apiVersion: "operators.coreos.com/v1alpha1",
    kind: "CatalogSource",
    metadata: { name: "gko-catalog", namespace: OLM_NS },
    spec: {
      sourceType: "grpc",
      image: INDEX_IMG,
      displayName: "GKO Test Catalog",
    },
  });

  const operatorGroup = YAML.stringify({
    apiVersion: "operators.coreos.com/v1",
    kind: "OperatorGroup",
    metadata: { name: "gko-og", namespace: TEST_NS },
    spec: { targetNamespaces: [TEST_NS] },
  });

  const subscription = YAML.stringify({
    apiVersion: "operators.coreos.com/v1alpha1",
    kind: "Subscription",
    metadata: { name: "gko-sub", namespace: TEST_NS },
    spec: {
      channel: "alpha",
      name: "gko",
      source: "gko-catalog",
      sourceNamespace: OLM_NS,
      installPlanApproval: "Automatic",
    },
  });

  for (const manifest of [catalogSource, operatorGroup, subscription]) {
    const tmpFile = path.join(
      os.tmpdir(),
      `olm-test-${Date.now()}-${Math.random().toString(36).slice(2)}.yaml`,
    );
    await fs.writeFile(tmpFile, manifest);
    try {
      await $`kubectl apply -f ${tmpFile}`;
    } finally {
      await fs.unlink(tmpFile).catch(() => {});
    }
  }

  LOG.green("  OLM resources applied");
}

async function waitForCSV() {
  LOG.blue("  Waiting for CSV to reach Succeeded phase ...");

  const maxAttempts = 40;
  const intervalMs = 5000;

  for (let i = 0; i < maxAttempts; i++) {
    try {
      const result =
        await $`kubectl get csv -n ${TEST_NS} -o jsonpath='{.items[0].status.phase}'`.quiet();
      const phase = result.stdout.replace(/'/g, "").trim();
      if (phase === "Succeeded") {
        LOG.green(`  CSV phase: ${phase}`);
        return;
      }
      if (phase === "Failed") {
        const reason =
          await $`kubectl get csv -n ${TEST_NS} -o jsonpath='{.items[0].status.message}'`.quiet();
        throw new Error(`CSV failed: ${reason.stdout}`);
      }
      LOG.log(`  CSV phase: ${phase || "(not yet available)"} (attempt ${i + 1}/${maxAttempts})`);
    } catch (e) {
      if (e.message?.startsWith("CSV failed")) throw e;
      LOG.log(`  Waiting for CSV to appear (attempt ${i + 1}/${maxAttempts}) ...`);
    }
    await sleep(intervalMs);
  }
  throw new Error(`CSV did not reach Succeeded phase within ${(maxAttempts * intervalMs) / 1000}s`);
}

async function verifyOperatorPod() {
  LOG.blue("  Verifying operator pod ...");
  await $`kubectl wait --for=condition=ready pod -l control-plane=controller-manager -n ${TEST_NS} --timeout=120s`;
  const pods =
    await $`kubectl get pods -n ${TEST_NS} -l control-plane=controller-manager -o wide`.quiet();
  LOG.green(`  Operator pod is running:\n${pods.stdout}`);
}

async function verifyCRDs() {
  LOG.blue("  Verifying Gravitee CRDs ...");
  const result = await $`kubectl get crd`.quiet();
  const graviteeCRDs = result.stdout
    .split("\n")
    .filter((line) => line.includes("gravitee.io"));

  if (graviteeCRDs.length !== EXPECTED_CRD_COUNT) {
    LOG.red(
      `  Expected ${EXPECTED_CRD_COUNT} Gravitee CRDs, found ${graviteeCRDs.length}`,
    );
    LOG.red(graviteeCRDs.join("\n"));
    throw new Error("CRD count mismatch");
  }

  LOG.green(`  All ${EXPECTED_CRD_COUNT} Gravitee CRDs registered`);
}

async function printDiagnostics() {
  LOG.yellow("\n  === Diagnostics ===");
  try {
    const events = await $`kubectl get events -n ${TEST_NS} --sort-by=.lastTimestamp`.nothrow().quiet();
    LOG.log(events.stdout);
  } catch { /* best effort */ }
  try {
    const csvs = await $`kubectl get csv -n ${TEST_NS} -o yaml`.nothrow().quiet();
    LOG.log(csvs.stdout);
  } catch { /* best effort */ }
  try {
    const subs = await $`kubectl get subscription -n ${TEST_NS} -o yaml`.nothrow().quiet();
    LOG.log(subs.stdout);
  } catch { /* best effort */ }
  try {
    const pods = await $`kubectl get pods -n ${TEST_NS} -o wide`.nothrow().quiet();
    LOG.log(pods.stdout);
  } catch { /* best effort */ }
  try {
    const catalog = await $`kubectl get catalogsource -n ${OLM_NS} -o yaml`.nothrow().quiet();
    LOG.log(catalog.stdout);
  } catch { /* best effort */ }
  try {
    const olmPods = await $`kubectl get pods -n ${OLM_NS} -o wide`.nothrow().quiet();
    LOG.log(olmPods.stdout);
  } catch { /* best effort */ }
}

async function teardown() {
  if (KEEP_CLUSTER) {
    LOG.yellow(
      `  --keep-cluster set, cluster "${CLUSTER_NAME}" left running for debugging`,
    );
    return;
  }
  LOG.blue(`  Tearing down cluster "${CLUSTER_NAME}" ...`);
  await $`kind delete cluster --name ${CLUSTER_NAME}`;
  LOG.green("  Cluster deleted");
  await stopRegistry();
}

// ── Main ──

let failed = false;
try {
  LOG.blue("Checking prerequisites ...");
  await time(async () => {
    await ensureTool("kind");
    await ensureTool("kubectl");
    await ensureTool("docker");
    LOG.green("  All required tools available");
  });

  const opmPath = await ensureOpm();

  LOG.blue("\nStep 1/8: Start local registry");
  await time(startRegistry);

  LOG.blue("\nStep 2/8: Create Kind cluster");
  await time(async () => {
    await deleteClusterIfExists();
    await createCluster();
    await connectRegistryToKind();
  });

  LOG.blue("\nStep 3/8: Install OLM");
  await time(installOLM);

  LOG.blue("\nStep 4/8: Build images");
  await time(async () => {
    await buildBundleImage();
    await buildOperatorImage();
  });

  LOG.blue("\nStep 5/8: Build catalog index");
  await time(() => buildCatalogIndex(opmPath));

  LOG.blue("\nStep 6/8: Load images into Kind");
  await time(loadImagesIntoKind);

  LOG.blue("\nStep 7/8: Apply OLM resources");
  await time(applyOLMResources);

  LOG.blue("\nStep 8/8: Verify installation");
  await time(async () => {
    await waitForCSV();
    await verifyOperatorPod();
    await verifyCRDs();
  });

  LOG.green(`
  OLM bundle test PASSED.

  The operator was successfully installed via OLM in the "${TEST_NS}" namespace.

  CSV:      gko.v${VERSION} (phase: Succeeded)
  CRDs:     ${EXPECTED_CRD_COUNT} Gravitee CRDs registered
  Cluster:  ${CLUSTER_NAME}
  `);
} catch (e) {
  failed = true;
  LOG.red(`\n  OLM bundle test FAILED: ${e.message}\n`);
  await printDiagnostics();
} finally {
  await teardown();
}

if (failed) {
  process.exit(1);
}
