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
 *   GKO-2125: clientCertificates entry missing content/ref entirely
 *   GKO-2131: Cert with both dates in the past (already expired)
 *   GKO-2135: Cert with endsAt before startsAt
 *   GKO-2146: Cert with explicit empty content (variant of 2125)
 *   GKO-2148: Cert with invalid date-format strings
 *
 * Each test uses kubectl.applyExpectFailure and asserts the rejection
 * stderr contains a relevant keyword. If a scenario turns out to be
 * accepted by GKO+APIM (i.e. the test fails because admission did NOT
 * reject), document it as APIM-only and adjust during live-cluster
 * validation.
 *
 * Deferred (not covered here):
 *   GKO-2124: name exceeds APIM length limit
 *   GKO-2133: name with invalid characters
 *   GKO-2143: empty name
 *   The GKO admission webhook rejects on PEM validity ("certificate is not
 *   a valid pem") before any cert-name validator runs, so these name-rule
 *   scenarios aren't reachable through the Application CRD path with the
 *   current stub PEM. Exercising them requires a valid PEM plus confirmation
 *   that APIM's dryRun actually enforces these name rules — tracked for a
 *   follow-up batch alongside per-scenario PKI fixtures.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const F = (name: string) => fixture(`crds/mtls-certificates/invalid/${name}.yaml`);

// Every rejection regex is deliberately scoped to substrings that should
// appear for the specific failure mode under test. Generic terms like
// "invalid" or "denied" are intentionally excluded — they appear in almost
// every admission/webhook error and would turn a real regression (e.g. the
// admission path erroring for an unrelated reason) into a false pass.
test.describe("mTLS — Application cert admission rejections (batch 8)", () => {
  test.afterEach(async () => {
    // Admission-rejected applies don't persist state, but be defensive in
    // case a scenario ever slips past admission — the safety net matches the
    // codebase's cleanup pattern.
    await kubectlSafe.del(F("2117-bad-pem")).catch(() => {});
    await kubectlSafe.del(F("2118-content-and-ref")).catch(() => {});
    await kubectlSafe.del(F("2122-start-eq-end")).catch(() => {});
    await kubectlSafe.del(F("2125-no-content-no-ref")).catch(() => {});
    await kubectlSafe.del(F("2131-expired")).catch(() => {});
    await kubectlSafe.del(F("2135-end-before-start")).catch(() => {});
    await kubectlSafe.del(F("2146-missing-cert-content")).catch(() => {});
    await kubectlSafe.del(F("2148-invalid-date-format")).catch(() => {});
  });

  test(`Deprecated clientCertificate with non-PEM content is rejected ${XRAY.MTLS_CERTIFICATES.CRD_BAD_PEM} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2117-bad-pem"));
    // validateSingleClientCertificate → "failed to parse TLS client certificate"
    expect(stderr.toLowerCase()).toMatch(/parse|pem|tls client certificate/);
  });

  test(`clientCertificates with both content and ref is rejected ${XRAY.MTLS_CERTIFICATES.CRD_FORBIDDEN_FIELD_UPDATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2118-content-and-ref"));
    // validateClientCertificates → "content and ref cannot both be set"
    expect(stderr.toLowerCase()).toMatch(/content and ref|both|cannot/);
  });

  test(`clientCertificate with startsAt equal to endsAt is rejected ${XRAY.MTLS_CERTIFICATES.CRD_START_EQ_END} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2122-start-eq-end"));
    expect(stderr.toLowerCase()).toMatch(/startsat|endsat|validity|window|never/);
  });

  test(`clientCertificates entry missing content and ref is rejected ${XRAY.MTLS_CERTIFICATES.CRD_MISSING_FIELDS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2125-no-content-no-ref"));
    // validateClientCertificates → "either content or ref must be set"
    expect(stderr.toLowerCase()).toMatch(/either content or ref|must be set/);
  });

  test(`Already-expired clientCertificate is rejected ${XRAY.MTLS_CERTIFICATES.CRD_EXPIRED_REJECTED} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2131-expired"));
    expect(stderr.toLowerCase()).toMatch(/expired|past|endsat/);
  });

  test(`clientCertificate with endsAt before startsAt is rejected ${XRAY.MTLS_CERTIFICATES.CRD_END_BEFORE_START} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2135-end-before-start"));
    expect(stderr.toLowerCase()).toMatch(/startsat|endsat|before|after/);
  });

  test(`clientCertificate with explicit empty content is rejected ${XRAY.MTLS_CERTIFICATES.CRD_MISSING_CERT_FIELD} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2146-missing-cert-content"));
    // validateClientCertificates → "either content or ref must be set"
    expect(stderr.toLowerCase()).toMatch(/either content or ref|must be set/);
  });

  test(`clientCertificate with invalid date-format strings is rejected ${XRAY.MTLS_CERTIFICATES.CRD_INVALID_DATA_DATES} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(F("2148-invalid-date-format"));
    expect(stderr.toLowerCase()).toMatch(/date|rfc3339|parse|format/);
  });
});
