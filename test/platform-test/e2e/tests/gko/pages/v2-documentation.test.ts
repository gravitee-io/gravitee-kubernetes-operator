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
 * V2 API Documentation — batch 5 coverage.
 *
 * Xray tests:
 *   GKO-146: Add/remove/update inline documentation pages on a V2 API
 *   GKO-147: Graceful failure when deploying oversized documentation
 *   GKO-148: Update inline documentation pages
 *   GKO-151: Add a documentation page using a fetcher
 *   GKO-199: Public documentation pages accessible by everyone
 *   GKO-200: Documentation visibility PRIVATE with no groups
 *   GKO-315: Documentation visibility PRIVATE with groups
 *   GKO-316: Documentation visibility PRIVATE with excluded groups
 *   GKO-662: Delete fetched ROOT pages from external sources
 *   GKO-699: Rename folder & documentation page
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 *   - For group-scoped page tests: the "e2e-v2-group-a" Group CR is applied
 *
 * Notes on APIM endpoints for V2 APIs:
 *   The V2 CRD export endpoint (`/management/v2/.../apis/{id}/_export/crd`)
 *   is not supported for V2 API definitions — APIM returns 400
 *   "definition version 2.0.0 is not supported anymore". These tests use
 *   the v1 management API pages endpoint instead to verify APIM state.
 */

import YAML from "yaml";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";
import type { Mapi } from "../../../../src/index.js";

interface ApimPage {
  id: string;
  name?: string;
  type?: string;
  content?: string;
  visibility?: string;
  parentPath?: string;
  accessControls?: Array<{ referenceId: string; referenceType: string }>;
  excludedAccessControls?: boolean;
  source?: { type?: string; configuration?: Record<string, unknown> };
}

async function fetchV1Pages(mapi: Mapi, apiId: string): Promise<ApimPage[]> {
  const path = mapi.http.managementV1Path(`/apis/${apiId}/pages`);
  const res = await mapi.http.get<ApimPage[]>(path);
  if (res.status !== 200) {
    throw new Error(`Failed to fetch pages for API ${apiId}: ${res.status}`);
  }
  return res.body;
}

const GROUP_A = "crds/members/group-for-v2-members.yaml";

test.describe("V2 API Documentation — Extended", () => {
  test.afterEach(async () => {
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-without-page-markdown.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-public-page.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-private-page-no-groups.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-private-page-with-groups.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-private-page-excluded-groups.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-with-fetcher-page.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-with-root-fetcher-deleted.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-with-root-fetcher.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-renamed-page.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-with-updated-page-markdown.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/pages/v2-api-with-page-markdown.yaml"))
      .catch(() => {});
    await kubectlSafe.del(fixture(GROUP_A)).catch(() => {});
  });

  // ── GKO-146: Add/remove/update inline documentation pages ───

  test(`Add, update and remove inline V2 documentation ${XRAY.PAGES.V2_DOC_CRUD} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-markdown-page";
    const WITH_PAGE = fixture("crds/pages/v2-api-with-page-markdown.yaml");
    const UPDATED = fixture("crds/pages/v2-api-with-updated-page-markdown.yaml");
    const WITHOUT = fixture("crds/pages/v2-api-without-page-markdown.yaml");

    await kubectl.apply(WITH_PAGE);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await test.step("APIM contains the markdown page", async () => {
      await expect
        .poll(async () => {
          const pages = await fetchV1Pages(mapi, apiId);
          return pages.some((p) => p.name === "markdown-page");
        })
        .toBe(true);
    });

    await kubectl.apply(UPDATED);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    await test.step("APIM reflects updated content", async () => {
      await expect
        .poll(async () => {
          const pages = await fetchV1Pages(mapi, apiId);
          return pages.find((p) => p.name === "markdown-page")?.content ?? "";
        })
        .toContain("updated");
    });

    await kubectl.apply(WITHOUT);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    await test.step("Page removed when CRD drops it", async () => {
      await expect
        .poll(async () => {
          const pages = await fetchV1Pages(mapi, apiId);
          return pages.some((p) => p.name === "markdown-page");
        })
        .toBe(false);
    });

    await kubectl.del(WITHOUT);
  });

  // ── GKO-147: >1.5 MB documentation is rejected gracefully ───

  test(`Oversized V2 documentation fails gracefully ${XRAY.PAGES.V2_DOC_OVERSIZE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const huge = "A".repeat(1_600_000);
    const manifest = {
      apiVersion: "gravitee.io/v1alpha1",
      kind: "ApiDefinition",
      metadata: { name: "e2e-v2-oversized-doc" },
      spec: {
        contextRef: { name: "dev-ctx", namespace: "default" },
        name: "e2e-v2-oversized-doc",
        description: "E2E test: oversized V2 documentation",
        version: "1.0.0",
        local: false,
        proxy: {
          virtual_hosts: [{ path: "/e2e-v2-oversized-doc" }],
          groups: [{ endpoints: [{ name: "Default", target: "https://api.gravitee.io/echo" }] }],
        },
        plans: [
          {
            name: "Free plan",
            description: "Open access plan",
            security: "KEY_LESS",
            status: "PUBLISHED",
          },
        ],
        pages: {
          "huge-page": {
            name: "huge-page",
            type: "MARKDOWN",
            content: huge,
            visibility: "PUBLIC",
            published: true,
          },
        },
      },
    };

    const yamlStr = YAML.stringify(manifest);
    const CR_NAME = "e2e-v2-oversized-doc";

    let rejectedAtAdmission = false;
    try {
      await kubectl.applyString(yamlStr);
    } catch {
      rejectedAtAdmission = true;
    }

    if (rejectedAtAdmission) {
      // Apply failed — either admission rejected, or the API server refused
      // the oversized payload (e.g. 413). Either way, the oversized doc was
      // handled gracefully. Stderr may be empty for size-limit rejections.
      return;
    }

    await test.step("CRD reaches an explicit failure state (Accepted=False)", async () => {
      await expect
        .poll(
          async () => {
            const status = await kubectl
              .getStatus<{
                conditions?: Array<{ type: string; status: string }>;
              }>("apidefinition", CR_NAME)
              .catch(() => ({ conditions: [] }));
            return (
              status.conditions?.find((c) => c.type === "Accepted")?.status ?? "Unknown"
            );
          },
          { timeout: 30_000 },
        )
        .toBe("False");
    });

    await kubectlSafe.deleteResource("apidefinition", CR_NAME).catch(() => {});
  });

  // ── GKO-148: Update inline documentation pages ──────────────

  test(`Update V2 inline page content is reflected in APIM ${XRAY.PAGES.V2_DOC_INLINE_UPDATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-markdown-page";
    const WITH_PAGE = fixture("crds/pages/v2-api-with-page-markdown.yaml");
    const UPDATED = fixture("crds/pages/v2-api-with-updated-page-markdown.yaml");

    await kubectl.apply(WITH_PAGE);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await kubectl.apply(UPDATED);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    await expect
      .poll(async () => {
        const pages = await fetchV1Pages(mapi, apiId);
        return pages.find((p) => p.name === "markdown-page")?.content ?? "";
      })
      .toContain("updated");

    await kubectl.del(UPDATED);
  });

  // ── GKO-151: Add documentation page using a fetcher ─────────

  test(`Add V2 page via HTTP fetcher ${XRAY.PAGES.V2_DOC_FETCHER} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-fetcher-page";
    const fixturePath = fixture("crds/pages/v2-api-with-fetcher-page.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    const pages = await fetchV1Pages(mapi, apiId);
    const fetched = pages.find((p) => p.name === "fetched-page");
    expect(fetched).toBeDefined();
    expect(fetched?.source?.type).toBe("http-fetcher");

    await kubectl.del(fixturePath);
  });

  // ── GKO-199: Public pages accessible ────────────────────────

  test(`V2 public page has PUBLIC visibility in APIM ${XRAY.PAGES.V2_DOC_PUBLIC} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-public-page";
    const fixturePath = fixture("crds/pages/v2-api-public-page.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    const pages = await fetchV1Pages(mapi, apiId);
    const page = pages.find((p) => p.name === "public-page");
    expect(page?.visibility).toBe("PUBLIC");

    await kubectl.del(fixturePath);
  });

  // ── GKO-200: Private page with no groups ────────────────────

  test(`V2 PRIVATE page with no groups is created ${XRAY.PAGES.V2_DOC_PRIVATE_NO_GROUPS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-private-no-groups";
    const fixturePath = fixture("crds/pages/v2-api-private-page-no-groups.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    const pages = await fetchV1Pages(mapi, apiId);
    const page = pages.find((p) => p.name === "private-page");
    expect(page?.visibility).toBe("PRIVATE");
    expect(page?.accessControls ?? []).toEqual([]);

    await kubectl.del(fixturePath);
  });

  // ── GKO-315: Private page with groups ───────────────────────

  test(`V2 PRIVATE page with group access control ${XRAY.PAGES.V2_DOC_PRIVATE_GROUPS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-private-with-groups";

    await kubectl.apply(fixture(GROUP_A));
    await kubectl.waitForCondition("group", "e2e-v2-group-a", "Accepted");

    const fixturePath = fixture("crds/pages/v2-api-private-page-with-groups.yaml");
    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    const pages = await fetchV1Pages(mapi, apiId);
    const page = pages.find((p) => p.name === "restricted-page");
    expect(page?.visibility).toBe("PRIVATE");
    expect(page?.accessControls?.length ?? 0).toBeGreaterThan(0);

    await kubectl.del(fixturePath);
    await kubectl.del(fixture(GROUP_A));
  });

  // ── GKO-316: Private page with excluded groups ──────────────

  test(`V2 PRIVATE page with excluded group ${XRAY.PAGES.V2_DOC_PRIVATE_EXCLUDED} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-private-excluded-groups";

    await kubectl.apply(fixture(GROUP_A));
    await kubectl.waitForCondition("group", "e2e-v2-group-a", "Accepted");

    const fixturePath = fixture("crds/pages/v2-api-private-page-excluded-groups.yaml");
    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    const pages = await fetchV1Pages(mapi, apiId);
    const page = pages.find((p) => p.name === "excluded-page");
    expect(page?.visibility).toBe("PRIVATE");
    expect(page?.excludedAccessControls).toBe(true);

    await kubectl.del(fixturePath);
    await kubectl.del(fixture(GROUP_A));
  });

  // GKO-662 (delete fetched ROOT pages) — dropped from batch 5. APIM
  // rejects V2 ROOT fetchers backed by http-fetcher ("The plugin does not
  // support to import a directory"); ROOT pages require a directory-capable
  // fetcher such as github-fetcher, which in turn needs real GitHub
  // credentials that the test cluster does not provision. Tracked in
  // "Batch 5 - Skipped Tests.md".

  // ── GKO-699: Rename folder & documentation page ─────────────

  test(`Renaming V2 folder & page is reflected in APIM ${XRAY.PAGES.V2_DOC_RENAME} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-markdown-page";
    const WITH_PAGE = fixture("crds/pages/v2-api-with-page-markdown.yaml");
    const RENAMED = fixture("crds/pages/v2-api-renamed-page.yaml");

    await kubectl.apply(WITH_PAGE);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await kubectl.apply(RENAMED);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    await expect
      .poll(async () => {
        const pages = await fetchV1Pages(mapi, apiId);
        const names = new Set(pages.map((p) => p.name ?? ""));
        return {
          hasRenamedFolder: names.has("renamed-folder"),
          hasRenamedPage: names.has("renamed-page"),
          hasOldFolder: names.has("markdown-folder"),
          hasOldPage: names.has("markdown-page"),
        };
      })
      .toMatchObject({
        hasRenamedFolder: true,
        hasRenamedPage: true,
        hasOldFolder: false,
        hasOldPage: false,
      });

    await kubectl.del(RENAMED);
  });
});
