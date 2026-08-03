# Journey: add an inline markdown page to an API

**As an API producer, I ship documentation alongside my API definition.**

Add an inline markdown page to a V4 API, then strip it. A page is an **inline
attribute** of the API (`apim_apiv4.pages[]` / `spec.pages`) — there is no
standalone Page resource, but the journey is still fully expressible through both
provisioners.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/api-with-page.yaml`](./gko/api-with-page.yaml) + [`gko/api-without-page.yaml`](./gko/api-without-page.yaml) | `ApiV4Definition.spec.pages` (map keyed by hrid); strip = re-apply without it. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4.pages`; `with_page = false` empties the list. |

**What it proves:** an inline markdown page set through either driver lands in
APIM (`GET /apis/{id}/pages`); stripping it removes it. The payload is a nested
object, so the fixtures and the assertion compare a page object
(`{ name, type, content, published }`), not a plain string.

**Adjacent, not covered here:** inline page **fetchers** (`pages[].source`, e.g.
an `http-fetcher` pulling a Swagger spec) are also expressible on **both** drivers
(the TF `apim_apiv4.pages[].source` block mirrors `spec.pages.<x>.source`), so a
fetcher-parity journey is a feasible follow-up. Only **V2** documentation has no
Terraform path (there is no `apim_apiv2` resource) and stays GKO-only under
[`tests/gko/pages/`](../../gko/pages/).

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-1470
npm --prefix test/platform-test run e2e -- --grep @GKO-1470 --provision-with terraform
```
