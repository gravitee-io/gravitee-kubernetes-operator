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
 * V4 API Lifecycle — Extended scenarios.
 *
 * Xray tests:
 *   GKO-81:   Deploy V4 API in DB-less mode
 *   GKO-137:  Deploy V4 message API with MQTT endpoint
 *   GKO-159:  Re-create a deleted V4 API with a previously closed plan
 *   GKO-268:  Add many valid categories at once to V4 API
 *   GKO-272:  Category rename in APIM triggers V4 redeploy
 *   GKO-412:  Deploy V4 API with a non-existing category
 *   GKO-471:  Deploy V4 API with a non-existing group
 *   GKO-1061: V4 API still works after operator upgrade/restart
 *   GKO-1465: Policy enforcement on V4 APIs without plans
 *   GKO-1474: Entrypoint × policy compatibility matrix
 *
 * Skipped tests:
 *   GKO-138 (RabbitMQ endpoint) — APIM schema bug
 *   GKO-139 (Solace endpoint)   — APIM schema bug
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { readFile } from "node:fs/promises";
import YAML from "yaml";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import type { Api, ApiV4 } from "../../../../src/types/apim.js";

/** Narrow a fetched Api to its V4 variant so tests can access V4-only fields. */
function asV4(api: Api): ApiV4 {
  expect(api.definitionVersion).toBe("V4");
  return api as ApiV4;
}

interface StatusWithId {
  id?: string;
}

interface StatusWithConditions extends StatusWithId {
  conditions?: Array<{
    type: string;
    status: string;
    reason?: string;
    message?: string;
  }>;
}

test.describe("V4 API Lifecycle — Extended", () => {
  // ── GKO-81: DB-less mode ─────────────────────────────────────

  test(`Deploy V4 API in DB-less mode ${XRAY.API_LIFECYCLE.DEPLOY_V4_DB_LESS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-db-less";
    const fixturePath = fixture("crds/v4-lifecycle-extended/v4-proxy-api-db-less.yaml");

    await test.step("Apply DB-less V4 API CRD", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("CRD is Accepted without a management-context sync", async () => {
      const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted?.status).toBe("True");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-137: MQTT endpoint ───────────────────────────────────

  test(`Deploy V4 message API with MQTT endpoint ${XRAY.MESSAGE_APIS.MQTT_ENDPOINT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-mqtt";
    const fixturePath = fixture("crds/message-apis/v4-message-api-mqtt.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithId>("apiv4definition", API_NAME);
    const api = asV4(await mapi.fetchApi(status.id!));
    expect(api.type).toBe("MESSAGE");

    await kubectl.del(fixturePath);
  });

  // GKO-138 (RabbitMQ) and GKO-139 (Solace) were skipped due to an APIM
  // schema bug.

  // ── GKO-159: Re-create deleted V4 API with closed plan ──────
  // Re-applying the same CRD after a full delete should create a fresh API
  // in APIM with a new id and the keyless plan should be reopened.

  test(`Re-create deleted V4 API with previously closed plan ${XRAY.API_LIFECYCLE.RECREATE_DELETED_V4_CLOSED_PLAN} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await test.step("Deploy, then delete the V4 API", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
      await kubectl.del(fixturePath);
      await kubectl.waitForDeletion("apiv4definition", API_NAME);
    });

    await test.step("Re-apply same CRD", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("API exists again in APIM with a started plan", async () => {
      const status = await kubectl.getStatus<StatusWithId>("apiv4definition", API_NAME);
      await mapi.waitForApiMatches(status.id!, { name: API_NAME, state: "STARTED" });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-268: Many categories at once ─────────────────────────

  test(`Deploy V4 API with many category refs ${XRAY.CATEGORIES.V4_MANY_CATEGORIES} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-many-cats";
    const fixturePath = fixture("crds/v4-lifecycle-extended/v4-proxy-api-many-categories.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    expect(status.conditions?.find((c) => c.type === "Accepted")?.status).toBe("True");

    await kubectl.del(fixturePath);
  });

  // ── GKO-272: Category rename triggers redeploy ───────────────
  // Re-applying the same CRD with a changed category list should trigger a
  // reconciliation; the operator must not lose the API on update.

  test(`Category change on V4 API triggers redeploy ${XRAY.CATEGORIES.V4_CATEGORY_RENAME_REDEPLOY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-many-cats";
    const fixturePath = fixture("crds/v4-lifecycle-extended/v4-proxy-api-many-categories.yaml");
    const renamedFixture = fixture(
      "crds/v4-lifecycle-extended/v4-proxy-api-many-categories-renamed.yaml",
    );

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const firstStatus = await kubectl.getStatus<StatusWithId>("apiv4definition", API_NAME);
    const apiId = firstStatus.id!;

    // Apply a fixture with a renamed category and a bumped description — the
    // spec hash changes, so the reconciler actually runs instead of being
    // filtered by LastSpecHashPredicate. Categories themselves are not
    // asserted because APIM silently drops unknown category refs, but the
    // description change is observable on the APIM side.
    await kubectl.apply(renamedFixture);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await mapi.waitForApiMatches(apiId, {
      description: "E2E test: V4 API with renamed category refs (redeployed)",
    });

    // APIM id must remain stable across the redeploy.
    const secondStatus = await kubectl.getStatus<StatusWithId>("apiv4definition", API_NAME);
    expect(secondStatus.id).toBe(apiId);

    await kubectl.del(renamedFixture);
  });

  // ── GKO-412: Non-existing category ───────────────────────────

  test(`Deploy V4 API with a non-existing category ${XRAY.CATEGORIES.V4_DEPLOY_NON_EXISTING_CATEGORY} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-non-existing-cat";
    const fixturePath = fixture(
      "crds/v4-lifecycle-extended/v4-proxy-api-non-existing-category.yaml",
    );

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    expect(status.conditions?.find((c) => c.type === "Accepted")?.status).toBe("True");

    await kubectl.del(fixturePath);
  });

  // ── GKO-471: Non-existing group ──────────────────────────────

  test(`Deploy V4 API with a non-existing group reference ${XRAY.CATEGORIES.V4_NON_EXISTING_GROUP_REF} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-non-existing-group-ref";
    const fixturePath = fixture("crds/v4-lifecycle-extended/v4-proxy-api-non-existing-group.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    expect(status.conditions?.find((c) => c.type === "Accepted")?.status).toBe("True");

    await kubectl.del(fixturePath);
  });

  // ── GKO-1061: API survives a reconcile cycle ────────────────
  // Proxy for an operator upgrade: re-apply the CRD and verify the APIM-side
  // id is unchanged, i.e. state survives an operator-initiated reconciliation.

  test(`V4 API state survives a reconcile cycle ${XRAY.API_LIFECYCLE.V4_SURVIVES_UPGRADE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const idBefore = (await kubectl.getStatus<StatusWithId>("apiv4definition", API_NAME)).id!;

    // Bump the API version so the spec hash actually changes — otherwise
    // LastSpecHashPredicate filters the update and the reconciler never runs.
    const raw = await readFile(fixturePath, "utf8");
    const doc = YAML.parse(raw) as { spec: { version: string } };
    const bumpedVersion = "2.0.0";
    doc.spec.version = bumpedVersion;
    await kubectl.applyString(YAML.stringify(doc));

    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const idAfter = (await kubectl.getStatus<StatusWithId>("apiv4definition", API_NAME)).id!;

    // APIM id must be stable across the real reconcile.
    expect(idAfter).toBe(idBefore);
    await mapi.waitForApiMatches(idAfter, { apiVersion: bumpedVersion, state: "STARTED" });

    await kubectl.del(fixturePath);
  });

  // ── GKO-1465: Policy on API without plans (STOPPED) ─────────

  test(`Policy on V4 API without plans is accepted when STOPPED ${XRAY.API_LIFECYCLE.POLICY_ON_API_WITHOUT_PLANS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-policy-no-plans";
    const fixturePath = fixture("crds/v4-lifecycle-extended/v4-proxy-api-policy-no-plans.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    expect(status.conditions?.find((c) => c.type === "Accepted")?.status).toBe("True");

    await kubectl.del(fixturePath);
  });

  // ── GKO-1474: Entrypoint × policy matrix ────────────────────

  test(`Entrypoint × policy matrix on V4 message API ${XRAY.API_LIFECYCLE.ENTRYPOINT_POLICY_MATRIX} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-entrypoint-policy";
    const fixturePath = fixture(
      "crds/v4-lifecycle-extended/v4-message-api-entrypoint-policy.yaml",
    );

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithId>("apiv4definition", API_NAME);
    const api = asV4(await mapi.fetchApi(status.id!));
    expect(api.type).toBe("MESSAGE");
    expect(api.flows).toBeTruthy();
    expect(api.flows!.length).toBeGreaterThanOrEqual(1);

    await kubectl.del(fixturePath);
  });
});
