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

const tfCoreMocks = vi.hoisted(() => ({
  showState: vi.fn(),
}));

vi.mock("../../../src/provisioners/engines/terraform-core.js", () => tfCoreMocks);

const { buildTerraformView } = await import(
  "../../../src/provisioners/terraform/terraform-provisioner.js"
);

const ws = { dir: "/tmp/does-not-matter", env: {} };

function spec() {
  return { fixtureDir: "/tmp/fixture", env: {}, addresses: { api: "apim_apiv4.test" } };
}

describe("buildTerraformView", () => {
  beforeEach(() => {
    tfCoreMocks.showState.mockReset();
  });

  it("returns gone when the address is absent from state", async () => {
    tfCoreMocks.showState.mockResolvedValue({ resources: [] });
    const result = await buildTerraformView(ws, spec()).read("api");
    expect(result.state).toBe("gone");
  });

  it("returns applied when the resource is present and not tainted", async () => {
    tfCoreMocks.showState.mockResolvedValue({
      resources: [{ address: "apim_apiv4.test" }],
    });
    const result = await buildTerraformView(ws, spec()).read("api");
    expect(result.state).toBe("applied");
  });

  it("returns failed when the resource is present and tainted", async () => {
    tfCoreMocks.showState.mockResolvedValue({
      resources: [{ address: "apim_apiv4.test", tainted: true }],
    });
    const result = await buildTerraformView(ws, spec()).read("api");
    expect(result.state).toBe("failed");
  });

  it("throws (never returns an unknown state) when no address is declared for the role", async () => {
    tfCoreMocks.showState.mockResolvedValue({ resources: [] });
    await expect(buildTerraformView(ws, { fixtureDir: "/tmp/fixture", env: {} }).read("api")).rejects.toThrow(
      /no resource address mapped/,
    );
    expect(tfCoreMocks.showState).not.toHaveBeenCalled();
  });
});
