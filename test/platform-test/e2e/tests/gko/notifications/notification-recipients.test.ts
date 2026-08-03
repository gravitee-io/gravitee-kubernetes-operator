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
 * Notification recipients & visibility.
 *
 * Xray tests:
 *   GKO-1196: Default recipient is the owner — no custom groups attached
 *             when an API is deployed without a notification ref.
 *
 * Configuring a console notification and its target groups
 * (GKO-1194/1195/1219/1239) is the shared journey
 * tests/user-journeys/api-producer/configure-api-notifications. This file keeps
 * the case that has no notification at all.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, expect, fixture } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

// Members referenced by the Group CRs (source=gravitee, sourceId=e2e-sa-…)
// are declared only. APIM does not require a pre-existing user record for
// group member declarations — the sourceId is resolved at auth time. An
// earlier implementation POSTed to /users here but the call consistently
// returned 400 (no matching IDP entry for source=gravitee in the test
// cluster), which is harmless but also useless; it has been removed.

const BASE_API = "notifications/v4-api-notif-base/crd.yaml";

test.describe(`Notifications — recipients & visibility ${PROVISIONER.GKO}`, () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(BASE_API)).catch(() => {});
  });

  // ── GKO-1196: Default has no extra groups ────────────────────

  test(`Default PORTAL notification setting has no extra groups ${XRAY.NOTIFICATIONS.DEFAULT_RECIPIENT_OWNER} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-notif-base";
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
});
