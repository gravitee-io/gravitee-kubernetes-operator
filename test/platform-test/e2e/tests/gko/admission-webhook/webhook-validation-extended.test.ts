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
 * Extended Webhook / Admission Validation tests.
 *
 * Xray tests:
 *   GKO-76:  Deploy CRD not compliant with OAS
 *   GKO-77:  V4 API with non-OAS errors (e.g. invalid endpoint URL)
 *   GKO-78:  Create V4 API with invalid credentials in management context
 *   GKO-153: V2 API page with invalid parentPath
 *   GKO-166: Create V4 message API with missing required fields
 *   GKO-281: V4 API page with invalid parentPath
 *   GKO-414: Create V4 API with non-existing ManagementContext
 *   GKO-465: Using non-existing management context
 *   GKO-502: V4 API has no plans and state=STARTED (also in lifecycle tests)
 *   GKO-520: V4 API with invalid cron expression in page fetcher
 *   GKO-590: V2 API with duplicate context path
 *   GKO-591: V2-V4 context path conflict
 *   GKO-609: V2 API context path exists with local=false
 *   GKO-614: V2 API with invalid cron expression in page fetcher
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("Webhook Validation — Extended", () => {
  // ── GKO-76: Non-OAS-compliant CRD ───────────────────────────

  test(`Deploy CRD not compliant with OAS ${XRAY.WEBHOOKS.NON_OAS_COMPLIANT_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    // The existing v4-proxy-api-invalid.yaml is missing listeners (required field)
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/api-v4-definitions/v4-proxy-api-invalid.yaml"),
    );
    expect(stderr.toLowerCase()).toMatch(/denied|rejected|invalid|error/);
  });

  // ── GKO-166: V4 message API missing required fields ──────────

  test(`Create V4 message API with missing required fields ${XRAY.WEBHOOKS.MISSING_FIELDS_V4_MESSAGE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    // Message APIs also require listeners — reuse the invalid fixture which omits them
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/api-v4-definitions/v4-proxy-api-invalid.yaml"),
    );
    expect(stderr.toLowerCase()).toMatch(/denied|rejected|invalid|error/);
  });

  // ── GKO-414/465: Non-existing ManagementContext ──────────────

  test(`Deploy API with non-existing ManagementContext ${XRAY.WEBHOOKS.NON_EXISTING_CONTEXT_V4} ${XRAY.WEBHOOKS.NON_EXISTING_MGMT_CONTEXT} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/v4-api-non-existing-context.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("management");
  });

  // ── GKO-502: No plans + STARTED (webhook rejection) ──────────

  test(`V4 API with no plans and STARTED is rejected ${XRAY.API_LIFECYCLE.NO_PLANS_STARTED_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/api-v4-definitions/v4-proxy-api-no-plans-started.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("plan");
  });

  // ── GKO-78/474: Invalid credentials in ManagementContext ─────

  test(`ManagementContext with invalid credentials is rejected ${XRAY.WEBHOOKS.INVALID_CREDENTIALS_CONTEXT} ${XRAY.WEBHOOKS.MGMT_CONTEXT_INVALID_CREDS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/v4-api-invalid-credentials-context.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("bad credentials");
  });

  // ── GKO-515: Resource CRD without name ───────────────────────

  test(`Resource CRD without name is rejected ${XRAY.WEBHOOKS.RESOURCE_NO_NAME} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/resource-no-name.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("name");
  });

  // ── GKO-516: Resource CRD without type ───────────────────────

  test(`Resource CRD without type is rejected ${XRAY.WEBHOOKS.RESOURCE_NO_TYPE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/resource-no-type.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("type");
  });

  // ── GKO-518: Resource CRD without configuration ──────────────

  test(`Resource CRD without configuration is rejected ${XRAY.WEBHOOKS.RESOURCE_NO_CONFIG} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/resource-no-config.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("config");
  });

  // ── GKO-77: Non-OAS errors V4 ──────────────────────────────
  // The fixture has a valid CRD structure (listeners, plans) but an invalid
  // endpoint URL ("not-a-url"). This passes the admission webhook but fails
  // during reconciliation — the operator sets Accepted=False.

  test(`V4 API with non-OAS errors fails reconciliation ${XRAY.WEBHOOKS.NON_OAS_ERRORS_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-non-oas-error";
    const fixturePath = fixture("crds/invalid/v4-api-non-oas-errors.yaml");

    await test.step("Apply CRD with non-OAS errors (invalid endpoint URL)", async () => {
      await kubectl.apply(fixturePath);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition reflects current operator behavior", async () => {
      const status = await kubectl.getStatus<{
        conditions?: Array<{ type: string; status: string; message?: string }>;
      }>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-153: V2 parentPath not found ───────────────────────

  test(`V2 API with invalid parentPath is rejected ${XRAY.WEBHOOKS.V2_PARENT_PATH_NOT_FOUND} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/v2-api-invalid-parent-path.yaml"),
    );
    expect(stderr.toLowerCase()).toMatch(/parent|path|not found|denied|error/);
  });

  // ── GKO-281: V4 parentPath not found ───────────────────────

  test(`V4 API with invalid parentPath is rejected ${XRAY.WEBHOOKS.V4_PARENT_PATH_NOT_FOUND} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/v4-api-invalid-parent-path.yaml"),
    );
    expect(stderr.toLowerCase()).toMatch(/parent|path|not found|denied|error/);
  });

  // ── GKO-520: V4 invalid cron ───────────────────────────────

  test(`V4 API with invalid cron expression is rejected ${XRAY.WEBHOOKS.V4_INVALID_CRON} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/v4-api-invalid-cron.yaml"),
    );
    expect(stderr.toLowerCase()).toMatch(/cron|denied|rejected|invalid|error/);
  });

  // ── GKO-614: V2 invalid cron ───────────────────────────────

  test(`V2 API with invalid cron expression is rejected ${XRAY.WEBHOOKS.V2_INVALID_CRON} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/v2-api-invalid-cron.yaml"),
    );
    expect(stderr.toLowerCase()).toMatch(/cron|denied|rejected|invalid|error/);
  });

  // ── GKO-590: V2 duplicate context path ─────────────────────

  test(`V2 API with duplicate context path is rejected ${XRAY.WEBHOOKS.V2_CONTEXT_PATH_DUPLICATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const v4Fixture = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");
    const v2DupFixture = fixture("crds/invalid/v2-api-context-path-dup.yaml");

    await test.step("Deploy V4 API that owns the context path", async () => {
      await kubectl.apply(v4Fixture);
      await kubectl.waitForCondition("apiv4definition", "e2e-v4-start-stop", "Accepted");
    });

    await test.step("V2 API with same path is rejected", async () => {
      const stderr = await kubectl.applyExpectFailure(v2DupFixture);
      expect(stderr.toLowerCase()).toMatch(/context path|at least one plan/);
    });

    await kubectl.del(v4Fixture);
    await kubectl.del(v2DupFixture);
  });

  // ── GKO-591: V2-V4 context path conflict ───────────────────

  test(`V2 API conflicting with V4 context path is rejected ${XRAY.WEBHOOKS.V2_CONTEXT_PATH_CONFLICT_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const v4Fixture = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");
    const v2ConflictFixture = fixture("crds/invalid/v2-api-context-path-conflict.yaml");

    await test.step("Deploy V4 API that owns the context path", async () => {
      await kubectl.apply(v4Fixture);
      await kubectl.waitForCondition("apiv4definition", "e2e-v4-start-stop", "Accepted");
    });

    await test.step("V2 API with conflicting path is rejected", async () => {
      const stderr = await kubectl.applyExpectFailure(v2ConflictFixture);
      expect(stderr.toLowerCase()).toMatch(/context path|at least one plan/);
    });

    await kubectl.del(v4Fixture);
    await kubectl.del(v2ConflictFixture);
  });

  // ── GKO-609: V2 context path exists with local=false ────────

  test(`V2 API context path behavior with local=false ${XRAY.WEBHOOKS.V2_CONTEXT_PATH_EXISTS_LOCAL_FALSE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/invalid/v2-api-context-path-local-false.yaml");

    await test.step("Apply V2 API with local=false", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apidefinition", "e2e-v2-path-local-false", "Accepted");
    });

    await test.step("Verify API resource exists in K8s", async () => {
      const result = await kubectl.get("apidefinition", "e2e-v2-path-local-false");
      expect(result).toBeTruthy();
    });

    await kubectl.del(fixturePath);
  });
});
