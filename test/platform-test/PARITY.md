# Coverage parity: GKO ↔ Terraform

**What this document is for:** the *status* view of the e2e suite — what the
Terraform provider can express, which areas are covered by a shared user journey,
and what is left to migrate. It answers "where do we stand?".

It is **not** the authoring guide. How to write a test (persona folders, the
variant-matrix rule, cleanup patterns, `view` vs `checks`) lives in
[AGENTS.md](./AGENTS.md) and the
[`write-e2e-test`](./.claude/skills/write-e2e-test/SKILL.md) skill; the journey
list lives in the [catalog](./e2e/tests/user-journeys/README.md). Keep those
boundaries — this file drifted before precisely because nothing said what it was
for.

> Regenerate the provider tables below from source rather than editing by hand;
> see [Regenerating this document](#regenerating-this-document).

---

## What the Terraform provider can express

The provider (`gravitee-io/apim`) is **Speakeasy-generated** from
`automation-api-oas.yaml` — APIM's Automation API OpenAPI spec — see
`.speakeasy/gen.yaml` in the provider repo. Its resources map **1:1 onto the
Automation API paths**: 10 resources, 10 paths.

That has a consequence worth stating plainly, because it was previously
misread as a hard limit:

> **"No Terraform resource" means "the Automation API has no path for it yet."**
> GKO writes through the *same* Automation API, so a gap is a Gravitee backlog
> item (add the path → regenerate the provider), not an external constraint. File
> it rather than recording the area as permanently GKO-only.
>
> The one genuinely permanent exclusion is **V2**: the Automation API is V4-only
> by design and V2 is legacy.

### Resources (verified on `origin/main`)

`apim_apiv4` · `apim_application` · `apim_subscription` · `apim_group` ·
`apim_dictionary` · `apim_shared_policy_group` · `apim_documentation_api` ·
`apim_documentation_portal` · `apim_portal` · `apim_portal_listing`

Each also has a matching data source.

### Inline attributes (where most parity actually lives)

A missing standalone resource does not mean a missing capability — `apim_apiv4`
in particular carries a rich inline surface:

| Resource | Attributes |
|---|---|
| `apim_apiv4` | `type` `state` `visibility` `lifecycle_state` `version` `description` `listeners` `endpoint_groups` `plans` `flows` `flow_execution` `labels` `categories` `groups` `members` `notify_members` `metadata` `pages` `properties` `resources` `response_templates` `services` `tags` `analytics` `failover` `console_notification` `portal_navigation` `allow_multi_jwt_oauth2_subscriptions` `allowed_in_api_products` |
| `apim_application` | `settings{app, oauth, tls.client_certificate}` `members` `notify_members` `groups` `metadata` `name` `description` `domain` `picture_url` `background` |
| `apim_subscription` | `api_hrid` `application_hrid` `plan_hrid` `api_keys` `consumer_configuration` `starting_at` `ending_at` `metadata` |
| `apim_group` | `name` `members` `notify_members` |
| `apim_dictionary` | `type` `dynamic` `manual` `deployed` `description` |
| `apim_shared_policy_group` | `api_type` `phase` `steps` `prerequisite_message` `description` |
| `apim_documentation_api` | `api_hrid` `name` `type` `content` `location` `order` |
| `apim_portal` / `apim_portal_listing` | `name` `navigation` / `apis` |

---

## Where coverage stands

`tests/gko/` holds **230** tests, and the buckets below partition exactly that
number. The 24 journeys in the [catalog](./e2e/tests/user-journeys/README.md) are
a *separate* population: their GKO arms carry the ids already migrated out of
this folder, so they are not counted here.

| Bucket | Tests | Notes |
|---|--:|---|
| **Migratable, not yet migrated** | **27** | the backlog below, area by area |
| Blocked on a framework gap | ~12 | APIM-side rejections — see [Rejection tests](#rejection-tests-a-framework-gap-not-a-provider-gap) |
| Genuinely GKO-only | 191 | operator mechanics; 64 of them V2 |

### How to read these numbers

**A falling test count is this work succeeding, not regressing.** Migration
merges duplicates: two single-driver tests of one behaviour become one journey
run against both drivers, so the count drops while the evidence grows. Before
consolidation the suite held copies that asserted *different* things about the
same behaviour (group rename checked the id on GKO and the name on Terraform),
which counted as two and covered neither properly.

Four numbers that move the right way while the count falls:

| Metric | Why it survives dedup |
|---|---|
| Behaviours covered (journeys × stages) | merging two copies leaves the behaviour count unchanged, not halved |
| Driver coverage (arms per journey) | the number that makes a missing Terraform arm visible |
| Assertion depth (fields asserted per test) | catches the `Accepted=True`-only tests listed [below](#migrate-these-ones-carefully--they-do-not-assert-their-premise), which inflate the count while proving nothing |
| Xray Test count | the external repository; it only grows |

So do not read `tests/gko` shrinking as lost coverage. Read it against the
[catalog](./e2e/tests/user-journeys/README.md) growing.

### Migrated

Each row is a journey now running against **both** provisioners. The GKO ids are
carried on the journey's GKO arm and the source tests are deleted; the split rule
(admission/rejection cases move to the matching `tests/gko/` domain folder rather
than into the journey) applied throughout.

| Area | Journey [persona] | Notes |
|---|---|---|
| V4 API visibility / lifecycle / delete | `configure-visibility-and-lifecycle` [producer] | the three portal combinations, in place, then retire |
| V4 API failover | `configure-endpoint-failover` [producer] | every field of the nested block, then a retry-budget change |
| Message-API entrypoints | variants inside `publish-a-message-api` [producer] | http-get / http-post / sse / websocket, each with its own version + description |
| V4 API members | `manage-api-members` [admin] | whole member set per stage, so a duplicated entry fails |
| API ↔ group association | `associate-groups-with-an-api` [admin] | group provisioned first; APIM drops unknown refs silently |
| Group members | folded into `create-group-with-member` [admin] | the member is now a resolvable service account and is asserted |
| Application members | `manage-application-members` [consumer] | as above, on `apim_application.members` |
| Application settings + metadata | folded into `register-and-retire-application` [consumer] | metadata read from its own endpoint; the detail response omits it |
| API notifications (V4) | `configure-api-notifications` [producer] | GKO's Notification CR vs TF's inline block, same PORTAL setting |
| Pages — V4 inline | `document-an-api` [producer] | ship, revise (rename + rewrite + visibility), remove |
| Pages — V4 inline fetchers | `sync-documentation-from-a-url` [producer] | `pages[].source`: the declared `http-fetcher` **and** the content APIM pulls from the URL, which a rename must not discard. `github-fetcher` stays GKO-only — admission pre-fetches the repository at apply time and the cluster has no GitHub credentials |
| Policies on a flow | `apply-policies-to-a-flow` [producer] | asserted at the **gateway**, not only in the definition |
| Subscriptions — JWT / OAuth2 | `subscribe-to-a-secured-plan` [consumer] | auto-validated despite `MANUAL` plans |
| mTLS plan + client certificates | `authenticate-with-client-certificate` [consumer] | issue, rotate, retire, revoke — asserted at the **gateway**; also covers the deprecated single-certificate field |
| Analytics: OTel logs + tracing, and the native reporter | `configure-api-observability` [producer] | two scenarios, because `tracing` is not for native APIs and `reporterMetricsEnabled` is native-only. The reporter is asserted **off** first: the provider defaults it to `true`, so `false` is what proves the declaration was transmitted |
| API consumption flags (`allowMultiJwtOauth2Subscriptions`, `allowedInApiProducts`) | folded into `subscribe-to-a-secured-plan` [consumer] | both arms already declared the first flag to make the journey legal and never asserted it. `allowedInApiProducts` can only be asserted as a round-trip: APIM exposes no product surface through either driver |
| Shared Policy Group create + update | `update-a-shared-policy-group` [producer] | the whole group per stage — step configuration, description, `lifecycleState` — where the source test read only the surviving `.status.id`. Id stability falls out of polling the original id for the new content |

Eight of these journeys run a Terraform arm but still show **TF TBD** in the
catalog. The coverage exists; only the Jira Test ids are missing, so Xray
under-reports Terraform. Fix with `/xray-sync-tests`.

### Migratable, not yet migrated

Every Terraform path below was verified against `origin/main` of the provider —
see [Regenerating this document](#regenerating-this-document). Ordered by batch:
rows 1 and 9 need a new journey folder each, rows 2-8 are variant-table additions
to journeys that already exist, 10-11 are the residue.

| # | Area | Tests | Terraform path | Journey [persona] |
|--:|---|--:|---|---|
| 1 | Export an API definition whatever created it | 1 | none needed — `mapi.exportApiCrd(id)` takes any APIM id | **new** `export-an-api-definition` [producer] — GKO-3116 |
| 2 | Category *references*: many at once, changed set redeploys with a stable id, unknown ref dropped, none declared → none in APIM | 6 | `apim_apiv4.categories` | variants in `assign-categories-to-api` |
| 3 | Application platform semantics: long name preserved, optional `client_id`, ARCHIVED on removal, `client_id` uniqueness | 4 | `apim_application.settings.app` | variants in `register-and-retire-application` |
| 4 | Message-API MQTT endpoint, entrypoint × policy matrix, message API with a policy | 3 | `endpoint_groups`, `flows` | variants in `publish-a-message-api` / `apply-policies-to-a-flow` |
| 5 | Unknown group reference is tolerated and dropped | 2 | `apim_apiv4.groups` | variant in `associate-groups-with-an-api` |
| 6 | Re-create an API after deleting it — previously closed plan reopens | 1 | `destroy` then re-apply | variant in `publish-api-and-serve-traffic` |
| 7 | Documentation **folder** + page rename, no duplicates left behind | 1 | `pages[].type = "FOLDER"` | variant in `document-an-api` |
| 8 | Primary owner is the identity the automation authenticated as | 2 | implicit (provider credentials) | fold into `manage-api-members` |
| 9 | `origin: KUBERNETES` marks a resource read-only in the console (API, application, notification settings) | 3 | any resource | **new** `automation-managed-resources-are-read-only` [admin] |
| 10 | Plan lifecycle (publish/close) | 1 | `apim_apiv4.plans` | `manage-plan-lifecycle` [producer] |
| 11 | Subscription slices that duplicate a journey (V4 JWT gateway call, mTLS plan) or never assert their premise (delete API with another plan) | 3 | — | retire against `subscribe-to-a-secured-plan` / `authenticate-with-client-certificate` rather than migrate |
| — | API metadata | 0 | `apim_apiv4.metadata` | `manage-api-metadata` [admin] — **new coverage**, not a migration: nothing covers V4 API metadata today |

Two readings of this document produced the old, much shorter backlog. Both were
wrong:

- **"Category CRUD has no Automation API path."** True for *creating* a category,
  which is why that row stays GKO-only. Every test in row 2 only *references* an
  existing one, and `apim_apiv4.categories` expresses that.
- **"Shared Policy Groups are blocked."** Only *reuse at the gateway* is (the
  `crossId` gap below). Creating and updating an SPG was untouched by it, and is
  now migrated as `update-a-shared-policy-group`.

Row 9 is confirmed, not assumed: a Terraform-created `apim_apiv4` reports
`originContext.origin: "KUBERNETES"`, measured on the cluster while building
`configure-api-observability`. Both drivers write through the Automation API, so
the marker is about automation rather than about Kubernetes.

### Rejection tests: a framework gap, not a provider gap

[`forEachProvisioner`](./e2e/helpers/for-each-provisioner.ts) always expects
`provision()` to succeed, so every negative test is GKO-only *by construction*
rather than by analysis. They are not all alike:

| Rejected by | Verdict |
|---|---|
| the CRD schema or GKO's own webhook logic | genuinely GKO-only |
| APIM's dry-run through the Automation API | `terraform apply` gets the identical error — shared behaviour, zero Terraform coverage today |

The second row is ~12 tests: a plan whose general conditions reference a missing
page, native-Kafka port ranges, duplicate certificate fingerprints, application
OAuth grant rules. Unlocking them needs an `expectFailure` scenario shape
(provision expected to throw, with a per-arm message matcher). That is worth more
than the tests it unblocks: it is the only thing that would prove the provider
*surfaces* APIM's validation instead of swallowing it.

### Migrate these ones carefully — they do not assert their premise

Found while auditing the backlog above. Each is inside a migratable area, so the
migration is also the fix; porting the assertion as-is would carry the hole into
the journey.

| Test | What it actually does |
|---|---|
| `DELETE_API_WITH_OTHER_PLAN` | never creates the subscription — the "despite an active subscription" premise is untested |
| `APP_CLIENT_ID_UNIQUE` | never applies a second application, so uniqueness is untested. Two `apim_application` resources express it directly |
| `GROUPS.CREATE_NON_EXISTING_USER` | `try/catch` both outcomes, so it cannot fail. `SPG_LIFECYCLE` was the same and has been deleted rather than ported |
| `PREVENT_PO_GROUP_AS_MEMBER` | asserts only that a group exists |
| ~10 others | stop at `Accepted=True` and never read APIM |

### Not expressible through Terraform (found while migrating)

| Case | Why |
|---|---|
| A member declared with **no role** (API or application) | the Automation API's `Member` schema marks `role` non-nullable, so the generated provider rejects omitting it. GKO can omit it and APIM defaults to `USER` — kept as a GKO-only test |
| Turning a console notification **off** | GKO-3085: the Automation API answers HTTP 500 (`organizationId must not be empty`) to an empty `events` list, and `console_notification = null` is a Terraform no-op. GKO removes the reference; kept as `@GKO-1238` |

Blocked: **`reuse-shared-policy-group`** — GKO-3001 (admission rejects the
documented `sharedPolicyGroupRef` flow form because it does not resolve the ref
before APIM's dry-run) *and* Terraform (`apim_shared_policy_group` exposes no
`cross_id`, and only the crossId executes the SPG at the gateway). The journey
documents the correct form as a `pending` fixme rather than being green-washed.
This blocks **reuse only** — creating and updating an SPG runs against both
drivers in `update-a-shared-policy-group`.

### Genuinely GKO-only

These have the **operator itself** as the system under test, so there is nothing
for a second provisioner to do. Counted by folder, so the rows partition the 192
exactly; V2 is a cross-cutting dimension, counted separately below.

| Folder | Tests | Why |
|---|--:|---|
| `admission-webhook/` | 28 | CRD schema / dry-run validation at the Kubernetes admission layer |
| `members/` | 27 | member and group resolution failures surfaced as reconciliation errors, plus V2 |
| `mtls-certificates/` | 24 | admission rejections; certificates resolved from cluster Secrets/ConfigMaps and `[[ … ]]` templating; base64-encoded content (the provider's `content` is PEM-only); certificate date-window and fingerprint-reuse semantics |
| `pages/` | 19 | V2 documentation and both versions' fetcher validation |
| `subscriptions/` | 19 | admission rules (immutability, plan matching, `syncFrom`) and V2 |
| `deployment-reconciliation/` | 13 | CR `.status` conditions, observedGeneration, operator restart, audit events |
| `api-lifecycle/` | 13 | V2 lifecycle, DB-less mode, reconcile-cycle identity |
| `applications/` | 13 | admission rejections and settings the test cluster cannot exercise (DCR is off) |
| `import-export/` | 9 | YAML CRD round-trips |
| `templating/` | 8 | `[[ … ]]` resolution from cluster ConfigMaps/Secrets |
| `categories/` | 7 | category CRUD and V2 |
| `management-context/` | 5 | Kubernetes custom-resource lifecycle |
| `notifications/` | 5 | Notification CR wiring and export |
| `defaults/` | 4 | CRD field defaulting |
| `dictionaries/`, `groups/` | 6 | admission rejections and CRD defaulting |
| `local-configmap/` | 2 | in-cluster ConfigMap locality |
| `policies/` | 4 | plan security types and plan lifecycle driven from the CR |

Two things have **no Automation API path at all**, so they stay here permanently
rather than becoming backlog items:

| Area | Tests | Why |
|---|--:|---|
| V2 — spread across the folders above | 64 | the Automation API is V4-only; V2 is legacy |
| Category CRUD (create/rename a category) | ~6 | no `/categories` path; an API can only *reference* categories inline (referencing is backlog row 2) |

### Terraform-only

| Behaviour | Where | Why |
|---|---|---|
| Drift detection, import, data sources, plan exit codes | `tests/terraform/groups.test.ts` | HCL desired-state vs server reconciliation; `terraform plan -detailed-exitcode` and `terraform import` semantics |
| Provider-side schema validation | `tests/terraform/groups.test.ts` | rejected before any API call |
| Sensitive redaction | `subscribe-and-call/apikey-tf-only.test.ts` | provider `Sensitive` attribute handling |
| Portal management | *not yet covered* | `apim_portal`, `apim_portal_listing`, `apim_documentation_portal` have no GKO CRD — the first TF-only *product* surface |

---

## Regenerating this document

The resource and attribute tables come from the provider source. The local
checkout drifts, so always read `origin/main`:

```sh
cd ~/dev/src/terraform-provider-apim && git fetch origin

# Resources + data sources
git show origin/main:internal/provider/provider.go | sed -n '/func.*Resources(ctx/,/^}/p'

# Top-level attributes of one resource
git show origin/main:internal/provider/apiv4_resource.go |
  grep -oE '^\t\t\t"[a-z_0-9]+": ' | tr -d '\t":' | sort

# The Automation API paths the provider is generated from
git show origin/main:automation-api-oas.yaml | grep -E "^  /"
```

Equivalently, against an initialised workspace: `terraform providers schema -json`.
