---
name: release
description: Release exactly one GKO version end-to-end - Jira Done/no-open-PR gate, GitHub Trigger Release (dry-run off, latest only if asked), wait until the GitHub tag and gravitee-platform-docs changelog PR exist (CircleCI canceled after those steps is still success), mark the Jira version released, and create the next patch with no dates. On success tell the user only that it worked plus the docs PR URL. Use when asked to release one GKO version, trigger a GKO release, run Trigger Release, or close a GKO Jira version after CI. Do not chain multiple versions.
user-invocable: true
---

# release

The steps for this skill live in `.agent/skills/release.md`, relative to the repository
root. Read that file now and follow it; this wrapper carries no instructions of its own.

@.agent/skills/release.md
