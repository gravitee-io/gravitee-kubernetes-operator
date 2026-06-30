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
 * Use case: register, update, and retire an application.
 *
 * An API consumer onboards an application, edits it, then offboards it. The
 * shared invariant is provisioner-agnostic: whichever driver creates the
 * application, APIM records it via the Automation API (origin KUBERNETES),
 * reflects a description update, and ARCHIVES it on retire.
 *
 * Fixtures (the discoverable demo) live in
 * fixtures/use-cases/register-and-retire-application/. GKO-only application
 * behaviour (members, OAuth settings, admission) stays under
 * tests/gko/applications.
 */

import { test, fixture } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";

const FIXTURE = "use-cases/register-and-retire-application";
const REGISTERED_DESCRIPTION = "E2E use-case: application registered";
const UPDATED_DESCRIPTION = "E2E use-case: application updated";

// The GKO "updated" state is a sibling manifest re-applied on update(); resolve
// its absolute path once at load time so the apply closure stays trivial.
const GKO_UPDATED_MANIFEST = fixture(`${FIXTURE}/gko/application-updated.yaml`);

/** The single knob the journey re-provisions with: the create vs updated state. */
interface AppLifecycleParams {
  updated: boolean;
}

forEachProvisioner<AppLifecycleParams>(
  {
    title: "Register, update, and retire an application",
    provisioners: {
      gko: gkoScenario<AppLifecycleParams>({
        manifests: [`${FIXTURE}/gko/application.yaml`],
        roles: { application: "e2e-uc-app-lifecycle" },
        // provision applies the "created" manifest; update() re-applies the
        // "updated" one. At provision params.updated is false, so this is a no-op.
        applyParams: async (k, params) => {
          if (params.updated) await k.apply(GKO_UPDATED_MANIFEST);
        },
      }),
      terraform: tfScenario<AppLifecycleParams>({
        fixture: `${FIXTURE}/terraform`,
        toVars: (params) => ({
          description: params.updated ? UPDATED_DESCRIPTION : REGISTERED_DESCRIPTION,
        }),
        // remove("application") drops the resource from the desired state and
        // re-applies, which APIM treats as a soft-delete (ARCHIVED).
        removeVars: { application: { create_application: false } },
      }),
    },
    xray: {
      gko: [XRAY.APPLICATIONS.CREATE_APP, XRAY.APPLICATIONS.UPDATE_APP, XRAY.APPLICATIONS.DELETE_APP],
      terraform: [XRAY.TERRAFORM.DELETE_APPLICATION_TF, XRAY.TERRAFORM.APPLICATION_LIFECYCLE_TF],
    },
    tags: [TAGS.REGRESSION],
    // apim_application + the Automation API application endpoints both ship in 4.12.
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const appId = await provisioned.applicationId();

    await test.step("Registered application lands in APIM (origin KUBERNETES)", async () => {
      await mapi.waitForApplicationMatches(
        appId,
        { description: REGISTERED_DESCRIPTION, origin: "KUBERNETES" },
        { timeoutMs: 30_000 },
      );
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
      await mapi.waitForApplicationMatches(appId, { status: "ARCHIVED" }, { timeoutMs: 30_000 });
    });
  },
  { updated: false },
);
