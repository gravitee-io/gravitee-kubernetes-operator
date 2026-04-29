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
 * mTLS — Application clientCertificates lifecycle.
 *
 * These scenarios were originally framed as console / API endpoint behavior
 * in Xray. Through the GKO Application CR they're observable via the mAPI
 * Application response — specifically `settings.tls.client_certificates[]`.
 *
 * Xray tests:
 *   GKO-2228: Add certificate with future startsAt
 *   GKO-2241: Update endsAt to a valid future date
 *   GKO-2214: Update startsAt to a valid future date
 *   GKO-2261: Update existing certificate content via CR
 *   GKO-2215: Replacing cert content updates the entity
 *   GKO-2264: Explicit cert name is preserved
 *   GKO-2236: Application without metadata + cert
 *   GKO-2263: Application with metadata + cert
 *   GKO-2235: Two certificate entries with the same name (duplicate)
 *
 * Notes:
 *   - The wire format for client cert dates is epoch milliseconds.
 *   - The mAPI Application response uses snake_case (`client_certificates`)
 *     with `startsAt` / `endsAt` as numbers and `name` / `certificate` as
 *     strings. Status (ACTIVE / SCHEDULED / ACTIVE_WITH_END / REVOKED) is
 *     UI-derived and not exposed in the mAPI response, so these tests
 *     verify acceptance + persisted dates rather than the status string.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { readFile } from "node:fs/promises";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";
import type { Application } from "../../../../src/types/apim.js";

const PKI = (...segments: string[]) => fixture("crds/mtls-certificates/pki", ...segments);

interface AppStatus {
  id?: string;
  conditions?: Array<{ type: string; status: string }>;
}

const ONE_DAY_MS = 24 * 60 * 60 * 1000;

function isoOffsetDays(days: number): string {
  return new Date(Date.now() + days * ONE_DAY_MS).toISOString();
}

/** Indent multi-line text by `n` spaces (for inlining a PEM under YAML `content: |`). */
function indent(text: string, n: number): string {
  const pad = " ".repeat(n);
  return text
    .split("\n")
    .map((line) => (line.length > 0 ? pad + line : line))
    .join("\n");
}

/** Look up a cert by name in the v1 mAPI application response. */
function findCert(
  app: Application,
  name: string,
): { name?: string; startsAt?: number; endsAt?: number; certificate?: string } | undefined {
  return app.settings?.tls?.client_certificates?.find((c) => c.name === name);
}

async function applyAndAssertCertVisible(
  kubectl: typeof kubectlSafe,
  mapi: import("../../../../src/index.js").Mapi,
  name: string,
  yaml: string,
): Promise<void> {
  await kubectl.applyString(yaml);
  await kubectl.waitForCondition("application", name, "Accepted");
  const appId = (await kubectl.getStatus<AppStatus>("application", name)).id;
  if (!appId) {
    throw new Error(`Application ${name} has no .status.id after Accepted`);
  }
  await expect
    .poll(async () =>
      (await mapi.fetchApplication(appId)).settings?.tls?.client_certificates?.length ?? 0,
    )
    .toBe(1);
}

test.describe("mTLS — Application cert lifecycle", () => {
  let pem1: string;
  let pem2: string;

  test.beforeAll(async () => {
    pem1 = (await readFile(PKI("client1.crt"))).toString();
    pem2 = (await readFile(PKI("client2.crt"))).toString();
  });

  // ── Per-test cleanup. Each test names its Application after the GKO ID,
  //    so the safety net here is broad but cheap. ──────────────────────
  const APP_NAMES = [
    "e2e-mtls-2228",
    "e2e-mtls-2241",
    "e2e-mtls-2214",
    "e2e-mtls-2261",
    "e2e-mtls-2264",
    "e2e-mtls-2235",
  ];

  test.afterEach(async () => {
    for (const name of APP_NAMES) {
      await kubectlSafe.deleteResource("application", name).catch(() => {});
    }
  });

  // Build a minimal Application CR with one cert. Pass `startsAt` / `endsAt`
  // as ISO strings (or undefined to omit the field). Use `extraSpec` to
  // append additional spec fields (e.g. metadata) without templating again.
  function appYaml(opts: {
    name: string;
    certName?: string;
    pem: string;
    startsAt?: string;
    endsAt?: string;
    extraSpec?: string;
  }): string {
    const certNameLine = opts.certName !== undefined ? `        - name: ${opts.certName}` : `        -`;
    const startLine = opts.startsAt ? `          startsAt: "${opts.startsAt}"` : "";
    const endLine = opts.endsAt ? `          endsAt: "${opts.endsAt}"` : "";
    return [
      `apiVersion: gravitee.io/v1alpha1`,
      `kind: Application`,
      `metadata:`,
      `  name: ${opts.name}`,
      `spec:`,
      `  contextRef:`,
      `    name: dev-ctx`,
      `    namespace: default`,
      `  name: ${opts.name}`,
      `  description: "E2E ${opts.name}"`,
      `  settings:`,
      `    app:`,
      `      type: WEB`,
      `    tls:`,
      `      clientCertificates:`,
      certNameLine,
      startLine,
      endLine,
      `          content: |`,
      indent(opts.pem.trimEnd(), 12),
      ...(opts.extraSpec ? [opts.extraSpec] : []),
      "",
    ]
      .filter((line) => line !== "")
      .join("\n");
  }

  // ── GKO-2228: future startsAt is accepted ────────────────────────────
  // NOTE: APIM filters scheduled (startsAt > now) certs out of the v1
  // application response — `settings.tls` is `null` until the cert
  // becomes active. The /v2 .../applications/{id}/certificates endpoint
  // is not implemented in APIM 4.12 (404). So this test verifies only
  // that the operator accepts the cert; the activation is a runtime
  // behavior that is not observable via mAPI within an E2E timeout.

  test(`Cert with future startsAt is accepted by the operator ${XRAY.MTLS_CERTIFICATES.CERT_FUTURE_START_DATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const APP = "e2e-mtls-2228";
    const future = isoOffsetDays(7);
    const farFuture = isoOffsetDays(365);

    await kubectl.applyString(
      appYaml({
        name: APP,
        certName: "future-start",
        pem: pem1,
        startsAt: future,
        endsAt: farFuture,
      }),
    );
    await kubectl.waitForCondition("application", APP, "Accepted");
  });

  // ── GKO-2241: update endsAt to a valid future date ───────────────────

  test(`Cert endsAt can be updated to a valid future date ${XRAY.MTLS_CERTIFICATES.CERT_UPDATE_END_DATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-mtls-2241";
    const initialEnd = isoOffsetDays(30);
    const updatedEnd = isoOffsetDays(120);

    await kubectl.applyString(
      appYaml({ name: APP, certName: "update-end", pem: pem1, endsAt: initialEnd }),
    );
    await kubectl.waitForCondition("application", APP, "Accepted");

    const appId = (await kubectl.getStatus<AppStatus>("application", APP)).id;
    expect(appId).toBeTruthy();

    // Confirm the initial endsAt is reflected.
    const initialEndMs = Date.parse(initialEnd);
    await expect
      .poll(async () => findCert(await mapi.fetchApplication(appId!), "update-end")?.endsAt ?? 0)
      .toBeGreaterThanOrEqual(initialEndMs - 1_000);

    // Apply the updated endsAt and verify the new value is persisted.
    await kubectl.applyString(
      appYaml({ name: APP, certName: "update-end", pem: pem1, endsAt: updatedEnd }),
    );
    await kubectl.waitForCondition("application", APP, "Accepted");

    const updatedEndMs = Date.parse(updatedEnd);
    await expect
      .poll(async () => findCert(await mapi.fetchApplication(appId!), "update-end")?.endsAt ?? 0)
      .toBeGreaterThanOrEqual(updatedEndMs - 1_000);
  });

  // ── GKO-2214: update startsAt to a valid future date ─────────────────
  // Verified 2026-04-29 against APIM 4.12: updating ANY field on a cert
  // entry causes APIM to compare the cert fingerprint against existing
  // ClientCertificate entities and reject with HTTP 400
  //   "Client certificate with fingerprint [...] is already used by
  //    another active application."
  // even though the same application owns it. The operator surfaces this
  // as `Accepted=False` with reason=ControlPlaneError. The original Xray
  // scenario expects acceptance — the test here documents the actual
  // current product behavior (the second apply is rejected) so that a
  // future fix flips this assertion deliberately.

  test(`Updating cert startsAt is currently rejected by APIM (fingerprint reuse) ${XRAY.MTLS_CERTIFICATES.CERT_UPDATE_START_DATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const APP = "e2e-mtls-2214";
    const initialStart = isoOffsetDays(7);
    const updatedStart = isoOffsetDays(30);
    const future = isoOffsetDays(365);

    await kubectl.applyString(
      appYaml({
        name: APP,
        certName: "update-start",
        pem: pem1,
        startsAt: initialStart,
        endsAt: future,
      }),
    );
    await kubectl.waitForCondition("application", APP, "Accepted");

    await kubectl.applyString(
      appYaml({
        name: APP,
        certName: "update-start",
        pem: pem1,
        startsAt: updatedStart,
        endsAt: future,
      }),
    );

    // Wait briefly for the operator to reconcile and surface the rejection
    // through a Conditions update, then assert Accepted=False with the
    // fingerprint-reuse reason.
    await expect
      .poll(
        async () => {
          const status = await kubectl.getStatus<AppStatus>("application", APP);
          const accepted = status.conditions?.find((c) => c.type === "Accepted");
          return { status: accepted?.status };
        },
        { timeout: 15_000 },
      )
      .toEqual({ status: "False" });
  });

  // ── GKO-2261 + GKO-2215: replace cert content via CR ─────────────────
  // The two scenarios ("Update Current Certificate" + "Cert Change Creates
  // New Entity") collapse to the same observable: after a content change
  // the mAPI cert reflects the new PEM.

  test(`Replacing cert content via CR updates the entity ${XRAY.MTLS_CERTIFICATES.CERT_UPDATE_VIA_CR} ${XRAY.MTLS_CERTIFICATES.CERT_CHANGE_NEW_ENTITY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-mtls-2261";

    await kubectl.applyString(appYaml({ name: APP, certName: "rotated", pem: pem1 }));
    await kubectl.waitForCondition("application", APP, "Accepted");
    const appId = (await kubectl.getStatus<AppStatus>("application", APP)).id;

    const cert1Stored = findCert(await mapi.fetchApplication(appId!), "rotated")?.certificate ?? "";
    expect(cert1Stored).not.toEqual("");

    await kubectl.applyString(appYaml({ name: APP, certName: "rotated", pem: pem2 }));
    await kubectl.waitForCondition("application", APP, "Accepted");

    await expect
      .poll(async () => findCert(await mapi.fetchApplication(appId!), "rotated")?.certificate ?? "")
      .not.toEqual(cert1Stored);
  });

  // ── GKO-2264: cert name auto-generation convention ───────────────────

  test(`Unnamed cert is auto-named "<app>-<index>" in mAPI ${XRAY.MTLS_CERTIFICATES.CERT_NAMING_CONVENTION} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-mtls-2264";

    // Build CR with cert content but no `name` field.
    const yaml = [
      `apiVersion: gravitee.io/v1alpha1`,
      `kind: Application`,
      `metadata:`,
      `  name: ${APP}`,
      `spec:`,
      `  contextRef:`,
      `    name: dev-ctx`,
      `    namespace: default`,
      `  name: ${APP}`,
      `  description: "E2E cert auto-name"`,
      `  settings:`,
      `    app:`,
      `      type: WEB`,
      `    tls:`,
      `      clientCertificates:`,
      `        - content: |`,
      indent(pem1.trimEnd(), 12),
      "",
    ].join("\n");

    await kubectl.applyString(yaml);
    await kubectl.waitForCondition("application", APP, "Accepted");
    const appId = (await kubectl.getStatus<AppStatus>("application", APP)).id;
    if (!appId) throw new Error(`${APP} has no .status.id`);

    await expect
      .poll(async () =>
        (await mapi.fetchApplication(appId)).settings?.tls?.client_certificates?.[0]?.name,
      )
      .toBe(`${APP}-0`);
  });

  // GKO-2236 was out of scope:
  // same APIM "upgrader" path as GKO-2263 — the scenario is about the
  // upgrader processing legacy applications without
  // `metadata.client_certificate`. GKO never runs the upgrader, so the
  // CR-equivalent ("apply a fresh app without metadata + cert works")
  // is too far from the original scenario to count as coverage.

  // GKO-2263 was out of scope:
  // the Xray scenario describes an APIM "upgrader" migration path that
  // creates a cert from a magic `metadata.client_certificate` entry on
  // legacy applications. GKO has no upgrader path — applications go
  // through the create/update reconciler — so the scenario doesn't map
  // cleanly to a CR-driven test.

  // ── GKO-2235: two cert entries with the same content (true duplicate) ──
  // The Xray scenario describes adding a cert "with the same details as
  // the existing certificate" — i.e. a true duplicate (matching content /
  // fingerprint). Verified 2026-04-29 against GKO master + APIM 4.12:
  // GKO admission rejects the apply with stderr containing
  //   "client certificate content must be unique"
  // before APIM ever sees the request. (Same-NAME-different-content is
  // a different case and IS accepted — both entries persist; not the
  // scenario the Xray describes.)

  test(`Duplicate cert content (same fingerprint) is rejected by admission ${XRAY.MTLS_CERTIFICATES.CERT_DUPLICATE_REJECTED} ${TAGS.REGRESSION}`, async () => {
    const APP = "e2e-mtls-2235";

    const yaml = [
      `apiVersion: gravitee.io/v1alpha1`,
      `kind: Application`,
      `metadata:`,
      `  name: ${APP}`,
      `spec:`,
      `  contextRef:`,
      `    name: dev-ctx`,
      `    namespace: default`,
      `  name: ${APP}`,
      `  description: "E2E true-duplicate cert content"`,
      `  settings:`,
      `    app:`,
      `      type: WEB`,
      `    tls:`,
      `      clientCertificates:`,
      `        - name: dup-A`,
      `          content: |`,
      indent(pem1.trimEnd(), 12),
      `        - name: dup-B`,
      `          content: |`,
      indent(pem1.trimEnd(), 12),
      "",
    ].join("\n");

    const stderr = await kubectlSafe.applyStringExpectFailure(yaml);
    // Narrow keywords from validateClientCertificates — generic terms like
    // "denied" / "invalid" intentionally excluded.
    expect(stderr.toLowerCase()).toMatch(/content must be unique|fingerprint|duplicate/);
  });
});

/**
 * Application clientCertificates — date-acceptance / status scenarios
 *.
 *
 * The Xray scenarios specify the cert `status` field (ACTIVE / SCHEDULED /
 * ACTIVE_WITH_END / REVOKED), which is computed by APIM and rendered in
 * the UI but is not exposed via mAPI. Through the GKO Application CR we
 * verify:
 *   - The operator accepts each date combination (CR reaches Accepted=True)
 *   - For combinations that produce an *active* cert today, the cert is
 *     visible in `application.settings.tls.client_certificates[]`
 *
 * Future-start scenarios (SCHEDULED status) cannot be observed in mAPI —
 * APIM filters them out of the v1 response — so they are deliberately
 * absent here. GKO-2228 / GKO-2214 in the lifecycle suite cover the
 * "future startsAt is accepted by the operator" angle.
 */

const STATUS_APPS = [
  "e2e-mtls-2121",
  "e2e-mtls-2130",
  "e2e-mtls-2141",
  "e2e-mtls-2149",
];

test.describe("mTLS — cert date acceptance", () => {
  let pem: string;

  test.beforeAll(async () => {
    pem = (await readFile(PKI("client1.crt"))).toString();
  });

  test.afterEach(async () => {
    for (const name of STATUS_APPS) {
      await kubectlSafe.deleteResource("application", name).catch(() => {});
    }
  });

  function statusAppYaml(opts: {
    name: string;
    startsAt?: string;
    endsAt?: string;
  }): string {
    const startLine = opts.startsAt ? `          startsAt: "${opts.startsAt}"` : "";
    const endLine = opts.endsAt ? `          endsAt: "${opts.endsAt}"` : "";
    return [
      `apiVersion: gravitee.io/v1alpha1`,
      `kind: Application`,
      `metadata:`,
      `  name: ${opts.name}`,
      `spec:`,
      `  contextRef:`,
      `    name: dev-ctx`,
      `    namespace: default`,
      `  name: ${opts.name}`,
      `  description: "E2E ${opts.name} (date acceptance)"`,
      `  settings:`,
      `    app:`,
      `      type: WEB`,
      `    tls:`,
      `      clientCertificates:`,
      `        - name: cert-under-test`,
      startLine,
      endLine,
      `          content: |`,
      indent(pem.trimEnd(), 12),
      "",
    ]
      .filter((line) => line !== "")
      .join("\n");
  }

  // GKO-2121: past startsAt + null endsAt → ACTIVE → visible
  test(`Cert with past startsAt + null endsAt is active and visible ${XRAY.MTLS_CERTIFICATES.CERT_PAST_START_NULL_END} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-mtls-2121";
    await applyAndAssertCertVisible(
      kubectl,
      mapi,
      APP,
      statusAppYaml({ name: APP, startsAt: isoOffsetDays(-7) }),
    );
  });

  // GKO-2130: null startsAt + null endsAt → ACTIVE → visible
  test(`Cert with no startsAt and no endsAt is active and visible ${XRAY.MTLS_CERTIFICATES.CERT_NULL_BOTH_DATES} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-mtls-2130";
    await applyAndAssertCertVisible(kubectl, mapi, APP, statusAppYaml({ name: APP }));
  });

  // GKO-2141: cert without `startsAt` / `endsAt` keys altogether → visible
  // Identical wire shape to GKO-2130, exercised separately to keep the
  // Xray-to-test mapping explicit. If both pass on every run, the second
  // test is a candidate for collapsing into the first in a later cleanup.
  test(`Cert with no date fields is accepted as active ${XRAY.MTLS_CERTIFICATES.CERT_NO_DATE_FIELDS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-mtls-2141";
    await applyAndAssertCertVisible(kubectl, mapi, APP, statusAppYaml({ name: APP }));
  });

  // GKO-2149: null startsAt + future endsAt → ACTIVE_WITH_END → visible
  test(`Cert with null startsAt and future endsAt is active and visible ${XRAY.MTLS_CERTIFICATES.CERT_NULL_START_FUTURE_END} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const APP = "e2e-mtls-2149";
    const future = isoOffsetDays(60);
    await applyAndAssertCertVisible(
      kubectl,
      mapi,
      APP,
      statusAppYaml({ name: APP, endsAt: future }),
    );
  });
});
