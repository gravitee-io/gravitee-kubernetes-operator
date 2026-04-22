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
 * V4 documentation structure changes — batch 7 + 8.
 *
 * Xray tests:
 *   GKO-1469 / GKO-700: Renaming a folder + page in the CR reconciles on the
 *             APIM side, no duplicate pages are left behind, and the new
 *             names are reflected in the exported CRD. (GKO-700 covers the
 *             same scenario as GKO-1469 — both are tagged on the same test.)
 *   GKO-1467: Cross-version confirmation that a page marked
 *             visibility=PUBLIC is exposed as public — extends the existing
 *             V2 (GKO-199) and V4 (GKO-282) coverage by verifying both
 *             versions in a single run.
 */

import YAML from "yaml";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const ORIGINAL = "crds/pages/v4-api-with-page-markdown.yaml";
const RENAMED = "crds/pages/v4-api-renamed-page.yaml";
const V2_PUBLIC = "crds/pages/v2-api-public-page.yaml";
const V4_PUBLIC = "crds/pages/v4-api-public-page.yaml";
const API_NAME = "e2e-v4-markdown-page";

interface ExportedPage {
  name?: string;
  type?: string;
  visibility?: string;
}

interface ExportedCrd {
  spec?: {
    pages?: Record<string, ExportedPage>;
  };
}

test.describe("V4 Documentation — Rename & cross-version visibility", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(RENAMED)).catch(() => {});
    await kubectlSafe.del(fixture(ORIGINAL)).catch(() => {});
    await kubectlSafe.del(fixture(V2_PUBLIC)).catch(() => {});
    await kubectlSafe.del(fixture(V4_PUBLIC)).catch(() => {});
  });

  // ── GKO-1469: Folder + page rename reconciles ────────────────

  test(`V4 folder and page rename reconciles without duplicates ${XRAY.PAGES.V4_DOC_RENAME} ${XRAY.PAGES.V4_DOC_FOLDER_RENAME} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("Deploy with original names", async () => {
      await kubectl.apply(fixture(ORIGINAL));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const apiId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;

    await test.step("Apply renamed folder and page", async () => {
      await kubectl.apply(fixture(RENAMED));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Exported CRD reflects renamed names with no duplicates", async () => {
      await expect
        .poll(
          async () => {
            const crd = YAML.parse(await mapi.exportApiCrd(apiId)) as ExportedCrd;
            const pages = crd.spec?.pages ?? {};
            return {
              folderName: pages["renamed-folder"]?.name,
              pageName: pages["renamed-page"]?.name,
              oldFolderPresent: pages["markdown-folder"] !== undefined,
              oldPagePresent: pages["markdown-page"] !== undefined,
            };
          },
          { timeout: 30_000, intervals: [1_000] },
        )
        .toEqual({
          folderName: "Renamed Folder",
          pageName: "Renamed Page",
          oldFolderPresent: false,
          oldPagePresent: false,
        });
    });

    await kubectl.del(fixture(RENAMED));
  });

  // ── GKO-1467: Cross-version PUBLIC visibility ────────────────

  test(`PUBLIC visibility is reported by APIM for both V2 and V4 ${XRAY.PAGES.DOC_PUBLIC_ACCESS_CROSS_VERSION} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("V2 API with PUBLIC page reports PUBLIC", async () => {
      await kubectl.apply(fixture(V2_PUBLIC));
      await kubectl.waitForCondition("apidefinition", "e2e-v2-public-page", "Accepted");
      const v2Id = (
        await kubectl.getStatus<{ id: string }>("apidefinition", "e2e-v2-public-page")
      ).id;
      // APIM does not support CRD export for V2 APIs (see batch 5 notes),
      // so query the v1 Management pages endpoint directly for visibility.
      const response = await mapi.http.get<Array<{ visibility?: string }>>(
        mapi.http.managementV1Path(`/apis/${v2Id}/pages`),
      );
      expect(response.status).toBe(200);
      const publicPage = response.body.find((p) => p.visibility === "PUBLIC");
      expect(publicPage, "expected a V2 page with visibility=PUBLIC").toBeTruthy();
    });

    await test.step("V4 API with PUBLIC page reports PUBLIC", async () => {
      await kubectl.apply(fixture(V4_PUBLIC));
      await kubectl.waitForCondition("apiv4definition", "e2e-v4-public-page", "Accepted");
      const v4Id = (
        await kubectl.getStatus<{ id: string }>("apiv4definition", "e2e-v4-public-page")
      ).id;
      const crd = YAML.parse(await mapi.exportApiCrd(v4Id)) as ExportedCrd;
      expect(crd.spec?.pages?.["public-page"]?.visibility).toBe("PUBLIC");
    });

    await kubectl.del(fixture(V2_PUBLIC));
    await kubectl.del(fixture(V4_PUBLIC));
  });
});
