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
 * V2 API Documentation — Fetchers.
 *
 * Xray tests:
 *   GKO-620: Web fetcher requires URL being set
 *   GKO-621: Web fetcher shows warning when using invalid configuration parameters
 *   GKO-622: Github fetcher requires owner, repository, username, PAT
 *   GKO-623: Github fetcher shows warning on invalid configuration
 *
 * Dropped:
 *   GKO-626, 675, 689, 692 — the GKO admission webhook pre-fetches
 *   github-fetcher pages at apply time, and the test cluster has no real
 *   GitHub credentials. Any positive github-fetcher test is rejected with
 *   "Page cannot be fetched, ... invalid / missing github credentials or
 *   an invalid file path". Re-enable once the test env provisions a
 *   service-account GitHub PAT pointing at a public fixture repo.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("V2 API Documentation — Fetchers", () => {
  // ── GKO-620: Web fetcher requires URL ───────────────────────

  test(`V2 web fetcher without URL is rejected ${XRAY.PAGES.V2_WEB_FETCHER_NO_URL} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/pages/v2-api-web-fetcher-no-url.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/url|required|denied|invalid/);
  });

  // ── GKO-621: Web fetcher warning on invalid parameters ──────

  test(`V2 web fetcher with invalid cron is rejected ${XRAY.PAGES.V2_WEB_FETCHER_WARNING} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/pages/v2-api-web-fetcher-invalid-param.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/cron|fetchcron|invalid|denied/);
  });

  // ── GKO-622: Github fetcher requires core fields ────────────

  test(`V2 github fetcher missing required fields is rejected ${XRAY.PAGES.V2_GITHUB_FETCHER_REQUIRED_FIELDS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/pages/v2-api-github-fetcher-missing-fields.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(
      /owner|repository|username|personalaccesstoken|required|denied|invalid/,
    );
  });

  // ── GKO-623: Github fetcher warning on invalid config ───────

  test(`V2 github fetcher with invalid cron is rejected ${XRAY.PAGES.V2_GITHUB_FETCHER_WARNING} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/pages/v2-api-github-fetcher-invalid-param.yaml");

    const stderr = await kubectl.applyExpectFailure(fixturePath);
    expect(stderr.toLowerCase()).toMatch(/cron|fetchcron|invalid|denied/);
  });
});
