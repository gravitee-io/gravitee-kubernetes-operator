# Journey: provision an API referencing a dictionary property

**As an API producer, I provision an API that references a dictionary property, resolved at the gateway.**

A MANUAL dictionary value is injected into a response header by the API's
transform-headers flow and asserted at the gateway echo. The same journey runs
against both provisioners, with two scenarios:

1. **Resolve** — the dictionary value resolves to `"test"` in the header.
2. **Update propagates** — changing the property value re-applies the dictionary
   and the new value reaches the gateway.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/dictionary.yaml`](./gko/dictionary.yaml) + [`gko/api.yaml`](./gko/api.yaml) + [`params.ts`](./params.ts) | `Dictionary` (MANUAL, deployed) + a V4 API referencing it by HRID via `{#dictionaries[...]}`. The update scenario renders the dictionary from params (applied via `applyParams`). |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_dictionary` + `apim_apiv4` with an inline `flows` block; the property value is the `env_value` tfvar. |

**What it proves:** a MANUAL dictionary created through either driver is deployed
and resolvable at the gateway — the API's flow resolves
`{#dictionaries['<hrid>']['env']}` in the echo response header — and an in-place
property update propagates to the gateway. DYNAMIC dictionaries and the
undeploy/delete lifecycle are covered by the `manage-dynamic-dictionary` journey.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-2903
npm --prefix test/platform-test run e2e -- --grep @GKO-2903 --provision-with terraform
```
