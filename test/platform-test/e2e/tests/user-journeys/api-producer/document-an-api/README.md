# Journey: document an API

**As an API producer, I ship documentation alongside the API definition, revise
it as the API changes, and remove it when it is obsolete.**

Documentation is an **inline attribute** of the API (`spec.pages` /
`apim_apiv4.pages[]`) — there is no standalone Page resource, but the journey is
still fully expressible through both provisioners.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `ApiV4Definition.spec.pages` (map keyed by hrid); strip = re-apply without it. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4.pages`; `with_page = false` empties the list. |

**What it proves:** the page lands in APIM (`GET /apis/{id}/pages`) with its
name, type, content, published flag and visibility; revising it — rename,
rewrite, narrow the visibility — updates the **existing** page rather than adding
a second one, because it keeps its hrid across the revision; and stripping it
removes it.

**Adjacent, not covered here:** inline page **fetchers** (`pages[].source`, e.g.
an `http-fetcher` pulling a Swagger spec) are also expressible on **both** drivers
(the TF `apim_apiv4.pages[].source` block mirrors `spec.pages.<x>.source`), so a
fetcher-parity journey is a feasible follow-up. Fetcher rejection cases, folder
renames and all **V2** documentation (there is no `apim_apiv2` resource) stay
GKO-only under [`tests/gko/pages/`](../../gko/pages/).

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-1470
npm --prefix test/platform-test run e2e -- --grep @GKO-1470 --provision-with terraform
```
