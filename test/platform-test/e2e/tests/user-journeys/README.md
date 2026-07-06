# User-journey catalog

Each folder here is a **self-contained, runnable demo of one customer journey**,
authored once and run against every provisioner (GKO + Terraform) via
`forEachProvisioner`. Everything a journey needs is co-located in its folder:

```
<journey>/
  <journey>.scenario.ts   # the shared intent, run against both provisioners
  gko/        *.yaml       # the GKO custom resources
  terraform/  main.tf      # the Terraform equivalent
  README.md                # what it demonstrates ("As a <persona>, I ...") + how to run it
```

Point someone at a folder to show them how to provision X through either driver.
GKO CRs carry a `gravitee.io/e2e: "true"` label (sweep with
`kubectl delete <kind> -l gravitee.io/e2e=true`); Terraform resources are cleaned
up by `terraform destroy`.

| Journey | Demonstrates | Xray (GKO / TF) |
|---|---|---|
| [`subscribe-and-call`](./subscribe-and-call/) | Subscribe an application to an api-key plan and call it (auto/custom/expiry/rotation/…) | GKO-2825… · TF GKO-2879… |
| [`api-references-dictionary-property`](./api-references-dictionary-property/) | An API resolves a dictionary property at the gateway | GKO-2903 · TF GKO-2998 |
| [`create-group-with-member`](./create-group-with-member/) | Create a group with a member | GKO-983 · TF GKO-2865 |
| [`register-and-retire-application`](./register-and-retire-application/) | Register, update, and retire an application | GKO-335/336/337 · TF GKO-3002 |
| [`publish-api-and-serve-traffic`](./publish-api-and-serve-traffic/) | Publish a V4 API, start/stop, gateway 200↔404 | GKO-69/1464 · TF GKO-3003 |
| [`secure-api-with-plan`](./secure-api-with-plan/) | Secure an API with a JWT plan and an OAuth2 plan | GKO-162/163 · TF GKO-3004 |
| [`consume-message-api`](./consume-message-api/) | Stand up a V4 MESSAGE (event) API | GKO-72/73 · TF GKO-3006 |
| [`label-an-api`](./label-an-api/) | Label a V4 API (inline `apim_apiv4.labels`) | GKO-1473 · TF GKO-3007 |
| [`assign-categories-to-api`](./assign-categories-to-api/) | Assign a portal category to a V4 API (inline `apim_apiv4.categories`; category pre-created via mAPI) | GKO-267/270 · TF GKO-3031 |
| [`reuse-shared-policy-group`](./reuse-shared-policy-group/) | Reuse a Shared Policy Group across an API — ⛔ pending (GKO-3001 + TF crossId gap) | GKO-976/980 · TF GKO-3005 |

Run any journey by its Xray tag (both arms), or pin a driver:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-335
npm --prefix test/platform-test run e2e -- --grep @GKO-335 --provision-with terraform
```

## Authoring a new journey

See [AGENTS.md → Adding a cross-provisioner parity scenario](../../../AGENTS.md#adding-a-cross-provisioner-parity-scenario)
and the prioritised backlog + scorecard in [PARITY.md](../../../PARITY.md).
