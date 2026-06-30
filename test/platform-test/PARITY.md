# GKO ↔ Terraform e2e test parity

Status of e2e coverage across the two provisioners the suite drives: the **GKO**
operator (Kubernetes CRs) and the **Terraform** APIM provider (`gravitee-io/apim`).
This is the analysis deliverable for GKO-2907 and the living backlog for GKO-2918.

> Counts are regenerated from the tree, not hand-maintained. Re-count with the
> snippet at the bottom.

## Provider-resource reality (GKO-2918)

The Terraform provider `gravitee-io/apim` registers exactly **6 standalone
resource types** (verified in the provider source `internal/provider/provider.go`
→ `Resources()`, not just the docs): `apim_apiv4`, `apim_application`,
`apim_subscription`, `apim_group`, `apim_dictionary`, `apim_shared_policy_group`.
A provider can only manage resource types it implements, so there is genuinely no
`apim_apiv2`, `apim_category`, `apim_notification`, or standalone `apim_page`.

**But standalone-resource ≠ no parity.** `apim_apiv4` carries rich **inline
attributes**, so several areas with no standalone resource are still fully
parity-testable at the API level:

| Inline attribute on `apim_apiv4` | Enables parity for |
|---|---|
| `type` (incl. MESSAGE), `state`, `visibility`, `lifecycle_state` | API lifecycle, message APIs |
| `plans[]` (with `security.type`) | plans (keyless/JWT/OAuth2/api-key) |
| `flows[]` | policies, SPG reuse |
| `labels` | API labels ✅ (see `label-an-api` journey) |
| `categories` | assigning categories to an API (can't *create* a category) |
| `groups` | associating existing groups with an API |
| `metadata` | API metadata |
| `pages[]` | inline markdown documentation (not standalone pages / fetchers) |

Genuinely impossible (no resource and no inline path): **V2 API lifecycle** (no
`apim_apiv2` at all), standalone **page fetchers**, environment-level **category
CRUD** and **notification configuration**.

## Use-case journeys (where parity lives)

Parity is organised as **customer-journey** scenarios, authored once and run
against both provisioners. Each journey is a self-contained, documented fixture
folder under [`e2e/fixtures/use-cases/`](./e2e/fixtures/use-cases/) (the catalog
lists them all), with the scenario under `e2e/tests/scenarios/<journey>/`.

| Journey | Resources | GKO / TF Xray |
|---|---|---|
| register-and-retire-application | `apim_application` | GKO-335/336/337 · TF GKO-1383 + new |
| publish-api-and-serve-traffic | `apim_apiv4` | GKO-69/1464 · TF new |
| secure-api-with-plan | `apim_apiv4.plans[]` | GKO-162/163 · TF new |
| reuse-shared-policy-group | `apim_shared_policy_group` | GKO-976/980 · TF new |
| consume-message-api | `apim_apiv4` (MESSAGE) | GKO-72/73 · TF new |
| label-an-api | `apim_apiv4.labels` (inline) | GKO-1473 · TF new |
| subscribe-and-call (api-key) | `apim_subscription` | existing |
| api-references-dictionary-property | `apim_dictionary` | existing |
| create-group-with-member | `apim_group` | existing |

> **Follow-up gap:** `reuse-shared-policy-group` asserts SPG reuse at the APIM
> config level (the API flow invokes the SPG). An end-to-end gateway check (the
> SPG's injected header reflected by the echo backend) did **not** resolve for
> EITHER provisioner, which points at an SPG deployment-lifecycle gap rather than
> a provisioner difference. Worth a focused investigation / possible product bug.

We do **not** want full parity. A large share of GKO coverage exercises
Kubernetes-only mechanics (admission, status conditions, templating, operator
restart) that have no Terraform equivalent, and Terraform has its own surface
(drift, redaction, plan exit codes). Parity means: **the APIM resources every
customer touches should be exercised through both drivers**, preferably once, as a
shared scenario in the provisioner layer (`tests/scenarios/`).

---

## The provisioner layer (where parity lives)

A shared scenario is authored once and run against every provisioner via
`forEachProvisioner` (`e2e/helpers/for-each-provisioner.ts`) with
`gkoScenario()` / `tfScenario()` (`e2e/helpers/provisioner-env.ts`). The body is
provisioner-agnostic (`provisioned` + `mapi`/`gateway`); each arm carries its own
Xray id. See the worked example + authoring guide in
[AGENTS.md](./AGENTS.md#adding-a-cross-provisioner-parity-scenario).

Today only **subscriptions/apikey** (10 scenarios), **groups** (1), and
**dictionaries** (1, this change) live there. Everything else is per-driver.

---

## Parity matrix (APIM resources — parity candidates)

| Feature area | Journey | Status |
|---|---|---|
| V4 API lifecycle (start/stop, visibility, lifecycle) | publish-api-and-serve-traffic | 🟢 done via journey |
| Applications (CRUD) | register-and-retire-application | 🟢 done via journey |
| Subscriptions — api-key | subscribe-and-call (apikey) | 🟢 done via journey |
| Plans (JWT / OAuth2 security types) | secure-api-with-plan | 🟢 done via journey |
| Shared Policy Groups | reuse-shared-policy-group | 🟢 done via journey (config-level; gateway gap noted above) |
| Dictionaries | api-references-dictionary-property | 🟢 done via journey |
| Message APIs (V4) | consume-message-api | 🟢 done via journey |
| Groups + members | create-group-with-member | 🟡 TF-led; journey covers create |
| Labels | label-an-api | 🟢 done via journey (inline `apim_apiv4.labels`) |
| Categories (assign to API) | — | 🟡 expressible inline (`apim_apiv4.categories`); next journey |
| Pages — inline markdown | — | 🟡 expressible inline (`apim_apiv4.pages[]`); next journey |
| Pages — standalone + fetchers | — | ⛔ no standalone `apim_page` |
| Notifications | — | ⛔ no `apim_notification` (no inline path) |
| V2 API lifecycle | — | ⛔ no `apim_apiv2` (no inline path) |
| Applications — members / OAuth | — | GKO-only (admission + member reconciliation) |

Legend: 🟢 parity met via journey · 🟡 expressible, journey pending · ⛔ no TF path (stays GKO-only).

---

## Intentionally GKO-only (no Terraform parity expected)

These exercise Kubernetes/operator mechanics the Terraform provider has no concept
of. They stay in `tests/gko/` and are **out of scope** for parity.

| Area | `tests/gko` | Why GKO-only |
|---|---:|---|
| Admission webhooks | 27 | CRD schema/dry-run validation at the K8s admission layer |
| Deployment & reconciliation | 15 | CR `.status` conditions, observedGeneration, operator restart |
| mTLS certificates (Application CRs) | 29 | Secret-backed Application TLS settings resolved from the cluster |
| Import / export | 10 | YAML CRD export/import round-trips |
| ConfigMap/Secret templating | 8 | `{{ … }}` resolution from cluster ConfigMaps/Secrets |
| ManagementContext CRD | 5 | Kubernetes custom resource lifecycle |
| CRD defaults | 4 | CRD spec field defaulting |
| Local ConfigMap | 2 | In-cluster ConfigMap locality/cleanup |

These exercise APIM behaviour with no Terraform path at all (no standalone
resource AND no inline `apim_apiv4` attribute), so they stay GKO-only:

| Area | `tests/gko` | Why no TF path |
|---|---:|---|
| V2 API lifecycle | ~12 | no `apim_apiv2` resource; provider is V4-only |
| Notifications | ~11 | no `apim_notification` resource, no inline attribute |
| Standalone pages / fetchers | ~20 | only inline `apim_apiv4.pages[]`; no `apim_page`, no fetcher support |
| Category CRUD (create/rename a category) | ~6 | no `apim_category`; an API can only *reference* categories inline |

Partially TF-expressible at the API level (assign-to-API only, no standalone
CRUD) and good candidates for follow-up journeys: **categories**
(`apim_apiv4.categories`), **inline markdown pages** (`apim_apiv4.pages[]`),
**group association** (`apim_apiv4.groups`), **metadata** (`apim_apiv4.metadata`).

## Intentionally Terraform-only (no GKO parity expected)

| Behaviour | Where | Why TF-only |
|---|---|---|
| Drift detection | `terraform plan` / group drift | HCL desired-state vs server reconciliation |
| Sensitive redaction | apikey sensitive-in-plan | provider `Sensitive` attribute handling |
| Plan exit codes / idempotency | post-apply, group import | `terraform plan -detailed-exitcode` semantics |

---

## Parity backlog (GKO-2918)

This batch (GKO-2918) delivered five new journeys (applications, V4 lifecycle,
plans, SPGs, message APIs) on top of the three existing scenarios (apikey,
dictionary, groups). What remains:

| # | Area | Status |
|---|---|---|
| 1 | Dictionaries / apikey / groups | ✅ done (pre-existing scenarios) |
| 2 | Applications (CRUD) | ✅ done — register-and-retire-application |
| 3 | V4 API lifecycle / visibility | ✅ done — publish-api-and-serve-traffic |
| 4 | Plans (JWT / OAuth2) | ✅ done — secure-api-with-plan |
| 5 | Shared Policy Groups | ✅ done — reuse-shared-policy-group (config-level) |
| 6 | Message APIs (V4) | ✅ done — consume-message-api |
| 7 | Labels | ✅ done — label-an-api (inline `apim_apiv4.labels`) |
| 8 | Categories (assign) · inline pages · group-assoc · metadata | ⏳ future journeys (inline `apim_apiv4` attrs, no standalone resource) |
| 9 | mTLS plans, gateway JWT/OAuth2 enforcement | ⏳ future (subscription + token orchestration) |
| 10 | SPG end-to-end gateway execution | 🔎 investigate (gap noted above; possible product bug) |
| 11 | Notifications · standalone pages/fetchers · category CRUD · V2 lifecycle | ⛔ no TF path — stay GKO-only |
| 12 | Relocate existing 3 scenarios into `fixtures/use-cases/` | ⏳ follow-up (consolidation only) |

When picking up a row, prefer **a use-case journey** (one intent, two arms) over a
standalone TF test, and follow the de-dup rule (remove the now-shared assertion
from the GKO-only file). Confirm the provider exposes the resource first — the 6
resources are `apim_apiv4`, `apim_application`, `apim_subscription`, `apim_group`,
`apim_dictionary`, `apim_shared_policy_group`.

---

## Re-generating the counts

```sh
cd test/platform-test/e2e/tests
# per-area GKO counts
for d in gko/*/; do printf "%-32s %s\n" "$d" \
  "$(grep -rhoE '\btest(\.(skip|fixme|only))?\(' "$d" --include='*.ts' | wc -l)"; done
# TF + scenario totals
grep -rhoE '\btest(\.(skip|fixme|only))?\(' terraform --include='*.ts' | wc -l
grep -rhoE 'forEachProvisioner(<[^>]*>)?\(' scenarios --include='*.ts' | wc -l   # ×2 drivers
```
