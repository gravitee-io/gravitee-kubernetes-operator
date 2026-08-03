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
 * Journey: register, update, and retire an application.
 *
 * As an application developer, I register an application with its client
 * settings and the metadata my organisation tags it with, edit it, then retire
 * it. The shared invariant is provisioner-agnostic: whichever driver creates the
 * application, APIM records it via the Automation API (origin KUBERNETES) with
 * the declared settings and metadata, reflects a description update, and
 * ARCHIVES it on retire.
 *
 * Metadata has to be read from its OWN endpoint: the application detail
 * response omits it entirely, which is why declaring metadata used to pass
 * without ever being checked.
 *
 * Fixtures are co-located in this folder (gko/ + terraform/ + README.md).
 * Application admission rules (both `app` and `oauth` settings at once, a
 * missing/duplicate client id, invalid redirect URIs or grant types) stay in
 * tests/gko/applications.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";
import { assertProvisioner } from "../../../../../src/provisioners/index.js";

const here = path.dirname(fileURLToPath(import.meta.url));
const REGISTERED_DESCRIPTION = "Application registered via the register/update/retire journey";
const UPDATED_DESCRIPTION = "Application updated via the register/update/retire journey";

/** The single knob the journey re-provisions with: the create vs updated state. */
interface AppLifecycleParams {
  updated: boolean;
}

forEachProvisioner<AppLifecycleParams>(
  {
    title: "Register, update, and retire an application",
    provisioners: {
      gko: gkoScenario<AppLifecycleParams>({
        manifests: [path.join(here, "gko/application.yaml")],
        roles: { application: "lifecycle-application" },
        // provision applies the "created" manifest; update() re-applies the
        // "updated" one. At provision params.updated is false, so this is a no-op.
        applyParams: async (k, params) => {
          if (params.updated) await k.apply(path.join(here, "gko/application-updated.yaml"));
        },
      }),
      terraform: tfScenario<AppLifecycleParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({
          description: params.updated ? UPDATED_DESCRIPTION : REGISTERED_DESCRIPTION,
        }),
        // remove("application") drops the resource from the desired state and
        // re-applies, which APIM treats as a soft-delete (ARCHIVED).
        removeVars: { application: { create_application: false } },
        addresses: { application: "apim_application.app" },
      }),
    },
    xray: {
      gko: [
        XRAY.APPLICATIONS.CREATE_APP,
        XRAY.APPLICATIONS.UPDATE_APP,
        XRAY.APPLICATIONS.DELETE_APP,
        XRAY.APPLICATIONS.APP_WITH_METADATA,
        XRAY.APPLICATIONS.APP_CONFIGURE_SETTINGS,
      ],
      terraform: [XRAY.TERRAFORM.DELETE_APPLICATION_TF, XRAY.TERRAFORM.APPLICATION_LIFECYCLE_TF],
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const appId = await provisioned.applicationId();

    await test.step("Registered application lands in APIM (origin KUBERNETES)", async () => {
      await mapi.waitForApplicationMatches(
        appId,
        {
          description: REGISTERED_DESCRIPTION,
          origin: "KUBERNETES",
          settings: { app: { type: "SIMPLE" } },
        },
        { timeoutMs: 30_000 },
      );
    });

    await test.step("Declared metadata is recorded against the application", async () => {
      await expect
        .poll(
          async () =>
            (await mapi.listApplicationMetadata(appId)).map((m) => ({
              name: m.name,
              value: m.value,
              format: m.format,
            })),
          { timeout: 30_000, message: "application metadata reaches APIM" },
        )
        .toEqual([{ name: "owner-team", value: "payments", format: "STRING" }]);
    });

    await test.step("Description update propagates to APIM", async () => {
      await provisioned.update({ updated: true });
      await mapi.waitForApplicationMatches(
        appId,
        { description: UPDATED_DESCRIPTION },
        { timeoutMs: 30_000 },
      );
    });

    await test.step("Retiring the application archives it in APIM", async () => {
      await provisioned.remove("application");
      await assertProvisioner(provisioned, "application", "gone");
      await mapi.waitForApplicationMatches(appId, { status: "ARCHIVED" }, { timeoutMs: 30_000 });
    });
  },
  { updated: false },
);
