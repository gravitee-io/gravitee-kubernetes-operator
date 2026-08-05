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
 * mTLS certificate content encodings the CRD accepts.
 *
 * Xray tests:
 *   GKO-2212: Add Single Certificate to Application
 *
 * Base64-encoded certificate content is a CRD input convenience: the operator
 * decodes it before writing to APIM. The Terraform provider's `content` attribute
 * is documented as PEM only, so there is nothing to compare a second driver
 * against and this stays GKO-only. Presenting a certificate at the gateway,
 * rotating it and revoking it are covered for BOTH drivers by the
 * `authenticate-with-client-certificate` journey.
 *
 * Preconditions:
 *   - APIM, Gateway (HTTP + mTLS), and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { readFile } from "node:fs/promises";
import { test, fixture } from "../../../setup.js";
import { XRAY, PROVISIONER } from "../../../helpers/tags.js";
import { createTlsFetch } from "../../../../src/utils/http/tls.js";
import * as kubectl from "../../../helpers/kubectl.js";

const PKI = (...segments: string[]) => fixture("mtls-certificates/pki", ...segments);

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

test.describe(`mTLS Certificates — encoded content ${PROVISIONER.GKO}`, () => {
  test.afterAll(async () => {
    const files = [
      "subscription-encoded", "application-encoded", "api-mtls-encoded", "tls-secrets-encoded",
    ];
    for (const f of files) {
      await kubectl.del(fixture(`mtls-certificates/${f}/crd.yaml`)).catch(() => {});
    }
  });

  test(`Base64-encoded cert ${XRAY.MTLS_CERTIFICATES.ADD_SINGLE_CERT}`, async ({
    kubectl,
    mapi,
    mtlsGatewayBaseUrl,
  }) => {
    const API_NAME = "e2e-mtls-encoded";
    const API_PATH = "/e2e-mtls-encoded";

    await test.step("Deploy API, Application, and Subscription", async () => {
      await kubectl.apply(fixture("mtls-certificates/api-mtls-encoded/crd.yaml"));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
      await kubectl.apply(fixture("mtls-certificates/application-encoded/crd.yaml"));
      await kubectl.waitForCondition("application", "e2e-mtls-encoded-app", "Accepted");
      await kubectl.apply(fixture("mtls-certificates/subscription-encoded/crd.yaml"));
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
    await kubectl.del(fixture("mtls-certificates/subscription-encoded/crd.yaml"));
    await kubectl.del(fixture("mtls-certificates/application-encoded/crd.yaml"));
    await kubectl.del(fixture("mtls-certificates/api-mtls-encoded/crd.yaml"));
  });

});
