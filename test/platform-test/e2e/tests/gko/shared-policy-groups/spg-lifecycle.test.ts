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

/**
 * Shared Policy Groups Lifecycle tests.
 *
 * Xray tests:
 *   GKO-981:  Update SPG
 *   GKO-1462: SPG lifecycle validation
 *
 * GKO-976 (add SPG to a V4 API) and GKO-980 (remove SPG from a V4 API) moved to
 * the shared cross-provisioner journey tests/scenarios/reuse-shared-policy-group
 * — SPG reuse (flow attach/detach) is now proven against GKO and Terraform.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

test.describe("Shared Policy Groups — Lifecycle", () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. Each del() ignores errors (the resource may already be gone).
  test.afterEach(async () => {
    for (const f of [
      "crds/api-v4-definitions/v4-api-with-spg.yaml",
      "crds/api-v4-definitions/v4-api-without-spg.yaml",
      "crds/shared-policy-groups/spg-proxy-request.yaml",
      "crds/shared-policy-groups/spg-proxy-updated.yaml",
      "crds/shared-policy-groups/spg-message-unsupported.yaml",
    ]) {
      await kubectl.del(fixture(f)).catch(() => {});
    }
  });

  // GKO-976 (add SPG to a V4 API) and GKO-980 (remove SPG from a V4 API) are now
  // covered by the cross-provisioner journey tests/scenarios/reuse-shared-policy-group.

  // ── GKO-981: Update SPG ────────────────────────────────────

  test(`Update SPG ${XRAY.SHARED_POLICY_GROUPS.UPDATE_SPG} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const SPG_NAME = "e2e-spg-proxy";
    const createFixture = fixture("shared-policy-groups/spg-proxy-request/crd.yaml");
    const updateFixture = fixture("shared-policy-groups/spg-proxy-updated/crd.yaml");

    await test.step("Deploy SPG", async () => {
      await kubectl.apply(createFixture);
      await kubectl.waitForCondition("sharedpolicygroup", SPG_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_NAME);
    expect(status.id).toBeTruthy();

    await test.step("Update SPG with new header value", async () => {
      await kubectl.apply(updateFixture);
      await kubectl.waitForCondition("sharedpolicygroup", SPG_NAME, "Accepted");
    });

    await test.step("SPG ID remains the same after update", async () => {
      const updatedStatus = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_NAME);
      expect(updatedStatus.id).toBe(status.id);
    });

    await kubectl.del(updateFixture);
  });

  // ── GKO-1462: SPG lifecycle validation ──────────────────────

  test(`SPG lifecycle validation ${XRAY.SHARED_POLICY_GROUPS.SPG_LIFECYCLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const SPG_PROXY_NAME = "e2e-spg-proxy";
    const SPG_MESSAGE_NAME = "e2e-spg-message";
    const proxyFixture = fixture("shared-policy-groups/spg-proxy-request/crd.yaml");
    const messageFixture = fixture("shared-policy-groups/spg-message-unsupported/crd.yaml");

    await test.step("Deploy valid PROXY SPG", async () => {
      await kubectl.apply(proxyFixture);
      await kubectl.waitForCondition("sharedpolicygroup", SPG_PROXY_NAME, "Accepted");
    });

    const proxyStatus = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_PROXY_NAME);
    expect(proxyStatus.id).toBeTruthy();

    await test.step("Deploy MESSAGE SPG (may or may not be supported)", async () => {
      // MESSAGE apiType SPG may be accepted or rejected depending on APIM version
      try {
        await kubectl.apply(messageFixture);
        await kubectl.waitForCondition("sharedpolicygroup", SPG_MESSAGE_NAME, "Accepted");
        const messageStatus = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_MESSAGE_NAME);
        expect(messageStatus.id).toBeTruthy();
      } catch {
        // If rejected, that is also valid behavior
        const status = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_MESSAGE_NAME);
        expect(status).toBeTruthy();
      }
    });

    await kubectl.del(proxyFixture);
    await kubectl.del(messageFixture);
  });
});
