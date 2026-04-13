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
 * V4 API Documentation — Extended scenarios (batch 4).
 *
 * Xray tests:
 *   GKO-236:  CRUD on existing V4 API operations (documentation context)
 *   GKO-280:  Documentation created by GKO is read-only when re-imported
 *   GKO-282:  Inline documentation with PUBLIC visibility
 *   GKO-1470: Documentation managed by GKO is reconciled end-to-end
 *
 * Skipped tests (see "Batch 4 - Skipped Tests.md" in hermesVault):
 *   GKO-283 (V4 spec.visibility PUBLIC-only) — GKO product bug
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import YAML from "yaml";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

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

test.describe("V4 API Documentation — Extended", () => {
  // ── GKO-236: Documentation CRUD on existing V4 operations ───

  test(`Documentation CRUD on existing V4 operations ${XRAY.PAGES.V4_DOC_OPERATIONS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-markdown-page";
    const WITH_PAGE = fixture("crds/pages/v4-api-with-page-markdown.yaml");
    const UPDATED = fixture("crds/pages/v4-api-with-updated-page-markdown.yaml");

    await test.step("Create API with markdown page", async () => {
      await kubectl.apply(WITH_PAGE);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Update markdown page content", async () => {
      await kubectl.apply(UPDATED);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Exported CRD reflects update", async () => {
      const crdYaml = await mapi.exportApiCrd(apiId);
      const crd = YAML.parse(crdYaml) as ExportedCrd;
      expect(crd.spec?.pages?.["markdown-page"]?.content).toContain("update");
    });

    await kubectl.del(UPDATED);
  });

  // ── GKO-280: GKO-created documentation is read-only ─────────
  // Re-applying the same CRD must not mutate the existing page content.

  test(`GKO-managed documentation is stable across re-apply ${XRAY.PAGES.V4_READ_ONLY_DOC} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-markdown-page";
    const WITH_PAGE = fixture("crds/pages/v4-api-with-page-markdown.yaml");

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

  // ── GKO-282: Inline documentation with PUBLIC visibility ────

  test(`Inline page with PUBLIC visibility is exposed ${XRAY.PAGES.V4_DOC_VISIBILITY_PUBLIC} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-public-page";
    const fixturePath = fixture("crds/pages/v4-api-public-page.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id;

    const crd = YAML.parse(await mapi.exportApiCrd(apiId)) as ExportedCrd;
    expect(crd.spec?.pages?.["public-page"]?.visibility).toBe("PUBLIC");

    await kubectl.del(fixturePath);
  });

  // GKO-283 (V4 spec.visibility only accepts PUBLIC) was skipped due to a
  // GKO product bug and is documented in "Batch 4 - Skipped Tests.md".

  // ── GKO-1470: Documentation is fully reconciled by the operator ─

  test(`Documentation is reconciled end-to-end ${XRAY.PAGES.V4_DOC_RECONCILED} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-markdown-page";
    const WITH_PAGE = fixture("crds/pages/v4-api-with-page-markdown.yaml");
    const WITHOUT_PAGE = fixture("crds/pages/v4-api-without-page-markdown.yaml");

    await kubectl.apply(WITH_PAGE);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id;

    await test.step("Page present after create", async () => {
      const crd = YAML.parse(await mapi.exportApiCrd(apiId)) as ExportedCrd;
      expect(crd.spec?.pages?.["markdown-page"]).toBeDefined();
    });

    await kubectl.apply(WITHOUT_PAGE);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    await test.step("Page removed after without-page re-apply", async () => {
      const crd = YAML.parse(await mapi.exportApiCrd(apiId)) as ExportedCrd;
      expect(crd.spec?.pages?.["markdown-page"]).toBeUndefined();
    });

    await kubectl.del(WITHOUT_PAGE);
  });
});
