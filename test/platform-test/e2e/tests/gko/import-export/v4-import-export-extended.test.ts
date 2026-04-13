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
 * V4 Import/Export — Extended scenarios (batch 4).
 *
 * Xray tests:
 *   GKO-237:  No email sent on exporting CRD for v4 APIs
 *   GKO-239:  API metadata handled by v4 import endpoint
 *   GKO-1472: Export → import round-trip integrity for V4 APIs
 *   GKO-1927: Terraform import/export parity
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import YAML from "yaml";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

interface ExportedCrd {
  spec?: {
    name?: string;
    labels?: string[];
    flows?: unknown[];
  };
}

test.describe("V4 Import/Export — Extended", () => {
  // ── GKO-237: No email on export ─────────────────────────────
  // Exporting a V4 CRD must not trigger notification emails. We can't assert
  // the absence of an email from here, but we can assert the export endpoint
  // returns a well-formed CRD without side effects on the API.

  test(`Export V4 CRD does not disturb API state ${XRAY.IMPORT_EXPORT.V4_NO_EMAIL_ON_EXPORT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-export";
    const fixturePath = fixture("crds/import-export/v4-proxy-api-export.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id;

    const stateBefore = (await mapi.fetchApi(apiId)).state;
    const crdYaml = await mapi.exportApiCrd(apiId);
    expect(crdYaml.length).toBeGreaterThan(0);
    const stateAfter = (await mapi.fetchApi(apiId)).state;
    expect(stateAfter).toBe(stateBefore);

    await kubectl.del(fixturePath);
  });

  // ── GKO-239: Metadata handled by v4 import endpoint ─────────

  test(`API metadata survives V4 import ${XRAY.IMPORT_EXPORT.V4_METADATA_ON_IMPORT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-export";
    const fixturePath = fixture("crds/import-export/v4-proxy-api-export.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id;

    const crd = YAML.parse(await mapi.exportApiCrd(apiId)) as ExportedCrd;
    expect(crd.spec?.name).toBe(API_NAME);
    expect(crd.spec?.labels).toContain("e2e-export-label");

    await kubectl.del(fixturePath);
  });

  // ── GKO-1472: Export → import round-trip ────────────────────
  // Apply the EXPORTED YAML on top of the live API (not the original fixture).
  // This is what actually validates the round-trip: if the export drops fields
  // or produces an unimportable CRD, the apply will fail or the reconciliation
  // will move off Accepted — neither of which the previous re-apply-original
  // formulation could catch.

  test(`V4 CRD export/import round-trip preserves identity ${XRAY.IMPORT_EXPORT.V4_EXPORT_IMPORT_ROUND_TRIP} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-export";
    const fixturePath = fixture("crds/import-export/v4-proxy-api-export.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const firstId = (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id;

    const exportedYaml = await mapi.exportApiCrd(firstId);
    const exportedCrd = YAML.parse(exportedYaml) as ExportedCrd;
    expect(exportedCrd.spec?.name).toBe(API_NAME);
    expect(Array.isArray(exportedCrd.spec?.flows)).toBe(true);
    expect((exportedCrd.spec?.flows ?? []).length).toBeGreaterThanOrEqual(1);

    // Apply the EXPORTED YAML over the live API. The operator reconciles it
    // into the existing APIM record, so the id must be stable end-to-end.
    await kubectl.applyString(exportedYaml);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    // The exported YAML pins `spec.id`, so the operator reuses it via pickID().
    // status.id is only populated when spec.id was absent, so read either.
    const reimported = await kubectl.get<{
      spec: { id?: string };
      status: { id?: string };
    }>("apiv4definition", API_NAME);
    const secondId = reimported.spec.id ?? reimported.status.id;
    expect(secondId).toBe(firstId);

    // And APIM must still know the same API under that id.
    await mapi.waitForApiMatches(firstId, { name: API_NAME });

    await kubectl.del(fixturePath);
  });

  // ── GKO-1927: Terraform import/export parity ────────────────
  // Ensure the exported CRD from a Terraform-created API matches a straight
  // CRD-managed one. We approximate by checking that export returns a CRD
  // with flows preserved (parity with what Terraform would produce).

  test(`Exported V4 CRD preserves flows (parity target) ${XRAY.IMPORT_EXPORT.V4_TERRAFORM_IMPORT_EXPORT_PARITY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-export";
    const fixturePath = fixture("crds/import-export/v4-proxy-api-export.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id;

    const crd = YAML.parse(await mapi.exportApiCrd(apiId)) as ExportedCrd;
    expect(Array.isArray(crd.spec?.flows)).toBe(true);
    expect((crd.spec?.flows ?? []).length).toBeGreaterThanOrEqual(1);

    await kubectl.del(fixturePath);
  });
});
