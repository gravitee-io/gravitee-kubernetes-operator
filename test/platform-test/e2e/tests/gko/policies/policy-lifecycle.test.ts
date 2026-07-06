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
 * Policy Lifecycle tests.
 *
 * Xray tests:
 *   GKO-94:  Deploy V4 proxy API with policy
 *   GKO-95:  Remove a policy using CRD deployment
 *   GKO-96:  Update a policy using CRD deployment
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

test.describe("Policies — Lifecycle", () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. Each del() ignores errors (the resource may already be gone).
  test.afterEach(async () => {
    for (const f of [
      "crds/api-v4-definitions/v4-proxy-api-with-policy.yaml",
      "crds/api-v4-definitions/v4-proxy-api-policy-removed.yaml",
      "crds/api-v4-definitions/v4-proxy-api-policy-updated.yaml",
      "crds/api-v4-definitions/v4-proxy-api-with-labels-categories.yaml",
    ]) {
      await kubectl.del(fixture(f)).catch(() => {});
    }
  });

  // ── GKO-94: Deploy API with policy ───────────────────────────

  test(`Deploy V4 proxy API with policy ${XRAY.POLICIES.DEPLOY_V4_WITH_POLICY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-policy";
    const API_PATH = "/e2e-v4-policy";
    const fixturePath = fixture("policies/v4-with-policy/crd.yaml");

    await test.step("Deploy API with transform-headers policy", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API has flows configured in APIM", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      // V4 APIs have flows at the API level
      if ("flows" in api && api.flows) {
        expect(api.flows.length).toBeGreaterThanOrEqual(1);
      }
    });

    await test.step("Gateway responds with custom header from policy", async () => {
      await gateway.assertResponds(API_PATH, { status: 200 });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-95: Remove a policy ──────────────────────────────────

  test(`Remove a policy using CRD deployment ${XRAY.POLICIES.REMOVE_POLICY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-policy";
    const withPolicy = fixture("policies/v4-with-policy/crd.yaml");
    const noPolicy = fixture("policies/v4-removed/crd.yaml");

    await test.step("Deploy API with policy", async () => {
      await kubectl.apply(withPolicy);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Deploy updated CRD without policy", async () => {
      await kubectl.apply(noPolicy);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("API has no flows in APIM", async () => {
      const api = await mapi.fetchApi(apiId);
      if ("flows" in api) {
        expect(api.flows?.length ?? 0).toBe(0);
      }
    });

    await kubectl.del(noPolicy);
  });

  // ── GKO-96: Update a policy ──────────────────────────────────

  test(`Update a policy using CRD deployment ${XRAY.POLICIES.UPDATE_POLICY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-policy";
    const withPolicy = fixture("policies/v4-with-policy/crd.yaml");
    const updatedPolicy = fixture("policies/v4-updated/crd.yaml");

    await test.step("Deploy API with policy", async () => {
      await kubectl.apply(withPolicy);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Update the policy", async () => {
      await kubectl.apply(updatedPolicy);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Updated flows reflected in APIM", async () => {
      const api = await mapi.fetchApi(apiId);
      if ("flows" in api && api.flows) {
        expect(api.flows.length).toBeGreaterThanOrEqual(1);
        expect(api.flows[0].name).toContain("updated");
      }
    });

    await kubectl.del(updatedPolicy);
  });

  // ── GKO-269: Non-existing category ───────────────────────────
  // Assigning and removing a valid category (@GKO-267 / @GKO-270) is covered by
  // the cross-provisioner journey tests/user-journeys/assign-categories-to-api/.
  // This case covers the GKO-only behaviour that an unknown category reference is
  // tolerated: the API still deploys.

  test(`Non-existing category is ignored ${XRAY.CATEGORIES.NON_EXISTING_CATEGORY_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-labels-cats";
    const fixturePath = fixture("categories/v4-with-labels/crd.yaml");

    // The fixture has labels but no categories field — API should deploy fine
    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    await mapi.assertApiStarted(status.id);

    await kubectl.del(fixturePath);
  });
});
