# Hand-rolled single-provisioner tests

For `e2e/tests/gko/<area>/` (the **operator** under test) and
`e2e/tests/terraform/` (the **provider** under test). Everything a customer could
also do through the other driver belongs in a user journey instead: see the main
skill's routing table.

Legitimate `tests/gko/` subjects: admission webhooks and CRD schema validation,
CR `.status` conditions and reconciliation, ConfigMap/Secret templating,
ManagementContext lifecycle, CRD defaulting, import/export round-trips, V2 APIs
(the Automation API is V4-only by design).

Legitimate `tests/terraform/` subjects: drift detection, `terraform import`,
data sources, `plan -detailed-exitcode`, provider-side schema validation,
sensitive-attribute redaction.

---

## Shape

```ts
<LICENCE>

/**
 * <What is under test and why it is provisioner-specific.>
 *
 * Xray tests:
 *   GKO-NNNN: <one line per test in this file>
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, expect } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

const RESOURCE_NAME = "e2e-<area>-<what>";

test.describe(`<Area> — <behaviour> ${PROVISIONER.GKO}`, () => {
  // Safety-net cleanup: inline cleanup does not run when a test times out.
  // Reverse dependency order; each delete ignores errors.
  test.afterEach(async () => {
    await kubectl.deleteResource("apiv4definition", RESOURCE_NAME).catch(() => {});
  });

  test(`<what it proves> ${XRAY.AREA.ID} ${TAGS.REGRESSION}`, async ({ kubectl, mapi }) => {
    // …
  });
});
```

**The provisioner tag on the describe title is mandatory.** Lane selection is by
title tag alone: a missing tag means the test runs in *every* lane, a wrong one
means it runs in the wrong lane, two means it runs in none. All three still pass
an untagged full run, which is why `npm run check:lanes` exists. Run it.

## Cleanup

Two ways to clean up, both needed:

```ts
// 1. By resource kind + name, when the manifest was built inline:
await kubectl.deleteResource("subscription", SUB_NAME).catch(() => {});
await kubectl.deleteResource("application", APP_NAME).catch(() => {});
await kubectl.deleteResource("apiv4definition", API_NAME).catch(() => {});

// 2. By fixture file, when the manifest came from e2e/fixtures/:
for (const f of ["<sub>/crd.yaml", "<app>/crd.yaml", "<api>/crd.yaml"]) {
  await kubectl.del(fixture(f)).catch(() => {});
}
```

- **Reverse dependency order: subscriptions, applications, APIs.** The admission
  webhook blocks deleting an application that still has subscriptions.
- **Always `.catch(() => {})`** on each call so one missing resource does not
  abort the rest of the teardown.
- **Never delete shared preconditions**: the `dev-ctx` ManagementContext, or the
  APIM org/environment Terraform authenticates against.
- Even a test whose resource should have been *rejected* needs the safety net:
  if a regression makes admission accept it, it must not leak downstream.

For Terraform, track the workspace and call `destroyWorkspace(ws)` in
`afterAll`; it re-runs `destroy` as a no-op when a test already destroyed
inline, so it is always safe.

## kubectl engine

`e2e/helpers/kubectl.ts` re-exports `src/provisioners/engines/kubectl.ts`
unchanged. It is also available as a Playwright fixture (`async ({ kubectl }) =>`).

| Function | Use |
|---|---|
| `apply(yamlPath, ns?)` / `applyString(yaml, ns?)` | create or update |
| `applyExpectFailure(path)` / `applyStringExpectFailure(yaml)` | **negative tests**: returns stderr for assertion |
| `applyDryRun(yamlPath, ns?)` | exercise admission without persisting |
| `del(yamlPath, ns?)` / `deleteResource(kind, name, ns?)` | teardown |
| `delExpectFailure(yamlPath)` | assert a deletion is blocked |
| `get<T>(kind, name, ns?)` / `getStatus<T>()` / `getField<T>()` / `getCondition()` | read CR state |
| `exists()` / `waitForExists()` / `waitForDeletion()` / `waitForCondition()` | wait on cluster state |
| `findPod()` / `rolloutRestart()` / `waitForRollout()` | operator-lifecycle tests |
| `writeConfigMap()` / `readConfigMap()` | cross-process handshakes (upgrade survival) |

**To trigger a reconcile, re-`apply` a modified CR file.** Not `kubectl patch`,
not `kubectl annotate`.

## Negative tests: pin the rejection to its cause

A rejection test that only asserts "apply failed" passes for the wrong reason as
soon as the manifest drifts. Assert the value **and** its field path:

```ts
const stderr = await kubectl.applyStringExpectFailure(v4Manifest({ name, type: "EDGE" }));
expect(stderr.toLowerCase()).toContain("unsupported value");
expect(stderr).toContain("EDGE");
expect(stderr).toContain("spec.type");
```

Keep the rest of the manifest valid, so the offending field is the only possible
cause. `metadata.name` must stay RFC 1123 DNS-safe, so map underscores to
hyphens when deriving a name from an enum value.

## Fixtures

Hand-rolled tests use the central `e2e/fixtures/<domain>/<scenario>/` tree via
`fixture("domain/scenario/crd.yaml")`. Conventions:

- Domain folders mirror the `tests/gko/` folders (`admission-webhook/`,
  `api-lifecycle/`, `categories/`, `policies/`).
- Scenario folders describe **what is being tested**, not what kind of CR sits at
  the top of the manifest: a V4 API with a JWT plan is `plans/v4-jwt/`, not
  `api-v4-definitions/`.
- Prefix with `v2-` / `v4-` inside domains holding both.
- Build the manifest inline (a template function plus `applyString`) when the
  test needs to vary one field across many cases; a fixture file per enum value
  is noise.

## Docstring

Every file opens with: what is under test, why it is provisioner-specific, the
`Xray tests:` list (one line per test, id plus intent), and `Preconditions:`.
Reference **Xray Test ids only**, never a Story or bug key.
