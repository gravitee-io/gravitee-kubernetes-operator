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
 * Deployment & Reconciliation: CRD reconciliation tests.
 *
 * Xray tests:
 *   GKO-1444: Verify custom resource reconciliation applies API configuration correctly
 *     WHEN a CR is applied with a valid API configuration
 *     THEN the operator creates or updates the corresponding API in APIM
 *     AND the APIM API configuration matches the CR spec
 *     AND status conditions on the CR reflect successful reconciliation
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, expect, fixture } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

const API_NAME = "e2e-v4-reconcile";
const API_PATH = "/e2e-v4-reconcile";

test(`Deployment & Reconciliation ${XRAY.DEPLOYMENT_RECONCILIATION.RECONCILE_API_CONFIG} ${TAGS.REGRESSION}`, async ({
  kubectl,
  mapi,
  gateway,
}) => {
  await test.step("Initial CRD apply creates API in APIM", async () => {
    await kubectl.apply(fixture("crds/api-v4-definitions/v4-proxy-api-reconcile.yaml"));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
  });

  const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
  const apiId = status.id;

  await test.step("API is created and reachable", async () => {
    await mapi.waitForApiMatches(apiId, {
      name: API_NAME,
      state: "STARTED",
    });
    await gateway.assertResponds(API_PATH, { status: 200 });
  });

  await test.step("Updated CRD is reconciled in APIM", async () => {
    await kubectl.apply(fixture("crds/api-v4-definitions/v4-proxy-api-reconcile-updated.yaml"));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    await mapi.waitForApiMatches(apiId, {
      name: "e2e-v4-reconcile-updated",
    });
  });

  await test.step("Status conditions reflect successful reconciliation", async () => {
    const status = await kubectl.getStatus<{
      conditions?: Array<{ type: string; status: string }>;
    }>("apiv4definition", API_NAME);

    expect(status.conditions).toBeDefined();
    const acceptedCondition = status.conditions!.find((c) => c.type === "Accepted");
    expect(acceptedCondition).toBeDefined();
    expect(acceptedCondition!.status).toBe("True");
  });

  await kubectl.del(fixture("crds/api-v4-definitions/v4-proxy-api-reconcile.yaml"));
});
