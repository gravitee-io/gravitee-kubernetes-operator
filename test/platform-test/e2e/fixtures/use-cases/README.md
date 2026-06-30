# Use-case catalog

Each folder here is a **self-contained, runnable demo of one customer journey**,
authored once and run against every provisioner (GKO + Terraform) by the matching
scenario under `e2e/tests/scenarios/<journey>/`. The fixtures are the
documentation: point someone at a folder to show them how to provision X.

A folder contains:
- `gko/*.yaml` — the GKO custom resources
- `terraform/main.tf` — the Terraform equivalent
- `README.md` — what it demonstrates and how to run it

| Journey | Demonstrates | Xray (GKO / TF) |
|---|---|---|
| [`register-and-retire-application`](./register-and-retire-application/) | Register, update, and retire (archive) an application | GKO-335/336/337 · TF GKO-1383 + new |
| [`publish-api-and-serve-traffic`](./publish-api-and-serve-traffic/) | Publish a V4 API (visibility + lifecycle), start/stop, gateway 200↔404 | GKO-69/1464 · TF new |
| [`secure-api-with-plan`](./secure-api-with-plan/) | Secure an API with a JWT plan and an OAuth2 plan | GKO-162/163 · TF new |
| [`reuse-shared-policy-group`](./reuse-shared-policy-group/) | Reuse a Shared Policy Group across a V4 API (attach/detach) — ⛔ pending (GKO-3001 + TF crossId gap) | GKO-976/980 · TF GKO-3005 |
| [`consume-message-api`](./consume-message-api/) | Stand up a V4 MESSAGE (event) API | GKO-72/73 · TF new |
| [`label-an-api`](./label-an-api/) | Label a V4 API and strip the labels (inline `apim_apiv4.labels`) | GKO-1473 · TF new |

Run any journey by its Xray tag (both arms), or pin a driver:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-335
npm --prefix test/platform-test run e2e -- --grep @GKO-335 --provision-with terraform
```

## Authoring a new journey

See [AGENTS.md → Adding a cross-provisioner parity scenario](../../AGENTS.md#adding-a-cross-provisioner-parity-scenario)
and the prioritised backlog in [PARITY.md](../../PARITY.md).

## Earlier cross-provisioner scenarios (not yet relocated here)

Three journeys predate this catalog and still live with their fixtures under the
original resource-named folders. Consolidating them into `use-cases/` is a
follow-up (they already run as shared scenarios):

- **subscribe-and-call (api-key)** — `tests/scenarios/subscriptions/apikey/`
- **api-references-dictionary-property** — `tests/scenarios/dictionaries/` (fixtures under `fixtures/dictionaries/{dictionary-manual,api-with-dictionary,manual-resolve}/`)
- **create-group-with-member** — `tests/scenarios/groups/` (fixtures under `fixtures/groups/lifecycle/`)
