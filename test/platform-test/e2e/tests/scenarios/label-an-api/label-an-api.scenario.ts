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
 * Use case: label a V4 API.
 *
 * Labels are an inline attribute of apim_apiv4 (no standalone Terraform
 * resource), yet the journey is fully expressible through both provisioners:
 * labels set through either driver land in APIM and are removed when stripped.
 * This is the reference pattern for the other inline apim_apiv4 attributes
 * (categories, groups, metadata, inline pages[]) that have no standalone
 * Terraform resource.
 *
 * Fixtures live in fixtures/use-cases/label-an-api/.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";

const FIXTURE = "use-cases/label-an-api";
const GKO_WITH_LABELS = fixture(`${FIXTURE}/gko/api-with-labels.yaml`);
const GKO_WITHOUT_LABELS = fixture(`${FIXTURE}/gko/api-without-labels.yaml`);
const LABELS = ["e2e-uc-label-1", "e2e-uc-label-2"];

/** The single knob: whether the API carries labels. */
interface LabelParams {
  withLabels: boolean;
}

forEachProvisioner<LabelParams>(
  {
    title: "Label a V4 API",
    provisioners: {
      gko: gkoScenario<LabelParams>({
        manifests: [],
        roles: { api: "e2e-uc-label" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          await k.apply(params.withLabels ? GKO_WITH_LABELS : GKO_WITHOUT_LABELS);
        },
      }),
      terraform: tfScenario<LabelParams>({
        fixture: `${FIXTURE}/terraform`,
        toVars: (params) => ({ with_labels: params.withLabels }),
      }),
    },
    xray: {
      gko: XRAY.CATEGORIES.V4_LABELS_LIFECYCLE,
      terraform: XRAY.TERRAFORM.API_LABELS_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    await test.step("Labels set through the provisioner land in APIM", async () => {
      await expect
        .poll(async () => ((await mapi.fetchApi(apiId)).labels ?? []).slice().sort(), {
          timeout: 30_000,
          message: "API labels reach APIM",
        })
        .toEqual([...LABELS].sort());
    });

    await test.step("Stripping the labels removes them in APIM", async () => {
      await provisioned.update({ withLabels: false });
      await expect
        .poll(async () => (await mapi.fetchApi(apiId)).labels ?? [], {
          timeout: 30_000,
          message: "API labels removed",
        })
        .toEqual([]);
    });
  },
  { withLabels: true },
);
