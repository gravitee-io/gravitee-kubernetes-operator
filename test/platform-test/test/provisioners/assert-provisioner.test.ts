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

import { describe, it, expect, vi } from "vitest";
import { assertProvisioner } from "../../src/provisioners/view.js";
import type { Provisioned } from "../../src/provisioners/types.js";
import type { ProvisionerViewResult } from "../../src/provisioners/view.js";

/** Minimal handle: assertProvisioner only touches `view.read` and `provisionerId`. */
function handleReading(read: () => Promise<ProvisionerViewResult>): Provisioned<unknown> {
  return { provisionerId: "gko", view: { read } } as unknown as Provisioned<unknown>;
}

const FAST = { timeoutMs: 500, intervalMs: 10 };

describe("assertProvisioner", () => {
  it("passes once the record reaches the expected state", async () => {
    const read = vi.fn().mockResolvedValue({ state: "applied", detail: {} });
    await assertProvisioner(handleReading(read), "api", "applied", FAST);
    expect(read).toHaveBeenCalledTimes(1);
  });

  it("retries a state that has not settled yet", async () => {
    // The reason this polls at all: after update() the provisioner does not wait,
    // so the first read can legitimately still show the pre-update record.
    const read = vi
      .fn()
      .mockResolvedValueOnce({ state: "gone", detail: {} })
      .mockResolvedValueOnce({ state: "gone", detail: {} })
      .mockResolvedValue({ state: "applied", detail: {} });
    await assertProvisioner(handleReading(read), "api", "applied", FAST);
    expect(read).toHaveBeenCalledTimes(3);
  });

  it("retries while read() throws because it cannot yet tell", async () => {
    // GKO's read() throws until the CR reports the condition at all; that is a
    // not-yet answer, not a failure.
    const read = vi
      .fn()
      .mockRejectedValueOnce(new Error('has no "Accepted" condition yet'))
      .mockResolvedValue({ state: "applied", detail: {} });
    await assertProvisioner(handleReading(read), "api", "applied", FAST);
    expect(read).toHaveBeenCalledTimes(2);
  });

  it("reports the observed state and detail when it never settles", async () => {
    const read = vi
      .fn()
      .mockResolvedValue({ state: "failed", detail: { severeErrors: ["boom"] } });

    const err = await assertProvisioner(handleReading(read), "api", "applied", FAST).then(
      () => undefined,
      (e: Error) => e,
    );
    if (!err) throw new Error("expected assertProvisioner to reject");

    expect(err.message).toMatch(/gko record for role "api" to be "applied"/);
    // The provisioner's own evidence must survive the timeout wrapper, or a failed
    // provision is indistinguishable from a slow one when triaging.
    const cause = (err as { cause?: Error }).cause;
    expect(cause?.message).toMatch(/got "failed".*severeErrors.*boom/);
  });
});
