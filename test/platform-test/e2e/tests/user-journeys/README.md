# User-journey catalog

Each folder here is a **self-contained, runnable demo of one customer journey**,
authored once and run against every provisioner (GKO + Terraform) via
`forEachProvisioner`. Point someone at a folder to show them how to do X through
either driver.

Journeys are grouped by **persona** — *who performs this?* — so the folder name
answers the only question you need to place a new one:

| Folder | Persona |
|---|---|
| [`api-producer/`](./api-producer/) | "As an API producer, I …" — designing, publishing, securing and documenting an API |
| [`api-consumer/`](./api-consumer/) | "As an API consumer, I …" — registering an application, subscribing, and calling the gateway |
| [`platform-admin/`](./platform-admin/) | "As a platform administrator, I …" — groups, members, dictionaries, environment-level config |

```
<persona>/<journey>/
  <journey>.scenario.ts   # the shared intent, run against every provisioner
  gko/        *.yaml      # the GKO custom resources
  terraform/  main.tf     # the Terraform equivalent
  README.md               # "As a <persona>, I ..." + how to run it
```

GKO CRs carry a `gravitee.io/e2e: "true"` label (sweep with
`kubectl delete <kind> -l gravitee.io/e2e=true`); Terraform resources are cleaned
up by `terraform destroy`.

## Journeys

### api-producer

| Journey | Demonstrates | Xray (GKO / TF) |
|---|---|---|
| [`publish-api-and-serve-traffic`](./api-producer/publish-api-and-serve-traffic/) | Publish a V4 API, start/stop, gateway 200↔404 | GKO-69/1464 · TF GKO-3003 |
| [`publish-a-message-api`](./api-producer/publish-a-message-api/) | Publish a V4 MESSAGE (event) API | GKO-72/73 · TF GKO-3006 |
| [`secure-api-with-plan`](./api-producer/secure-api-with-plan/) | Secure an API with a JWT plan and an OAuth2 plan | GKO-162/163 · TF GKO-3004 |
| [`label-an-api`](./api-producer/label-an-api/) | Label a V4 API (inline `apim_apiv4.labels`) | GKO-1473 · TF GKO-3007 |
| [`assign-categories-to-api`](./api-producer/assign-categories-to-api/) | Assign a portal category to a V4 API | GKO-267/270 · TF GKO-3031 |
| [`add-inline-markdown-page-in-api`](./api-producer/add-inline-markdown-page-in-api/) | Add an inline markdown documentation page | GKO-1470 · TF GKO-3034 |
| [`reuse-shared-policy-group`](./api-producer/reuse-shared-policy-group/) | Reuse a Shared Policy Group — ⛔ pending (GKO-3001 + TF crossId gap) | GKO-976/980 · TF GKO-3005 |

### api-consumer

| Journey | Demonstrates | Xray (GKO / TF) |
|---|---|---|
| [`register-and-retire-application`](./api-consumer/register-and-retire-application/) | Register, update, and retire an application | GKO-335/336/337 · TF GKO-1383/3002 |
| [`subscribe-and-call`](./api-consumer/subscribe-and-call/) | Subscribe to an api-key plan and call it (auto/custom/expiry/rotation) | GKO-2825… · TF GKO-2879… |

### platform-admin

| Journey | Demonstrates | Xray (GKO / TF) |
|---|---|---|
| [`create-group-with-member`](./platform-admin/create-group-with-member/) | Create, rename, and delete a group | GKO-983/986/985 · TF GKO-2865/2868/2869/2866 |
| [`api-references-dictionary-property`](./platform-admin/api-references-dictionary-property/) | An API resolves a MANUAL dictionary property at the gateway | GKO-2903 · TF GKO-2998 |
| [`manage-dynamic-dictionary`](./platform-admin/manage-dynamic-dictionary/) | Run a DYNAMIC dictionary through its lifecycle | GKO-2904/2910/2911/2909 · TF GKO-3014… |

## Running one

Every journey arm carries its own Xray id, so `--grep` selects it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-335
npm --prefix test/platform-test run e2e -- --grep @GKO-335 --provision-with terraform
```

## Authoring a new journey

See [AGENTS.md → Adding a user journey](../../../AGENTS.md#adding-a-user-journey).
Two rules decide where a test goes:

1. **Persona** picks the folder — who performs this journey?
2. **A new folder only when the _story_ differs**, not when the configuration
   does. Variants of one story (message-API entrypoint types, plan security
   types, api-key modes) are a variant table inside the journey's own
   `*.scenario.ts`.

Coverage status and the migration backlog live in
[PARITY.md](../../../PARITY.md).
