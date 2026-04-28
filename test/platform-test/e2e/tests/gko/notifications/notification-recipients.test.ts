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
 * Notification recipients & visibility — batch 8.
 *
 * Xray tests:
 *   GKO-1194: View notification settings for an API (PORTAL setting present)
 *   GKO-1195: Notification labelled "Console Notification" — verified at the
 *             data layer as a PORTAL config_type. The literal label is a UI
 *             constant and not exposed via mAPI; the underlying setting type
 *             is what GKO is responsible for.
 *   GKO-1196: Default recipient is the owner — no custom groups attached
 *             when an API is deployed without a notification ref.
 *   GKO-1219: Customise the target user(s) for in-app API-level notifications
 *             via groupRefs on a console-targeted Notification CR.
 *   GKO-1239: Members of groups are notified — groupRefs on the Notification
 *             CR propagate to the API's PORTAL notification setting groups.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, expect, fixture } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

// Members referenced by the Group CRs (source=gravitee, sourceId=e2e-sa-…)
// are declared only. APIM does not require a pre-existing user record for
// group member declarations — the sourceId is resolved at auth time. An
// earlier implementation POSTed to /users here but the call consistently
// returned 400 (no matching IDP entry for source=gravitee in the test
// cluster), which is harmless but also useless; it has been removed.

const BASE_API = "crds/notifications/v4-api-batch8-base.yaml";
const NOTIF_1219 = "crds/notifications/notification-1219-portal-target-user.yaml";
const GROUP_1219 = "crds/notifications/group-1219-portal-target.yaml";
const API_1219 = "crds/notifications/v4-api-1219-with-portal-target.yaml";
const NOTIF_1239 = "crds/notifications/notification-1239-group-members.yaml";
const GROUP_1239 = "crds/notifications/group-1239-group-members.yaml";
const API_1239 = "crds/notifications/v4-api-1239-with-group-members.yaml";

// Skipped: master-GKO ↔ APIM 4.11 payload mismatch. Master GKO sends embedded
// `consoleNotificationConfiguration` via the Automation API v4 import; APIM
// 4.11 silently drops it (PORTAL setting comes back with hooks/groups empty
// and origin=MANAGEMENT). Verified passing with GKO 4.11.4 against APIM 4.11.
// Re-enable when the test setup pins APIM ≥ 4.12.
test.describe.skip("Notifications — recipients & visibility (batch 8)", () => {
  test.afterEach(async () => {
    // Reverse dependency order: API → Notification → Group
    await kubectlSafe.del(fixture(API_1239)).catch(() => {});
    await kubectlSafe.del(fixture(NOTIF_1239)).catch(() => {});
    await kubectlSafe.del(fixture(GROUP_1239)).catch(() => {});
    await kubectlSafe.del(fixture(API_1219)).catch(() => {});
    await kubectlSafe.del(fixture(NOTIF_1219)).catch(() => {});
    await kubectlSafe.del(fixture(GROUP_1219)).catch(() => {});
    await kubectlSafe.del(fixture(BASE_API)).catch(() => {});
  });

  // ── GKO-1194: PORTAL setting exists for a deployed API ───────

  test(`A PORTAL notification setting is present for the API ${XRAY.NOTIFICATIONS.VIEW_NOTIFICATION_SETTINGS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-notif-batch8-base";
    await kubectl.apply(fixture(BASE_API));
    await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", NAME)).id;

    await expect
      .poll(
        async () => {
          const settings = await mapi.fetchApiNotificationSettings(apiId);
          return settings.find((s) => s.config_type === "PORTAL");
        },
        { timeout: 10_000, intervals: [1_000] },
      )
      .toBeTruthy();
  });

  // ── GKO-1195: PORTAL setting type matches the "Console Notification" label

  test(`PORTAL notification setting type backs the Console Notification label ${XRAY.NOTIFICATIONS.NOTIFICATION_LABEL} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-notif-batch8-base";
    await kubectl.apply(fixture(BASE_API));
    await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", NAME)).id;

    await expect
      .poll(
        async () => {
          const settings = await mapi.fetchApiNotificationSettings(apiId);
          // The console "Console Notification" label is rendered from the
          // PORTAL config_type. Verify the type field, which is what GKO
          // is responsible for populating.
          const portal = settings.find((s) => s.config_type === "PORTAL");
          return portal?.config_type;
        },
        { timeout: 10_000, intervals: [1_000] },
      )
      .toBe("PORTAL");
  });

  // ── GKO-1196: Default has no extra groups ────────────────────

  test(`Default PORTAL notification setting has no extra groups ${XRAY.NOTIFICATIONS.DEFAULT_RECIPIENT_OWNER} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-notif-batch8-base";
    await kubectl.apply(fixture(BASE_API));
    await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", NAME)).id;

    await expect
      .poll(
        async () => {
          const settings = await mapi.fetchApiNotificationSettings(apiId);
          return settings.find((s) => s.config_type === "PORTAL");
        },
        { timeout: 10_000, intervals: [1_000] },
      )
      .toMatchObject({
        // The Owner is implicit at the API-level. With no notification CR,
        // the "groups" recipient list should not contain any extra entries.
        groups: [],
      });
  });

  // ── GKO-1219: Customise the portal-notifier target user ──────

  test(`Custom group is set as portal-notifier recipient via Notification CR ${XRAY.NOTIFICATIONS.PORTAL_NOTIFIER_TARGET_USER} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-1219-portal-target";
    const GROUP_NAME = "e2e-group-1219-portal-target";

    await test.step("Apply group, notification, then API", async () => {
      await kubectl.apply(fixture(GROUP_1219));
      await kubectl.waitForCondition("group", GROUP_NAME, "Accepted");
      await kubectl.apply(fixture(NOTIF_1219));
      await kubectl.apply(fixture(API_1219));
      await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    });

    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", NAME)).id;
    const groupId = (await kubectl.getStatus<{ id: string }>("group", GROUP_NAME)).id;

    await expect
      .poll(
        async () => {
          const settings = await mapi.fetchApiNotificationSettings(apiId);
          return settings.find((s) => s.config_type === "PORTAL");
        },
        { timeout: 15_000, intervals: [1_000] },
      )
      .toMatchObject({
        groups: expect.arrayContaining([groupId]),
      });
  });

  // ── GKO-1239: Group members are notified via groupRefs ───────

  test(`Notification CR with groupRefs propagates to PORTAL groups ${XRAY.NOTIFICATIONS.GROUP_MEMBERS_NOTIFIED} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-1239-group-members";
    const GROUP_NAME = "e2e-group-1239-group-members";

    await test.step("Apply group, notification, then API", async () => {
      await kubectl.apply(fixture(GROUP_1239));
      await kubectl.waitForCondition("group", GROUP_NAME, "Accepted");
      await kubectl.apply(fixture(NOTIF_1239));
      await kubectl.apply(fixture(API_1239));
      await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    });

    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", NAME)).id;
    const groupId = (await kubectl.getStatus<{ id: string }>("group", GROUP_NAME)).id;

    await expect
      .poll(
        async () => {
          const settings = await mapi.fetchApiNotificationSettings(apiId);
          return settings.find((s) => s.config_type === "PORTAL");
        },
        { timeout: 15_000, intervals: [1_000] },
      )
      .toMatchObject({
        hooks: expect.arrayContaining(["API_STARTED"]),
        groups: expect.arrayContaining([groupId]),
      });
  });
});
