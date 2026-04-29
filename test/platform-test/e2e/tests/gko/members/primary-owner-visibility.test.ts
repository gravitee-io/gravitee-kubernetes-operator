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
 * Primary owner + visibility.
 *
 * Xray tests:
 *   GKO-1457: A CR that specifies a PRIMARY_OWNER user and PRIVATE
 *             visibility reconciles correctly — APIM shows the right
 *             primary owner and the visibility flag matches.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const FIXTURE = "crds/members/v4-api-po-user-private.yaml";
const API_NAME = "e2e-v4-po-private";

test.describe("Primary owner & visibility", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(FIXTURE)).catch(() => {});
  });

  test(`Primary owner + PRIVATE visibility reconcile correctly ${XRAY.MEMBERS.PRIMARY_OWNER_VISIBILITY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await kubectl.apply(fixture(FIXTURE));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const apiId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;
    const api = await mapi.fetchApi(apiId);

    expect(api.visibility).toBe("PRIVATE");
    expect(api.primaryOwner?.displayName).toBe("admin");

    await kubectl.del(fixture(FIXTURE));
  });
});
