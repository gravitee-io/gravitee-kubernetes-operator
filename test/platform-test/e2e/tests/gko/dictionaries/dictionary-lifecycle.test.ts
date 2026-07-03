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
 * Dictionaries Lifecycle tests.
 *
 * Xray tests:
 *   GKO-2905: Delete a dictionary
 *   GKO-2912: Admission webhook rejects DYNAMIC dictionary with manual field set
 *
 * GKO-2903 (manual resolve) and the DYNAMIC resolve + lifecycle tests
 * (GKO-2904/2909/2910/2911) are now cross-provisioner journeys under
 * tests/user-journeys/ (api-references-dictionary-property and
 * manage-dynamic-dictionary), so they are intentionally not duplicated here.
 * What remains here is GKO-specific and has no Terraform counterpart: plain CR
 * deletion (finalizer release) and the K8s admission webhook.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

const DICT_NAME = "e2e-dict-manual";

test.describe("Dictionaries — Lifecycle @since-4.12", () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. The manual dictionary is the only resource these tests create
  // (the admission test's CR is rejected at apply). del() ignores errors.
  test.afterEach(async () => {
    await kubectl.del(fixture("dictionaries/dictionary-manual/crd.yaml")).catch(() => {});
  });

  // ── GKO-2563: Delete a dictionary ──────────────────────────────

  test(`Delete a dictionary ${XRAY.DICTIONARIES.DELETE_DICTIONARY} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const dictFixture = fixture("dictionaries/dictionary-manual/crd.yaml");

    await test.step("Create dictionary", async () => {
      await kubectl.apply(dictFixture);
      await kubectl.waitForCondition("dictionary", DICT_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("dictionary", DICT_NAME);
    expect(status.id).toBeTruthy();

    await test.step("Delete the dictionary CRD", async () => {
      await kubectl.del(dictFixture);
      await kubectl.waitForDeletion("dictionary", DICT_NAME);
    });
  });

  // ── Admission: DYNAMIC + manual field is rejected ─────────────

  test(`Admission webhook rejects DYNAMIC dictionary with manual field set ${XRAY.DICTIONARIES.ADMISSION_REJECTS_DYNAMIC_WITH_MANUAL} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const invalidFixture = fixture("dictionaries/dictionary-dynamic-invalid/crd.yaml");

    const stderr = await kubectl.applyExpectFailure(invalidFixture);
    expect(stderr).toMatch(/dictionary type is DYNAMIC but 'manual' field is set/);
  });
});
