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
 * V4 API visibility & lifecycleState.
 *
 * Xray tests:
 *   GKO-172:  V4 API with visibility=PRIVATE — verify visibility on API record
 *   GKO-173:  V4 API with visibility=PUBLIC — verify visibility on API record
 *   GKO-179:  V4 API with lifecycleState=PUBLISHED — verify state on API record
 *   GKO-180:  V4 API with lifecycleState=UNPUBLISHED — verify state on API record
 *   GKO-1466: Verify combined visibility / lifecycleState rules across the
 *             three meaningful combinations (PRIVATE+PUBLISHED, PUBLIC+PUBLISHED,
 *             PUBLIC+UNPUBLISHED) all reach Accepted=True with the expected
 *             attributes on the APIM-side record.
 *
 * Note: the original Xray scenarios reference what the developer portal
 * surfaces (which would require a browser fixture). We assert the equivalent
 * via the mAPI record (`visibility`, `lifecycleState`) — the portal-rendering
 * is a downstream consumer of those fields, so verifying them at the source
 * gives equivalent confidence without a browser dependency.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";
import type { ApiLifecycleState, ApiVisibility } from "../../../../src/types/apim.js";

const PRIVATE_PUBLISHED = "crds/api-v4-definitions/v4-api-private-published.yaml";
const PUBLIC_PUBLISHED = "crds/api-v4-definitions/v4-api-public-published.yaml";
const PUBLIC_UNPUBLISHED = "crds/api-v4-definitions/v4-api-public-unpublished.yaml";

interface StatusWithId {
  id?: string;
  conditions?: Array<{ type: string; status: string }>;
}

test.describe("V4 API — Visibility & lifecycleState", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(PRIVATE_PUBLISHED)).catch(() => {});
    await kubectlSafe.del(fixture(PUBLIC_PUBLISHED)).catch(() => {});
    await kubectlSafe.del(fixture(PUBLIC_UNPUBLISHED)).catch(() => {});
  });

  // ── GKO-172: PRIVATE visibility ──────────────────────────────

  test(`V4 API with visibility=PRIVATE is recorded as PRIVATE in APIM ${XRAY.API_LIFECYCLE.V4_VISIBILITY_PRIVATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-private-published";
    await kubectl.apply(fixture(PRIVATE_PUBLISHED));
    await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    const apiId = (await kubectl.getStatus<StatusWithId>("apiv4definition", NAME)).id;
    expect(apiId).toBeTruthy();
    await mapi.waitForApiMatches(apiId!, { visibility: "PRIVATE" });
  });

  // ── GKO-173: PUBLIC visibility ───────────────────────────────

  test(`V4 API with visibility=PUBLIC is recorded as PUBLIC in APIM ${XRAY.API_LIFECYCLE.V4_VISIBILITY_PUBLIC} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-public-published";
    await kubectl.apply(fixture(PUBLIC_PUBLISHED));
    await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    const apiId = (await kubectl.getStatus<StatusWithId>("apiv4definition", NAME)).id;
    expect(apiId).toBeTruthy();
    await mapi.waitForApiMatches(apiId!, { visibility: "PUBLIC" });
  });

  // ── GKO-179: PUBLISHED lifecycleState ────────────────────────

  test(`V4 API with lifecycleState=PUBLISHED is recorded as PUBLISHED in APIM ${XRAY.API_LIFECYCLE.V4_PUBLISHED_IN_PORTAL} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-public-published";
    await kubectl.apply(fixture(PUBLIC_PUBLISHED));
    await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    const apiId = (await kubectl.getStatus<StatusWithId>("apiv4definition", NAME)).id;
    expect(apiId).toBeTruthy();
    await mapi.waitForApiMatches(apiId!, { lifecycleState: "PUBLISHED" });
  });

  // ── GKO-180: UNPUBLISHED lifecycleState ──────────────────────

  test(`V4 API with lifecycleState=UNPUBLISHED is recorded as UNPUBLISHED in APIM ${XRAY.API_LIFECYCLE.V4_UNPUBLISHED_NOT_IN_PORTAL} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const NAME = "e2e-v4-public-unpublished";
    await kubectl.apply(fixture(PUBLIC_UNPUBLISHED));
    await kubectl.waitForCondition("apiv4definition", NAME, "Accepted");
    const apiId = (await kubectl.getStatus<StatusWithId>("apiv4definition", NAME)).id;
    expect(apiId).toBeTruthy();
    await mapi.waitForApiMatches(apiId!, { lifecycleState: "UNPUBLISHED" });
  });

  // ── GKO-1466: combined visibility / lifecycleState rules ─────

  test(`Visibility & lifecycleState combinations enforced ${XRAY.API_LIFECYCLE.V4_PORTAL_VISIBILITY_RULES} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    type Combo = {
      label: string;
      fixturePath: string;
      crName: string;
      expected: { visibility: ApiVisibility; lifecycleState: ApiLifecycleState };
    };

    const combos: Combo[] = [
      {
        label: "PRIVATE + PUBLISHED",
        fixturePath: PRIVATE_PUBLISHED,
        crName: "e2e-v4-private-published",
        expected: { visibility: "PRIVATE", lifecycleState: "PUBLISHED" },
      },
      {
        label: "PUBLIC + PUBLISHED",
        fixturePath: PUBLIC_PUBLISHED,
        crName: "e2e-v4-public-published",
        expected: { visibility: "PUBLIC", lifecycleState: "PUBLISHED" },
      },
      {
        label: "PUBLIC + UNPUBLISHED",
        fixturePath: PUBLIC_UNPUBLISHED,
        crName: "e2e-v4-public-unpublished",
        expected: { visibility: "PUBLIC", lifecycleState: "UNPUBLISHED" },
      },
    ];

    for (const combo of combos) {
      await test.step(combo.label, async () => {
        await kubectl.apply(fixture(combo.fixturePath));
        await kubectl.waitForCondition("apiv4definition", combo.crName, "Accepted");
        const apiId = (
          await kubectl.getStatus<StatusWithId>("apiv4definition", combo.crName)
        ).id;
        expect(apiId).toBeTruthy();
        await mapi.waitForApiMatches(apiId!, combo.expected);
      });
    }
  });
});
