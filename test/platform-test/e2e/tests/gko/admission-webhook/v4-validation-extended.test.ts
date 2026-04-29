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
 * V4 API Validation & Reconciliation — Extended scenarios.
 *
 * Xray tests:
 *   GKO-1476: Context path conflict validation on V4 APIs
 *   GKO-1479: OAS compliance webhook rejects invalid endpoints
 *   GKO-1480: Default values are applied during reconciliation
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("V4 API Validation & Reconciliation — Extended", () => {
  // ── GKO-1476: Context path conflict ─────────────────────────

  test(`Context path conflict between two V4 APIs ${XRAY.VALIDATION.V4_CONTEXT_PATH_CONFLICT} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const firstFixture = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");
    const conflictFixture = fixture("crds/v4-lifecycle-extended/v4-proxy-api-started-conflict.yaml");

    await kubectl.apply(firstFixture);
    await kubectl.waitForCondition("apiv4definition", "e2e-v4-start-stop", "Accepted");

    try {
      await test.step("Second API with same path is rejected", async () => {
        const stderr = await kubectl.applyExpectFailure(conflictFixture);
        expect(stderr.toLowerCase()).toMatch(/context.?path|already exists|denied/);
      });
    } finally {
      // Safety-net: delete the second CRD if admission unexpectedly accepted it,
      // so it doesn't leak into downstream tests.
      await kubectl.del(conflictFixture).catch(() => {});
      await kubectl.del(firstFixture).catch(() => {});
    }
  });

  // ── GKO-1479: OAS compliance webhook rejection ─────────────

  test(`Non-OAS-compliant V4 API is rejected by the webhook ${XRAY.VALIDATION.V4_OAS_COMPLIANCE_WEBHOOK} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-invalid.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/denied|invalid|error/);
  });

  // ── GKO-1480: Default values applied during reconcile ──────

  test(`Operator applies default values during reconciliation ${XRAY.VALIDATION.V4_DEFAULT_VALUES} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const api = await mapi.fetchApi(status.id);

    // Fields not declared in the CRD are still present in APIM (defaulted).
    expect(api.visibility).toBe("PRIVATE");
    expect(api.lifecycleState).toBeTruthy();
    expect(api.createdAt).toBeTruthy();

    await kubectl.del(fixturePath);
  });
});
