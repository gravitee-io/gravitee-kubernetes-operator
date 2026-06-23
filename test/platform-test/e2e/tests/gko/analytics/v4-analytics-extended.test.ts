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
 * V4 API Analytics — Extended scenarios.
 *
 * This suite is primarily about validating that analytics-related fields
 * are accepted by the CRD/admission layer.
 */

import { test, fixture, expect } from "../../../setup.js";
import { TAGS } from "../../../helpers/tags.js";

test.describe("V4 API Analytics — Extended", () => {
  test(`V4 API analytics fields are accepted by admission ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-analytics-otel";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-analytics-otel.yaml");
    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const api = await mapi.fetchApi(status.id);

    await mapi.waitForApiMatches(api.id, {
      analytics: {
        enabled: true,
        otelLogs: { enabled: true },
        tracing: { enabled: true, verbose: true },
      },
    })

    await kubectl.del(fixturePath).catch(() => {});
  });
});

