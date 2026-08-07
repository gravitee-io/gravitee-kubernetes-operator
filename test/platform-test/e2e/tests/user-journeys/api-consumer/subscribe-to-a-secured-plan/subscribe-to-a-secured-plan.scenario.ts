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
 * Journey: subscribe an application to a token-secured plan.
 *
 * As an API consumer, I subscribe my application to an API's JWT plan and to its
 * OAuth2 plan.
 *
 * Both plans declare MANUAL validation on purpose: a subscription written
 * through the Automation API is auto-validated regardless of the plan's
 * validation mode, and that is the behaviour worth pinning — a regression there
 * would leave every declarative subscription stuck PENDING.
 *
 * The two plan security types are a variant of ONE story, so they are two
 * subscription roles in the same scenario rather than two folders; that also
 * keeps them on one API, which is how a real consumer meets them.
 *
 * There is no gateway assertion here. Minting a token the gateway would accept
 * needs a JWT/OAuth2 resource wired to a real issuer, which the test cluster
 * does not provide, and the no-token path cannot be asserted on THIS API either:
 * an OAuth2 plan with no `oauthResource` configured makes the gateway throw an
 * NPE and answer 500 instead of 401 (GKO-3086). The 401 on a
 * JWT-only API is @GKO-817 in tests/gko/subscriptions.
 *
 * Fixtures are co-located in this folder. Subscription rules that live in the
 * operator — immutability of an accepted Subscription CR, admission on a
 * mismatched plan, a syncFrom=KUBERNETES API, an unstarted API, deleting a
 * subscribed plan, and everything V2 — stay in tests/gko/subscriptions.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The plan security types the journey subscribes to, by role label. */
const PLANS = [
  { label: "jwt", title: "JWT" },
  { label: "oauth2", title: "OAuth2" },
] as const;

forEachProvisioner(
  {
    title: "Subscribe an application to a JWT plan and an OAuth2 plan",
    provisioners: {
      gko: gkoScenario<void>({
        // API and application are static and reconcile first; the subscriptions
        // are dynamic so they are applied only afterwards. APIM's subscription
        // endpoint rejects an application it still considers archived, which a
        // single multi-document apply runs straight into.
        manifests: [path.join(here, "gko/api.yaml"), path.join(here, "gko/application.yaml")],
        roles: {
          api: "secured-plans-api",
          application: "secured-plans-app",
          "subscription:jwt": { kind: "subscription", name: "secured-plans-sub-jwt" },
          "subscription:oauth2": { kind: "subscription", name: "secured-plans-sub-oauth2" },
        },
        dynamicRoles: ["subscription:jwt", "subscription:oauth2"],
        contextPath: "/secured-plans-api",
        applyParams: async (k) => {
          await k.apply(path.join(here, "gko/subscription-jwt.yaml"));
          await k.apply(path.join(here, "gko/subscription-oauth2.yaml"));
        },
      }),
      terraform: tfScenario<void>({
        fixture: path.join(here, "terraform"),
        outputs: { "subscription:jwt": "sub_jwt_id", "subscription:oauth2": "sub_oauth2_id" },
      }),
    },
    xray: {
      gko: [
        XRAY.SUBSCRIPTIONS.V4_JWT_SUBSCRIPTION,
        XRAY.SUBSCRIPTIONS.V4_OAUTH2_SUBSCRIPTION,
        XRAY.SUBSCRIPTIONS.AUTO_VALIDATE_V4,
      ],
      terraform: XRAY.TERRAFORM.SUBSCRIBE_SECURED_PLAN_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 120_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    await test.step("The API permits multiple token subscriptions per application", async () => {
      // The flag is what makes the rest of this journey legal: without it APIM
      // rejects the second subscription for the same application. Asserting it
      // here pins the permission the two subscriptions below depend on.
      await mapi.waitForApiMatches(apiId, {
        allowMultiJwtOauth2Subscriptions: true,
        allowedInApiProducts: true,
      });
    });

    for (const plan of PLANS) {
      await test.step(`The ${plan.title} subscription is accepted despite manual validation`, async () => {
        const subscriptionId = await provisioned.subscriptionId(plan.label);
        await mapi.assertSubscriptionAccepted(apiId, subscriptionId);
      });
    }
  },
);
