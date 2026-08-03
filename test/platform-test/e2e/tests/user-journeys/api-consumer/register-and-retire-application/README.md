# Journey: register, update, and retire an application

**As an application developer, I register an application, edit it, then retire it.**

The same journey runs against both provisioners and asserts the same APIM outcome.
Everything needed lives in this folder.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/application.yaml`](./gko/application.yaml) + [`gko/application-updated.yaml`](./gko/application-updated.yaml) | Application CR (created + updated states); retire = delete the CR. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_application`; `description` re-applied for the update, `create_application = false` to retire. |

**What it proves:** an application created through either driver lands in APIM via
the Automation API (`origin: KUBERNETES`), reflects a description update, and is
ARCHIVED when retired.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-335
npm --prefix test/platform-test run e2e -- --grep @GKO-335 --provision-with terraform
```
