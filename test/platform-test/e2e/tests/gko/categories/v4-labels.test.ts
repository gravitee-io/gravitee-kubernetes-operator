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
 * V4 API Labels lifecycle.
 *
 * Xray tests:
 *   GKO-1473: Labels on V4 API CRs are set initially and removed when the
 *             CR no longer declares them (sync-down).
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const WITH_LABELS = "crds/api-v4-definitions/v4-proxy-api-with-labels-categories.yaml";
const WITHOUT_LABELS = "crds/api-v4-definitions/v4-proxy-api-without-labels.yaml";
const API_NAME = "e2e-v4-labels-cats";

test.describe("V4 API Labels — Lifecycle", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(WITHOUT_LABELS)).catch(() => {});
    await kubectlSafe.del(fixture(WITH_LABELS)).catch(() => {});
  });

  test(`Labels are set, then removed when stripped from CR ${XRAY.CATEGORIES.V4_LABELS_LIFECYCLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("Deploy with labels", async () => {
      await kubectl.apply(fixture(WITH_LABELS));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const apiId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;

    await test.step("APIM reports both labels", async () => {
      await mapi.waitForApiMatches(apiId, {
        labels: ["e2e-label-1", "e2e-label-2"],
      });
    });

    await test.step("Remove labels via CR update", async () => {
      await kubectl.apply(fixture(WITHOUT_LABELS));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("APIM no longer reports the labels", async () => {
      await expect
        .poll(
          async () => {
            const api = await mapi.fetchApi(apiId);
            return api.labels ?? [];
          },
          { timeout: 30_000, intervals: [1_000] },
        )
        .toEqual([]);
    });

    await kubectl.del(fixture(WITHOUT_LABELS));
  });
});
