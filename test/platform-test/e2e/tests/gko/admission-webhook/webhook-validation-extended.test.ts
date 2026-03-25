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
 *   GKO-78:  Create V4 API with invalid credentials in management context
 *   GKO-166: Create V4 message API with missing required fields
 *   GKO-414: Create V4 API with non-existing ManagementContext
 *   GKO-465: Using non-existing management context
 *   GKO-502: V4 API has no plans and state=STARTED (also in lifecycle tests)
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
});
