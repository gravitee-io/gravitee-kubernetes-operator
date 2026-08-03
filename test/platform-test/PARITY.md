# Coverage parity: GKO ↔ Terraform

**What this document is for:** the *status* view of the e2e suite — what the
Terraform provider can express, which areas are covered by a shared user journey,
and what is left to migrate. It answers "where do we stand?".

It is **not** the authoring guide. How to write a test (persona folders, the
variant-matrix rule, cleanup patterns, `view` vs `checks`) lives in
[AGENTS.md](./AGENTS.md); the journey list lives in the
[catalog](./e2e/tests/user-journeys/README.md). Keep those boundaries — this file
drifted before precisely because nothing said what it was for.

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

Of the ~290 tests in `tests/gko/`:

| Bucket | Approx | Notes |
|---|--:|---|
| Covered by a shared journey | ~35 | the 13 journeys in the [catalog](./e2e/tests/user-journeys/README.md) |
| **Migratable, not yet migrated** | **~100** | see the backlog below |
| Genuinely GKO-only | ~155 | Kubernetes/operator mechanics + V2 |

### Migratable, not yet migrated

Each row is a journey to author (persona folder in brackets). Confirm the
attribute on `origin/main` before writing the Terraform arm.

| Area | Tests | Terraform path | Journey [persona] |
|---|--:|---|---|
| V4 API lifecycle / visibility | ~32 | `state` `visibility` `lifecycle_state` `failover` | `configure-visibility-and-lifecycle` [producer] |
| V4 API members | ~25 | `apim_apiv4.members` `notify_members` | `manage-api-members` [admin] |
| Subscriptions — JWT / OAuth2 / mTLS / manual | ~15 | `apim_subscription` | `subscribe-with-*`, `request-manual-approval` [consumer] |
| Application members + OAuth/app settings | ~12 | `apim_application.members` `settings.oauth` | `manage-application-members`, `configure-application-oauth-settings` [consumer/admin] |
| API notifications (V4) | ~11 | `apim_apiv4.console_notification` | `configure-api-notifications` [producer] |
| Pages — V4 inline + standalone | ~10 | `apim_apiv4.pages` + `apim_documentation_api` | `document-an-api` [producer] |
| Message-API entrypoints | ~8 | `apim_apiv4.listeners` | variants inside `publish-a-message-api` [producer] |
| Plans / policies lifecycle | ~8 | `apim_apiv4.plans` `flows` | `apply-policies-to-a-flow`, `manage-plan-lifecycle` [producer] |
| mTLS certs — inline content only | subset of 28 | `apim_application.settings.tls.client_certificate` | `authenticate-with-client-certificate` [consumer] |
| API metadata / group association | few | `apim_apiv4.metadata` `groups` | `manage-api-metadata`, `associate-groups-with-an-api` [admin] |
| Group members | ~8 | `apim_group.members` | fold into `create-group-with-member` [admin] |

Blocked: **`reuse-shared-policy-group`** — GKO-3001 (admission rejects the
documented `sharedPolicyGroupRef` flow form because it does not resolve the ref
before APIM's dry-run) *and* Terraform (`apim_shared_policy_group` exposes no
`cross_id`, and only the crossId executes the SPG at the gateway). The journey
documents the correct form as a `pending` fixme rather than being green-washed.

### Genuinely GKO-only

These have the **operator itself** as the system under test, so there is nothing
for a second provisioner to do. They live in `tests/gko/`.

| Area | Tests | Why |
|---|--:|---|
| Admission webhooks | 27 | CRD schema / dry-run validation at the Kubernetes admission layer |
| mTLS certificates via Secret refs | subset of 29 | resolved from cluster Secrets, not inline content |
| Deployment & reconciliation | 15 | CR `.status` conditions, observedGeneration, operator restart |
| Import / export | 10 | YAML CRD round-trips |
| ConfigMap/Secret templating | 8 | `[[ … ]]` resolution from cluster ConfigMaps/Secrets |
| ManagementContext CRD | 5 | Kubernetes custom-resource lifecycle |
| CRD defaults | 4 | CRD field defaulting |
| Local ConfigMap | 2 | in-cluster ConfigMap locality |

And these have **no Automation API path at all**:

| Area | Tests | Why |
|---|--:|---|
| V2 API lifecycle, V2 pages/fetchers, V2 members, V2 subscriptions | ~39 | the Automation API is V4-only; V2 is legacy |
| Category CRUD (create/rename a category) | ~6 | no `/categories` path; an API can only *reference* categories inline |

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
