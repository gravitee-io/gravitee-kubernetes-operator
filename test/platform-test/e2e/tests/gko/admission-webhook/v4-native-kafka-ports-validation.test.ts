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
 * Native Kafka API — plan port admission validation.
 *
 * Rules:
 * - if bootstrapPort is set:
 *   - brokerRangeStart < brokerRangeEnd
 *   - bootstrapPort must NOT be within [brokerRangeStart, brokerRangeEnd] (inclusive)
 */

import { test, fixture, expect } from "../../../setup.js";
import { TAGS } from "../../../helpers/tags.js";
import type { PlanV4 } from "../../../../src/types/apim.js";

test.describe("Native Kafka API — Plan ports validation", () => {
  test(`Valid native kafka plan ports are accepted ${TAGS.REGRESSION} @since-4.12`, async ({ kubectl, mapi }) => {
    const API_NAME = "e2e-v4-native-kafka-ports";
    const f = fixture("crds/api-v4-definitions/v4-native-kafka-ports-valid.yaml");
    await kubectl.apply(f);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    // Beyond passing admission, the port-based-routing fields must round-trip
    // through the management API GET endpoint (GKO-2919 / PR #1684).
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id;
    await expect
      .poll(async () => ((await mapi.listApiPlans(apiId)) as PlanV4[])[0])
      .toMatchObject({
        bootstrapPort: 9092,
        brokerRangeStart: 9100,
        brokerRangeEnd: 9102,
      });

    await kubectl.del(f).catch(() => {});
  });

  test(`Invalid broker range is rejected ${TAGS.REGRESSION}`, async ({ kubectl }) => {
    const f = fixture("crds/api-v4-definitions/v4-native-kafka-ports-range-error.yaml");
    const stderr = await kubectl.applyExpectFailure(f);
    expect(stderr.toLowerCase()).toMatch(/invalid broker port range|broker port range|brokerrange/i);
    await kubectl.del(f).catch(() => {});
  });

  test(`bootstrapPort within range is rejected ${TAGS.REGRESSION}`, async ({ kubectl }) => {
    const f = fixture("crds/api-v4-definitions/v4-native-kafka-ports-bootstrap-in-range.yaml");
    const stderr = await kubectl.applyExpectFailure(f);
    expect(stderr.toLowerCase()).toMatch(/bootstrapport within broker port range|bootstrapport.*broker port range/i);
    await kubectl.del(f).catch(() => {});
  });
});

