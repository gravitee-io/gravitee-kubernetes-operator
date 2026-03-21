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
 * Terraform / HCL: Post-apply verification and idempotency tests.
 *
 * Xray tests:
 *   GKO-1926: Terraform apply creates the complex API configuration successfully
 *     WHEN terraform init && terraform apply is run
 *     THEN Terraform completes without errors AND the API exists in APIM
 *     AND is reachable on the gateway
 *
 *   GKO-1932: Terraform configuration is idempotent for complex configuration
 *     WHEN terraform plan is run after a successful apply
 *     THEN no changes are detected AND no resources are recreated or modified
 *
 * Preconditions:
 *   - APIM and Gateway are running
 *   - terraform CLI is installed
 */

import { test, expect } from "../../setup.js";
import { XRAY, TAGS } from "../../helpers/tags.js";
import * as terraform from "../../helpers/terraform.js";
import type { TfWorkspace } from "../../helpers/terraform.js";

let ws: TfWorkspace;
let apiId: string;
const apiPath = "/e2e-tf-post-apply/";

test.describe("Terraform — Post-apply & Idempotency", () => {
  test.beforeAll(async () => {
    ws = await terraform.initWorkspace("terraform");
    await terraform.apply(ws);
    apiId = await terraform.output(ws, "api_id");
  });

  test.afterAll(async () => {
    if (ws) await terraform.destroyWorkspace(ws);
  });

  test(`Terraform apply creates API in APIM ${XRAY.TERRAFORM.APPLY_COMPLEX_CONFIG} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    await mapi.assertApiMatches(apiId, { state: "STARTED" });
  });

  test(`API is reachable on the gateway ${XRAY.TERRAFORM.APPLY_COMPLEX_CONFIG} ${TAGS.REGRESSION}`, async ({
    gateway,
  }) => {
    await gateway.assertResponds(apiPath, { status: 200 });
  });

  test(`terraform plan shows no changes ${XRAY.TERRAFORM.IDEMPOTENT_CONFIG} ${TAGS.REGRESSION}`, async () => {
    const result = await terraform.plan(ws);
    expect(result.hasChanges).toBe(false);
  });
});
