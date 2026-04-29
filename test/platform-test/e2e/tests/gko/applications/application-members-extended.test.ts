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
 *   GKO-534: Remove member from an application on re-apply
 *   GKO-538: Add member with a role name
 *   GKO-539: Change member role on an application
 *
 * Skipped: GKO-537 (remove member from APIM env) requires console UI.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

interface StatusWithConditions {
  id?: string;
  conditions?: Array<{ type: string; status: string }>;
}

function acceptedTrue(status: StatusWithConditions): boolean {
  return status.conditions?.find((c) => c.type === "Accepted")?.status === "True";
}

test.describe("Applications — Members Extended", () => {
  test.afterEach(async () => {
    await kubectlSafe
      .del(fixture("crds/applications/application-member-non-existing-role-531.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/applications/application-with-member.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/applications/application-with-member-removed.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/applications/application-with-member-reviewer.yaml"))
      .catch(() => {});
  });

  // ── GKO-531: Non-existing role ──────────────────────────────
  // Xray contract: "the application is created, the default role is added".

  test(`Application with member using non-existing role is created ${XRAY.APPLICATIONS_MEMBERS.APP_NON_EXISTING_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const APP_NAME = "e2e-app-member-non-existing-role";
    const fixturePath = fixture(
      "crds/applications/application-member-non-existing-role-531.yaml",
    );

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("application", APP_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-534: Remove member ──────────────────────────────────

  test(`Remove member from Application on re-apply ${XRAY.APPLICATIONS_MEMBERS.APP_REMOVE_MEMBER} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP_NAME = "e2e-app-with-member";
    const withMember = fixture("crds/applications/application-with-member.yaml");
    const removed = fixture("crds/applications/application-with-member-removed.yaml");

    await kubectl.apply(withMember);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    const appId = (await kubectl.getStatus<{ id: string }>("application", APP_NAME)).id;

    await kubectl.apply(removed);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await mapi.waitForApplicationMatches(appId, {
      description: "E2E test: application with member removed",
    });

    await kubectl.del(removed);
  });

  // ── GKO-538: Add member with role name ──────────────────────

  test(`Application with member by role name is created ${XRAY.APPLICATIONS_MEMBERS.APP_ADD_MEMBER_ROLE_NAME} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP_NAME = "e2e-app-with-member";
    const fixturePath = fixture("crds/applications/application-with-member.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const appId = (await kubectl.getStatus<{ id: string }>("application", APP_NAME)).id;
    await mapi.waitForApplicationMatches(appId, { name: APP_NAME });

    await kubectl.del(fixturePath);
  });

  // ── GKO-539: Change member role ─────────────────────────────

  test(`Change Application member role ${XRAY.APPLICATIONS_MEMBERS.APP_CHANGE_MEMBER_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP_NAME = "e2e-app-with-member";
    const withMember = fixture("crds/applications/application-with-member.yaml");
    const reviewer = fixture("crds/applications/application-with-member-reviewer.yaml");

    await kubectl.apply(withMember);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    const appId = (await kubectl.getStatus<{ id: string }>("application", APP_NAME)).id;

    await kubectl.apply(reviewer);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await mapi.waitForApplicationMatches(appId, {
      description: "E2E test: application with member role changed",
    });

    await kubectl.del(reviewer);
  });
});
