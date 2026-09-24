# Skill: Release GKO

## When to Use

Use this skill when asked to release **one** GKO version, trigger a GKO release, run
"Trigger Release", or close a GKO Jira version after CI.

This skill releases **exactly one version**. If the user names more than one, ask which
single version to run and do not start until they pick one.

Do **not** dispatch GitHub or mutate Jira until the gate in step 1 passes.

## Constants

- Jira project: `GKO` on `gravitee.atlassian.net`
- GitHub repo: `gravitee-io/gravitee-kubernetes-operator`
- Docs repo: `gravitee-io/gravitee-platform-docs`
- Workflow name: `Trigger Release`
- CircleCI API: `https://circleci.com/api/v2`
- Tools: `twg` + `gh`. Do not use GitHub MCP.

## Hard rules

1. Release exactly one version, then stop. Do not chain a second version in the same run.
2. Real release: `dry-run=false`. Never leave dry-run at its default `true`.
3. `latest=false` unless the user explicitly ticks latest for this version.
4. If the Jira gate fails, **stop**. Do not dispatch.
5. Success is **artifacts**, not CircleCI job status:
   - GitHub release tag `<version>` exists
   - docs PR `[GKO] Changelog for version <version>` exists (open or merged)
6. CircleCI `Release` may end `canceled` after those steps (later steps such as OLM).
   Treat `canceled` as success when both artifacts exist. Treat it as failure only when
   they do not.
7. Mark Jira as **released** (the Released checkbox). Do **not** generate Jira release notes.
8. Next Jira version: patch + 1, **no dates**, not released. Do this even if CircleCI
   ended `canceled`, as long as artifacts exist.
9. User-facing status: short. On success, only "it worked" + the docs PR URL. On failure,
   give details. Do not dump CircleCI, GitHub release, or Jira ids on success.

## Inputs

From the user:

- **version** (required), e.g. `4.9.33`
- **latest** (optional, default false)

Next patch: bump the last number (`4.9.33` → `4.9.34`).

## Steps

### 1. Jira gate

Resolve the version:

```bash
twg jira space versions query --key GKO --search-string "<version>" --limit 5
```

Use the row whose `name` equals `<version>` exactly. Record `id`, `released`, `archived`.

If missing, archived, or already `released: true`: **stop** and report failure.

List bound tickets:

```bash
twg jira workitem query --jql "project = GKO AND fixVersion = '<version>' ORDER BY key"
```

Empty ticket list is allowed (sync release). Non-empty: every ticket `status.statusCategory.key`
must be `done`.

For each ticket key, search open PRs:

```bash
gh search prs --repo gravitee-io/gravitee-kubernetes-operator "<KEY>" --state open
```

Any open PR: **stop**. Report key, status, PR URL.

Gate pass: continue.

### 2. Dispatch Trigger Release

```bash
gh workflow run "Trigger Release" -R gravitee-io/gravitee-kubernetes-operator \
  -f version=<version> -f dry-run=false -f latest=<true|false>
```

Find the run and watch it:

```bash
gh run list -R gravitee-io/gravitee-kubernetes-operator --workflow "Trigger Release" --limit 3
gh run watch <run-id> -R gravitee-io/gravitee-kubernetes-operator
```

GH job is short (~20s). If it fails: **stop**. Dump the failing step.

### 3. Parse CircleCI URL

```bash
gh run view <run-id> -R gravitee-io/gravitee-kubernetes-operator --log
```

Look for:

```
Pipeline is running at https://app.circleci.com/pipelines/github/gravitee-io/gravitee-kubernetes-operator/<N>
```

If missing: **stop**.

### 4. Wait until the release exists (or CircleCI fails for real)

Poll every 20–30s (typical ~4–5 minutes). On each tick:

1. CircleCI jobs:

```bash
curl -sS "https://circleci.com/api/v2/project/github/gravitee-io/gravitee-kubernetes-operator/pipeline/<N>"
curl -sS "https://circleci.com/api/v2/pipeline/<uuid>/workflow"
curl -sS "https://circleci.com/api/v2/workflow/<workflow-id>/job"
```

2. Artifacts:

```bash
gh release view <version> -R gravitee-io/gravitee-kubernetes-operator
gh pr list -R gravitee-io/gravitee-platform-docs \
  --search "[GKO] Changelog for version <version>" --state all --limit 5
```

**Proceed to step 5** as soon as the GitHub release exists. Do not wait for CircleCI
`status: success`. The `Release` job often keeps running (OLM) and can be marked
`canceled` after the tag and docs PR are already published.

**Fail** only if the CircleCI `Release` job is `failed` **and** the GitHub release
does not exist. Then report CircleCI URL, job name, status.

If CircleCI is `canceled` and the GitHub release does not exist yet, wait ~1 minute
and check artifacts once more before failing.

### 5. Docs PR

Expected title: `[GKO] Changelog for version <version>`.
Expected branch: `release-gko-<version>`.

```bash
gh pr list -R gravitee-io/gravitee-platform-docs \
  --search "[GKO] Changelog for version <version>" --state all --limit 5
```

Prefer `--json url,title,state,number` (or `gh pr view`) so you have the URL.

If the GitHub release exists but the docs PR is missing, wait ~1 minute and retry once.
If still missing: **stop** as failure (release shipped, docs PR did not). Do still close
Jira in step 6.

### 6. Close Jira version and open the next patch

Do this when the GitHub release exists, including when CircleCI ended `canceled`.

```bash
twg jira space version update --id <version-id> --released true
twg jira space version create --key GKO --name <next-patch>
```

No dates. No release notes. If `<next-patch>` already exists, do not create a duplicate;
record the existing id.

### 7. Report

**Success** (GitHub release + docs PR):

```
Released <version>
Docs PR: <url>
```

Nothing else.

**Failure**:

```
Failed <version>
<one-line reason>
<details as needed: job status, missing artifact, open PRs, Jira gate>
```
