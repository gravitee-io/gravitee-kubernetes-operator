# Journey: sync API documentation from a URL

**As an API producer, I keep my OpenAPI spec at a URL rather than pasting it into
the API definition, and let APIM pull it in and keep it fresh.**

The fetcher is an **inline attribute** of the page (`pages[].source`) on both
drivers, so neither needs a standalone page resource. The fixtures declare no
`content` at all: every character of it comes from the fetch.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `ApiV4Definition.spec.pages.<hrid>.source` (`http-fetcher` + its configuration); strip = re-apply without the page. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4.pages[].source`, whose `configuration` is a JSON string; `with_page = false` empties the list. |

**What it proves:** the page and its `source` (`type` plus the fetcher
configuration: `url`, `autoFetch`, `fetchCron`, `useSystemProxy`) reach APIM;
APIM really fetches the spec, so the page content parses as a spec document with
a title and paths; renaming the page keeps its hrid, so the rename updates the
existing page and the fetched content is **not** discarded by the sync; and
stripping the page removes it.

**Adjacent, not covered here:** documentation authored **inline** is
[`document-an-api`](../document-an-api/). `github-fetcher` stays GKO-only —
admission pre-fetches the repository at apply time and the test cluster has no
GitHub credentials. Fetcher rejection cases (a web fetcher with no URL, an
invalid cron) and all **V2** documentation live in
[`tests/gko/pages/`](../../gko/pages/).

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-279
npm --prefix test/platform-test run e2e -- --grep @GKO-3109 --provision-with terraform
```
