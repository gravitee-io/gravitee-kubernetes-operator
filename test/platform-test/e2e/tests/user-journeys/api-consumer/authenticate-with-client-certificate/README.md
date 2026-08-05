# Authenticate with a client certificate

**As an API consumer, I register a client certificate on my application, call an
mTLS-secured API with it, roll onto a replacement before the old one expires, and
revoke it when I am done.**

The journey stands up an API whose only plan is mTLS-secured, an application, and a
subscription binding them, then drives the application's certificate set through
four stages. Every stage asserts at the **gateway**, not only in the definition:
the certificate list APIM reports is necessary but not sufficient evidence that a
consumer can or cannot get through.

| Stage | Certificates | What it proves |
|---|---|---|
| 1 | `client1` | an unauthenticated caller gets 401, `client1` gets 200, and an unrelated certificate is refused |
| 2 | `client1` + `client2` | during a rollover both certificates work, so the consumer never loses access |
| 3 | `client2` | the retired certificate **stops** being accepted |
| 4 | none | revoking the last certificate leaves the application with an empty list and no access |

Stage 3 is the one carrying the most regression value. The test it replaces removed
a certificate and then only checked the custom resource still reported `Accepted`,
which would have passed even if the retired certificate kept working, which is the
entire point of rotating one.

A second scenario covers the deprecated single-certificate field
(`clientCertificate` / `client_certificate`), which both drivers still expose, and
the migration off it onto the list form. It swaps in a different certificate
deliberately: re-registering the same content under the new field trips APIM's
fingerprint-reuse rule, which has its own test in `tests/gko/mtls-certificates`.

Both arms present the **same** certificates, read from the shared PKI in
[`fixtures/mtls-certificates/pki/`](../../../../fixtures/mtls-certificates/pki/),
so any difference between the arms is a real difference in how the two drivers
registered them.

| Arm | Fixtures | Xray |
|---|---|---|
| GKO | [`gko/`](./gko/) | @GKO-2243 @GKO-2231 @GKO-2247 @GKO-2246 @GKO-2219, and @GKO-2244 for the deprecated field |
| Terraform | [`terraform/`](./terraform/) | @GKO-3103, and @GKO-3104 for the deprecated field |

## Staying in `tests/gko/mtls-certificates`

These have the operator, not the platform, as the system under test:

- admission rejections (bad PEM, already expired, `endsAt` before `startsAt`,
  content and ref together, missing fields)
- certificates resolved from cluster `Secret`s and `ConfigMap`s, and `[[ … ]]`
  templating
- base64-encoded certificate content, a CRD input convenience with no equivalent
  in the provider's PEM-only `content` attribute
- the certificate date-window and fingerprint-reuse edge cases

## Running it

```sh
npm --prefix test/platform-test run e2e -- --grep "client certificate"
npm --prefix test/platform-test run e2e -- --grep "Call an mTLS API with a client certificate.*@gko"
```

The Terraform arm needs the gateway's mTLS listener reachable at
`gateway.mtlsBaseUrl` in `config.yaml` (`https://localhost:30084` locally).
