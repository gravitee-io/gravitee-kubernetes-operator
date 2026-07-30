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

import { describe, it, expect, vi, beforeEach } from "vitest";
import type { K8sCondition } from "../../../src/provisioners/engines/kubectl.js";

const kubectlMocks = vi.hoisted(() => ({
  exists: vi.fn(),
  getCondition: vi.fn(),
  getStatus: vi.fn(),
}));

vi.mock("../../../src/provisioners/engines/kubectl.js", () => kubectlMocks);

const { buildGkoView } = await import("../../../src/provisioners/gko/gko-provisioner.js");

function spec() {
  return {
    manifests: [],
    roles: { api: { kind: "apiv4definition", name: "my-api" } },
    dynamicRoles: [],
  };
}

function condition(overrides: Partial<K8sCondition> = {}): K8sCondition {
  return {
    type: "Accepted",
    status: "True",
    reason: "Ready",
    message: "",
    lastTransitionTime: "2026-01-01T00:00:00Z",
    ...overrides,
  };
}

describe("buildGkoView", () => {
  beforeEach(() => {
    kubectlMocks.exists.mockReset();
    kubectlMocks.getCondition.mockReset();
    kubectlMocks.getStatus.mockReset();
  });

  it("returns gone when the CR does not exist", async () => {
    kubectlMocks.exists.mockResolvedValue(false);
    const result = await buildGkoView(spec()).read("api");
    expect(result.state).toBe("gone");
    expect(kubectlMocks.getCondition).not.toHaveBeenCalled();
  });

  it("returns applied when the condition is True", async () => {
    kubectlMocks.exists.mockResolvedValue(true);
    kubectlMocks.getCondition.mockResolvedValue(condition({ status: "True" }));
    const result = await buildGkoView(spec()).read("api");
    expect(result.state).toBe("applied");
  });

  it("returns failed with severe errors when the condition is False", async () => {
    kubectlMocks.exists.mockResolvedValue(true);
    kubectlMocks.getCondition.mockResolvedValue(
      condition({ status: "False", reason: "ValidationFailed" }),
    );
    kubectlMocks.getStatus.mockResolvedValue({ errors: { severe: ["endpoint target is invalid"] } });
    const result = await buildGkoView(spec()).read("api");
    expect(result.state).toBe("failed");
    expect((result.detail as { severeErrors?: string[] }).severeErrors).toEqual([
      "endpoint target is invalid",
    ]);
  });

  it("throws (never returns an unknown state) when the condition has not been reported yet", async () => {
    kubectlMocks.exists.mockResolvedValue(true);
    kubectlMocks.getCondition.mockResolvedValue(undefined);
    await expect(buildGkoView(spec()).read("api")).rejects.toThrow(/cannot determine state/);
  });

  it("throws (never returns an unknown state) when the condition is Unknown", async () => {
    kubectlMocks.exists.mockResolvedValue(true);
    kubectlMocks.getCondition.mockResolvedValue(condition({ status: "Unknown" }));
    await expect(buildGkoView(spec()).read("api")).rejects.toThrow(/cannot determine state/);
  });
});
