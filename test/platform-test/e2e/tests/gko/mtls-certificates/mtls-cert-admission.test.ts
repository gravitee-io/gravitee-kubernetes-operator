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
 * mTLS — Application clientCertificates admission rejections (batch 8).
 *
 * GKO does not have a standalone MTLSCertificate CRD; the original Xray tests
 * (GKO-2117 etc.) describe APIM cert-management API behavior. We exercise the
 * equivalent paths through the Application CRD's spec.settings.app.tls
 * (clientCertificate / clientCertificates), which lands in:
 *   - internal/admission/application/validate.go (validateSingleClientCertificate,
 *     validateClientCertificates)
 *   - apim DryRunCreateOrUpdate (APIM-side validation)
 *
 * Xray tests:
 *   GKO-2117: Bad PEM in deprecated clientCertificate
 *   GKO-2118: Both content and ref set on a clientCertificates entry
 *   GKO-2122: Cert startsAt == endsAt
 *   GKO-2124: Cert name exceeds APIM length limit
 *   GKO-2125: clientCertificates entry missing content/ref entirely
 *   GKO-2131: Cert with both dates in the past (already expired)
 *   GKO-2133: Cert name with invalid characters
 *   GKO-2135: Cert with endsAt before startsAt
 *   GKO-2143: Cert with empty name
 *   GKO-2146: Cert with explicit empty content (variant of 2125)
 *   GKO-2148: Cert with invalid date-format strings
 *
 * Each test uses kubectl.applyExpectFailure and asserts the rejection
 * stderr contains a relevant keyword. If a scenario turns out to be
 * accepted by GKO+APIM (i.e. the test fails because admission did NOT
 * reject), document it as APIM-only and adjust during live-cluster
 * validation.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

const F = (name: string) => fixture(`crds/mtls-certificates/invalid/${name}.yaml`);

test.describe("mTLS — Application cert admission rejections (batch 8)", () => {
  test(`Deprecated clientCertificate with non-PEM content is rejected ${XRAY.MTLS_CERTIFICATES.CRD_BAD_PEM} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2117-bad-pem"));
    expect(stderr.toLowerCase()).toMatch(/parse|certificate|pem|tls|invalid|denied/);
  });

  test(`clientCertificates with both content and ref is rejected ${XRAY.MTLS_CERTIFICATES.CRD_FORBIDDEN_FIELD_UPDATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2118-content-and-ref"));
    expect(stderr.toLowerCase()).toMatch(/content|ref|both|cannot|denied/);
  });

  test(`clientCertificate with startsAt equal to endsAt is rejected ${XRAY.MTLS_CERTIFICATES.CRD_START_EQ_END} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2122-start-eq-end"));
    expect(stderr.toLowerCase()).toMatch(/date|window|valid|invalid|certificate|denied/);
  });

  test(`clientCertificate with overly long name is rejected ${XRAY.MTLS_CERTIFICATES.CRD_NAME_TOO_LONG} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2124-name-too-long"));
    expect(stderr.toLowerCase()).toMatch(/name|length|too long|maximum|invalid|denied/);
  });

  test(`clientCertificates entry missing content and ref is rejected ${XRAY.MTLS_CERTIFICATES.CRD_MISSING_FIELDS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2125-no-content-no-ref"));
    expect(stderr.toLowerCase()).toMatch(/content|ref|either|must be set|denied/);
  });

  test(`Already-expired clientCertificate is rejected ${XRAY.MTLS_CERTIFICATES.CRD_EXPIRED_REJECTED} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2131-expired"));
    expect(stderr.toLowerCase()).toMatch(/expired|date|past|invalid|certificate|denied/);
  });

  test(`clientCertificate name with invalid characters is rejected ${XRAY.MTLS_CERTIFICATES.CRD_INVALID_CHARS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2133-invalid-chars"));
    expect(stderr.toLowerCase()).toMatch(/name|character|invalid|denied/);
  });

  test(`clientCertificate with endsAt before startsAt is rejected ${XRAY.MTLS_CERTIFICATES.CRD_END_BEFORE_START} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2135-end-before-start"));
    expect(stderr.toLowerCase()).toMatch(/date|before|startsAt|endsAt|invalid|denied/);
  });

  test(`clientCertificate with empty name is rejected ${XRAY.MTLS_CERTIFICATES.CRD_MISSING_NAME} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2143-missing-name"));
    expect(stderr.toLowerCase()).toMatch(/name|required|invalid|denied/);
  });

  test(`clientCertificate with explicit empty content is rejected ${XRAY.MTLS_CERTIFICATES.CRD_MISSING_CERT_FIELD} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2146-missing-cert-content"));
    expect(stderr.toLowerCase()).toMatch(/content|ref|either|must be set|denied/);
  });

  test(`clientCertificate with invalid date-format strings is rejected ${XRAY.MTLS_CERTIFICATES.CRD_INVALID_DATA_DATES} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2148-invalid-date-format"));
    expect(stderr.toLowerCase()).toMatch(/date|format|parse|invalid|denied/);
  });
});
