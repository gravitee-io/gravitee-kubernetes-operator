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
 * API version compatibility.
 *
 * Xray tests:
 *   GKO-1448: Minimal V2 and V4 CRs each reach Accepted=True and show
 *             the expected definition version on the APIM side.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const V2_MIN = "crds/api-definitions/v2-api-min.yaml";
const V4_MIN = "crds/api-v4-definitions/v4-api-min.yaml";

test.describe("API version compatibility — V2 & V4", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(V2_MIN)).catch(() => {});
    await kubectlSafe.del(fixture(V4_MIN)).catch(() => {});
  });

  test(`Minimal V2 and V4 CRs both reconcile ${XRAY.API_LIFECYCLE.V2_V4_COMPATIBILITY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("Deploy V4 minimal API", async () => {
      await kubectl.apply(fixture(V4_MIN));
      await kubectl.waitForCondition("apiv4definition", "e2e-v4-min", "Accepted");
      const v4Id = (
        await kubectl.getStatus<{ id: string }>("apiv4definition", "e2e-v4-min")
      ).id;
      const v4Api = await mapi.fetchApi(v4Id);
      expect(v4Api.definitionVersion).toBe("V4");
    });

    await test.step("Deploy V2 minimal API", async () => {
      await kubectl.apply(fixture(V2_MIN));
      await kubectl.waitForCondition("apidefinition", "e2e-v2-min", "Accepted");
      const v2Id = (
        await kubectl.getStatus<{ id: string }>("apidefinition", "e2e-v2-min")
      ).id;
      const v2Api = await mapi.fetchApi(v2Id);
      expect(v2Api.definitionVersion).toBe("V2");
    });

    await kubectl.del(fixture(V2_MIN));
    await kubectl.del(fixture(V4_MIN));
  });
});
