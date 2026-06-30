# GKO ↔ Terraform e2e test parity

Status of e2e coverage across the two provisioners the suite drives: the **GKO**
operator (Kubernetes CRs) and the **Terraform** APIM provider (`gravitee-io/apim`).
This is the analysis deliverable for GKO-2907 and the living backlog for GKO-2918.

> Counts are regenerated from the tree, not hand-maintained. As of **2026-06-29**:
> `tests/gko/` = **310** `test()` blocks across 21 areas, `tests/terraform/` = **20**,
> `tests/scenarios/` = **11** cross-provisioner scenarios (→ 22 runtime tests) + 7
> driver-specific = **29** runtime tests. Re-count with the snippet at the bottom.

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

| Feature area | GKO `tests/gko` | TF `tests/terraform` | Scenario layer | Status |
|---|---:|---:|---|---|
| V4/V2 API lifecycle (start/stop, visibility) | 37 | 1 (post-apply create) | — | ❌ gap |
| Applications (CRUD + members) | 26 | 1 (delete) | — | ❌ gap |
| Subscriptions — api-key | 25 | 2 (errors) | ✅ apikey (10) | 🟢 done via scenario |
| Subscriptions — other plan types | (incl. above) | 0 | — | ❌ gap |
| Plans & policies | 10 | 1 (general conditions) | — | ❌ gap |
| Groups + members | 5 | 13 | ✅ groups (1) | 🟡 TF-led; scenario covers create |
| Categories & labels | 10 | 0 | — | ❌ gap |
| **Dictionaries** | 7 | 1 (new) | ✅ dictionary (1) | 🟡 in progress (GKO-2997) |
| Pages / documentation | 25 | ~1 (hierarchy) | — | ❌ gap |
| Shared Policy Groups | 4 | 0 | — | ❌ gap |
| Notifications | 11 | 0 | — | ❌ gap |
| Message APIs (V4) | 9 | 0 | — | ❌ gap |
| Analytics | 1 | 0 | — | ❌ gap |

Legend: ✅ covered · 🟢 parity met · 🟡 partial / in progress · ❌ no TF coverage.

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

## Intentionally Terraform-only (no GKO parity expected)

| Behaviour | Where | Why TF-only |
|---|---|---|
| Drift detection | `terraform plan` / group drift | HCL desired-state vs server reconciliation |
| Sensitive redaction | apikey sensitive-in-plan | provider `Sensitive` attribute handling |
| Plan exit codes / idempotency | post-apply, group import | `terraform plan -detailed-exitcode` semantics |

---

## Parity backlog (GKO-2918 follow-up scope)

Prioritised by customer impact × effort. **This change delivers dictionaries only;
the rest is explicitly deferred to future tickets**, tracked here.

| # | Area | Approx. new TF/scenario tests | Notes | Status |
|---|---|---:|---|---|
| 1 | **Dictionaries** | 1+ | MANUAL resolve via gateway. Pattern reference. | 🟡 this change (GKO-2997) |
| 2 | Applications (CRUD + members) | ~6 | `apim_application` exists; mirror app lifecycle | ⏳ future ticket |
| 3 | Categories & labels | ~4 | check provider resource availability first | ⏳ future ticket |
| 4 | Plans (keyless, JWT, OAuth2, mTLS) | ~6 | `apim_apiv4.plans[]` already used | ⏳ future ticket |
| 5 | V4/V2 API lifecycle (start/stop, visibility) | ~6 | extend `post-apply` into scenarios | ⏳ future ticket |
| 6 | Shared Policy Groups | ~4 | `apim_shared_policy_group` exists in provider | ⏳ future ticket |
| 7 | Pages / documentation | ~6 | provider page/doc resources | ⏳ future ticket |
| 8 | Notifications | ~4 | confirm provider support | ⏳ future ticket |
| 9 | Message APIs (V4) | ~4 | confirm provider support | ⏳ future ticket |

When picking up a backlog row, prefer **migrating the GKO test into a shared
scenario** (one intent, two arms) over writing a standalone TF test, and follow the
de-dup rule (remove the now-shared assertion from the GKO-only file). Confirm the
provider exposes the resource first — the registry README lists `apim_apiv4`,
`apim_application`, `apim_subscription`, `apim_group`, `apim_dictionary`,
`apim_shared_policy_group`.

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
