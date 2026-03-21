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
 * Webhooks: Admission webhook validation tests.
 *
 * Xray tests:
 *   GKO-1447: Validate webhook rejects invalid/unsupported CRs (general)
 *     WHEN an invalid or unsupported CR is submitted
 *     THEN the admission webhook rejects it with a clear error message
 *     AND the resource is never created in APIM
 *     Test cases: missing required fields, invalid enum values, deprecated fields
 *
 * Preconditions:
 *   - GKO operator is running with admission webhooks enabled
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, expect, fixture } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

const INVALID_FIXTURE = fixture("crds/api-v4-definitions/v4-proxy-api-invalid.yaml");

test.describe("Webhooks — Admission validation", () => {
  test.afterAll(async ({ kubectl }) => {
    // Guard: clean up in case the webhook failed to reject and the resource was created
    await kubectl.del(INVALID_FIXTURE);
  });

  test(`Webhook rejects invalid CR and does not create API in APIM ${XRAY.WEBHOOKS.REJECT_INVALID_CRS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    let stderr: string;

    await test.step("Apply invalid CR — expect webhook rejection", async () => {
      stderr = await kubectl.applyExpectFailure(INVALID_FIXTURE);
      expect(stderr.toLowerCase()).toMatch(/denied|rejected|invalid|error/);
    });

    await test.step("Verify API was not created in APIM", async () => {
      await mapi.assertApiHttpStatus("e2e-v4-invalid-api", 404);
    });
  });
});
