# Assertions cheatsheet

The provisioner-agnostic surface a scenario body asserts through. Source of
truth: [`src/assertions/apim/mapi.ts`](../../../../src/assertions/apim/mapi.ts)
and [`gateway.ts`](../../../../src/assertions/apim/gateway.ts). If a method you
need is missing, add it via the `extend-platform-assertions` skill rather than
inlining a `fetch` in the test.

## Naming contract

| Prefix | Behaviour |
|---|---|
| `assertX` | single shot; throws `AssertionError` on mismatch |
| `waitForX` | polled variant; takes a trailing `PollOptions` (`{ timeoutMs, intervalMs, description }`) |
| `checkX` | non-throwing; returns an `AssertionReport` for the caller to inspect |
| `listX` / `createX` / `deleteX` | plain data access to build assertions on |

**Anything read after an apply must be a `waitFor*`.** Both `kubectl apply` and
`terraform apply` return before APIM has converged.

## mAPI: APIs

```ts
await mapi.assertApiMatches(apiId, { name, definitionVersion: "V4", type: "PROXY", state: "STARTED" });
await mapi.waitForApiMatches(apiId, { visibility: "PUBLIC", lifecycleState: "PUBLISHED" }, { timeoutMs: 30_000 });
await mapi.checkApiMatches(apiId, {…});          // returns { pass, failures }

await mapi.assertApiState(apiId, "STARTED");
await mapi.assertApiStarted(apiId);   await mapi.waitForApiStarted(apiId, opts);
await mapi.assertApiStopped(apiId);   await mapi.waitForApiStopped(apiId, opts);
await mapi.assertApiHttpStatus(apiId, 404);      // e.g. after a delete
await mapi.waitForApiAbsent(apiId, opts);

await mapi.deleteApi(apiId, /* closePlans */ true);
```

Matching is **partial and positional**: only fields present in `expected` are
checked, objects match recursively, arrays match by index, primitives are strict
equality.

## mAPI: plans, subscriptions, api-keys

```ts
await mapi.listApiPlans(apiId);
await mapi.assertPlanMatches(apiId, planId, { security: { type: "JWT" } });
await mapi.assertPlanStatus(apiId, planId, "PUBLISHED");
await mapi.assertPlanPublished(apiId, planId);

await mapi.assertSubscriptionMatches(apiId, subId, {…});
await mapi.assertSubscriptionStatus(apiId, subId, "ACCEPTED");
await mapi.assertSubscriptionAccepted(apiId, subId);

await mapi.listSubscriptionApiKeys(apiId, subId);        // includes revoked/expired
await mapi.listActiveSubscriptionApiKeys(apiId, subId);  // filtered client-side
await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 2, opts);
```

APIM's api-key listing endpoint has **no server-side filter**, so it returns
revoked and expired keys; use the `Active` variant or filter on
`revoked`/`expired` yourself. API-key **values are unique per API** including
already-revoked entries, so custom-key tests need a per-run unique value.

## mAPI: applications, groups, categories, members, pages

```ts
await mapi.assertApplicationMatches(appId, {…});
await mapi.waitForApplicationMatches(appId, {…}, opts);
await mapi.assertApplicationHttpStatus(appId, 404);
await mapi.listApplicationMembers(appId);
await mapi.listApplicationMetadata(appId);    // the detail response omits metadata

await mapi.listGroups();
await mapi.assertGroupMatches(hrid, {…});        await mapi.waitForGroupMatches(hrid, {…}, opts);
await mapi.assertGroupMatchesById(id, {…});      await mapi.waitForGroupMatchesById(id, {…}, opts);
await mapi.assertGroupAbsent(hrid);              await mapi.waitForGroupAbsent(hrid, opts);
await mapi.assertGroupAbsentById(id);            await mapi.waitForGroupAbsentById(id, opts);
await mapi.deleteGroup(groupId);
await mapi.createServiceAccount(name);           // a resolvable member for group/member journeys

await mapi.listCategories();
await mapi.createCategory({ key, name, description });   // APIM drops unknown category refs silently
await mapi.deleteCategory(categoryId);

await mapi.listApiMembers(apiId);
await mapi.listApiPages(apiId);
```

Assert the **whole member set per stage**, not just the added entry, so a
duplicated or stale member fails the test.

## Gateway

```ts
await gateway.assertResponds("/ctx", { status: 200 });
await gateway.assertResponds("/ctx", { status: 200, headers: { Authorization: `Bearer ${token}` } });
await gateway.assertNotResponds("/ctx", { notStatus: 200 });
await gateway.assertRespondsWithHeaders("/ctx", { "X-Added": "yes", "X-Removed": null }, { status: 200 });
```

- All three **poll internally** (500ms interval, 30s window by default).
- `assertNotResponds` treats a connection or TLS error as "not responding",
  which is how a rejected mTLS handshake is asserted.
- `assertRespondsWithHeaders` matches header names case-insensitively; a `null`
  expectation asserts the header is **absent**, which is how "the policy stopped
  adding it" is checked.
- A policy is only really proven **at the gateway**: APIM will happily record a
  flow the gateway never applies.

mTLS:

```ts
import { createTlsFetch } from "@gravitee/platform-test/utils/http";
const secureGw = mapi.gateway({ baseUrl: mtlsGatewayBaseUrl }, createTlsFetch({ cert, key, ca }));
```

## Provisioner readout

```ts
import { assertProvisioner } from "@gravitee/platform-test/provisioners";

await provisioned.remove("subscription");
await assertProvisioner(provisioned, "subscription", "gone");
```

Three states only: `applied`, `failed`, `gone`. Use after `remove()` and on
failure paths. **Not** as a convergence wait after `update()`: GKO's condition
`observedGeneration` can lag after a re-apply (GKO-2940), so a post-update read
can return the pre-update `Accepted=True` and pass without asserting anything.
`mapi` is the convergence signal after an update.

The Terraform arm needs `addresses: { role: "apim_x.name" }` in its scenario
spec (with a `[0]` suffix for a `count`-gated resource) and throws without it.

## Polling utilities

```ts
import { poll, deepPartialMatch } from "@gravitee/platform-test";

await poll(() => mapi.assertApiStarted(apiId), {
  timeoutMs: 15_000,
  intervalMs: 1_000,
  description: "API to reach STARTED state",
});

// Playwright-native, for anything not covered by a waitFor*:
await expect.poll(() => fetchSomething()).toMatchObject({ a: 1, b: 2 });
```

**Combine polled checks atomically.** One `expect.poll(...).toMatchObject({...})`
over every field, not one poll per field followed by a fresh fetch for the rest:
the second read can observe a different state.

## Anti-patterns

| Do not | Do |
|---|---|
| bump `timeout` to make a test pass | find the real convergence signal and poll it |
| `test.skip` a failing test | `since()` if it is a version gate, `pending` if the arm is unimplemented, otherwise fix it |
| assert immediately after `kubectl apply` / `terraform apply` | `waitFor*` or `expect.poll` |
| poll one field, then re-fetch for the rest | one atomic `toMatchObject` |
| branch the body on `provisionerId` | move the difference into fixtures, `params.ts`, or a `-gko-only` / `-tf-only` file |
| `kubectl patch` / `kubectl annotate` to trigger a reconcile | re-`kubectl apply -f` a modified CR file |
| assert only that a resource exists | assert the fields that carry regression value |
