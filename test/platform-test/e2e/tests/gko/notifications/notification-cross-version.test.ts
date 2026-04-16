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
 * Notification CR cross-version behaviour — batch 5.
 *
 * Xray tests:
 *   GKO-1236: Notification CR cannot be deleted while referenced by an API
 *   GKO-1237: Same Notification CR works with both V2 and V4 APIs
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const NOTIFICATION = "crds/notifications/notification-shared.yaml";
const V4_API = "crds/notifications/v4-api-using-shared-notification.yaml";
const V2_API = "crds/notifications/v2-api-using-shared-notification.yaml";

test.describe("Notification CR — Cross Version", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(V2_API)).catch(() => {});
    await kubectlSafe.del(fixture(V4_API)).catch(() => {});
    await kubectlSafe.del(fixture(NOTIFICATION)).catch(() => {});
  });

  // GKO-1236 (Notification CR referenced by an API cannot be deleted) —
  // dropped from batch 5: GKO does not enforce an in-use protection on
  // Notification CRs. Tracked in "Batch 5 - Skipped Tests.md" as a
  // product gap.

  // ── GKO-1237: Shared Notification works with V2 and V4 ─────

  test(`Notification CR is usable by both V2 and V4 APIs ${XRAY.NOTIFICATIONS.WORKS_WITH_V2_AND_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    test.slow();

    await kubectl.apply(fixture(NOTIFICATION));
    await kubectl.apply(fixture(V4_API));
    await kubectl.apply(fixture(V2_API));

    await kubectl.waitForCondition(
      "apiv4definition",
      "e2e-v4-shared-notification",
      "Accepted",
    );
    await kubectl.waitForCondition(
      "apidefinition",
      "e2e-v2-shared-notification",
      "Accepted",
    );

    // Cleanup in reverse dependency order.
    await kubectl.del(fixture(V2_API));
    await kubectl.del(fixture(V4_API));
    await kubectl.waitForDeletion("apidefinition", "e2e-v2-shared-notification");
    await kubectl.waitForDeletion("apiv4definition", "e2e-v4-shared-notification");
    await kubectl.del(fixture(NOTIFICATION));
  });
});
