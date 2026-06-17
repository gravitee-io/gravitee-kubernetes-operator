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

const execFileAsync = promisify(execFile);

const NAMESPACE = process.env["K8S_NAMESPACE"] ?? "default";

interface ExecResult {
  stdout: string;
  stderr: string;
}

const EXEC_TIMEOUT_MS = 15_000;

async function run(args: string[], timeoutMs = EXEC_TIMEOUT_MS): Promise<ExecResult> {
  return execFileAsync("kubectl", args, { timeout: timeoutMs });
}

/** Apply a YAML manifest file. */
export async function apply(yamlPath: string, namespace = NAMESPACE): Promise<void> {
  await run(["apply", "-f", yamlPath, "-n", namespace]);
}

/** Server-side dry-run apply — validates through admission without persisting. */
export async function applyDryRun(yamlPath: string, namespace = NAMESPACE): Promise<void> {
  await run(["apply", "-f", yamlPath, "-n", namespace, "--dry-run=server"]);
}

/** Delete resources defined in a YAML manifest file. */
export async function del(yamlPath: string, namespace = NAMESPACE): Promise<void> {
  await run(["delete", "-f", yamlPath, "-n", namespace, "--ignore-not-found"]);
}

/** Delete a single resource by kind and name. */
export async function deleteResource(
  kind: string,
  name: string,
  namespace = NAMESPACE,
): Promise<void> {
  await run(["delete", kind, name, "-n", namespace, "--ignore-not-found"]);
}

/** Wait for a condition on a resource. */
export async function waitForCondition(
  kind: string,
  name: string,
  condition: string,
  timeoutSeconds = 60,
  namespace = NAMESPACE,
): Promise<void> {
  // Use kubectl's own timeout + a small buffer for the exec timeout
  const execTimeoutMs = (timeoutSeconds + 5) * 1_000;
  await run(
    [
      "wait",
      `--for=condition=${condition}`,
      `${kind}/${name}`,
      `--timeout=${timeoutSeconds}s`,
      "-n",
      namespace,
    ],
    execTimeoutMs,
  );
}

/** Get the full resource as parsed JSON. */
export async function get<T = unknown>(
  kind: string,
  name: string,
  namespace = NAMESPACE,
): Promise<T> {
  const { stdout } = await run(["get", `${kind}/${name}`, "-n", namespace, "-o", "json"]);
  return JSON.parse(stdout) as T;
}

/** Get just the .status field of a resource. */
export async function getStatus<T = unknown>(
  kind: string,
  name: string,
  namespace = NAMESPACE,
): Promise<T> {
  const resource = await get<{ status: T }>(kind, name, namespace);
  return resource.status;
}

/** Assert that a Kubernetes event for the given resource contains a message substring. */
export async function assertEventContains(
  kind: string,
  name: string,
  message: string,
  namespace = NAMESPACE,
): Promise<void> {
  const { stdout } = await run([
    "get",
    "events",
    `--field-selector=involvedObject.name=${name}`,
    "-n",
    namespace,
    "-o",
    "json",
  ]);
  const events = JSON.parse(stdout) as { items: Array<{ message: string }> };
  const found = events.items.some((e) => e.message.includes(message));
  if (!found) {
    throw new Error(
      `No event for ${kind}/${name} contains "${message}". ` +
        `Events: ${events.items.map((e) => e.message).join("; ")}`,
    );
  }
}

/** Extract a specific field using jsonpath. */
export async function getField<T = string>(
  kind: string,
  name: string,
  jsonpath: string,
  namespace = NAMESPACE,
): Promise<T> {
  const { stdout } = await run([
    "get",
    `${kind}/${name}`,
    "-n",
    namespace,
    "-o",
    `jsonpath=${jsonpath}`,
  ]);
  try {
    return JSON.parse(stdout) as T;
  } catch {
    return stdout as unknown as T;
  }
}

/** Poll until a resource is deleted (no longer found). */
export async function waitForDeletion(
  kind: string,
  name: string,
  timeoutSeconds = 60,
  namespace = NAMESPACE,
): Promise<void> {
  const deadline = Date.now() + timeoutSeconds * 1_000;
  while (Date.now() < deadline) {
    try {
      await run(["get", `${kind}/${name}`, "-n", namespace]);
    } catch {
      return; // Resource not found = deleted
    }
    await new Promise((r) => setTimeout(r, 1_000));
  }
  throw new Error(`${kind}/${name} still exists after ${timeoutSeconds}s`);
}

/** Try to delete and expect it to fail (e.g., webhook blocks deletion). Returns stderr. */
export async function delExpectFailure(
  yamlPath: string,
  namespace = NAMESPACE,
): Promise<string> {
  try {
    await run(["delete", "-f", yamlPath, "-n", namespace]);
    throw new Error(`Expected kubectl delete to fail for ${yamlPath}, but it succeeded`);
  } catch (err: unknown) {
    if (err != null && typeof err === "object" && "stderr" in err && typeof (err as Record<string, unknown>).stderr === "string") {
      return (err as Record<string, unknown>).stderr as string;
    }
    throw err;
  }
}

/** Check if a resource exists (returns true/false without throwing). */
export async function exists(
  kind: string,
  name: string,
  namespace = NAMESPACE,
): Promise<boolean> {
  try {
    await run(["get", `${kind}/${name}`, "-n", namespace]);
    return true;
  } catch {
    return false;
  }
}

/** Trigger a rolling restart of a workload (deployment, daemonset, statefulset). */
export async function rolloutRestart(
  kind: string,
  name: string,
  namespace = NAMESPACE,
): Promise<void> {
  await run(["rollout", "restart", `${kind}/${name}`, "-n", namespace]);
}

/** Wait for a rolling update to finish. */
export async function waitForRollout(
  kind: string,
  name: string,
  timeoutSeconds = 120,
  namespace = NAMESPACE,
): Promise<void> {
  const execTimeoutMs = (timeoutSeconds + 5) * 1_000;
  await run(
    [
      "rollout",
      "status",
      `${kind}/${name}`,
      `--timeout=${timeoutSeconds}s`,
      "-n",
      namespace,
    ],
    execTimeoutMs,
  );
}

/** Apply a YAML string via stdin (useful for dynamically generated manifests). */
export async function applyString(
  yamlContent: string,
  namespace = NAMESPACE,
): Promise<void> {
  return new Promise<void>((resolve, reject) => {
    const child = execFile(
      "kubectl",
      ["apply", "-f", "-", "-n", namespace],
      { timeout: EXEC_TIMEOUT_MS },
      (err: Error | null) => (err ? reject(err) : resolve()),
    );
    child.stdin?.end(yamlContent);
  });
}

/**
 * Try to apply a YAML string via stdin and expect it to fail.
 * Returns the stderr output for further assertions.
 * Throws if the apply unexpectedly succeeds.
 */
export async function applyStringExpectFailure(
  yamlContent: string,
  namespace = NAMESPACE,
): Promise<string> {
  return new Promise<string>((resolve, reject) => {
    const child = execFile(
      "kubectl",
      ["apply", "-f", "-", "-n", namespace],
      { timeout: EXEC_TIMEOUT_MS },
      (err: Error | null, _stdout, stderr) => {
        if (err) {
          resolve(stderr ?? "");
        } else {
          reject(new Error("Expected kubectl apply to fail, but it succeeded"));
        }
      },
    );
    child.stdin?.end(yamlContent);
  });
}

/**
 * Try to apply a manifest and expect it to fail (e.g., admission webhook rejection).
 * Returns the stderr output for further assertions.
 * Throws if the apply unexpectedly succeeds.
 */
export async function applyExpectFailure(
  yamlPath: string,
  namespace = NAMESPACE,
): Promise<string> {
  try {
    await run(["apply", "-f", yamlPath, "-n", namespace]);
    throw new Error(`Expected kubectl apply to fail for ${yamlPath}, but it succeeded`);
  } catch (err: unknown) {
    if (err != null && typeof err === "object" && "stderr" in err && typeof (err as Record<string, unknown>).stderr === "string") {
      return (err as Record<string, unknown>).stderr as string;
    }
    throw err;
  }
}
