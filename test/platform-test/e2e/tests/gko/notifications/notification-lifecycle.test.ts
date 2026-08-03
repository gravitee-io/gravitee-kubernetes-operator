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
 * Notification lifecycle tests (V4 only).
 *
 * Xray tests:
 *   GKO-1238: Removing the notification reference in an API removes related notifications in UI
 *
 * Configuring the console notification (GKO-1231/1232/1461/1194/1195/1219/1239)
 * is the shared journey
 * tests/user-journeys/api-producer/configure-api-notifications. Removal stays
 * here: the Automation API answers HTTP 500 to a console notification with an
 * empty `events` list (GKO-3085), so the Terraform arm cannot express turning it
 * off.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, expect, fixture } from "../../../setup.js";
import { XRAY, PROVISIONER } from "../../../helpers/tags.js";

/** Create a service account in APIM (needed for group members). */
async function createServiceAccount(mapi: { http: { managementV1Path(r: string): string; post<T>(p: string, b: unknown): Promise<{ status: number; body: T }> } }, name: string): Promise<void> {
  const path = mapi.http.managementV1Path("/users");
  const res = await mapi.http.post(path, {
    firstname: name,
    lastname: "Service",
    email: `${name}@gravitee.io`,
    source: "gravitee",
    sourceId: name,
  });
  // 200 = created, 400/409 = already exists (OK for idempotent setup)
  if (res.status !== 200 && res.status !== 201 && res.status !== 400 && res.status !== 409) {
    throw new Error(`Failed to create service account ${name}: ${res.status}`);
  }
}


test.describe(`Notification Lifecycle @since-4.12 ${PROVISIONER.GKO}`, () => {
  test(`Remove notification from API ${XRAY.NOTIFICATIONS.REMOVE_NOTIFICATION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-remove-notification";

    await test.step("Create service account for group member", async () => {
      await createServiceAccount(mapi, "e2e-sa-remove-notification");
    });

    await test.step("Create Group and Notification resources", async () => {
      await kubectl.apply(fixture("notifications/group-for-remove-notification/crd.yaml"));
      await kubectl.waitForCondition("group", "e2e-group-remove-notification", "Accepted");
      await kubectl.apply(fixture("notifications/notification-for-remove/crd.yaml"));
    });

    await test.step("Deploy API with notification reference", async () => {
      await kubectl.apply(fixture("notifications/v4-api-with-notification-remove/crd.yaml"));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Remove notification reference from API", async () => {
      await kubectl.apply(fixture("notifications/v4-api-with-removed-notification/crd.yaml"));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("PORTAL notification has empty hooks and groups", async () => {
      await expect.poll(async () => {
        const settings = await mapi.fetchApiNotificationSettings(apiId);
        return settings.find((s) => s.config_type === "PORTAL");
      }, { timeout: 10_000 }).toMatchObject({
        hooks: [],
        groups: [],
      });
    });

    // Cleanup
    await kubectl.del(fixture("notifications/v4-api-with-removed-notification/crd.yaml"));
    await kubectl.del(fixture("notifications/notification-for-remove/crd.yaml"));
    await kubectl.del(fixture("notifications/group-for-remove-notification/crd.yaml"));
  });
});
