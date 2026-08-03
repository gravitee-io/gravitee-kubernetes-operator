# Journey: operate a dynamic dictionary through its lifecycle

**As an API producer, I run a DYNAMIC dictionary whose values an API resolves at the gateway, and I manage it through its lifecycle.**

A DYNAMIC dictionary's HTTP provider polls the Gravitee echo endpoint every 5s; a
JOLT spec maps the echoed `X-Test-Specific` header into a property of the same
name. A keyless PROXY API injects that property into the `X-Env` response header
via `transform-headers`. The same four scenarios run against both provisioners.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/api.yaml`](./gko/api.yaml) + [`params.ts`](./params.ts) (`dictionaryYaml`) | Static V4 API referencing the dictionary by HRID (`default-dyn-dictionary`); the `Dictionary` CR is rendered from params and applied via `applyParams`, so `update()` re-applies a changed CR. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_dictionary` (`type = "DYNAMIC"`, count-gated) + `apim_apiv4`; lifecycle knobs come from tfvars (`header_value`, `deployed`, `create_dictionary`). |

**Scenarios (each an arm per provisioner):**

1. **Resolve** — the dictionary value resolves to `ABCDEF` in the `X-Env` header.
2. **Update propagates** — changing the provider header value to `ZYXWVU` (`update()`) reaches the gateway.
3. **`deployed=false` stops** — flipping `deployed=false` (`update()`) stops resolution without deleting the resource.
4. **Delete stops** — removing only the dictionary (`remove("dictionary")`) stops resolution while the API stays up.

The stop assertion is provisioner-agnostic: once the dictionary is gone or
stopped, the gateway either returns **500** (the EL expression can no longer
resolve) or **200** with the value gone.

**What it proves:** a DYNAMIC dictionary created through either driver is started,
resolvable at the gateway, and its lifecycle transitions (provider update,
undeploy, delete) reach the gateway identically. Provisioner-specific dictionary
behaviour with no cross-provisioner meaning (GKO admission validation, GKO
secret-templating, plain CR delete) stays under `tests/gko/dictionaries`.

Run it:

```sh
# both arms
npm --prefix test/platform-test run e2e -- --grep "dynamic dictionary"
# Terraform arm only
npm --prefix test/platform-test run e2e -- --grep "dynamic dictionary" --provision-with terraform
```
