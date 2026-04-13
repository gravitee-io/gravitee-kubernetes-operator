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
 * Application Members & OAuth — batch 4 coverage.
 *
 * Xray tests:
 *   GKO-533: Add non-existing member to application
 *   GKO-535: Application member without role (defaults applied)
 *   GKO-536: Application member without source (rejected)
 *   GKO-548: Application referencing non-existing group
 *   GKO-555: Application member with non-existing role
 *   GKO-581: WEB application must include authorization_code
 *
 * Skipped tests (see "Batch 4 - Skipped Tests.md" in hermesVault):
 *   GKO-553 (Configure OAuth settings)        — DCR not enabled in test env
 *   GKO-574 (BACKEND_TO_BACKEND grant)        — DCR not enabled in test env
 *   GKO-576 (B2B redirectURIs optional)       — DCR not enabled in test env
 *   GKO-580 (WEB accepts multiple grants)     — DCR not enabled in test env
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

interface StatusWithConditions {
  id?: string;
  conditions?: Array<{
    type: string;
    status: string;
  }>;
}

function acceptedTrue(status: StatusWithConditions): boolean {
  return status.conditions?.find((c) => c.type === "Accepted")?.status === "True";
}

test.describe("Applications — Members & OAuth", () => {
  // ── GKO-533: Non-existing member ─────────────────────────────

  test(`Application with non-existing member ${XRAY.APPLICATIONS_MEMBERS.APP_NON_EXISTING_MEMBER} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const APP_NAME = "e2e-app-bad-member";
    const fixturePath = fixture("crds/applications/application-non-existing-member.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("application", APP_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-535: Member without role ────────────────────────────

  test(`Application member without role gets default ${XRAY.APPLICATIONS_MEMBERS.APP_MEMBER_NO_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const APP_NAME = "e2e-app-member-no-role";
    const fixturePath = fixture("crds/applications/application-member-no-role.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("application", APP_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-536: Member without source ──────────────────────────

  test(`Application member without source is rejected ${XRAY.APPLICATIONS_MEMBERS.APP_MEMBER_NO_SOURCE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/applications/application-member-no-source.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/source|required|denied|invalid/);
  });

  // ── GKO-548: Non-existing group ──────────────────────────────

  test(`Application with non-existing group ${XRAY.APPLICATIONS_MEMBERS.APP_NON_EXISTING_GROUP} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const APP_NAME = "e2e-app-bad-group";
    const fixturePath = fixture("crds/applications/application-non-existing-group.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("application", APP_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // GKO-553 (Configure OAuth settings) was skipped because the test APIM
  // environment does not have Dynamic Client Registration enabled. See
  // "Batch 4 - Skipped Tests.md" in hermesVault.

  // ── GKO-555: Non-existing role ──────────────────────────────

  test(`Application member with non-existing role ${XRAY.APPLICATIONS_MEMBERS.APP_MEMBER_NON_EXISTING_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const APP_NAME = "e2e-app-bad-role";
    const fixturePath = fixture("crds/applications/application-member-bad-role.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("application", APP_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // GKO-574, GKO-576, GKO-580 were skipped for the same DCR precondition as
  // GKO-553. See "Batch 4 - Skipped Tests.md" in hermesVault.

  // ── GKO-581: WEB must include authorization_code ────────────

  test(`WEB application without authorization_code is rejected ${XRAY.APPLICATIONS_MEMBERS.APP_WEB_REQUIRES_AUTH_CODE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/applications/application-oauth-web-no-auth-code.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/authorization_code|grant|denied|invalid/);
  });
});
