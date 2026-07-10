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
 * Upgrade survival - BEFORE phase.
 *
 * Runs while the cluster is still on the OLD line (released GKO 4.11 + APIM 4.11).
 * Provisions the survival scenario, asserts it is healthy, and DELIBERATELY does
 * not tear it down: the resources must survive into the after-phase, which runs
 * in a separate process once both GKO and APIM have been upgraded.
 *
 * This is a `*.spec.ts` file run only via the `e2e:upgrade:before` script and the
 * dedicated upgrade config; the normal suite (testMatch *.test.ts / *.scenario.ts)
 * never collects it.
 */

import { test, expect } from "../../setup.js";
import * as kubectl from "../../helpers/kubectl.js";
import { signJwt } from "../../helpers/jwt.js";
import {
  DATASTORE,
  SURVIVAL_NON_HRID,
  survivalNonHridScenario,
  survivalScenario,
  survivalV2Scenario,
} from "./survival-scenario.js";

test("upgrade survival: provision and verify on the old line @upgrade @before", async ({
  mapi,
  gateway,
}) => {
  const provisioned = await survivalScenario().provision({});

  const apiId = await provisioned.apiId();
  const subId = await provisioned.subscriptionId();
  const ctx = await provisioned.contextPath();

  // Healthy on the OLD line: the API is started and the subscription is accepted.
  await mapi.waitForApiStarted(apiId, { timeoutMs: 30_000 });
  await mapi.assertSubscriptionAccepted(apiId, subId);

  // Reachable through the gateway on the OLD line: rejected without a token,
  // served with a valid JWT (client_id = the subscribed app's clientId).
  await gateway.assertResponds(ctx, { status: 401 });
  await gateway.assertResponds(ctx, {
    status: 200,
    headers: { Authorization: `Bearer ${signJwt("legacy-client")}` },
  });

  // No teardown on purpose - the after-phase re-attaches to these resources.
});

test("upgrade survival (V2): keyless V2 API reachable on the old line @upgrade @before", async ({
  gateway,
}) => {
  const provisioned = await survivalV2Scenario().provision({});
  const ctx = await provisioned.contextPath();

  // provision() already waited for the V2 API CR to reconcile (Accepted); the
  // keyless plan means the gateway serves it without a token. Leave it running.
  await gateway.assertResponds(ctx, { status: 200 });
});

test("upgrade survival (non-HRID names): provision and verify on the old line @upgrade @before", async ({
  mapi,
  gateway,
}) => {
  const provisioned = await survivalNonHridScenario().provision({});

  const apiId = await provisioned.apiId();
  const subId = await provisioned.subscriptionId();
  const ctx = await provisioned.contextPath();

  // Healthy on the OLD line despite the spaced plan/page keys and the lowercase
  // flow mode: the API is started and the subscription (which references the
  // plan by its raw spaced key) is accepted.
  await mapi.waitForApiStarted(apiId, { timeoutMs: 30_000 });
  await mapi.assertSubscriptionAccepted(apiId, subId);

  // Reachable through the gateway: rejected without a token, served with a
  // valid JWT (client_id = the subscribed app's clientId).
  await gateway.assertResponds(ctx, { status: 401 });
  await gateway.assertResponds(ctx, {
    status: 200,
    headers: { Authorization: `Bearer ${signJwt(SURVIVAL_NON_HRID.clientId)}` },
  });

  // The pages declared under spaced map keys made it into APIM, once each.
  await expect
    .poll(
      async () => {
        const pages = await mapi.listApiPages(apiId);
        return {
          folders: pages.filter((p) => p.name === SURVIVAL_NON_HRID.folderName).length,
          markdowns: pages.filter((p) => p.name === SURVIVAL_NON_HRID.pageName).length,
        };
      },
      { timeout: 30_000 },
    )
    .toEqual({ folders: 1, markdowns: 1 });

  // No teardown on purpose - the after-phase re-attaches to these resources.
});

test("upgrade survival: record the datastore identity before the upgrade @upgrade @before", async () => {
  // Snapshot the MongoDB pod's identity so the after-phase can prove the same
  // pod survived the in-place upgrade (the APIM data lives there). Recorded in a
  // ConfigMap because the two phases run in separate processes and only share
  // state through the cluster. If no MongoDB pod is found, record nothing - the
  // after-phase then skips rather than failing on a cluster shaped differently.
  const mongo = await kubectl.findPod(DATASTORE.podNamePattern);
  if (!mongo) {
    console.warn(
      `[upgrade] no pod matching "${DATASTORE.podNamePattern}" found; ` +
        "skipping datastore-continuity recording",
    );
    return;
  }
  await kubectl.writeConfigMap(DATASTORE.stateConfigMap, {
    uid: mongo.uid,
    restartCount: String(mongo.restartCount),
    name: mongo.name,
    namespace: mongo.namespace,
  });
});
