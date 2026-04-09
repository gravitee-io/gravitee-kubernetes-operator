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
 * Local ConfigMap behaviour tests.
 *
 * Xray tests:
 *   GKO-765:  Make sure local=false is the default setting
 *   GKO-1452: Validate deletion and finalizer cleanup removes all dependent resources
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, expect, fixture } from "../../../setup.js";
import { XRAY } from "../../../helpers/tags.js";

test.describe("Local ConfigMap", () => {
  test(`Local=false means no ConfigMap is created ${XRAY.LOCAL_CONFIGMAP.LOCAL_FALSE_NO_CONFIGMAP}`, async ({
    kubectl,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-local-false";
    const API_PATH = "/e2e-v4-local-false";
    const FIXTURE = fixture("crds/local-configmap/v4-api-local-false.yaml");

    await test.step("Deploy V4 API with syncFrom MANAGEMENT", async () => {
      await kubectl.apply(FIXTURE);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("ConfigMap should not exist", async () => {
      const cmExists = await kubectl.exists("configmap", API_NAME);
      expect(cmExists).toBe(false);
    });

    await test.step("Gateway responds 200", async () => {
      await gateway.assertResponds(API_PATH, { status: 200 });
    });

    await kubectl.del(FIXTURE);
  });

  test(`Delete when missing in APIM ${XRAY.LOCAL_CONFIGMAP.DELETION_FINALIZER_CLEANUP}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-delete-when-missing";
    const API_PATH = "/e2e-v4-delete-when-missing";
    const FIXTURE = fixture("crds/local-configmap/v4-api-delete-when-missing.yaml");

    await test.step("Deploy V4 API with syncFrom KUBERNETES", async () => {
      await kubectl.apply(FIXTURE);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Gateway responds 200", async () => {
      await gateway.assertResponds(API_PATH, { status: 200 });
    });

    await test.step("Delete API directly from APIM", async () => {
      await mapi.deleteApi(apiId);
    });

    await test.step("Delete CRD from cluster", async () => {
      await kubectl.del(FIXTURE);
    });

    await test.step("ConfigMap should not exist anymore", async () => {
      // Give operator time to clean up
      await kubectl.waitForDeletion("configmap", API_NAME, 30);
    });

    await test.step("Gateway no longer serves traffic", async () => {
      await gateway.assertResponds(API_PATH, { status: 404 });
    });
  });
});
