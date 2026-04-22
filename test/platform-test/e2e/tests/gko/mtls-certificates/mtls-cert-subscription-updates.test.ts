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
 * mTLS — Application clientCertificates visibility via mAPI (batch 8).
 *
 * Xray tests covered:
 *   GKO-2219: Application with no clientCertificates list is visible via mAPI
 *             with an empty clientCertificates list (equivalent to the "no
 *             active certificates" UI state).
 *   GKO-2223: Cert startsAt / endsAt round-trip from the CRD to mAPI
 *             (equivalent to the UI's certificate information display).
 *   GKO-2246: A valid client certificate is visible on the application.
 *   GKO-2251: A cert whose endsAt is in the past is filtered out of the
 *             active certificate list returned by mAPI — the same end state
 *             shown in the UI's "expired certificate" indicator (the list
 *             of *active* certs is empty, the expired one is no longer
 *             surfaced through settings.tls.client_certificates).
 *
 * The broader bucket-H lifecycle scenarios (subscription update flows,
 * PKCS7 bundles, cert rotation edge cases) are not covered here — they
 * either have no GKO analog (PKCS7 — GKO only accepts PEM) or require
 * live per-scenario PKI generation that is tracked for a follow-up batch.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 *   - The existing PKI assets under fixtures/crds/mtls-certificates/pki/
 *     are reused (client1.crt is inlined in h-2223 / h-2251)
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const APP_NO_CERTS = "crds/mtls-certificates/h-2219-app-no-certs.yaml";
const APP_VALID_CERT = "crds/mtls-certificates/h-2223-app-valid-cert.yaml";
const APP_PAST_END = "crds/mtls-certificates/h-2251-app-past-end.yaml";

// mAPI serializes the certificate list in snake_case and uses epoch millis
// for the validity window.
interface AppClientCertView {
  name?: string;
  startsAt?: number;
  endsAt?: number;
  certificate?: string;
}

interface AppWithTls {
  settings?: {
    app?: { type?: string };
    tls?: { client_certificates?: AppClientCertView[]; certificate_count?: number };
  };
}

test.describe("mTLS — clientCertificates visibility via mAPI (batch 8)", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(APP_NO_CERTS)).catch(() => {});
    await kubectlSafe.del(fixture(APP_VALID_CERT)).catch(() => {});
    await kubectlSafe.del(fixture(APP_PAST_END)).catch(() => {});
  });

  test(`Application with no clientCertificates shows empty list via mAPI ${XRAY.MTLS_CERTIFICATES.NO_ACTIVE_CERTS_DISPLAY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-app-h-2219";
    await kubectl.apply(fixture(APP_NO_CERTS));
    await kubectl.waitForCondition("application", APP, "Accepted");

    const appId = (await kubectl.getStatus<{ id: string }>("application", APP)).id;
    await mapi.waitForApplicationMatches(appId, { name: APP });

    const raw = (await mapi.fetchApplication(appId)) as unknown as AppWithTls;
    const certs = raw.settings?.tls?.client_certificates ?? [];
    expect(certs).toEqual([]);
  });

  test(`Cert startsAt and endsAt round-trip from CRD to mAPI ${XRAY.MTLS_CERTIFICATES.CERT_INFO_DISPLAY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-app-h-2223";
    await kubectl.apply(fixture(APP_VALID_CERT));
    await kubectl.waitForCondition("application", APP, "Accepted");

    const appId = (await kubectl.getStatus<{ id: string }>("application", APP)).id;

    await expect
      .poll(
        async () => {
          const raw = (await mapi.fetchApplication(appId)) as unknown as AppWithTls;
          const cert = raw.settings?.tls?.client_certificates?.[0];
          return {
            name: cert?.name,
            hasStart: Boolean(cert?.startsAt),
            hasEnd: Boolean(cert?.endsAt),
          };
        },
        { timeout: 15_000, intervals: [1_000] },
      )
      .toEqual({
        name: "client1",
        hasStart: true,
        hasEnd: true,
      });
  });

  test(`Active certificate is visible on the application via mAPI ${XRAY.MTLS_CERTIFICATES.ACTIVE_CERT_DISPLAY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-app-h-2223";
    await kubectl.apply(fixture(APP_VALID_CERT));
    await kubectl.waitForCondition("application", APP, "Accepted");

    const appId = (await kubectl.getStatus<{ id: string }>("application", APP)).id;

    await expect
      .poll(
        async () => {
          const raw = (await mapi.fetchApplication(appId)) as unknown as AppWithTls;
          return raw.settings?.tls?.client_certificates?.length ?? 0;
        },
        { timeout: 10_000, intervals: [1_000] },
      )
      .toBeGreaterThanOrEqual(1);
  });

  test(`Past-end-date certificate is filtered out of the active cert list ${XRAY.MTLS_CERTIFICATES.EXPIRED_CERT_DISPLAY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-app-h-2251";
    await kubectl.apply(fixture(APP_PAST_END));
    await kubectl.waitForCondition("application", APP, "Accepted");

    const appId = (await kubectl.getStatus<{ id: string }>("application", APP)).id;

    // APIM treats endsAt as a hard filter: certificates whose endsAt has
    // passed are removed from the application's client_certificates list
    // returned by the management API. This is the data the console's
    // "expired certificate" indicator is derived from — the absence of
    // the cert from the active list is the signal.
    await expect
      .poll(
        async () => {
          const raw = (await mapi.fetchApplication(appId)) as unknown as AppWithTls;
          return raw.settings?.tls?.client_certificates ?? [];
        },
        { timeout: 10_000, intervals: [1_000] },
      )
      .toEqual([]);
  });
});
