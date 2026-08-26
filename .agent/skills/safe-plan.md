# safe-plan

Plan first. The plan file is the contract. Code comes after an
unambiguous execute.

## Modes

Pick one mode per user message. Do not mix.

| Mode | User intent | Agent does | Then |
|---|---|---|---|
| Draft | plan / design / explore | write `plans/<slug>.md` | stop |
| Challenge | questions a Case / Issue / decision | answer that question | **stop** |
| Revise | compromise agreed | snapshot `revN`, edit original, sync | stop |
| Execute | "implement the plan" / "execute" / "go ahead and implement" | follow the plan, do not edit it | — |

If Draft vs Challenge is unclear: **Challenge**.
If Challenge vs Execute is unclear: **Challenge**.

## Hard rules

1. Do not implement in Draft, Challenge, or Revise.
2. Challenge: answer. Do not edit the plan. Do not start the next Case. Stop.
3. Revise only after the user agrees the new shape (not when you agree with yourself).
4. Execute: do not edit the plan file.
5. Canonical path is `<repo>/plans/<slug>.md`. Not `~/.cursor/plans/`.
6. Every implementable chunk is a numbered Case the user can name.

## File names

- Live plan: `plans/<slug>.md` (kebab-case, no `.rev` in the live name).
- Snapshot before a Revise: `plans/<slug>.revN.md`
- `N` starts at `0`. Never overwrite a `revN` file.
- Next `N` = highest existing `rev*` + 1. If none, use `0`.

If more than one repo is in scope, the live file and the new `revN`
must be byte-identical in each repo's `plans/`.

`plans/` may be gitignored. Still write there. Do not commit plans
unless the user asks.

## Draft

1. List the repos in scope. If two, you will sync both.
2. Read enough code to make the plan actionable. Do not start the work.
3. Write `plans/<slug>.md` using the template below.
4. If two repos: write the same file to both `plans/` directories.
5. Show the user the Cases. Stop. Wait.

### Template

```markdown
# <title>

Repos: <repo-a>[, <repo-b>]
File: plans/<slug>.md

## Goal
<one short paragraph, the outcome>

## Case 01 — <short name>
<actionable steps: paths, symbols, snippets, commands>

## Case 02 — <short name>
...

## Out of scope
- <explicit non-goals>

## Done when
- <verification commands or gates>
```

### Actionable means

Each Case must be implementable without guessing:

- concrete paths (`internal/core/interface.go`)
- symbols to add/remove
- call sites
- commands that prove the Case (`go build`, `git diff -- <crd.yaml>`)

Forbidden in a Case: "consider", "we might", "clean up as needed".
Open questions are `## Issue 01 — ...`, not buried in a Case.

IDs are stable: `Case 01`, `Case 02`, `Issue 01`. The user will
challenge by ID (`Case 02`, `Issue 05`).

## Challenge

Triggers: "why", "should we", "I don't think", a Case/Issue id,
or any push-back on a decision.

1. Name the Case or Issue in the first line.
2. Answer the question. Be direct. Agree or disagree with a reason.
3. If the answer implies a plan change, say what would change.
   Do **not** write it yet.
4. Stop.

No Write/Edit/StrReplace on the plan in this mode.
No implementation in this mode.

Same message also says to update the plan? That is Revise, below.

## Revise

Triggers: user agrees the new shape ("do that", "update the plan",
"drop X", "yes, CredentialAuth only").

1. Copy live file → `plans/<slug>.revN.md` (next unused N).
2. Edit `plans/<slug>.md` to match the agreed shape.
3. Sync both files to every repo listed under `Repos:`.
4. Tell the user: what changed, new rev id, which Cases moved.
5. Stop. Do not implement unless the same message is Execute.

## Execute

Only when the user is unambiguous about implementing.

Follow the live plan. Stop and ask if reality does not match a Case.
Do not edit the plan file.

## Sync

Two-repo work (GKO + e2e, or any pair listed in `Repos:`):

- After every Draft write and every Revise, both `plans/` copies
  match.
- Same `<slug>`, same `revN`.
- If `plans/` is missing, create it.

One-repo work: no sync.
