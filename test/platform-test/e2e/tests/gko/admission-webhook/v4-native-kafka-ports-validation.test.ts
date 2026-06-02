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

test.describe("Native Kafka API — Plan ports validation", () => {
  test(`Valid native kafka plan ports are accepted ${TAGS.REGRESSION}`, async ({ kubectl }) => {
    const f = fixture("crds/api-v4-definitions/v4-native-kafka-ports-valid.yaml");
    await kubectl.apply(f);
    await kubectl.waitForCondition("apiv4definition", "e2e-v4-native-kafka-ports", "Accepted");
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

