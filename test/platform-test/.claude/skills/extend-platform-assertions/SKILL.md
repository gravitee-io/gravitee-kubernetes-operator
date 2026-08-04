---
name: extend-platform-assertions
description: Add or change an assertion in the @gravitee/platform-test library - a new mapi or gateway method, a new APIM entity type, or a whole new product surface such as Access Management, Alert Engine or Cockpit. Use when a test needs a platform assertion that does not exist yet, when a test is about to inline a raw fetch against APIM, or when extending the assertion library's public API.
allowed-tools: Bash, Read, Edit, Write, Glob, Grep
---

# extend-platform-assertions

`src/assertions/` is the **provisioner-agnostic** half of the framework: the
assertions that stay identical no matter how a resource was created. Anything a
test inlines as a raw `fetch` against APIM is a missing method here.

Invariants: [`.agent/rules/e2e-framework-development.md`](../../../.agent/rules/e2e-framework-development.md).
Current surface: [`references` in the write-e2e-test skill](../write-e2e-test/references/assertions-cheatsheet.md).

Work from `test/platform-test/`.

---

## Step 1: check it does not already exist

```bash
grep -nE "^  async (assert|wait|check|list|create|delete)" src/assertions/apim/mapi.ts
grep -nE "^  async " src/assertions/apim/gateway.ts
```

`Mapi` already covers API, plan, subscription, api-key, application, group,
category, member and page families. Extend the right family rather than adding a
parallel one.

## Step 2: pick the method shape

One naming contract, no exceptions:

| Prefix | Behaviour | Returns |
|---|---|---|
| `assertX(...)` | single shot; throws `AssertionError` on mismatch | `Promise<void>` |
| `waitForX(..., options?: PollOptions)` | polled variant of the same assertion | `Promise<void>` |
| `checkX(...)` | non-throwing; the caller inspects the report | `Promise<AssertionReport>` |
| `listX` / `createX` / `deleteX` | plain data access for tests to build assertions on | the data |

**Add `waitForX` whenever the value converges asynchronously**, which is true of
anything written through the Automation API. A test that has to wrap your
`assertX` in its own `poll()` is a sign the `waitForX` is missing.

The body is always the same three lines:

```ts
async assertThingMatches(id: string, expected: DeepPartial<Thing>): Promise<void> {
  const thing = await this.fetchThing(id);
  const report = deepPartialMatch(thing, expected);
  throwIfFailed(report);
}

async waitForThingMatches(
  id: string,
  expected: DeepPartial<Thing>,
  options: PollOptions = {},
): Promise<void> {
  await poll(() => this.assertThingMatches(id, expected), {
    description: `thing ${id} to match`,
    ...options,
  });
}
```

Success is `void`; failure is a `node:assert` `AssertionError` produced by
`deepPartialMatch` + `throwIfFailed`. Never return a boolean, never `console.log`
a mismatch.

## Step 3: model only what is asserted

Add or extend the entity type in `src/types/apim.ts`.

- Matching is **partial**, so fields the suite never asserts do not need
  modelling. Mirroring the whole APIM response is churn that has to be
  maintained against every APIM release.
- Prefer the union of literal states (`"STARTED" | "STOPPED"`) over `string`, so
  a typo in a test fails at `npm run typecheck:e2e`.

## Step 4: know the endpoint's quirks

Encode the quirk in the method rather than leaving it for each test to
rediscover. Precedents worth copying:

- `listActiveSubscriptionApiKeys` exists because APIM's api-key listing has **no
  server-side filter** and returns revoked and expired keys.
- `listApplicationMetadata` is a separate call because the application detail
  response **omits** metadata.
- `assertRespondsWithHeaders` matches header names case-insensitively and treats
  a `null` expectation as "header absent", because that is how "the policy
  stopped adding it" is asserted.
- `assertNotResponds` counts a connection or TLS error as "not responding",
  because a rejected mTLS handshake produces no HTTP status at all.

Document the quirk in the method's doc comment, and add it to the
[cheatsheet](../write-e2e-test/references/assertions-cheatsheet.md).

## Step 5: a whole new product surface

`src/assertions/am/` is the placeholder for Access Management and shows the
intended shape. A new surface is a sibling of `src/assertions/apim/`:

```
src/assertions/<product>/
  index.ts        # public re-exports
  <client>.ts     # the class, built on HttpClient from src/utils/http/
```

Then:

- Add an `exports` entry in [`package.json`](../../../package.json)
  (`"./assertions/<product>"`), pointing at the built `dist/` paths, and mirror
  the pattern of the existing entries.
- Re-export from `src/index.ts` if it belongs on the top-level surface.
- Use `HttpClient` and native `fetch`. A new runtime dependency needs a reason
  that survives "can native `fetch` do this?".
- Add the config block and its env-var overrides to `src/config/`, `config.yaml`,
  and the override tables in `README.md` and `e2e/README.md`.

## Step 6: test it

Unit tests live under `test/` and run with vitest, **no cluster needed**:

```
test/apim/apim.test.ts        # Mapi
test/apim/gateway.test.ts     # Gateway
test/match/partial.test.ts    # deepPartialMatch
test/provisioners/            # provisioner layer
```

```bash
npm test
```

Cover the matching and report logic with a stubbed fetch: the pass case, the
mismatch case (assert the failure path and message), and the quirk from Step 4.
The cluster-dependent behaviour is proven by an e2e test that uses the method.

## Step 7: verify and document

```bash
npm run build
npm run typecheck && npm run typecheck:e2e && npm run typecheck:examples
npm test
npm run check:lanes
npm run e2e -- --grep @GKO-NNNN     # a test that actually calls the new method
```

Then:

- Add the method to the table in [`README.md`](../../../README.md).
- Add it to the
  [assertions cheatsheet](../write-e2e-test/references/assertions-cheatsheet.md)
  so test authors find it.
- Apache 2.0 header on every new file.

## Checklist

- [ ] Family checked for an existing method before adding one
- [ ] Naming triplet respected; `waitForX` added if the value converges
- [ ] void + throw via `deepPartialMatch` + `throwIfFailed`
- [ ] Type models only what is asserted; literal unions over `string`
- [ ] Endpoint quirks encoded in the method and documented
- [ ] `exports` map + `src/index.ts` updated for a new import path
- [ ] Vitest coverage for pass, mismatch, and the quirk
- [ ] Apache 2.0 header on new files
- [ ] Full verification gate green, including one real e2e test using the method
- [ ] `README.md` table and the cheatsheet updated
