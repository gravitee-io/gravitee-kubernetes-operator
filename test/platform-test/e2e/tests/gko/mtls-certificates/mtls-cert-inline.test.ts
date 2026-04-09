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
 * mTLS certificate tests: inline, encoded, and backward-compat variants.
 *
 * Xray tests:
 *   GKO-2243: Add Multiple Certificates to Application
 *   GKO-2212: Add Single Certificate to Application
 *   GKO-2244: Deprecated Field Functionality (backward-compat)
 *
 * Preconditions:
 *   - APIM, Gateway (HTTP + mTLS), and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { readFile } from "node:fs/promises";
import { test, fixture } from "../../../setup.js";
import { XRAY } from "../../../helpers/tags.js";
import { createTlsFetch } from "../../../../src/utils/http/tls.js";
import * as kubectl from "../../../helpers/kubectl.js";

const PKI = (...segments: string[]) => fixture("crds/mtls-certificates/pki", ...segments);

async function loadPki() {
  const [cert1, key1, cert2, key2, ca] = await Promise.all([
    readFile(PKI("client1.crt")),
    readFile(PKI("client1.key")),
    readFile(PKI("client2.crt")),
    readFile(PKI("client2.key")),
    readFile(PKI("ca.crt")),
  ]);
  return { cert1, key1, cert2, key2, ca };
}

test.describe("mTLS Certificates — Inline, Encoded, Backward-compat", () => {
  test.afterAll(async () => {
    const files = [
      "subscription-inline", "application-inline", "api-mtls-inline", "tls-secrets-inline",
      "subscription-encoded", "application-encoded", "api-mtls-encoded", "tls-secrets-encoded",
      "subscription-backward-compat", "application-compat-cert2-list", "api-mtls-backward-compat", "tls-secrets-backward-compat",
    ];
    for (const f of files) {
      await kubectl.del(fixture(`crds/mtls-certificates/${f}.yaml`)).catch(() => {});
    }
  });

  test(`Inline certs: both grant access ${XRAY.MTLS_CERTIFICATES.ADD_MULTIPLE_CERTS}`, async ({
    kubectl,
    mapi,
    mtlsGatewayBaseUrl,
  }) => {
    const API_NAME = "e2e-mtls-inline";
    const API_PATH = "/e2e-mtls-inline";

    await test.step("Deploy API, Application, and Subscription", async () => {
      await kubectl.apply(fixture("crds/mtls-certificates/api-mtls-inline.yaml"));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
      await kubectl.apply(fixture("crds/mtls-certificates/application-inline.yaml"));
      await kubectl.waitForCondition("application", "e2e-mtls-inline-app", "Accepted");
      await kubectl.apply(fixture("crds/mtls-certificates/subscription-inline.yaml"));
      await kubectl.waitForCondition("subscription", "e2e-mtls-inline-sub", "Accepted");
    });

    const pki = await loadPki();

    await test.step("Gateway rejects request without cert (401)", async () => {
      const gw = mapi.gateway({ baseUrl: mtlsGatewayBaseUrl }, createTlsFetch({ ca: pki.ca }));
      await gw.assertResponds(API_PATH, { status: 401 });
    });

    await test.step("Gateway accepts client1 cert (200)", async () => {
      const gw = mapi.gateway(
        { baseUrl: mtlsGatewayBaseUrl },
        createTlsFetch({ cert: pki.cert1, key: pki.key1, ca: pki.ca }),
      );
      await gw.assertResponds(API_PATH, { status: 200 });
    });

    await test.step("Gateway accepts client2 cert (200)", async () => {
      const gw = mapi.gateway(
        { baseUrl: mtlsGatewayBaseUrl },
        createTlsFetch({ cert: pki.cert2, key: pki.key2, ca: pki.ca }),
      );
      await gw.assertResponds(API_PATH, { status: 200 });
    });

    // Cleanup
    await kubectl.del(fixture("crds/mtls-certificates/subscription-inline.yaml"));
    await kubectl.del(fixture("crds/mtls-certificates/application-inline.yaml"));
    await kubectl.del(fixture("crds/mtls-certificates/api-mtls-inline.yaml"));
  });

  test(`Base64-encoded cert ${XRAY.MTLS_CERTIFICATES.ADD_SINGLE_CERT}`, async ({
    kubectl,
    mapi,
    mtlsGatewayBaseUrl,
  }) => {
    const API_NAME = "e2e-mtls-encoded";
    const API_PATH = "/e2e-mtls-encoded";

    await test.step("Deploy API, Application, and Subscription", async () => {
      await kubectl.apply(fixture("crds/mtls-certificates/api-mtls-encoded.yaml"));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
      await kubectl.apply(fixture("crds/mtls-certificates/application-encoded.yaml"));
      await kubectl.waitForCondition("application", "e2e-mtls-encoded-app", "Accepted");
      await kubectl.apply(fixture("crds/mtls-certificates/subscription-encoded.yaml"));
      await kubectl.waitForCondition("subscription", "e2e-mtls-encoded-sub", "Accepted");
    });

    const pki = await loadPki();

    await test.step("Gateway accepts client1 cert (200)", async () => {
      const gw = mapi.gateway(
        { baseUrl: mtlsGatewayBaseUrl },
        createTlsFetch({ cert: pki.cert1, key: pki.key1, ca: pki.ca }),
      );
      await gw.assertResponds(API_PATH, { status: 200 });
    });

    // Cleanup
    await kubectl.del(fixture("crds/mtls-certificates/subscription-encoded.yaml"));
    await kubectl.del(fixture("crds/mtls-certificates/application-encoded.yaml"));
    await kubectl.del(fixture("crds/mtls-certificates/api-mtls-encoded.yaml"));
  });

  test(`Backward compat: deprecated clientCertificate field ${XRAY.MTLS_CERTIFICATES.DEPRECATED_FIELD}`, async ({
    kubectl,
    mapi,
    mtlsGatewayBaseUrl,
  }) => {
    const API_NAME = "e2e-mtls-backward-compat";
    const API_PATH = "/e2e-mtls-backward-compat";

    await test.step("Deploy API, Application (deprecated field), and Subscription", async () => {
      await kubectl.apply(fixture("crds/mtls-certificates/api-mtls-backward-compat.yaml"));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
      await kubectl.apply(fixture("crds/mtls-certificates/application-backward-compat.yaml"));
      await kubectl.waitForCondition("application", "e2e-mtls-compat-app", "Accepted");
      await kubectl.apply(fixture("crds/mtls-certificates/subscription-backward-compat.yaml"));
      await kubectl.waitForCondition("subscription", "e2e-mtls-compat-sub", "Accepted");
    });

    const pki = await loadPki();

    await test.step("Gateway accepts client1 cert with deprecated field (200)", async () => {
      const gw = mapi.gateway(
        { baseUrl: mtlsGatewayBaseUrl },
        createTlsFetch({ cert: pki.cert1, key: pki.key1, ca: pki.ca }),
      );
      await gw.assertResponds(API_PATH, { status: 200 });
    });

    await test.step("Upgrade to clientCertificates list with client2 only", async () => {
      await kubectl.apply(fixture("crds/mtls-certificates/application-compat-cert2-list.yaml"));
      await kubectl.waitForCondition("application", "e2e-mtls-compat-app", "Accepted");
    });

    await test.step("Gateway rejects client1 after upgrade", async () => {
      const gw = mapi.gateway(
        { baseUrl: mtlsGatewayBaseUrl },
        createTlsFetch({ cert: pki.cert1, key: pki.key1, ca: pki.ca }),
      );
      await gw.assertNotResponds(API_PATH, { notStatus: 200 });
    });

    await test.step("Gateway accepts client2 after upgrade (200)", async () => {
      const gw = mapi.gateway(
        { baseUrl: mtlsGatewayBaseUrl },
        createTlsFetch({ cert: pki.cert2, key: pki.key2, ca: pki.ca }),
      );
      await gw.assertResponds(API_PATH, { status: 200 });
    });

    // Cleanup
    await kubectl.del(fixture("crds/mtls-certificates/subscription-backward-compat.yaml"));
    await kubectl.del(fixture("crds/mtls-certificates/application-compat-cert2-list.yaml"));
    await kubectl.del(fixture("crds/mtls-certificates/api-mtls-backward-compat.yaml"));
  });
});
