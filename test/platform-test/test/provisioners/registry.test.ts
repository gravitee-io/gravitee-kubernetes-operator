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

import { describe, it, expect } from "vitest";
import { PROVISIONER_ORDER, PROVISIONER_LANES } from "../../src/provisioners/registry.js";

describe("provisioner registry", () => {
  it("has exactly one lane per provisioner in PROVISIONER_ORDER", () => {
    const laneIds = PROVISIONER_LANES.map((lane) => lane.id);
    expect(new Set(laneIds)).toEqual(new Set(PROVISIONER_ORDER));
    expect(laneIds).toHaveLength(PROVISIONER_ORDER.length);
  });

  it("gives every lane a non-empty, unique tag", () => {
    const tags = PROVISIONER_LANES.map((lane) => lane.tag);
    expect(tags.every((tag) => tag.startsWith("@"))).toBe(true);
    expect(new Set(tags).size).toBe(tags.length);
  });

  it("gives no lane a tag that collides with the Xray tag namespace", () => {
    // Lane selection greps titles for these tags, and titles also carry Xray ids
    // like "@GKO-1234". A lane tag that is a case-insensitive prefix of that
    // namespace (e.g. "@gko") only stays selective because playwright.config.ts
    // greps case-SENSITIVELY with a trailing \b. Guard the assumption: a tag that
    // matched "@GKO-1234" case-insensitively AND at a word boundary would silently
    // select the entire suite for every lane.
    for (const lane of PROVISIONER_LANES) {
      const laneOnly = new RegExp(lane.tag + String.raw`\b`);
      expect(laneOnly.test("Some test @GKO-1234")).toBe(false);
    }
  });
});
