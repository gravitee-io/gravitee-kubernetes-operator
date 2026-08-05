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
 * Journey: call an mTLS-secured API with a client certificate, rotate it, revoke it.
 *
 * As an API consumer, I register a client certificate on my application, call an
 * mTLS-secured API with it, roll onto a replacement before the old one expires,
 * and revoke it when I am done. The certificate set is an inline attribute on
 * both drivers (`spec.settings.tls.clientCertificates` /
 * `apim_application.settings.tls.client_certificates`), so each stage asserts the
 * outcome where it actually matters: at the GATEWAY, not only in the definition.
 *
 * Both arms present the SAME certificates, read from the shared PKI under
 * fixtures/mtls-certificates/pki, so a difference between the arms is a real
 * difference in how the two drivers registered them.
 *
 * The retirement stage is the one that carries the most regression value: it
 * asserts the retired certificate STOPS being accepted. A rotation that only
 * checks the new certificate works would pass even if the old one never lost
 * access, which is the whole point of rotating.
 *
 * Fixtures are co-located in this folder. The certificate cases only the operator
 * can produce stay in tests/gko/mtls-certificates: admission rejections (bad PEM,
 * expired, end-before-start, content-and-ref together), certificates resolved
 * from cluster Secrets and ConfigMaps, `[[ … ]]` templating, and base64-encoded
 * content, which is a CRD input convenience the provider's PEM-only `content`
 * attribute has no equivalent for.
 */

import path from "node:path";
import { readFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";
import { createTlsFetch } from "../../../../../src/utils/http/tls.js";
import type { Mapi } from "../../../../../src/index.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** Shared with the tests still living in tests/gko/mtls-certificates. */
const PKI_DIR = path.resolve(here, "../../../../fixtures/mtls-certificates/pki");

/**
 * Which certificates the application currently declares. `deprecated` is the
 * legacy singular field; the rest are entries of the `clientCertificates` list.
 */
type CertMode = "client1" | "both" | "client2" | "none" | "deprecated";

/** The list-form modes, which are the only ones the rotation journey drives. */
type ListCertMode = Exclude<CertMode, "deprecated">;

interface CertParams {
  certMode: CertMode;
}

interface ListCertParams {
  certMode: ListCertMode;
}

/** The client identities the journey calls the gateway with. */
type Caller = "anonymous" | "client1" | "client2";

async function loadPki() {
  const [ca, cert1, key1, cert2, key2] = await Promise.all([
    readFile(path.join(PKI_DIR, "ca.crt")),
    readFile(path.join(PKI_DIR, "client1.crt")),
    readFile(path.join(PKI_DIR, "client1.key")),
    readFile(path.join(PKI_DIR, "client2.crt")),
    readFile(path.join(PKI_DIR, "client2.key")),
  ]);
  return { ca, cert1, key1, cert2, key2 };
}

type Pki = Awaited<ReturnType<typeof loadPki>>;

/**
 * A gateway client on the mTLS listener, presenting `caller`'s certificate. The
 * CA is always supplied so the gateway's own certificate verifies; only the
 * CLIENT half changes between callers.
 */
function gatewayAs(mapi: Mapi, baseUrl: string, pki: Pki, caller: Caller) {
  const clientCredentials: Record<Caller, { cert?: Buffer; key?: Buffer }> = {
    anonymous: {},
    client1: { cert: pki.cert1, key: pki.key1 },
    client2: { cert: pki.cert2, key: pki.key2 },
  };
  return mapi.gateway(
    { baseUrl },
    createTlsFetch({ ...clientCredentials[caller], ca: pki.ca }),
  );
}

/**
 * Assert who can reach the API right now. `expected` maps a caller to the status
 * it must get, or to "rejected" when the gateway must refuse it — an unknown
 * client certificate is turned away during the TLS handshake, so there is no HTTP
 * status to assert on.
 */
async function assertGatewayAccess(
  mapi: Mapi,
  baseUrl: string,
  pki: Pki,
  contextPath: string,
  expected: Partial<Record<Caller, number | "rejected">>,
): Promise<void> {
  for (const [caller, outcome] of Object.entries(expected) as [
    Caller,
    number | "rejected",
  ][]) {
    const gw = gatewayAs(mapi, baseUrl, pki, caller);
    if (outcome === "rejected") {
      await gw.assertNotResponds(contextPath, { notStatus: 200 });
    } else {
      await gw.assertResponds(contextPath, { status: outcome });
    }
  }
}

/** The certificate names APIM should report on the application, per stage. */
async function assertRegisteredCertificates(
  mapi: Mapi,
  appId: string,
  expected: string[],
): Promise<void> {
  await expect
    .poll(
      async () => {
        const app = await mapi.fetchApplication(appId);
        return (app.settings?.tls?.client_certificates ?? [])
          .map((c) => c.name ?? "(unnamed)")
          .sort((a, b) => a.localeCompare(b));
      },
      { timeout: 30_000, message: `registered certificates: [${expected.join(", ")}]` },
    )
    .toEqual([...expected].sort((a, b) => a.localeCompare(b)));
}

// ── Journey 1: issue, rotate and revoke through the certificate list ──────────

interface CertStage {
  label: string;
  params: ListCertParams;
  /** Certificate names APIM should report once the stage settles. */
  registered: string[];
  access: Partial<Record<Caller, number | "rejected">>;
}

const STAGES: CertStage[] = [
  {
    label: "issued to client1",
    params: { certMode: "client1" },
    registered: ["client1"],
    access: { anonymous: 401, client1: 200, client2: "rejected" },
  },
  {
    label: "rolling over, client1 and client2 both valid",
    params: { certMode: "both" },
    registered: ["client1", "client2"],
    access: { client1: 200, client2: 200 },
  },
  {
    label: "client1 retired, only client2 left",
    params: { certMode: "client2" },
    registered: ["client2"],
    access: { client1: "rejected", client2: 200 },
  },
  {
    label: "every certificate revoked",
    params: { certMode: "none" },
    registered: [],
    access: { client2: "rejected" },
  },
];

const GKO_APP_MANIFESTS: Record<ListCertMode, string> = {
  client1: "gko/application-client1.yaml",
  both: "gko/application-both.yaml",
  client2: "gko/application-client2.yaml",
  none: "gko/application-none.yaml",
};

forEachProvisioner<ListCertParams>(
  {
    title: "Call an mTLS API with a client certificate, rotate it, then revoke it",
    provisioners: {
      gko: gkoScenario<ListCertParams>({
        // The API reconciles first; the application carries the certificates and
        // is re-applied per stage, so the subscription is applied after both to
        // avoid APIM's application-archived race on a single multi-doc apply.
        manifests: [path.join(here, "gko/api.yaml")],
        roles: {
          api: "mtls-client-cert",
          application: "mtls-client-cert-app",
          subscription: "mtls-client-cert-sub",
        },
        // Teardown walks this list in order, and the admission webhook refuses to
        // delete an application that still has a subscription, so the subscription
        // must come first. Apply order is set by applyParams, not by this list.
        dynamicRoles: ["subscription", "application"],
        contextPath: "/mtls-client-cert",
        applyParams: async (k, params) => {
          await k.apply(path.join(here, GKO_APP_MANIFESTS[params.certMode]));
          // Gate on the application before touching the subscription. APIM's
          // subscription endpoint rejects an application it still considers
          // archived from an earlier run, and applying both back to back loses
          // that race: the Subscription CR then never reaches Accepted.
          await k.waitForCondition("application", "mtls-client-cert-app", "Accepted");
          await k.apply(path.join(here, "gko/subscription.yaml"));
        },
      }),
      terraform: tfScenario<ListCertParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({
          cert_mode: params.certMode,
          pki_dir: PKI_DIR,
          resource_prefix: "mtls-client-cert",
        }),
      }),
    },
    xray: {
      gko: [
        XRAY.MTLS_CERTIFICATES.ADD_MULTIPLE_CERTS,
        XRAY.MTLS_CERTIFICATES.CERT_ROTATION,
        XRAY.MTLS_CERTIFICATES.REMOVE_CERT,
        XRAY.MTLS_CERTIFICATES.ACTIVE_CERT_DISPLAY,
        XRAY.MTLS_CERTIFICATES.NO_ACTIVE_CERTS_DISPLAY,
      ],
      terraform: XRAY.TERRAFORM.CLIENT_CERTIFICATE_ROTATION_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 150_000, terraform: 240_000 },
  },
  async ({ provisioned, mapi, mtlsGatewayBaseUrl }) => {
    const appId = await provisioned.applicationId();
    const apiId = await provisioned.apiId();
    const contextPath = await provisioned.contextPath();
    const pki = await loadPki();

    await test.step("The API's only plan is mTLS-secured", async () => {
      await expect
        .poll(
          async () => (await mapi.listApiPlans(apiId)).map((p) => p.security?.type),
          { timeout: 30_000, message: "the API exposes an MTLS plan" },
        )
        .toEqual(["MTLS"]);
    });

    for (const [index, stage] of STAGES.entries()) {
      await test.step(`Certificates ${stage.label}`, async () => {
        // Stage 0 is what provision() already applied.
        if (index > 0) await provisioned.update(stage.params);

        await assertRegisteredCertificates(mapi, appId, stage.registered);
        await assertGatewayAccess(mapi, mtlsGatewayBaseUrl, pki, contextPath, stage.access);
      });
    }
  },
  STAGES[0].params,
);

// ── Journey 2: the deprecated single-certificate field ────────────────────────

/**
 * A separate scenario rather than a stage of the one above: it starts from the
 * legacy field, and the migration deliberately swaps in a DIFFERENT certificate.
 * Re-registering the same certificate content under the new field would trip
 * APIM's fingerprint-reuse rule, which is its own GKO-only test.
 */
forEachProvisioner<CertParams>(
  {
    title: "Migrate from the deprecated single client certificate to the list",
    provisioners: {
      gko: gkoScenario<CertParams>({
        manifests: [path.join(here, "gko/api-legacy.yaml")],
        roles: {
          api: "mtls-legacy-cert",
          application: "mtls-legacy-cert-app",
          subscription: "mtls-legacy-cert-sub",
        },
        // Teardown walks this list in order, and the admission webhook refuses to
        // delete an application that still has a subscription, so the subscription
        // must come first. Apply order is set by applyParams, not by this list.
        dynamicRoles: ["subscription", "application"],
        contextPath: "/mtls-legacy-cert",
        applyParams: async (k, params) => {
          await k.apply(
            path.join(
              here,
              params.certMode === "deprecated"
                ? "gko/application-legacy-field.yaml"
                : "gko/application-legacy-migrated.yaml",
            ),
          );
          // See the sibling scenario: the subscription has to wait for the
          // application, or it loses APIM's archived-application race.
          await k.waitForCondition("application", "mtls-legacy-cert-app", "Accepted");
          await k.apply(path.join(here, "gko/subscription-legacy.yaml"));
        },
      }),
      terraform: tfScenario<CertParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({
          cert_mode: params.certMode,
          pki_dir: PKI_DIR,
          resource_prefix: "mtls-legacy-cert",
        }),
      }),
    },
    xray: {
      gko: XRAY.MTLS_CERTIFICATES.DEPRECATED_FIELD,
      terraform: XRAY.TERRAFORM.CLIENT_CERTIFICATE_DEPRECATED_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 120_000, terraform: 240_000 },
  },
  async ({ provisioned, mapi, mtlsGatewayBaseUrl }) => {
    const contextPath = await provisioned.contextPath();
    const pki = await loadPki();

    await test.step("The deprecated field grants access to client1", async () => {
      await assertGatewayAccess(mapi, mtlsGatewayBaseUrl, pki, contextPath, {
        anonymous: 401,
        client1: 200,
      });
    });

    await test.step("Moving to the list with client2 revokes client1", async () => {
      await provisioned.update({ certMode: "client2" });
      await assertGatewayAccess(mapi, mtlsGatewayBaseUrl, pki, contextPath, {
        client1: "rejected",
        client2: 200,
      });
    });
  },
  { certMode: "deprecated" },
);
