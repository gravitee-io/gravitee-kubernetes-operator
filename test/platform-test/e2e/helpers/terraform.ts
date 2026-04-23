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
import { cp, mkdir, mkdtemp, rm } from "node:fs/promises";
import { homedir, tmpdir } from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { loadGraviteeConfig } from "../../src/cmd/config.js";
import { fixture } from "../setup.js";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const execFileAsync = promisify(execFile);
const TF_TIMEOUT_MS = 60_000;

export interface TfWorkspace {
  dir: string;
  env: Record<string, string>;
}

/**
 * Resolve the plugin-cache directory Terraform should use across all test
 * workspaces. Each initWorkspace call otherwise creates a fresh temp dir and
 * redownloads the APIM provider from GitHub, which is slow and flaky: it
 * regularly exceeds the 30 s beforeAll timeout and occasionally hits GitHub
 * rate-limiting or transient network errors ("context deadline exceeded").
 *
 * Respects TF_PLUGIN_CACHE_DIR if the caller already set one; otherwise uses
 * the standard ~/.terraform.d/plugin-cache location (created if missing).
 */
async function resolvePluginCacheDir(): Promise<string> {
  const existing = process.env["TF_PLUGIN_CACHE_DIR"];
  if (existing) return existing;
  const dir = path.join(homedir(), ".terraform.d", "plugin-cache");
  await mkdir(dir, { recursive: true });
  return dir;
}

/**
 * Create a Terraform workspace from a fixture directory.
 * Copies the fixture to a temp dir, loads APIM env vars, and runs `terraform init`.
 */
export async function initWorkspace(fixtureName: string): Promise<TfWorkspace> {
  const configPath = path.resolve(__dirname, "../../config.yaml");
  const config = await loadGraviteeConfig(configPath);
  const baseUrl = config.apim?.baseUrl ?? "http://localhost:30083";
  const pluginCache = await resolvePluginCacheDir();

  const env: Record<string, string> = {
    ...(process.env as Record<string, string>),
    APIM_SERVER_URL: `${baseUrl}/automation`,
    APIM_USERNAME: config.apim?.auth?.username ?? "admin",
    APIM_PASSWORD: config.apim?.auth?.password ?? "admin",
    TF_PLUGIN_CACHE_DIR: pluginCache,
  };

  const dir = await mkdtemp(path.join(tmpdir(), "e2e-tf-"));
  await cp(fixture(fixtureName), dir, { recursive: true });
  await tf({ dir, env }, ["init", "-no-color"]);

  return { dir, env };
}

/** Run an arbitrary terraform command in a workspace. */
export async function tf(
  ws: TfWorkspace,
  args: string[],
): Promise<{ stdout: string; stderr: string }> {
  return execFileAsync("terraform", args, {
    cwd: ws.dir,
    env: ws.env,
    timeout: TF_TIMEOUT_MS,
  });
}

/** Run `terraform apply -auto-approve`. Returns stdout. */
export async function apply(ws: TfWorkspace): Promise<string> {
  const { stdout } = await tf(ws, ["apply", "-auto-approve", "-no-color"]);
  return stdout;
}

/** Run `terraform plan -detailed-exitcode`. Handles exit code 2 (changes detected). */
export async function plan(ws: TfWorkspace): Promise<{ stdout: string; hasChanges: boolean }> {
  try {
    const { stdout } = await tf(ws, ["plan", "-detailed-exitcode", "-no-color"]);
    return { stdout, hasChanges: false };
  } catch (err: unknown) {
    const e = err as { code?: number; stdout?: string };
    if (e.code === 2) {
      return { stdout: e.stdout ?? "", hasChanges: true };
    }
    throw err;
  }
}

/** Get a terraform output value by name. */
export async function output(ws: TfWorkspace, name: string): Promise<string> {
  const { stdout } = await tf(ws, ["output", "-raw", name]);
  return stdout.trim();
}

/** Run `terraform destroy -auto-approve`. */
export async function destroy(ws: TfWorkspace): Promise<void> {
  await tf(ws, ["destroy", "-auto-approve", "-no-color"]);
}

/** Destroy resources and remove the temp directory. */
export async function destroyWorkspace(ws: TfWorkspace): Promise<void> {
  await destroy(ws).catch((err: unknown) => {
    console.error(
      `[terraform] destroy failed — APIM resources in workspace "${ws.dir}" may be orphaned.\n`,
      err,
    );
  });
  await rm(ws.dir, { recursive: true, force: true });
}
