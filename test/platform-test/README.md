# @gravitee/platform-test

Assertion library for Gravitee platform e2e testing. TypeScript, minimal dependencies, test-runner agnostic.

## Running the E2E test suite

The Playwright suite under [`e2e/`](e2e/) drives end-to-end tests for both the
**Gravitee Kubernetes Operator** and the **Terraform APIM provider** against a
live APIM + Gateway stack. See **[e2e/README.md](e2e/README.md)** for
prerequisites, cluster bootstrap, and how to run the suite.

Quick start (assumes a local cluster with APIM + GKO is already running):

```bash
npm install && npm run build
npm run e2e
```

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
// import { loadGraviteeConfig, createMapiFromConfig } from "@gravitee/platform-test/config";
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

## Configuration

`loadGraviteeConfig(path)` reads a `config.yaml`; with no path it looks in the
current working directory. The e2e suite loads the one at the package root.

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

This covers the core API/Plan/Subscription assertions. `Mapi` also exposes
`waitFor*` polling variants and Application, Group, Category, and Page
assertion families — see
[`src/assertions/apim/mapi.ts`](src/assertions/apim/mapi.ts) for the full
surface.

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
├── provisioners/ Pluggable creation paths (GKO, Terraform) behind one interface
│   ├── engines/  kubectl + terraform CLI wrappers
│   ├── gko/      GkoProvisioner + its ProvisionerView
│   └── terraform/ TerraformProvisioner + its ProvisionerView and checks
├── utils/
│   ├── http/     HTTP client (native fetch, swappable for undici/mTLS)
│   └── match/    Deep partial matching engine, poll utility
├── types/        TypeScript type definitions for APIM entities
└── config/       config.yaml loading + env-var overrides
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
import { GkoProvisioner, assertProvisioner } from "@gravitee/platform-test/provisioners";
```

## Development

```bash
npm install
npm run build               # TypeScript compilation
npm test                    # Unit tests (vitest)
npm run typecheck           # src/ — also typecheck:e2e, typecheck:examples
npm run check:lanes         # assert the provisioner lanes partition the suite
```

For the e2e suite itself see [e2e/README.md](e2e/README.md) (bootstrap) and
[AGENTS.md](AGENTS.md) (how to write a test).

AI agents contributing here have two rule sets in [`.agent/rules/`](.agent/rules/)
(writing tests · extending the framework) and four skills in
[`.claude/skills/`](.claude/skills/): `write-e2e-test`,
`investigate-e2e-failure`, `add-provisioner`, `extend-platform-assertions`.

## License

Apache-2.0
