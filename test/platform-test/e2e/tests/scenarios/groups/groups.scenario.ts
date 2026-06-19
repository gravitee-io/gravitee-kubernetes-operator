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
 * Group creation, shared across provisioners. A Group created through any
 * provisioner must land in APIM via the Automation API (origin KUBERNETES).
 *
 * This is a deliberately small, param-free scenario: no params.ts, no closures,
 * no gateway/contextPath. The provisioner-specific group behaviour (GKO member
 * reconciliation / admission; Terraform drift, import, data source, validation)
 * stays in the per-provisioner suites under tests/gko/groups and
 * tests/terraform/groups.test.ts.
 */

import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";

forEachProvisioner(
  {
    title: "Group created through a provisioner lands in APIM",
    provisioners: {
      gko: gkoScenario<void>({
        manifests: ["groups/lifecycle/crd.yaml"],
        roles: { group: "e2e-group-simple" },
      }),
      terraform: tfScenario<void>({ fixture: "groups/lifecycle" }),
    },
    xray: { gko: XRAY.GROUPS.CREATE_WITH_MEMBER, terraform: XRAY.TERRAFORM.GROUP_CREATE },
    tags: [TAGS.REGRESSION],
    // The Terraform apim_group resource ships in 4.12; the GKO Group CRD is older,
    // so only the Terraform arm is version-gated.
    since: { terraform: "4.12" },
  },
  async ({ provisioned, mapi }) => {
    // The shared invariant: both provisioners write the group through the
    // Automation API, so APIM records it with origin KUBERNETES.
    const groupId = await provisioned.groupId();
    await mapi.waitForGroupMatchesById(groupId, { origin: "KUBERNETES" });
  },
);
