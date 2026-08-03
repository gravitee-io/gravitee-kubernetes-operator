# Journey: assign a category to an API

**As an API producer, I organise my API in the portal by assigning it to a category.**

Assign a pre-existing portal category to a V4 API, then strip it. Categories are
an **inline attribute** of the API (`apim_apiv4.categories`): there is no
standalone Terraform resource and no GKO Category CRD, and an API can only
*reference* a category that already exists in the environment (**APIM silently
drops references to unknown categories**). The category itself is therefore
created once as a provisioner-agnostic precondition via mAPI
(`POST /configuration/categories`), and both drivers assign it by key.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/api-with-categories.yaml`](./gko/api-with-categories.yaml) + [`gko/api-without-categories.yaml`](./gko/api-without-categories.yaml) | `ApiV4Definition.spec.categories`; strip = re-apply without them. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4.categories`; `with_categories = false` empties the list. |

**Precondition:** the scenario `beforeEach` creates the referenced category
(`e2e-portal-category`) via mAPI and deletes it in `afterEach`. Neither GKO nor
Terraform can create a category; they can only assign an existing one.

**What it proves:** a category assigned through either driver lands on the API in
APIM; stripping it removes the assignment. This is the same inline-attribute
pattern as [`label-an-api`](../label-an-api/) and the template for the other
inline `apim_apiv4` attributes (`groups`, `metadata`, inline `pages[]`).

Run it:

```sh
# GKO arm (title carries @GKO-267 + @GKO-270)
npm --prefix test/platform-test run e2e -- --grep @GKO-267
# Terraform arm
npm --prefix test/platform-test run e2e -- --grep @GKO-3031
```
