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
 * GKO-only api-key behaviour: things whose *assertion* is Kubernetes-specific
 * (admission acceptance/rejection, Secret templating, push-based idempotent
 * reconcile, and Subscription-only deletion) and so have no Terraform analog.
 * These use the GKO provisioner handle and its `checks` surface, but stay out
 * of the shared matrix.
 */

import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../../helpers/tags.js";
import * as kubectl from "../../../../helpers/kubectl.js";
import { gkoScenario } from "../../../../helpers/provisioner-env.js";
import {
  isGko,
  subscriptionYaml,
  apiKeySecretYaml,
  type Provisioner,
  type Provisioned,
} from "../../../../../src/provisioners/index.js";
import { APIKEY_GKO, gkoApplyApiKeys, uniqueKey, RUN_ID, type ApiKeyParams } from "./params.js";

/** Full GKO provisioner (API + app + an api-key subscription named `subName`). */
function fullGko(subName: string): Provisioner<ApiKeyParams> {
  return gkoScenario<ApiKeyParams>({
    manifests: [APIKEY_GKO.apiManifest, APIKEY_GKO.appManifest],
    roles: {
      api: APIKEY_GKO.apiName,
      application: APIKEY_GKO.appName,
      subscription: subName,
    },
    dynamicRoles: ["subscription"],
    contextPath: APIKEY_GKO.contextPath,
    applyParams: gkoApplyApiKeys(subName),
  })();
}

test.describe(`GKO-only: api-key admission, templating, idempotency ${PROVISIONER.GKO}`, () => {
  let handle: Provisioned<ApiKeyParams> | undefined;
  let provisioner: Provisioner<ApiKeyParams> | undefined;
  let secretName: string | undefined;

  /**
   * Provision through `p`, tracking the provisioner BEFORE provision() resolves
   * so afterEach can tear down even when provision() itself fails partway after
   * applying some CRs (which would otherwise leak the API/app/subscription and
   * poison later tests).
   */
  const provisionTracked = async (
    p: Provisioner<ApiKeyParams>,
    params: ApiKeyParams,
  ): Promise<Provisioned<ApiKeyParams>> => {
    provisioner = p;
    handle = await p.provision(params);
    return handle;
  };

  test.afterEach(async () => {
    if (handle) {
      await handle.destroy().catch(() => {});
    } else if (provisioner?.cleanup) {
      await provisioner.cleanup().catch(() => {});
    }
    handle = undefined;
    provisioner = undefined;
    if (secretName) {
      await kubectl.deleteResource("secret", secretName).catch(() => {});
      secretName = undefined;
    }
  });

  // ── Webhook accepts API_KEY subscriptions (GKO-2547 regression) ──
  test(`Admission webhook accepts api-key plan subscriptions ${XRAY.SUBSCRIPTIONS.V4_APIKEY_WEBHOOK_ACCEPTED} ${TAGS.REGRESSION}`, async () => {
    // Before GKO-2547 the webhook rejected api-key plan subscriptions. The
    // Subscription reaching its Accepted condition during provision() IS the
    // assertion (provision waits for it).
    const h = await provisionTracked(fullGko("e2e-sub-apikey-webhook"), {});
    expect(isGko(h.checks)).toBe(true);
    if (isGko(h.checks)) {
      const status = await h.checks.status<{ id?: string }>("subscription");
      expect(status.id).toBeTruthy();
    }
  });

  // ── CRD schema rejects api-keys outside 32-256 char bounds (admission) ──
  test(`CRD schema rejects api-keys outside 32-256 char bounds ${XRAY.SUBSCRIPTIONS.V4_APIKEY_LENGTH_REJECTED} ${TAGS.REGRESSION}`, async () => {
    // Base API + app only (no subscription); then try bad subs that admission
    // must reject. Terraform's equivalent (server-side rejection) is a
    // structurally different assertion and lives in apikey-tf-only.test.ts.
    await provisionTracked(
      gkoScenario<ApiKeyParams>({
        manifests: [APIKEY_GKO.apiManifest, APIKEY_GKO.appManifest],
        roles: { api: APIKEY_GKO.apiName, application: APIKEY_GKO.appName },
        contextPath: APIKEY_GKO.contextPath,
      })(),
      {},
    );

    const badSub = (name: string, key: string): string =>
      subscriptionYaml({
        name,
        apiName: APIKEY_GKO.apiName,
        plan: APIKEY_GKO.plan,
        applicationName: APIKEY_GKO.appName,
        keys: [{ key }],
      });

    await test.step("31-char key rejected (below minLength=32)", async () => {
      const stderr = await kubectl.applyStringExpectFailure(badSub("e2e-sub-apikey-too-short", "a".repeat(31)));
      expect(stderr).toMatch(/spec\.apiKeys.*key|minLength|too short|invalid/i);
    });
    await test.step("257-char key rejected (above maxLength=256)", async () => {
      const stderr = await kubectl.applyStringExpectFailure(badSub("e2e-sub-apikey-too-long", "a".repeat(257)));
      expect(stderr).toMatch(/spec\.apiKeys.*key|maxLength|too long|invalid/i);
    });
  });

  // ── Custom api-key sourced from a Kubernetes Secret (GKO templating) ──
  test(`Custom api-key sourced from a Kubernetes Secret ${XRAY.SUBSCRIPTIONS.V4_APIKEY_SECRET_SOURCED} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(60_000);
    secretName = `e2e-apikey-secret-${RUN_ID}`;
    const secretValue = uniqueKey("secret-source-apikey");
    const templateRef = `[[ secret \`${secretName}/apiKey\` ]]`;
    const sub = "e2e-sub-apikey-secret";

    await kubectl.applyString(apiKeySecretYaml(secretName, secretValue));
    const h = await provisionTracked(fullGko(sub), { keys: [{ key: templateRef }] });

    const apiId = await h.apiId();
    const subId = await h.subscriptionId();

    // The discriminator: APIM stores the RESOLVED secret value, not the literal
    // `[[ secret ... ]]`. A templating regression would leave the literal (or
    // fail admission entirely).
    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    expect(active.key).toBe(secretValue);

    await gateway.assertResponds(APIKEY_GKO.contextPath, {
      status: 200,
      headers: { "X-Gravitee-Api-Key": secretValue },
    });

    // Release the GKO finalizer on the Secret before afterEach deletes it.
    await kubectl.deleteResource("subscription", sub);
    await kubectl.waitForDeletion("subscription", sub);
  });

  // ── Idempotent reconcile: re-applying the same spec creates no duplicates ──
  test(`Re-applying same custom-key spec does not create extra keys ${XRAY.SUBSCRIPTIONS.V4_APIKEY_IDEMPOTENT} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    test.setTimeout(60_000);
    const KEY = uniqueKey("idempotent-key");
    const h = await provisionTracked(fullGko("e2e-sub-apikey-idempotent"), { keys: [{ key: KEY }] });

    const apiId = await h.apiId();
    const subId = await h.subscriptionId();
    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    expect(active.key).toBe(KEY);

    // Re-apply the identical spec twice; a regression that recreates the key on
    // every reconcile would surface as extra entries sharing this run's value.
    await h.update({ keys: [{ key: KEY }] });
    await h.update({ keys: [{ key: KEY }] });
    await new Promise((r) => setTimeout(r, 5_000));

    const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
    expect(data.filter((k) => !k.revoked && !k.expired).length).toBe(1);
    expect(data.filter((k) => k.key === KEY)).toHaveLength(1);
  });

});
