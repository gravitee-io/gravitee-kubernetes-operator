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

import { execFile } from "node:child_process";
import { promisify } from "node:util";
import { loadGraviteeConfig } from "../src/cmd/config.js";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const execFileAsync = promisify(execFile);

/**
 * Warm up the APIM Automation API before any tests run.
 *
 * The Automation API is the path the Terraform provider writes through. On a
 * freshly (re)started APIM it is cold — the first heavyweight request pays JVM
 * JIT, lazy controller init, and connection-pool warm-up, and can take far
 * longer than a steady-state call. That cold cost otherwise lands on whichever
 * terraform test happens to run first, where it can blow the hook timeout.
 *
 * This issues authenticated GETs against the Automation API paths the writers
 * use (`.../apis` and `.../applications`), blocking until each answers quickly
 * so the first real request sees a warm server. The `applications` path matters
 * for the upgrade after-phase: the GKO admission webhook dry-run-validates
 * Applications against it, and after an in-place APIM upgrade the gateway is up
 * long before the Automation controllers are, so an un-warmed endpoint makes the
 * webhook -> APIM call time out. Best-effort: failures are logged, never thrown
 * (a genuinely unreachable APIM is already caught by the infrastructure checks),
 * and a global cap keeps a stubbornly-cold APIM from blocking the suite forever.
 */
async function warmUpApim(
  mapiUrl: string,
  envId: string,
  auth: { username: string; password: string },
): Promise<void> {
  // org is "DEFAULT" to match the terraform fixtures' Automation API path.
  const base = `${mapiUrl}/automation/organizations/DEFAULT/environments/${envId}`;
  const paths = ["apis", "applications"];
  const credentials = Buffer.from(`${auth.username}:${auth.password}`).toString("base64");
  const headers = { Authorization: `Basic ${credentials}` };
  const WARM_ENOUGH_MS = 1_500;
  const MAX_WAIT_MS = 120_000;
  const started = Date.now();

  // Block until every endpoint answers quickly (warm) or the cap elapses. When
  // APIM is already warm (the normal suite, the before-phase) each loop exits on
  // the first fast response, so this stays near-instant.
  for (const path of paths) {
    for (;;) {
      const t0 = Date.now();
      try {
        await fetch(`${base}/${path}`, { headers, signal: AbortSignal.timeout(20_000) });
      } catch {
        // A slow/failed warm-up request still exercised the path.
      }
      if (Date.now() - t0 < WARM_ENOUGH_MS || Date.now() - started >= MAX_WAIT_MS) {
        break;
      }
      await new Promise((r) => setTimeout(r, 500));
    }
  }

  console.log(`APIM Automation API warm-up completed in ${Date.now() - started}ms.`);
}

/**
 * Playwright globalSetup — runs once before all tests.
 * Checks that required infrastructure (APIM, Gateway, K8s cluster) is reachable.
 * If any check fails, the entire test suite aborts immediately with a clear message.
 */
export default async function globalSetup() {
  const configPath = path.resolve(__dirname, "../config.yaml");
  const config = await loadGraviteeConfig(configPath);

  const mapiUrl = config.apim?.baseUrl ?? "http://localhost:30083";
  const gatewayUrl = config.gateway?.baseUrl ?? "http://localhost:30082";

  // Run independent infrastructure checks in parallel
  const results = await Promise.allSettled([
    fetch(mapiUrl, { signal: AbortSignal.timeout(3_000) }),
    fetch(gatewayUrl, { signal: AbortSignal.timeout(3_000) }),
    execFileAsync("kubectl", ["cluster-info"], { timeout: 5_000 }),
  ]);

  const errors: string[] = [];

  if (results[0].status === "rejected") {
    errors.push(`Management API is not reachable at ${mapiUrl}`);
  }
  if (results[1].status === "rejected") {
    errors.push(`Gateway is not reachable at ${gatewayUrl}`);
  }
  if (results[2].status === "rejected") {
    errors.push("kubectl cannot reach a Kubernetes cluster");
  }

  // Check GKO CRDs are installed (requires kubectl to be reachable)
  if (errors.length === 0) {
    try {
      await execFileAsync("kubectl", ["get", "crd", "apiv4definitions.gravitee.io"], {
        timeout: 5_000,
      });
    } catch {
      errors.push(
        "GKO CRDs are not installed (apiv4definitions.gravitee.io not found). " +
          "Install the Gravitee Kubernetes Operator before running E2E tests.",
      );
    }
  }

  // Check GKO operator is running
  if (errors.length === 0) {
    try {
      const { stdout } = await execFileAsync(
        "kubectl",
        ["get", "deploy", "-A", "-l", "app.kubernetes.io/name=gko", "-o", "name"],
        { timeout: 5_000 },
      );
      if (!stdout.trim()) {
        // Fallback: search for any deployment with "gko" or "gravitee-kubernetes-operator" in the name
        const { stdout: fallback } = await execFileAsync(
          "kubectl",
          ["get", "deploy", "-A", "-o", "name"],
          { timeout: 5_000 },
        );
        const hasGko = fallback
          .split("\n")
          .some((line) => /gko|gravitee.*operator/i.test(line));
        if (!hasGko) {
          errors.push(
            "GKO operator deployment not found. " +
              "The CRDs are installed but the operator is not running.",
          );
        }
      }
    } catch {
      errors.push("Failed to check GKO operator deployment status");
    }
  }

  if (errors.length > 0) {
    const msg = [
      "",
      "=".repeat(70),
      " E2E INFRASTRUCTURE CHECK FAILED",
      "=".repeat(70),
      "",
      ...errors.map((e) => `  - ${e}`),
      "",
      " Make sure APIM, Gateway, and a K8s cluster are running",
      " before executing E2E tests.",
      "=".repeat(70),
      "",
    ].join("\n");

    throw new Error(msg);
  }

  console.log("Infrastructure check passed: Management API, Gateway, and K8s cluster are reachable.");

  // Prime the Automation API so the first terraform test doesn't absorb APIM's
  // cold-start cost. Best-effort — never blocks the suite.
  await warmUpApim(mapiUrl, config.apim?.envId ?? "DEFAULT", {
    username: config.apim?.auth?.username ?? "admin",
    password: config.apim?.auth?.password ?? "admin",
  });

  // Ensure the dev-ctx ManagementContext exists before any tests run.
  const ctxFixture = path.resolve(__dirname, "fixtures/management-context/dev-ctx/crd.yaml");
  try {
    await execFileAsync("kubectl", ["apply", "-f", ctxFixture, "-n", "default"], {
      timeout: 10_000,
    });
    console.log("ManagementContext 'dev-ctx' applied successfully.");
  } catch (err) {
    throw new Error(`Failed to apply ManagementContext 'dev-ctx': ${err}`);
  }
}
