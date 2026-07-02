# Journey: provision an API referencing a dictionary property

**As an API producer, I provision an API that references a dictionary property, resolved at the gateway.**

A MANUAL dictionary value is injected into a response header by the API's
transform-headers flow and asserted at the gateway echo. The same journey runs
against both provisioners.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/dictionary.yaml`](./gko/dictionary.yaml) + [`gko/api.yaml`](./gko/api.yaml) | `Dictionary` (MANUAL, deployed) + a V4 API referencing it by HRID via `{#dictionaries[...]}`. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_dictionary` + `apim_apiv4` with an inline `flows` block. |

**What it proves:** a MANUAL dictionary created through either driver is deployed
and resolvable at the gateway — the API's flow resolves
`{#dictionaries['<hrid>']['env']}` to `"test"` in the echo response header.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-2903
npm --prefix test/platform-test run e2e -- --grep @GKO-2903 --provision-with terraform
```
