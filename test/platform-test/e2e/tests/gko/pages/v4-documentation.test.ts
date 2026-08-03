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
 * V4 API Documentation — Extended scenarios.
 *
 * Xray tests:
 *   GKO-280:  Documentation created by GKO is read-only when re-imported
 *
 * Shipping, revising and removing an inline page (GKO-1470/277/278/236/1469/282)
 * is the shared journey tests/user-journeys/api-producer/document-an-api.
 *
 * Skipped tests:
 *   GKO-283 (V4 spec.visibility PUBLIC-only) — GKO product bug
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import YAML from "yaml";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

interface ExportedPage {
  type?: string;
  visibility?: string;
  content?: string;
  name?: string;
}

interface ExportedCrd {
  spec?: {
    pages?: Record<string, ExportedPage>;
  };
}

test.describe(`V4 API Documentation — Extended ${PROVISIONER.GKO}`, () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. Each del() ignores errors (the resource may already be gone).
  test.afterEach(async () => {
    await kubectl.del(fixture("pages/v4-api-with-page-markdown/crd.yaml")).catch(() => {});
  });

  // ── GKO-280: GKO-created documentation is read-only ─────────
  // Re-applying the same CRD must not mutate the existing page content.

  // FIXME(GKO-2940): re-applying a V4 CR that omits the endpoint `secondary`
  // field (post-GKO-2857) churns metadata.generation while the Accepted
  // condition's observedGeneration lags, so kubectl wait (observedGeneration-
  // aware in kubectl >=1.31) blocks. Re-enable once GKO-2940 is fixed.
  test.fixme(`GKO-managed documentation is stable across re-apply ${XRAY.PAGES.V4_READ_ONLY_DOC} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-markdown-page";
    const WITH_PAGE = fixture("pages/v4-api-with-page-markdown/crd.yaml");

    await kubectl.apply(WITH_PAGE);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id;

    const contentBefore = (
      YAML.parse(await mapi.exportApiCrd(apiId)) as ExportedCrd
    ).spec?.pages?.["markdown-page"]?.content;

    await kubectl.apply(WITH_PAGE);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const contentAfter = (
      YAML.parse(await mapi.exportApiCrd(apiId)) as ExportedCrd
    ).spec?.pages?.["markdown-page"]?.content;

    expect(contentAfter).toBe(contentBefore);

    await kubectl.del(WITH_PAGE);
  });
});
