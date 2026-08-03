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
| [`publish-a-message-api`](./api-producer/publish-a-message-api/) | Publish a V4 MESSAGE (event) API over each consumer entrypoint | GKO-72/73/129/130/132/133/134/136/141 · TF GKO-3006 |
| [`secure-api-with-plan`](./api-producer/secure-api-with-plan/) | Secure an API with a JWT plan and an OAuth2 plan | GKO-162/163 · TF GKO-3004 |
| [`label-an-api`](./api-producer/label-an-api/) | Label a V4 API (inline `apim_apiv4.labels`) | GKO-1473/83 · TF GKO-3007 |
| [`assign-categories-to-api`](./api-producer/assign-categories-to-api/) | Assign a portal category to a V4 API | GKO-267/270 · TF GKO-3031 |
| [`document-an-api`](./api-producer/document-an-api/) | Ship, revise and remove an inline markdown documentation page | GKO-1470/277/278/236/1469/282 · TF GKO-3034 |
| [`configure-visibility-and-lifecycle`](./api-producer/configure-visibility-and-lifecycle/) | Control portal visibility & lifecycle state, then retire the API | GKO-172/173/179/180/1466/140 · TF TBD |
| [`configure-endpoint-failover`](./api-producer/configure-endpoint-failover/) | Give an API a failover policy and tighten its retry budget | GKO-859 · TF TBD |
| [`configure-api-notifications`](./api-producer/configure-api-notifications/) | Notify a group about the API's events (CR ref vs inline block) | GKO-1231/1232/1461/1194/1195 · TF TBD |
| [`apply-policies-to-a-flow`](./api-producer/apply-policies-to-a-flow/) | Add, rewrite and remove a policy — asserted at the gateway | GKO-94/95/96 · TF TBD |
| [`reuse-shared-policy-group`](./api-producer/reuse-shared-policy-group/) | Reuse a Shared Policy Group — ⛔ pending (GKO-3001 + TF crossId gap) | GKO-976/980 · TF GKO-3005 |

### api-consumer

| Journey | Demonstrates | Xray (GKO / TF) |
|---|---|---|
| [`register-and-retire-application`](./api-consumer/register-and-retire-application/) | Register an application with its settings & metadata, update, retire | GKO-335/336/337/194/552 · TF GKO-1383/3002 |
| [`manage-application-members`](./api-consumer/manage-application-members/) | Share an application, re-role and revoke the member | GKO-538/539/534 · TF TBD |
| [`subscribe-to-a-secured-plan`](./api-consumer/subscribe-to-a-secured-plan/) | Subscribe to a JWT plan and an OAuth2 plan (auto-validated) | GKO-800/819/815 · TF TBD |
| [`subscribe-and-call`](./api-consumer/subscribe-and-call/) | Subscribe to an api-key plan and call it (auto/custom/expiry/rotation) | GKO-2825… · TF GKO-2879… |

### platform-admin

| Journey | Demonstrates | Xray (GKO / TF) |
|---|---|---|
| [`create-group-with-member`](./platform-admin/create-group-with-member/) | Create a group with a resolved member, rename it, delete it | GKO-983/986/985 · TF GKO-2865/2868/2869/2866/2872/2870 |
| [`manage-api-members`](./platform-admin/manage-api-members/) | Grant, re-role and revoke access to an API | GKO-247/259/253/213/402 · TF TBD |
| [`associate-groups-with-an-api`](./platform-admin/associate-groups-with-an-api/) | Put an API under a group's ownership, then detach it | GKO-257/314/1004 · TF TBD |
| [`api-references-dictionary-property`](./platform-admin/api-references-dictionary-property/) | An API resolves a MANUAL dictionary property at the gateway | GKO-2903 · TF GKO-2998 |
| [`manage-dynamic-dictionary`](./platform-admin/manage-dynamic-dictionary/) | Run a DYNAMIC dictionary through its lifecycle | GKO-2904/2910/2911/2909 · TF GKO-3014… |

## Running one

Every journey arm carries its own Xray id, so `--grep` selects it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-335
npm --prefix test/platform-test run e2e -- --grep @GKO-335 --provision-with terraform
```

## Authoring a new journey

See [AGENTS.md → Adding a user journey](../../../AGENTS.md#adding-a-user-journey),
or invoke the [`write-e2e-test`](../../../.claude/skills/write-e2e-test/SKILL.md)
skill for the full procedure. Two rules decide where a test goes:

1. **Persona** picks the folder — who performs this journey?
2. **A new folder only when the _story_ differs**, not when the configuration
   does. Variants of one story (message-API entrypoint types, plan security
   types, api-key modes) are a variant table inside the journey's own
   `*.scenario.ts`.

Coverage status and the migration backlog live in
[PARITY.md](../../../PARITY.md).
