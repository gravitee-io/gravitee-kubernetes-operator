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
 * V4 API — Automation API fields.
 *
 * Validates that fields managed via the Automation API are applied in APIM.
 */

import { test, fixture } from "../../../setup.js";
import { TAGS } from "../../../helpers/tags.js";

test.describe("V4 API — Automation API fields @since-4.12", () => {
  test(`Automation API fields are applied in APIM ${TAGS.REGRESSION}`, async ({ kubectl, mapi }) => {
    const API_NAME = "e2e-v4-new-fields";
    const f = fixture("crds/api-v4-definitions/v4-proxy-api-new-fields.yaml");

    await kubectl.apply(f);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    await mapi.waitForApiMatches(status.id, {
      allowedInApiProducts: true,
      allowMultiJwtOauth2Subscriptions: true,
    });

    await kubectl.del(f).catch(() => {});
  });
});

