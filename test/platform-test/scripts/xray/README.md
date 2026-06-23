# test/platform-test/scripts/xray/

Wiring for **pushing CircleCI e2e results into Xray Cloud** so the
"Test Coverage" panel on each linked Jira Test (and on the parent Story it
`is tested by`) reflects the latest run instead of staying at `TO DO`.

This is the **write** counterpart to the read-side `xray-export` Claude Code
skill — they don't overlap.

## How it fits together

1. The Playwright e2e suite (`test/platform-test/`) is already configured to
   emit a JUnit XML report at `test/platform-test/playwright-results/results.xml`
   (see `playwright.config.ts`, `["junit", { outputFile: ... }]`).
2. Each test title already carries its Xray Test issue key as a tag, e.g.
   `terraform apply creates the group in APIM @GKO-2865 @regression`, via
   the `XRAY` constants in [`e2e/helpers/tags.ts`](../../e2e/helpers/tags.ts).
3. After the e2e job runs, the CircleCI step `Push results to Xray Cloud`
   invokes [`push-results.sh`](./push-results.sh). It authenticates to Xray
   Cloud, runs [`junit-to-xray.mjs`](./junit-to-xray.mjs) to convert the JUnit
   report into Xray's **JSON results format** keyed by the `@GKO-NNNN` id in
   each title, and POSTs it to `/api/v2/import/execution`.
4. Xray Cloud creates a new **Test Execution** in Jira and, because each result
   carries a `testKey`, attaches it to the *existing* Test issue — updating that
   Test's "Test Coverage" panel with `PASSED` / `FAILED`.
5. If a `XRAY_TEST_PLAN_KEY` is set, the Test Execution is also linked to
   that Test Plan, so the Plan view shows trend over time.

> **Why JSON, not raw JUnit?** Xray's JUnit importer matches Tests by
> `classname`+`name` (its "Generic Test Definition"), **not** by a Jira key in
> the title — so a raw JUnit import creates a *duplicate* Test issue per
> testcase instead of updating the existing one. The JSON format's `testKey`
> is the only thing that binds a result to a pre-existing Test.

Skipped / `test.fixme`'d cases are **omitted** from the payload so a
temporarily-disabled test doesn't overwrite a Test's last real result with
`TODO`.

The script is **best-effort** — auth failures, credential gaps, missing
`jq` / `node`, missing results file, network errors all log a warning and
exit 0. The pipeline's pass/fail gate is unaffected (CircleCI parses
`store_test_results` independently).

## Environment variables

| Variable | Required? | Default | Purpose |
|----------|-----------|---------|---------|
| `XRAY_CLIENT_ID` | yes | — | Xray Cloud API key client id |
| `XRAY_CLIENT_SECRET` | yes | — | Xray Cloud API key client secret |
| `XRAY_TEST_PLAN_KEY` | no | — | Existing Test Plan issue (e.g. `GKO-NNNN`) to roll runs up to |
| `CIRCLE_BUILD_URL` / `CIRCLE_BUILD_NUM` / `CIRCLE_BRANCH` | no | autodetected | Used to template the Test Execution summary + description; CircleCI sets them automatically |

### Test Execution title

The script defers the title to [`junit-to-xray.mjs`](./junit-to-xray.mjs) — it
has the per-test totals naturally. The default template is:

```
GKO Playwright e2e on <branch> — <n> passed, <n> failed[, <n> skipped] (CircleCI #<num>)
```

Examples:

- `GKO Playwright e2e on master — 355 passed, 0 failed (CircleCI #85230)`
- `GKO Playwright e2e on master — 354 passed, 1 failed, 12 skipped (CircleCI #85231)`
- `GKO Playwright e2e on master — 1 passed, 0 failed (local)` *(local dry-run)*

Pass an explicit `XRAY_SUMMARY` env var to override.

### Test Execution description

Deliberately small — the Xray-generated "Tests" panel inside the issue body
already shows totals and the per-test rundown, so duplicating that here is
noise. The default is:

```
Commit: <git log -1 --pretty=%s>
Automated import from CircleCI: <build URL>
```

The `Commit:` line is omitted when `XRAY_COMMIT_SUBJECT` isn't set
(e.g. local dry-runs in a worktree without a HEAD), and the second line
falls back to a generic "Automated import of Playwright e2e results."
when no `XRAY_BUILD_URL` is available. Pass an explicit `XRAY_DESCRIPTION`
to override entirely.

> **Test Execution workflow state.** The created issue stays in its default
> "To Do" state — that's intentional. The per-test results (PASSED / FAILED
> rows on each linked Test's "Test Coverage" panel) are what consumers
> actually look at; the Test Execution issue is a per-run log entry whose
> workflow column nobody reads. Don't add a Jira API token just to flip the
> badge.

The Test Execution is created in the project of the referenced Tests (GKO),
so no project key is needed.

In CI the two `XRAY_CLIENT_*` vars are injected by the `keeper/env-export`
orb from the Keeper record `so8_Jh2tP-AZSbtIYcbBvg` (custom fields
`CLIENT ID` / `CLIENT SECRET`) — see the two `keeper/env-export` steps in
`job-e2e-tests`, placed *before* the e2e run so they still run while the
pipeline is green (the orb command runs `on_success`). Generate a fresh
key pair from Jira if needed: *Apps → Xray → Global Settings → API Keys →
Create*.


## Local run

After a local e2e run:

```sh
export XRAY_CLIENT_ID=...
export XRAY_CLIENT_SECRET=...
# optional: export XRAY_TEST_PLAN_KEY=GKO-NNNN

./test/platform-test/scripts/xray/push-results.sh
# → [xray-push] reporting 6 test result(s) to Xray Cloud
# → [xray-push] created Test Execution GKO-NNNN
```

Spot-check by opening the created Test Execution in Jira, and any of the
linked Tests (e.g. [GKO-2865](https://gravitee.atlassian.net/browse/GKO-2865))
— the "Test Coverage" panel should now show the run with the right status.

You can pass a custom results path as the first argument if you've moved
the JUnit file:

```sh
./test/platform-test/scripts/xray/push-results.sh /path/to/results.xml
```

## CI invocation

The script is invoked from `.circleci/config.yml` inside `job-e2e-tests`,
right after the existing `store_artifacts` steps. The CI step uses
`when: always` so failed runs still report their failures back to Xray.

The credentials are loaded earlier in the job by two `keeper/env-export`
steps placed **before** the `Run Playwright E2E tests` step. That ordering
matters: `keeper/env-export` runs `on_success`, so if it were placed next to
the `when: always` push step and the e2e run failed, the export would be
skipped and the push would run without credentials.

## What this script intentionally doesn't do

- **No Playwright reporter swap.** A community `playwright-xray-cloud-reporter`
  exists but adds a runtime dep and pushes per-test rather than per-suite.
  Keeping the push as a post-step keeps the test container clean.
- **No Cucumber import.** The Gherkin scenarios live in the Jira Test
  descriptions, not in `.feature` files. The JUnit report (transformed to
  Xray JSON) is sufficient for per-Test status.
- **No PR comment / Slack callout.** The Test Execution key is written to
  `/tmp/xray-test-execution.txt` and stored as a CircleCI artifact —
  follow-up automation can pick it up from there.
