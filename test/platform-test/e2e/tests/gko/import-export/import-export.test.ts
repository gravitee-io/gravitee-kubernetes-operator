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
 * Import/Export: CRD round-trip tests.
 *
 * Xray tests:
 *   GKO-229: Export V4 API CRD — verify API fields in APIM
 *   GKO-218: Import exported V4 API CRD — verify creation
 *   GKO-301: Export V2 API CRD — verify API fields in APIM
 *   GKO-303: Import exported V2 API CRD — verify creation
 *   GKO-231: K8s-compliant names with special characters in spec.name
 *   GKO-93:  Exported V4 policies/flows are reflected in the CRD
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("Import/Export — CRD Round-trips", () => {
  // ── GKO-229: Export V4 API CRD ──────────────────────────────

  test(`Export V4 API CRD and verify fields in APIM ${XRAY.IMPORT_EXPORT.EXPORT_V4_CRD} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-export";
    const fixturePath = fixture("crds/import-export/v4-proxy-api-export.yaml");

    await test.step("Deploy V4 API with flows, plans, and labels", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API in APIM has expected fields", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api.name).toBe(API_NAME);
      expect(api.labels).toContain("e2e-export-label");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-218: Import exported V4 API CRD ────────────────────

  test(`Import V4 API CRD creates API in APIM ${XRAY.IMPORT_EXPORT.IMPORT_V4_CRD} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-export";
    const fixturePath = fixture("crds/import-export/v4-proxy-api-export.yaml");

    await test.step("Apply V4 API CRD (simulating import)", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API exists in APIM after import", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-301: Export V2 API CRD ──────────────────────────────

  test(`Export V2 API CRD and verify fields in APIM ${XRAY.IMPORT_EXPORT.EXPORT_V2_CRD} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-export";
    const fixturePath = fixture("crds/import-export/v2-api-export.yaml");

    await test.step("Deploy V2 API with flows and plans", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME);
    const apiId = status.id;

    await test.step("API in APIM has expected fields", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api.name).toBe(API_NAME);
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-303: Import exported V2 API CRD ────────────────────

  test(`Import V2 API CRD creates API in APIM ${XRAY.IMPORT_EXPORT.IMPORT_V2_CRD} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-export";
    const fixturePath = fixture("crds/import-export/v2-api-export.yaml");

    await test.step("Apply V2 API CRD (simulating import)", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME);
    const apiId = status.id;

    await test.step("API exists in APIM after import", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-231: K8s-compliant names ───────────────────────────

  test(`K8s-compliant metadata.name with special characters in spec.name ${XRAY.IMPORT_EXPORT.K8S_COMPLIANT_NAMES} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const K8S_NAME = "my-e2e-api-v4-test";
    const SPEC_NAME = "My E2E API (v4) Test";
    const fixturePath = fixture("crds/import-export/v4-proxy-api-with-special-name.yaml");

    await test.step("Deploy API with special characters in spec.name", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", K8S_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", K8S_NAME);
    const apiId = status.id;

    await test.step("API name in APIM matches spec.name (with special chars)", async () => {
      await mapi.assertApiMatches(apiId, { name: SPEC_NAME });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-93: Exported V4 policies reflected in CRD ──────────

  test(`Exported V4 API policies and flows are present ${XRAY.IMPORT_EXPORT.EXPORTED_POLICIES_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-export";
    const fixturePath = fixture("crds/import-export/v4-proxy-api-export.yaml");

    await test.step("Deploy V4 API with policy flows", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API in APIM has flows with policies", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api.flows).toBeTruthy();
      expect(api.flows!.length).toBeGreaterThanOrEqual(1);
    });

    await kubectl.del(fixturePath);
  });
});
