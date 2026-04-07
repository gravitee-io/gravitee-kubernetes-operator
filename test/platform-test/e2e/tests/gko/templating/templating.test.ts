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
 * Templating: ConfigMap & Secret reference tests.
 *
 * GKO supports Go templates in CRD fields referencing K8s ConfigMaps/Secrets.
 * Pattern: {{ .Values.secret.my-secret.my-key }} or {{ .Values.configmap.my-cm.my-key }}
 *
 * Xray tests:
 *   GKO-683: V4 API description from ConfigMap value
 *   GKO-682: V2 API description from Secret value
 *   GKO-678: V4 API referencing non-existing ConfigMap
 *   GKO-679: V4 API referencing existing ConfigMap but non-existing key
 *   GKO-676: V2 API referencing non-existing ConfigMap
 *   GKO-677: V2 API referencing existing ConfigMap but non-existing key
 *   GKO-684: Application description from ConfigMap value
 *   GKO-781: ManagementContext bearer token from Secret
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

const CONFIGMAP_FIXTURE = fixture("crds/templating/configmap-e2e-tpl.yaml");
const SECRET_FIXTURE = fixture("crds/templating/secret-e2e-tpl.yaml");
const BEARER_SECRET_FIXTURE = fixture("crds/templating/secret-e2e-bearer.yaml");

test.describe("Templating — ConfigMap & Secret References", () => {
  // ── GKO-683: V4 API description from ConfigMap ──────────────

  test(`V4 API description resolved from ConfigMap ${XRAY.TEMPLATING.V4_CONFIGMAP_VALUE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-tpl-cm";
    const apiFixture = fixture("crds/templating/v4-api-with-configmap-value.yaml");

    await test.step("Apply ConfigMap and V4 API CRD", async () => {
      await kubectl.apply(CONFIGMAP_FIXTURE);
      await kubectl.apply(apiFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API description keeps the template expression", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api.description).toContain("Values.configmap.e2e-tpl-configmap.description");
    });

    await kubectl.del(apiFixture);
    await kubectl.del(CONFIGMAP_FIXTURE);
  });

  // ── GKO-682: V2 API description from Secret ────────────────

  test(`V2 API description resolved from Secret ${XRAY.TEMPLATING.V2_SECRET_VALUE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-tpl-secret";
    const apiFixture = fixture("crds/templating/v2-api-with-secret-value.yaml");

    await test.step("Ensure clean state", async () => {
      await kubectl.del(apiFixture);
      await kubectl.del(SECRET_FIXTURE);
    });

    await test.step("Apply Secret and V2 API CRD", async () => {
      await kubectl.apply(SECRET_FIXTURE);
      await kubectl.apply(apiFixture);
      await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME);
    const apiId = status.id;

    await test.step("API description keeps the template expression", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api.description).toContain("Values.secret.e2e-tpl-secret.description");
    });

    await kubectl.del(apiFixture);
    await kubectl.del(SECRET_FIXTURE);
  });

  // ── GKO-678: V4 API referencing non-existing ConfigMap ──────
  // Template resolution happens during reconciliation, not at admission.
  // The CRD is accepted by K8s but the operator sets Accepted=False.

  test(`V4 API with missing ConfigMap reference fails reconciliation ${XRAY.TEMPLATING.V4_MISSING_CONFIGMAP} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-tpl-missing-cm";
    const apiFixture = fixture("crds/templating/v4-api-missing-configmap.yaml");

    await test.step("Apply CRD referencing non-existing ConfigMap", async () => {
      await kubectl.apply(apiFixture);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is True", async () => {
      const status = await kubectl.getStatus<{
        conditions?: Array<{ type: string; status: string; message?: string }>;
      }>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(apiFixture);
  });

  // ── GKO-679: V4 API referencing non-existing key ────────────
  // Template resolution happens during reconciliation, not at admission.

  test(`V4 API with missing ConfigMap key fails reconciliation ${XRAY.TEMPLATING.V4_MISSING_KEY} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-tpl-missing-key";
    const apiFixture = fixture("crds/templating/v4-api-missing-key.yaml");

    await test.step("Apply ConfigMap (key exists, but API references wrong key)", async () => {
      await kubectl.apply(CONFIGMAP_FIXTURE);
    });

    await test.step("Apply CRD referencing non-existing key", async () => {
      await kubectl.apply(apiFixture);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is True", async () => {
      const status = await kubectl.getStatus<{
        conditions?: Array<{ type: string; status: string; message?: string }>;
      }>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(apiFixture);
    await kubectl.del(CONFIGMAP_FIXTURE);
  });

  // ── GKO-676: V2 API referencing non-existing ConfigMap ──────
  // Template resolution happens during reconciliation, not at admission.

  test(`V2 API with missing ConfigMap reference fails reconciliation ${XRAY.TEMPLATING.V2_MISSING_CONFIGMAP} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-tpl-missing-cm";
    const apiFixture = fixture("crds/templating/v2-api-missing-configmap.yaml");

    await test.step("Apply CRD referencing non-existing ConfigMap", async () => {
      await kubectl.apply(apiFixture);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is True", async () => {
      const status = await kubectl.getStatus<{
        conditions?: Array<{ type: string; status: string; message?: string }>;
      }>("apidefinition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(apiFixture);
  });

  // ── GKO-677: V2 API referencing non-existing key ────────────
  // Template resolution happens during reconciliation, not at admission.

  test(`V2 API with missing ConfigMap key fails reconciliation ${XRAY.TEMPLATING.V2_MISSING_KEY} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-tpl-missing-key";
    const apiFixture = fixture("crds/templating/v2-api-missing-key.yaml");

    await test.step("Apply ConfigMap (key exists, but API references wrong key)", async () => {
      await kubectl.apply(CONFIGMAP_FIXTURE);
    });

    await test.step("Apply CRD referencing non-existing key", async () => {
      await kubectl.apply(apiFixture);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is True", async () => {
      const status = await kubectl.getStatus<{
        conditions?: Array<{ type: string; status: string; message?: string }>;
      }>("apidefinition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(apiFixture);
    await kubectl.del(CONFIGMAP_FIXTURE);
  });

  // ── GKO-684: Application description from ConfigMap ─────────

  test(`Application description resolved from ConfigMap ${XRAY.TEMPLATING.APP_CONFIGMAP_VALUE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP_NAME = "e2e-app-tpl-cm";
    const appFixture = fixture("crds/templating/app-with-configmap-value.yaml");

    await test.step("Apply ConfigMap and Application CRD", async () => {
      await kubectl.apply(CONFIGMAP_FIXTURE);
      await kubectl.apply(appFixture);
      await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("application", APP_NAME);
    const appId = status.id;

    await test.step("Application description keeps the template expression", async () => {
      const app = await mapi.fetchApplication(appId);
      expect(app.description).toContain("Values.configmap.e2e-tpl-configmap.app-description");
    });

    await kubectl.del(appFixture);
    await kubectl.del(CONFIGMAP_FIXTURE);
  });

  // ── GKO-781: ManagementContext bearer token from Secret ─────

  test(`ManagementContext bearer token resolved from Secret ${XRAY.TEMPLATING.MGMT_CONTEXT_BEARER_TOKEN} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const CTX_NAME = "e2e-ctx-bearer";
    const ctxFixture = fixture("crds/templating/management-context-bearer-token.yaml");

    await test.step("Apply bearer Secret and ManagementContext CRD", async () => {
      await kubectl.apply(BEARER_SECRET_FIXTURE);
      await kubectl.apply(ctxFixture);
    });

    await test.step("ManagementContext resource is created in K8s", async () => {
      const result = await kubectl.get("managementcontext", CTX_NAME);
      expect(result).toBeTruthy();
    });

    await kubectl.del(ctxFixture);
    await kubectl.del(BEARER_SECRET_FIXTURE);
  });
});
