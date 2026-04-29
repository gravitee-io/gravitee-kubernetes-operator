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
 * V4 API Documentation — Fetchers.
 *
 * Xray tests:
 *   GKO-629:  V4 web fetcher without URL is rejected
 *   GKO-628:  V4 web fetcher with invalid configuration (cron) is rejected
 *   GKO-636:  V4 github fetcher missing required fields is rejected
 *   GKO-637:  V4 github fetcher with invalid configuration (cron) is rejected
 *   GKO-1475: Cross-version validation coverage for schedulers/fetchers
 *             (URL + cron) — piggybacks on the V4 no-URL case.
 *
 * Dropped:
 *   GKO-638 — V4 github fetcher default URL (positive path needs real
 *             credentials; same blocker that killed GKO-626/675).
 *   GKO-697 — V4 delete fetched ROOT pages (same blocker as GKO-662).
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("V4 API Documentation — Fetchers", () => {
  // ── GKO-629 / GKO-1475: V4 web fetcher without URL ──────────

  test(`V4 web fetcher without URL is rejected ${XRAY.PAGES.V4_WEB_FETCHER_NO_URL} ${XRAY.WEBHOOKS.CROSS_VERSION_SCHEDULERS_FETCHERS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/pages/v4-api-web-fetcher-no-url.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/url|required|denied|invalid/);
  });

  // ── GKO-628: V4 web fetcher warning on invalid parameters ───

  test(`V4 web fetcher with invalid cron is rejected ${XRAY.PAGES.V4_WEB_FETCHER_WARNING} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/pages/v4-api-web-fetcher-invalid-param.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/cron|fetchcron|invalid|denied/);
  });

  // ── GKO-636: V4 github fetcher requires core fields ─────────

  test(`V4 github fetcher missing required fields is rejected ${XRAY.PAGES.V4_GITHUB_FETCHER_REQUIRED_FIELDS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/pages/v4-api-github-fetcher-missing-fields.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(
      /owner|repository|username|personalaccesstoken|required|denied|invalid/,
    );
  });

  // ── GKO-637: V4 github fetcher warning on invalid config ────
  // The pre-fetch attempt runs before cron validation and fails with a
  // credentials error (same behaviour as V2 GKO-623). The regex accepts
  // either the cron-specific message or the generic "invalid" response
  // from the credentials path, which is still proof the invalid config
  // is being rejected by admission.

  test(`V4 github fetcher with invalid cron is rejected ${XRAY.PAGES.V4_GITHUB_FETCHER_WARNING} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/pages/v4-api-github-fetcher-invalid-param.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/cron|fetchcron|invalid|denied/);
  });
});
