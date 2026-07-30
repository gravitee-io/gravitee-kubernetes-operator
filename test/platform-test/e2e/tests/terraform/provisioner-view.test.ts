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
 * Terraform / HCL: `ProvisionerView` detects a tainted resource.
 *
 * This is the permanent record of the empirical spike behind
 * `buildTerraformView` (src/provisioners/terraform/terraform-provisioner.ts):
 * `terraform show -json`'s STATE representation (no plan needed) exposes a
 * `tainted: true` field directly on a tainted resource, and the field is
 * simply absent otherwise. Confirmed live against `terraform taint`/`untaint`
 * with Terraform 1.14.8; the pinned CI version is 1.12.1 — if that pin moves,
 * re-run this test to confirm the field is still exposed the same way.
 *
 * Xray tests:
 *   GKO-TBD-tf-view-tainted: Terraform ProvisionerView reports "failed" for a tainted resource
 *     WHEN a resource is provisioned and then `terraform taint`ed
 *     THEN `assertProvisioner(provisioned, "api", "failed")` passes
 *     AND untainting it makes `assertProvisioner(provisioned, "api", "applied")` pass again
 *
 * Preconditions:
 *   - APIM and Gateway are running
 *   - terraform CLI is installed
 */

import { test } from "../../setup.js";
import { XRAY, TAGS } from "../../helpers/tags.js";
import { fixture } from "../../setup.js";
import { terraformEnv } from "../../helpers/terraform.js";
import {
  TerraformProvisioner,
  assertProvisioner,
  isTerraform,
  type Provisioned,
} from "../../../src/provisioners/index.js";
import * as terraform from "../../helpers/terraform.js";

let provisioned: Provisioned<void>;

test.describe("Terraform — ProvisionerView", () => {
  test.beforeAll(async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    const provisioner = new TerraformProvisioner<void>({
      fixtureDir: fixture("subscriptions/v4-full-stack"),
      env: await terraformEnv(),
      addresses: { api: "apim_apiv4.e2e_tf_test" },
    });
    provisioned = await provisioner.provision();
  });

  test.afterAll(async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    if (provisioned) await provisioned.destroy();
  });

  test(`terraform view reports failed for a tainted resource ${XRAY.TERRAFORM.VIEW_DETECTS_TAINTED_RESOURCE} ${TAGS.REGRESSION}`, async () => {
    await assertProvisioner(provisioned, "api", "applied");

    if (!isTerraform(provisioned.checks)) {
      throw new Error("expected a Terraform-provisioned handle");
    }
    await provisioned.checks.taint("api");
    await assertProvisioner(provisioned, "api", "failed");

    // Restore to a clean state so the afterAll destroy() is a normal teardown.
    await provisioned.checks.untaint("api");
    await assertProvisioner(provisioned, "api", "applied");
  });
});
