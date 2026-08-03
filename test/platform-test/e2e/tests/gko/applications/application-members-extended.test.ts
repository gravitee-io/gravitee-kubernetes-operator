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
 * Application Members (non-OAuth).
 *
 * Xray tests:
 *   GKO-531: Add member with non-existing role (default role applied)
 *
 * Granting, re-roling and revoking an application member (GKO-534/538/539) is
 * the shared journey
 * tests/user-journeys/api-consumer/manage-application-members.
 *
 * Skipped: GKO-537 (remove member from APIM env) requires console UI.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

interface StatusWithConditions {
  id?: string;
  conditions?: Array<{ type: string; status: string }>;
}

function acceptedTrue(status: StatusWithConditions): boolean {
  return status.conditions?.find((c) => c.type === "Accepted")?.status === "True";
}

test.describe(`Applications — Members Extended ${PROVISIONER.GKO}`, () => {
  test.afterEach(async () => {
    await kubectlSafe
      .del(fixture("applications/application-member-non-existing-role/crd.yaml"))
      .catch(() => {});
  });

  // ── GKO-531: Non-existing role ──────────────────────────────
  // Xray contract: "the application is created, the default role is added".

  test(`Application with member using non-existing role is created ${XRAY.APPLICATIONS_MEMBERS.APP_NON_EXISTING_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const APP_NAME = "e2e-app-member-non-existing-role";
    const fixturePath = fixture(
      "applications/application-member-non-existing-role/crd.yaml",
    );

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("application", APP_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

});
