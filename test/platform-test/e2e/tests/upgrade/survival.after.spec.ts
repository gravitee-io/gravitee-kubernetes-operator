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
 * Upgrade survival - AFTER phase.
 *
 * Runs after the cluster has been upgraded to the NEW line (branch GKO + APIM
 * 4.12), in a SEPARATE process from the before-phase. It re-attaches to the
 * resources the before-phase created (they persist as CRs across the in-place
 * upgrade), then:
 *   1. asserts they survived the operator + APIM upgrade,
 *   2. updates the carried-over API - the highest-risk step, where the NEW
 *      operator rewrites a resource the OLD operator created (and where CRD
 *      field evolution would surface),
 *   3. tears everything down in reverse dependency order and confirms removal.
 *
 * Run only via the `e2e:upgrade:after` script; the normal suite never collects it.
 */

import { readFile } from "node:fs/promises";
import { test, expect, fixture } from "../../setup.js";
import * as kubectl from "../../helpers/kubectl.js";
import { signJwt } from "../../helpers/jwt.js";
import { createTlsFetch } from "../../../src/utils/http/tls.js";
import { compareVersions } from "../../../src/utils/version/index.js";
import { XRAY } from "../../helpers/tags.js";
import {
  survivalScenario,
  SURVIVAL,
  survivalV2Scenario,
  SURVIVAL_V2,
  survivalV2SubScenario,
  SURVIVAL_V2_SUB,
} from "./survival-scenario.js";

/**
 * True when the upgrade TARGET (the NEW line, read from E2E_MAX_VERSION) is at
 * least `version`. Unset target -> assume latest, so the step runs. This lets
 * version-specific steps skip on older targets - e.g. a future 4.9 -> 4.11 run
 * skips the 4.12-only mTLS + V2-subscription steps. (Full pairing is GKO-2985.)
 */
function targetAtLeast(version: string): boolean {
  const target = process.env["E2E_MAX_VERSION"];
  if (!target) return true;
  return compareVersions(target, version) >= 0;
}

// Safety-net cleanup: if any step fails mid-test the inline teardown never runs,
// so delete every upgrade resource best-effort in dependency order (subscriptions
// -> applications -> APIs/secret) to leave the cluster clean for the next run.
test.afterAll(async () => {
  const order = [
    "v2-sub", "sub-legacy-weather", "sub-mobile-legacy", "sub-mobile-weather", "sub-legacy-mtls", "sub-legacy-legacy-jwt",
    "v2-sub-app", "app-mobile", "app-legacy",
    "v2-sub-api", "api-weather", "api-legacy", "v2-legacy-api", "jwt-secret",
  ];
  for (const f of order) {
    await kubectl.del(fixture(`upgrade/${f}.yaml`)).catch(() => {});
  }
});

test(`upgrade survival: reconnect, verify, update, delete on the new line ${XRAY.API_LIFECYCLE.V4_SURVIVES_UPGRADE} @upgrade @after`, async ({
  mapi,
  gateway,
  mtlsGatewayBaseUrl,
}) => {
  const provisioner = survivalScenario();
  if (!provisioner.attach) {
    throw new Error("the GKO provisioner does not implement attach(); cannot run the after-phase");
  }

  // Rebuild a handle to the resources created before the upgrade (a missing CR
  // here means it did not survive).
  const provisioned = await provisioner.attach({
    api: { hrid: SURVIVAL.apiName },
    application: { hrid: SURVIVAL.appName },
    subscription: { hrid: SURVIVAL.subName },
  });

  const apiId = await provisioned.apiId();
  const subId = await provisioned.subscriptionId();
  const ctx = await provisioned.contextPath();
  const bearer = () => ({ Authorization: `Bearer ${signJwt("legacy-client")}` });

  await test.step("resources survived the GKO + APIM upgrade", async () => {
    await mapi.waitForApiStarted(apiId, { timeoutMs: 30_000 });
    await mapi.assertSubscriptionAccepted(apiId, subId);
    // Still reachable through the gateway after the operator + APIM upgrade.
    await gateway.assertResponds(ctx, { status: 401 });
    await gateway.assertResponds(ctx, { status: 200, headers: bearer() });
  });

  await test.step("the branch operator updates the carried-over API", async () => {
    await kubectl.apply(fixture("upgrade/api-legacy-updated.yaml"));
    await mapi.waitForApiMatches(
      apiId,
      { description: SURVIVAL.updatedDescription },
      { timeoutMs: 30_000 },
    );
    // The update did not break reachability.
    await gateway.assertResponds(ctx, { status: 200, headers: bearer() });
  });

  // mTLS clientCertificate is a 4.12-only surface, so this step only runs when the
  // upgrade TARGET is >= 4.12. On an older target it is skipped (see GKO-2985).
  if (targetAtLeast("4.12")) {
    await test.step("the new operator adds an mTLS plan + cert to the carried-over API (4.12+)", async () => {
      // The upgraded operator extends the 4.11-created API with an mTLS plan, gives
      // the carried-over app a client cert, and subscribes it.
      await kubectl.apply(fixture("upgrade/api-legacy-with-mtls.yaml"));
      await kubectl.waitForCondition("apiv4definition", SURVIVAL.apiName, "Accepted");
      await kubectl.apply(fixture("upgrade/app-legacy-with-cert.yaml"));
      await kubectl.waitForCondition("application", SURVIVAL.appName, "Accepted");
      await kubectl.apply(fixture("upgrade/sub-legacy-mtls.yaml"));
      await kubectl.waitForCondition("subscription", "legacy-mtls", "Accepted");

      const [cert, key, ca] = await Promise.all([
        readFile(fixture("mtls-certificates/pki/client1.crt")),
        readFile(fixture("mtls-certificates/pki/client1.key")),
        readFile(fixture("mtls-certificates/pki/ca.crt")),
      ]);
      // mTLS data plane: rejected without a client cert, served with client1's cert.
      await mapi
        .gateway({ baseUrl: mtlsGatewayBaseUrl }, createTlsFetch({ ca }))
        .assertResponds(ctx, { status: 401 });
      await mapi
        .gateway({ baseUrl: mtlsGatewayBaseUrl }, createTlsFetch({ cert, key, ca }))
        .assertResponds(ctx, { status: 200 });
      // JWT on the HTTP port still works alongside the new mTLS plan.
      await gateway.assertResponds(ctx, { status: 200, headers: bearer() });
    });
  }

  await test.step("cross-subscriptions: new APIs/apps cross-subscribe on the new line", async () => {
    // A second API + app, plus the three cross-subscriptions (legacy->weather,
    // mobile->legacy, mobile->weather), alongside the carried-over resources.
    await kubectl.apply(fixture("upgrade/api-weather.yaml"));
    await kubectl.waitForCondition("apiv4definition", "weather-api", "Accepted");
    await kubectl.apply(fixture("upgrade/app-mobile.yaml"));
    await kubectl.waitForCondition("application", "mobile-app", "Accepted");
    for (const s of ["sub-legacy-weather", "sub-mobile-legacy", "sub-mobile-weather"]) {
      await kubectl.apply(fixture(`upgrade/${s}.yaml`));
    }
    await kubectl.waitForCondition("subscription", "legacy-weather-jwt", "Accepted");
    await kubectl.waitForCondition("subscription", "mobile-legacy-jwt", "Accepted");
    await kubectl.waitForCondition("subscription", "mobile-weather-jwt", "Accepted");

    const legacyTok = { Authorization: `Bearer ${signJwt("legacy-client")}` };
    const mobileTok = { Authorization: `Bearer ${signJwt("mobile-app")}` };
    await gateway.assertResponds("/weather", { status: 200, headers: legacyTok }); // legacy-app -> weather-api
    await gateway.assertResponds(ctx, { status: 200, headers: mobileTok }); // mobile-app -> legacy-api
    await gateway.assertResponds("/weather", { status: 200, headers: mobileTok }); // mobile-app -> weather-api
  });

  await test.step("clean teardown in reverse dependency order", async () => {
    // Subscriptions first (admission requires subs before apps), then the extra
    // apps/APIs, then the carried-over core via destroy().
    for (const s of ["sub-legacy-weather", "sub-mobile-legacy", "sub-mobile-weather", "sub-legacy-mtls"]) {
      await kubectl.del(fixture(`upgrade/${s}.yaml`)).catch(() => {});
    }
    await kubectl.del(fixture("upgrade/app-mobile.yaml")).catch(() => {});
    await kubectl.del(fixture("upgrade/api-weather.yaml")).catch(() => {});
    await provisioned.destroy();
    // destroy() issues deletes best-effort (errors swallowed); wait for the CR to
    // actually be gone (finalizer + APIM deletion) before asserting the mAPI 404.
    await kubectl.waitForDeletion("apiv4definition", SURVIVAL.apiName);
    await mapi.assertApiHttpStatus(apiId, 404);
    // The gateway no longer serves the deleted API.
    await gateway.assertNotResponds(ctx, { notStatus: 200, headers: bearer() });
  });
});

test(`upgrade survival (V2): keyless V2 API survives the upgrade ${XRAY.API_LIFECYCLE.V2_SURVIVES_UPGRADE} @upgrade @after`, async ({
  mapi,
  gateway,
}) => {
  const provisioner = survivalV2Scenario();
  if (!provisioner.attach) {
    throw new Error("the GKO provisioner does not implement attach(); cannot run the after-phase");
  }
  const provisioned = await provisioner.attach({ api: { hrid: SURVIVAL_V2.apiName } });
  const apiId = await provisioned.apiId();
  const ctx = await provisioned.contextPath();

  await test.step("V2 API survived the upgrade (control plane + gateway)", async () => {
    await mapi.assertApiHttpStatus(apiId, 200); // still present in the management API
    await gateway.assertResponds(ctx, { status: 200 });
  });

  await test.step("the new operator updates the carried-over V2 API", async () => {
    // Re-apply with a changed description + a second context path; the new path
    // being served proves the update reconciled (version-agnostic).
    await kubectl.apply(fixture("upgrade/v2-legacy-api-updated.yaml"));
    await kubectl.waitForCondition("apidefinition", SURVIVAL_V2.apiName, "Accepted");
    await gateway.assertResponds("/legacy-v2-updated", { status: 200 });
    await gateway.assertResponds(ctx, { status: 200 });
  });

  // V2 subscriptions go through the Automation API (4.12-only) -> gated (GKO-2985).
  if (targetAtLeast("4.12")) {
    await test.step("the new operator creates a V2 subscription (4.12+)", async () => {
      const sub = await survivalV2SubScenario().provision({});
      // provision() waited for the Subscription CR to reach Accepted; a resolvable
      // APIM id confirms the V2 subscription was created server-side.
      expect(await sub.subscriptionId()).toBeTruthy();
      await sub.destroy();
      await kubectl.waitForDeletion("apidefinition", SURVIVAL_V2_SUB.apiName);
    });
  }

  await test.step("clean teardown", async () => {
    await provisioned.destroy();
    await kubectl.waitForDeletion("apidefinition", SURVIVAL_V2.apiName);
    await gateway.assertNotResponds(ctx, { notStatus: 200 });
  });
});
