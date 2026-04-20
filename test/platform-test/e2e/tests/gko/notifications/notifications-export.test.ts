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
 * Notification — export & admission checks (batch 7).
 *
 * Xray tests:
 *   GKO-1233: CRD export does not include Notification configuration
 *   GKO-1235: API CR with two console-target notifications is rejected
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import YAML from "yaml";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const NOTIFICATION = "crds/notifications/notification-shared.yaml";
const V4_API_WITH_NOTIF = "crds/notifications/v4-api-using-shared-notification.yaml";
const CONSOLE_A = "crds/notifications/notification-console-a.yaml";
const CONSOLE_B = "crds/notifications/notification-console-b.yaml";
const DUPLICATE_API = "crds/notifications/v4-api-duplicate-console-notifications.yaml";

interface ExportedCrd {
  spec?: Record<string, unknown> & {
    notificationsRefs?: unknown;
    notifications?: unknown;
  };
}

test.describe("Notifications — export & validation", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(V4_API_WITH_NOTIF)).catch(() => {});
    await kubectlSafe.del(fixture(DUPLICATE_API)).catch(() => {});
    await kubectlSafe.del(fixture(CONSOLE_A)).catch(() => {});
    await kubectlSafe.del(fixture(CONSOLE_B)).catch(() => {});
    await kubectlSafe.del(fixture(NOTIFICATION)).catch(() => {});
  });

  // ── GKO-1233: Exported CRD does not include notifications ───

  test(`Exported V4 API CRD excludes Notification refs ${XRAY.NOTIFICATIONS.NOT_IN_EXPORT} ${TAGS.REGRESSION}`, async ({
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

    const crdYaml = await mapi.exportApiCrd(apiId);
    const crd = YAML.parse(crdYaml) as ExportedCrd;

    expect(
      crd.spec?.notificationsRefs,
      `expected notificationsRefs to be absent from export, got ${JSON.stringify(crd.spec?.notificationsRefs)}`,
    ).toBeUndefined();
    expect(
      crd.spec?.notifications,
      `expected notifications to be absent from export, got ${JSON.stringify(crd.spec?.notifications)}`,
    ).toBeUndefined();

    await kubectl.del(fixture(V4_API_WITH_NOTIF));
    await kubectl.del(fixture(NOTIFICATION));
  });

  // ── GKO-1235: Two console notifications on same API rejected

  test(`API with two console-target notifications is rejected ${XRAY.NOTIFICATIONS.DUPLICATE_CONSOLE_REJECTED} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    await kubectl.apply(fixture(CONSOLE_A));
    await kubectl.apply(fixture(CONSOLE_B));

    const stderr = await kubectl.applyExpectFailure(fixture(DUPLICATE_API));
    expect(stderr.toLowerCase()).toMatch(/console|notification|already|duplicate|denied/);

    await kubectl.del(fixture(CONSOLE_A));
    await kubectl.del(fixture(CONSOLE_B));
  });
});
