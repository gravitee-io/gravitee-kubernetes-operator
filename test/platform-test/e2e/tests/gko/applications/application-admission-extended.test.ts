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
 * Application admission & lifecycle edge cases.
 *
 * Xray tests:
 *   GKO-505:  Applications deployed via CR are read-only in APIM (origin=KUBERNETES)
 *   GKO-578:  BROWSER application type with invalid redirectUris is rejected
 *   GKO-579:  SPA application type only accepts authorization_code/implicit
 *   GKO-1382: Application name length edge case (exceeds APIM limit)
 *   GKO-1383: Successful deletion of an application via CR removal
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const FIXTURE_505 = "crds/applications/application-505-readonly.yaml";
const FIXTURE_578 = "crds/applications/application-578-browser-bad-uris.yaml";
const FIXTURE_579 = "crds/applications/application-579-spa-bad-grants.yaml";
const FIXTURE_1382 = "crds/applications/application-1382-long-name.yaml";
const FIXTURE_1383 = "crds/applications/application-1383-delete.yaml";

interface StatusWithId {
  id?: string;
  conditions?: Array<{ type: string; status: string }>;
}

test.describe("Applications — admission & lifecycle", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(FIXTURE_505)).catch(() => {});
    await kubectlSafe.del(fixture(FIXTURE_1383)).catch(() => {});
    // 578, 579, 1382 expected to fail admission; nothing to clean up.
  });

  // ── GKO-505: read-only via origin=KUBERNETES ─────────────────

  test(`Applications deployed via CR are read-only in APIM ${XRAY.APPLICATIONS.APP_READ_ONLY_IN_APIM} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP_NAME = "e2e-app-505-readonly";
    await kubectl.apply(fixture(FIXTURE_505));
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const appId = (await kubectl.getStatus<StatusWithId>("application", APP_NAME)).id;
    expect(appId).toBeTruthy();

    await mapi.waitForApplicationMatches(appId!, {
      // origin=KUBERNETES is APIM's marker that the resource is managed by GKO
      // and read-only in the console.
      origin: "KUBERNETES",
    });
  });

  // ── GKO-578: BROWSER + invalid redirectUris ──────────────────

  test(`BROWSER application with invalid redirectUris is rejected ${XRAY.APPLICATIONS.APP_BROWSER_VALID_URIS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(fixture(FIXTURE_578));
    // Scoped to redirectUri-specific terms; a generic webhook failure (e.g.
    // DNS/timeout) would no longer satisfy the assertion.
    expect(stderr.toLowerCase()).toMatch(/uri|redirect/);
  });

  // ── GKO-579: SPA grant types restricted ──────────────────────

  test(`SPA application with disallowed grant_type is rejected ${XRAY.APPLICATIONS.APP_SPA_GRANT_TYPES} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(fixture(FIXTURE_579));
    expect(stderr.toLowerCase()).toMatch(/grant|authorization_code|implicit/);
  });

  // ── GKO-1382: name length edge case ──────────────────────────
  // APIM does not enforce a hard length limit on application `name` (a 60-char
  // name is accepted via the CRD + dry-run path). The Xray scenario is a
  // boundary check: apply the CR with a long name and verify the application
  // is created with the full name preserved (no silent truncation / rename).

  test(`Application with a long name is created with the name preserved ${XRAY.APPLICATIONS.APP_NAME_LENGTH_EDGE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const CR_NAME = "e2e-app-1382-longname";
    const LONG_NAME = "e2e-1382-this-app-name-is-far-too-long-and-exceeds-fifty-chars";

    await kubectl.apply(fixture(FIXTURE_1382));
    await kubectl.waitForCondition("application", CR_NAME, "Accepted");

    const appId = (await kubectl.getStatus<StatusWithId>("application", CR_NAME)).id;
    expect(appId).toBeTruthy();

    await mapi.waitForApplicationMatches(appId!, { name: LONG_NAME });
  });

  // ── GKO-1383: successful deletion ────────────────────────────
  // Per existing GKO-337 coverage, deleting the Application CR marks the APIM
  // application as ARCHIVED (APIM's soft-delete behavior). We assert that.

  test(`Application is ARCHIVED in APIM when the CR is removed ${XRAY.APPLICATIONS.APP_DELETE_SUCCESS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP_NAME = "e2e-app-1383-delete";

    await test.step("Apply application CR", async () => {
      await kubectl.apply(fixture(FIXTURE_1383));
      await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    });

    const appId = (await kubectl.getStatus<StatusWithId>("application", APP_NAME)).id;
    expect(appId).toBeTruthy();

    await test.step("Application is fetchable via mAPI before delete", async () => {
      await mapi.assertApplicationMatches(appId!, { name: APP_NAME });
    });

    await test.step("Delete the application CR", async () => {
      await kubectl.del(fixture(FIXTURE_1383));
      await kubectl.waitForDeletion("application", APP_NAME);
    });

    await test.step("Application is ARCHIVED in APIM after delete", async () => {
      await mapi.waitForApplicationMatches(
        appId!,
        { status: "ARCHIVED" },
        { timeoutMs: 15_000 },
      );
    });
  });
});
