# Use case: stand up a message (event) API

Provision a V4 MESSAGE API (HTTP-GET + webhook subscription entrypoints over a
mock message endpoint). The same journey runs against both provisioners.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/message-api.yaml`](./gko/message-api.yaml) | `ApiV4Definition` `type: MESSAGE`. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4` `type = "MESSAGE"` with HTTP + subscription listeners. |

**What it proves:** a MESSAGE API created through either driver is recorded in
APIM as `type: MESSAGE` and reaches `state: STARTED`. The entrypoint-type matrix
(SSE/webhook/websocket consumption) stays GKO-only.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-72
npm --prefix test/platform-test run e2e -- --grep @GKO-72 --provision-with gko
npm --prefix test/platform-test run e2e -- --grep @GKO-72 --provision-with terraform
```
