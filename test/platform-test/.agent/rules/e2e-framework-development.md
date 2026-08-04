# Rules: extending the e2e framework

Audience: any agent (SDET) **changing the framework itself**: `src/` (the
`@gravitee/platform-test` library), `e2e/helpers/`, `e2e/setup.ts`,
`e2e/playwright*.config.ts`, `scripts/`, or the package surface.

Adding or changing a *test* is a different job: see
[`e2e-test-authoring.md`](./e2e-test-authoring.md).

Procedures live in the `add-provisioner` and `extend-platform-assertions`
skills. These are the invariants they rest on.

---

## Layering

```
src/            @gravitee/platform-test - runner-agnostic, publishable, no Playwright
e2e/            the Playwright binding - config, fixtures, helpers, tests
```

- **`src/` must not import from `e2e/`.** It is published as a standalone
  package and consumed by `runner-examples/jest` as well as by this suite.
- The two things `src/` cannot own are supplied by the e2e binding in
  [`e2e/helpers/provisioner-env.ts`](../../e2e/helpers/provisioner-env.ts): the
  APIM auth/server environment loaded from `config.yaml`, and fixture-path
  resolution via `fixture()`. Scenario authors pass fixture-relative paths;
  that module resolves them to the absolute paths the `src/` constructors expect.
- [`e2e/helpers/kubectl.ts`](../../e2e/helpers/kubectl.ts) and
  [`e2e/helpers/terraform.ts`](../../e2e/helpers/terraform.ts) are thin shims
  over `src/provisioners/engines/`. Put logic in the engine, not the shim.

## Library design principles

From [`README.md`](../../README.md#design-principles); do not erode them:

| Principle | Implementation |
|---|---|
| **void + throw** | Success = `void`, failure = a `node:assert` `AssertionError` |
| **Partial matching** | Only assert fields the caller specifies; ignore the rest |
| **Minimal deps** | Native `fetch`, `node:assert`, `yaml` for config only |
| **Runner-agnostic** | Works with Vitest, Jest, `node:test`, Playwright |
| **Extensible** | New product surfaces (AM, AE, Cockpit) follow the same shape |

## Assertion naming triplet

[`src/assertions/apim/mapi.ts`](../../src/assertions/apim/mapi.ts) uses one
naming contract. Match it when adding a method:

| Prefix | Behaviour | Signature shape |
|---|---|---|
| `assertX(...)` | single shot; throws `AssertionError` on mismatch | `Promise<void>` |
| `waitForX(...)` | polled variant of the same assertion | takes a trailing `PollOptions` |
| `checkX(...)` | non-throwing; returns the report for the caller to inspect | `Promise<AssertionReport>` |
| `listX(...)` / `createX(...)` / `deleteX(...)` | plain data access used by tests to build assertions | returns data |

Not every family needs all three. Add `waitForX` whenever the value converges
asynchronously, which is almost always true for anything written through the
Automation API.

## Provisioner contracts

- **[`registry.ts`](../../src/provisioners/registry.ts) is the single source of
  truth** for the set of provisioners. The `ProvisionerId` type,
  `forEachProvisioner`'s generation order, the `--provision-with` flag, and the
  Playwright lane `grepInvert` all derive from it. Never hardcode a provisioner
  list anywhere else.
- **`ProvisionerView.read()` has exactly three outcomes**: `applied`, `failed`,
  `gone`. There is deliberately no `unknown`. A provisioner that cannot yet tell
  what happened **must throw**, so a half-implemented `read()` fails loudly
  instead of passing silently. See
  [`view.ts`](../../src/provisioners/view.ts).
- **`view` vs `checks`:** anything expressible as "did MY layer land this role?"
  belongs in `view`, which shared scenario bodies call without narrowing.
  `checks` is the narrowed escape hatch (`isTerraform()`) and stays rare; a
  provisioner with none must not declare the property at all, so shared bodies
  cannot reach for it by habit. See
  [`types.ts`](../../src/provisioners/types.ts).
- **Id resolution is defined once.** `BaseProvisioned` implements `apiId()`,
  `applicationId()`, `subscriptionId()`, `groupId()` on top of a single
  `resolveId(role)` that each provisioner implements. Adding a kind is one line
  in [`base.ts`](../../src/provisioners/base.ts), not a new method per
  provisioner.
- **`destroy()` and `cleanup()` must never throw.** They run in `finally` and in
  the `afterEach` safety net; a teardown failure must not mask a test result.

## Lanes

Lane selection is **by title tag only**; nothing keys off the folder a test
lives in, which is what lets the tree be reorganised freely.

- [`scripts/check-lanes.mjs`](../../scripts/check-lanes.mjs) asserts the lanes
  **partition** the suite exactly: no test in two lanes, none in zero. Run it
  after any change to tags, lane logic, or the test tree.
- It imports from `dist/`, so **`npm run build` must run before
  `npm run check:lanes`**.
- The config's `grepInvert` is **case-sensitive** on purpose: Playwright's CLI
  `--grep` is case-insensitive, so a bare `@gko` would also match every
  `@GKO-NNNN` Xray tag.

## Package surface

- A new importable path needs an entry in the `exports` map of
  [`package.json`](../../package.json), and a matching `index.ts` re-export.
- Keep `dependencies` minimal. A new runtime dependency needs a reason that
  survives "can native `fetch` do this?".
- **Apache 2.0 licence header on every new `.ts` / `.mjs` file.** Enforced by CI
  (`job-lint-licenses`); copy the header from any existing source file.

## Docs contract

Each document has one job. Do not cross-write; PARITY.md records that drift
happened before precisely because nothing said what each file was for.

| File | Answers |
|---|---|
| [`README.md`](../../README.md) | the library API: what `@gravitee/platform-test` exposes |
| [`e2e/README.md`](../../e2e/README.md) | bootstrap and how to run the suite |
| [`AGENTS.md`](../../AGENTS.md) | the authoring model and where a test goes |
| [`PARITY.md`](../../PARITY.md) | coverage status: GKO to Terraform parity and the backlog |
| [`e2e/tests/user-journeys/README.md`](../../e2e/tests/user-journeys/README.md) | the journey catalog |
| `.agent/rules/*.md` | the invariants, per audience |
| `.claude/skills/*/SKILL.md` | the procedures |

## Verification gate

Any framework change runs the whole gate, not just the part you touched:

```bash
npm run build
npm run typecheck && npm run typecheck:e2e && npm run typecheck:examples
npm test                 # vitest unit tests, no cluster needed
npm run check:lanes      # needs the build above
npm run e2e -- --provision-with gko    # at least one real lane
```

`build`, `typecheck*`, `test` and `check:lanes` need no cluster and should be
run first. Never report a framework change as done on the cluster-free checks
alone: the layer exists to talk to a live platform.
