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
 * ManagementContext CRD tests.
 *
 * Xray tests:
 *   GKO-892: Cannot delete ManagementContext when referenced by V2 API
 *   GKO-893: Cannot delete ManagementContext when referenced by V4 API
 *   GKO-894: Cannot delete ManagementContext when referenced by Application
 *   GKO-895: ManagementContext deletion succeeds with no references
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("ManagementContext — Lifecycle", () => {
  // ── GKO-893: Cannot delete when referenced by V4 API ─────────

  test(`Cannot delete ManagementContext referenced by V4 API ${XRAY.MANAGEMENT_CONTEXT.DELETE_WITH_V4_API_REF} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const apiFixture = fixture("crds/api-v4-definitions/v4-proxy-api-sync-from-mgmt.yaml");

    await test.step("Deploy an API referencing dev-ctx", async () => {
      await kubectl.apply(apiFixture);
      await kubectl.waitForCondition("apiv4definition", "e2e-v4-sync-mgmt", "Accepted");
    });

    await test.step("Attempting to delete dev-ctx fails", async () => {
      const stderr = await kubectl.delExpectFailure(
        fixture("crds/management-context/dev-ctx.yaml"),
      );
      expect(stderr.toLowerCase()).toContain("cannot be deleted");
    });

    await kubectl.del(apiFixture);
  });

  // ── GKO-472: Non-existing environment ────────────────────────

  test(`ManagementContext with non-existing environment is rejected ${XRAY.MANAGEMENT_CONTEXT.NON_EXISTING_ENV} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/management-context-invalid-env.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("invalid organization or environment");
  });

  // ── GKO-473: Non-existing organization ───────────────────────

  test(`ManagementContext with non-existing organization is rejected ${XRAY.MANAGEMENT_CONTEXT.NON_EXISTING_ORG} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/invalid/management-context-invalid-org.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("invalid organization or environment");
  });

  // ── GKO-894: Cannot delete ManagementContext with Application ref ─

  test(`Cannot delete ManagementContext referenced by Application ${XRAY.MANAGEMENT_CONTEXT.DELETE_WITH_APP_REF} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const appFixture = fixture("crds/applications/application-simple.yaml");

    await test.step("Deploy an Application referencing dev-ctx", async () => {
      await kubectl.apply(appFixture);
      await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");
    });

    await test.step("Attempting to delete dev-ctx fails", async () => {
      const stderr = await kubectl.delExpectFailure(
        fixture("crds/management-context/dev-ctx.yaml"),
      );
      expect(stderr.toLowerCase()).toContain("cannot be deleted");
    });

    await kubectl.del(appFixture);
  });

  // ── GKO-895: Delete ManagementContext with no references ─────

  test(`ManagementContext deletion succeeds with no references ${XRAY.MANAGEMENT_CONTEXT.DELETE_NO_REFS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const ctxFixture = fixture("crds/management-context/temporary-ctx.yaml");

    await test.step("Create a temporary ManagementContext with valid creds", async () => {
      await kubectl.apply(ctxFixture);
      // ManagementContext does not expose a condition to wait on — give the operator
      // a moment to register it before attempting deletion.
      await new Promise((r) => setTimeout(r, 2_000));
    });

    await test.step("Delete succeeds (no APIs/Apps reference it)", async () => {
      await kubectl.del(ctxFixture);
      await kubectl.waitForDeletion("managementcontext", "e2e-temp-ctx");
    });
  });
});
