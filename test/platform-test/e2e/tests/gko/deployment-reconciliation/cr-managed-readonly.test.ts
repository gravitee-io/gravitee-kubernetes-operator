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
 * CR-managed resources are tagged read-only in APIM — batch 7.
 *
 * Xray tests:
 *   GKO-1456: APIM reports CR-managed resources (APIs, Applications) with
 *             definitionContext.origin=KUBERNETES — the signal the UI
 *             reads to render read-only banners and block edit actions.
 *   GKO-1234: Notification settings created through a GKO Notification CR
 *             are tagged origin=KUBERNETES in the notification-settings
 *             endpoint — same read-only signal, applied to notifications.
 *
 * The actual edit-blocking is enforced by the UI layer based on this flag.
 * Asserting the flag is the meaningful contract between GKO and the UI;
 * it verifies the read-only intent is propagated.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const V4_API = "crds/api-v4-definitions/v4-proxy-api-started.yaml";
const APP = "crds/applications/application-simple.yaml";
const NOTIFICATION = "crds/notifications/notification-shared.yaml";
const V4_API_WITH_NOTIF = "crds/notifications/v4-api-using-shared-notification.yaml";

test.describe("CR-managed resources are read-only in APIM", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(V4_API)).catch(() => {});
    await kubectlSafe.del(fixture(APP)).catch(() => {});
    await kubectlSafe.del(fixture(V4_API_WITH_NOTIF)).catch(() => {});
    await kubectlSafe.del(fixture(NOTIFICATION)).catch(() => {});
  });

  // ── GKO-1456: API and Application report origin=KUBERNETES ──

  test(`API and Application are tagged origin=KUBERNETES in APIM ${XRAY.DEPLOYMENT_RECONCILIATION.CR_MANAGED_READ_ONLY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("Deploy V4 API via CR", async () => {
      await kubectl.apply(fixture(V4_API));
      await kubectl.waitForCondition("apiv4definition", "e2e-v4-start-stop", "Accepted");
    });

    await test.step("Deploy Application via CR", async () => {
      await kubectl.apply(fixture(APP));
      await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");
    });

    await test.step("API reports origin=KUBERNETES", async () => {
      const apiId = (
        await kubectl.getStatus<{ id: string }>("apiv4definition", "e2e-v4-start-stop")
      ).id;
      const api = await mapi.fetchApi(apiId);
      expect(api.originContext?.origin).toBe("KUBERNETES");
    });

    await test.step("Application reports origin=KUBERNETES", async () => {
      const appId = (await kubectl.getStatus<{ id: string }>("application", "e2e-app-simple"))
        .id;
      const app = await mapi.fetchApplication(appId);
      expect(app.origin).toBe("KUBERNETES");
    });

    await kubectl.del(fixture(APP));
    await kubectl.del(fixture(V4_API));
  });

  // ── GKO-1234: Notification settings tagged origin=KUBERNETES
  //
  // Skipped: master-GKO ↔ APIM 4.11 payload mismatch. Master GKO sends embedded
  // `consoleNotificationConfiguration` via the Automation API v4 import; APIM
  // 4.11 silently drops it, so the default PORTAL setting comes back with
  // origin=MANAGEMENT. Verified passing with GKO 4.11.4 against APIM 4.11.
  // Re-enable when the test setup pins APIM ≥ 4.12.

  test.skip(`Notification settings created via CR are origin=KUBERNETES ${XRAY.NOTIFICATIONS.CR_READONLY_VIA_MAPI} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await kubectl.apply(fixture(NOTIFICATION));
    await kubectl.apply(fixture(V4_API_WITH_NOTIF));
    await kubectl.waitForCondition(
      "apiv4definition",
      "e2e-v4-shared-notification",
      "Accepted",
    );

    const apiId = (
      await kubectl.getStatus<{ id: string }>(
        "apiv4definition",
        "e2e-v4-shared-notification",
      )
    ).id;

    const settings = await mapi.fetchApiNotificationSettings(apiId);
    const consoleSetting = settings.find((s) => s.config_type === "PORTAL");
    expect(consoleSetting, "expected a PORTAL (console) notification setting").toBeTruthy();
    expect(consoleSetting?.origin).toBe("KUBERNETES");

    await kubectl.del(fixture(V4_API_WITH_NOTIF));
    await kubectl.del(fixture(NOTIFICATION));
  });
});
