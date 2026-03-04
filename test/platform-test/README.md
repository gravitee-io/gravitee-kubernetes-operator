# @gravitee/platform-test

Assertion library for Gravitee platform e2e testing. TypeScript, minimal dependencies, test-runner agnostic.

## Install

```bash
npm install @gravitee/platform-test
```

## Quick Start

```typescript
import { createMapi, poll } from "@gravitee/platform-test";

const mapi = createMapi({
  baseUrl: "http://localhost:8083",
  auth: { type: "basic", username: "admin", password: "admin" },
});

// Or load connection details from config.yaml:
// import { loadGraviteeConfig, createMapiFromConfig } from "@gravitee/platform-test/cmd";
// const config = await loadGraviteeConfig();
// const mapi = createMapiFromConfig(config);

// Assert an API matches an expected shape (partial deep match)
await mapi.assertApiMatches("api-id", {
  name: "Petstore API",
  definitionVersion: "V4",
  type: "PROXY",
  state: "STARTED",
  listeners: [{ type: "HTTP", entrypoints: [{ type: "http-proxy" }] }],
});

// Convenience shortcuts
await mapi.assertApiStarted("api-id");
await mapi.assertPlanPublished("api-id", "plan-id");
await mapi.assertSubscriptionAccepted("api-id", "sub-id");

// Poll for eventual consistency
await poll(() => mapi.assertApiStarted("api-id"), {
  timeoutMs: 15_000,
  description: "API to reach STARTED state",
});

// Gateway assertions (with built-in retry/polling)
const gateway = mapi.gateway({ baseUrl: "http://localhost:8082" });
await gateway.assertResponds("/my-api", { status: 200 });
await gateway.assertNotResponds("/my-api", { notStatus: 200 });
```

## CLI

The library ships a CLI for use in shell scripts and Chainsaw test steps.

```
platform-test <subcommand> [flags]
```

Exit codes: `0` assertion passed, `1` assertion failed, `2` config/network error.

### assert-api

Assert that an APIM API matches an expected partial shape.

```bash
# Assert API state
platform-test assert-api --api-id <id> --state STARTED

# Assert API listener path
platform-test assert-api --api-id <id> --path /petstore

# Assert HTTP status (useful for deletion checks — skips property matching)
platform-test assert-api --api-id <id> --status 404

# Assert arbitrary fields with --match (JSON partial)
platform-test assert-api --api-id <id> --match '{"categories":["finance"]}'
platform-test assert-api --api-id <id> --match '{"visibility":"PUBLIC","tags":["internal"]}'

# Assert from a YAML match file (same shape as the API object)
platform-test assert-api --api-id <id> --match-file expected-api.yaml

# Combine --match-file with flag overrides (flags take precedence)
platform-test assert-api --api-id <id> --match-file expected-api.yaml --state STARTED

# Combine --match-file with --match (--match overrides --match-file, flags override both)
platform-test assert-api --api-id <id> --match-file expected-api.yaml --match '{"visibility":"PUBLIC"}'

# Combine specific flags with --match
platform-test assert-api --api-id <id> --state STARTED --match '{"categories":["finance"]}'
```

| Flag | Description |
|------|-------------|
| `--api-id <id>` | API ID to assert (required) |
| `--status <code>` | Expected HTTP status code (e.g. 404); skips property checks |
| `--state <state>` | Expected lifecycle state (e.g. STARTED, STOPPED) |
| `--path <path>` | Expected listener path (e.g. /petstore) |
| `--match <json>` | Arbitrary JSON partial merged into the assertion |
| `--match-file <file>` | YAML file with expected partial API shape (merged before `--match`/`--state`/`--path`) |
| `--config <file>` | Path to `config.yaml` (default: CWD) |

**How `--match` works:** The JSON value is parsed and merged into the partial object
passed to `assertApiMatches()`. You can assert any field on the API object — categories,
tags, labels, properties, visibility, etc. Explicit flags (`--state`, `--path`) take
precedence over overlapping keys in `--match`.

**How `--match-file` works:** The YAML file is parsed as a plain object with the same shape
as the API (i.e. `DeepPartial<Api>`). This provides a file-based alternative to `--match`
for complex assertions. The merge order is: `--match-file` (base) → `--match` JSON
(overrides) → individual flags (`--state`, `--path`) (highest precedence).

Example `expected-api.yaml`:

```yaml
state: STARTED
categories:
  - finance
listeners:
  - type: HTTP
    paths:
      - path: /petstore
```

### assert-gateway

Assert that the APIM gateway responds with an expected HTTP status. Retries
automatically (500ms interval, 30s timeout by default) to handle eventual
consistency from operator reconciliation and gateway sync.

```bash
# Assert gateway returns 200
platform-test assert-gateway --path /petstore --status 200

# Assert gateway does NOT return 200 (e.g. after API stopped)
platform-test assert-gateway --path /petstore --not-status 200

# With authorization header
platform-test assert-gateway --path /jwt-demo --status 200 --authorization "Bearer <token>"

# Custom gateway URL and timeouts
platform-test assert-gateway --path /petstore --status 200 \
  --gateway http://localhost:9082 \
  --timeout 60000 --retry-interval 1000

# mTLS
platform-test assert-gateway --path /secure --status 200 \
  --cert client.crt --key client.key --cacert ca.crt
```

| Flag | Description |
|------|-------------|
| `--path <path>` | Gateway path to call, e.g. /petstore (required) |
| `--status <code>` | Expected HTTP status code (mutually exclusive with `--not-status`) |
| `--not-status <code>` | Status code that must NOT appear |
| `--gateway <url>` | Gateway base URL (default: `http://localhost:30082`) |
| `--authorization <value>` | Authorization header value |
| `--cert <file>` | Client certificate PEM (mTLS) |
| `--key <file>` | Client private key PEM (mTLS) |
| `--cacert <file>` | CA certificate PEM |
| `--timeout <ms>` | Total retry timeout in ms (default: 30000) |
| `--retry-interval <ms>` | Interval between retries in ms (default: 500) |
| `--config <file>` | Path to `config.yaml` (default: CWD) |

### Configuration File

Both subcommands accept a `--config` flag pointing to a `config.yaml` file.
If not provided, the CLI looks for one in the current working directory.

```yaml
apim:
  baseUrl: http://localhost:30083
  envId: DEFAULT
  auth:
    username: admin
    password: admin

gateway:
  baseUrl: http://localhost:30082
  mtlsBaseUrl: https://localhost:30084
```

Environment variables override config file values:

| Variable | Overrides |
|----------|-----------|
| `GRAVITEE_BASE_URL` | `apim.baseUrl` |
| `GRAVITEE_ENV_ID` | `apim.envId` |
| `GRAVITEE_USERNAME` | `apim.auth.username` |
| `GRAVITEE_PASSWORD` | `apim.auth.password` |
| `GRAVITEE_GATEWAY_URL` | `gateway.baseUrl` |
| `GRAVITEE_GATEWAY_MTLS_URL` | `gateway.mtlsBaseUrl` |

### Usage in Chainsaw Tests

The CLI can be invoked from Chainsaw `script:` steps. Since the package is not
published to a registry, invoke it directly via `node` with a relative path to the
built CLI entry point:

```yaml
bindings:
  - name: platformTestCli
    value: "../../../../../../platform-test/dist/cmd/cli.js"

steps:
  - name: Assert gateway returns 200
    try:
      - script:
          env:
            - name: PLATFORM_TEST_CLI
              value: ($platformTestCli)
            - name: API_NAME
              value: ($apiName)
          content: |
            node $PLATFORM_TEST_CLI assert-gateway --path $API_NAME --status 200

  - name: Assert API categories via management API
    try:
      - script:
          env:
            - name: PLATFORM_TEST_CLI
              value: ($platformTestCli)
            - name: API_VERSION_K8S
              value: ($apiVersionK8s)
            - name: API_NAME
              value: ($apiName)
            - name: CATEGORY_NAME
              value: ($categoryName)
          content: |
            API_ID=$(kubectl get $API_VERSION_K8S -n default $API_NAME -o jsonpath='{.status.id}')
            node $PLATFORM_TEST_CLI assert-api --api-id "$API_ID" --match '{"categories":["'$CATEGORY_NAME'"]}'

  - name: Assert API shape from YAML match file
    try:
      - script:
          env:
            - name: PLATFORM_TEST_CLI
              value: ($platformTestCli)
            - name: API_VERSION_K8S
              value: ($apiVersionK8s)
            - name: API_NAME
              value: ($apiName)
          content: |
            API_ID=$(kubectl get $API_VERSION_K8S -n default $API_NAME -o jsonpath='{.status.id}')
            node $PLATFORM_TEST_CLI assert-api --api-id "$API_ID" --match-file expected-api.yaml
```

## TypeScript API

### mAPI Assertions

```typescript
import { createMapi } from "@gravitee/platform-test";

const mapi = createMapi({
  baseUrl: "http://localhost:8083",
  auth: { type: "basic", username: "admin", password: "admin" },
});
```

| Method | Description |
|--------|-------------|
| `assertApiMatches(apiId, partial)` | Assert API matches a partial shape |
| `checkApiMatches(apiId, partial)` | Non-throwing variant (returns report) |
| `assertApiState(apiId, state)` | Assert specific lifecycle state |
| `assertApiStarted(apiId)` | Assert API is STARTED |
| `assertApiStopped(apiId)` | Assert API is STOPPED |
| `assertApiHttpStatus(apiId, status)` | Assert management API HTTP status (e.g. 404) |
| `assertPlanMatches(apiId, planId, partial)` | Assert plan matches a partial shape |
| `assertPlanPublished(apiId, planId)` | Assert plan is PUBLISHED |
| `assertSubscriptionMatches(apiId, subId, partial)` | Assert subscription matches |
| `assertSubscriptionAccepted(apiId, subId)` | Assert subscription is ACCEPTED |

### Gateway Assertions

```typescript
const gateway = mapi.gateway({ baseUrl: "http://localhost:8082" });

// Assert endpoint returns expected status (retries automatically)
await gateway.assertResponds("/my-api", { status: 200 });

// With auth header
await gateway.assertResponds("/jwt-demo", {
  status: 200,
  headers: { Authorization: "Bearer <token>" },
});

// Assert endpoint stops returning a specific status
await gateway.assertNotResponds("/my-api", { notStatus: 200 });

// mTLS
import { createTlsFetch } from "@gravitee/platform-test/utils/http";
const mtlsFetch = createTlsFetch({ cert, key, ca });
const secureGw = mapi.gateway({ baseUrl: "https://localhost:8443" }, mtlsFetch);
await secureGw.assertResponds("/mtls-demo", { status: 200 });
```

### Matching Engine

```typescript
import { deepPartialMatch } from "@gravitee/platform-test/utils/match";

const report = deepPartialMatch(actualObject, {
  name: "My API",
  state: "STARTED",
  listeners: [{ type: "HTTP" }],
});

if (!report.pass) {
  console.log(report.failures);
  // [{ path: "$.state", expected: "STARTED", actual: "STOPPED", message: "..." }]
}
```

Matching rules:
- Only fields in `expected` are checked; everything else is ignored
- Objects are matched recursively (partial)
- Arrays are matched positionally (`expected[0]` against `actual[0]`, etc.)
- Primitives use strict equality

### Poll Utility

```typescript
import { poll } from "@gravitee/platform-test";

await poll(() => mapi.assertApiStarted(apiId), {
  timeoutMs: 15_000,
  intervalMs: 1_000,
  description: "API to reach STARTED state",
});
```

### Error Output

```
AssertionError: Assertion failed (2 mismatches):
  path:     $.name
  expected: "Petstore API v2"
  actual:   "Petstore API"

  path:     $.state
  expected: "STARTED"
  actual:   "STOPPED"
```

## Architecture

```
@gravitee/platform-test
├── assertions/
│   ├── apim/     mAPI (Management API) & gateway assertions
│   └── am/       Access Management (placeholder)
├── utils/
│   ├── http/     HTTP client (native fetch, swappable for undici/mTLS)
│   └── match/    Deep partial matching engine, poll utility
├── types/        TypeScript type definitions for APIM entities
└── cmd/          CLI entry point and subcommands
```

### Design Principles

| Principle | Implementation |
|---|---|
| **void + throw** | Success = void, failure = `node:assert` AssertionError |
| **Partial matching** | Only assert fields you specify; others are ignored |
| **Minimal deps** | Native `fetch`, `node:assert`, `yaml` for config parsing |
| **Test-runner agnostic** | Works with Vitest, Jest, node:test, Playwright, anything |
| **Extensible** | Add AM, AE, Cockpit modules following the same pattern |

## Sub-package Imports

```typescript
import { Mapi } from "@gravitee/platform-test/assertions/apim";
import { deepPartialMatch, poll } from "@gravitee/platform-test/utils/match";
import { HttpClient } from "@gravitee/platform-test/utils/http";
```

## Development

```bash
npm install
npm run build      # TypeScript compilation
npm test           # Run tests (vitest)
npm run typecheck  # Type check without emitting
```

## License

Apache-2.0
