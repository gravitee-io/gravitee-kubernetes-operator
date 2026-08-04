# Journey templates

Copy-paste skeletons for a new `e2e/tests/user-journeys/<persona>/<journey>/`.
Derived from `api-producer/publish-api-and-serve-traffic/`, which is the
reference implementation to read when in doubt.

Every `.ts` and `.yaml` file starts with the Apache 2.0 licence header (CI job
`job-lint-licenses` enforces it for sources). Copy it from any existing file in
the suite; it is elided below for brevity as `<LICENCE>`.

---

## `<journey>.scenario.ts`

```ts
<LICENCE>

/**
 * Journey: <one line, the customer-visible outcome>.
 *
 * As a <persona>, I <story>. <What is asserted and why it carries regression
 * value.>
 *
 * Fixtures are co-located in this folder. <Anything that stays GKO-only or
 * Terraform-only, and why.>
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The single knob this journey varies. */
interface MyParams {
  state: "STARTED" | "STOPPED";
}

forEachProvisioner<MyParams>(
  {
    title: "<Imperative sentence describing the journey>",
    provisioners: {
      gko: gkoScenario<MyParams>({
        manifests: [path.join(here, "gko/api.yaml")],
        roles: { api: "<cr-metadata-name>" },   // role -> CR name; kind by convention
        contextPath: "/<context-path>",
        // Parameterized resources are applied here instead of via `manifests`.
        // Declare them in `dynamicRoles` so their ids resolve after this runs.
        // dynamicRoles: ["subscription"],
        // applyParams: async (k, params) => { await k.apply(path.join(here, `gko/${params.state}.yaml`)); },
      }),
      terraform: tfScenario<MyParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ state: params.state }),
        // addresses: { api: "apim_apiv4.api" },   // required for view/assertProvisioner
        // outputs: { api: "api_id" },             // only to override the defaults
        // removeVars: { subscription: { create_subscription: false } },
      }),
    },
    xray: {
      gko: XRAY.<AREA>.<ID>,               // a string[] is allowed when one arm covers several tickets
      terraform: XRAY.TERRAFORM.<ID>,
    },
    tags: [TAGS.REGRESSION],
    // since: { gko: "4.12", terraform: "4.12" },
    // timeoutMs: { gko: 90_000 },
    // pending: { terraform: "no `x` attribute on apim_apiv4 yet (GKO-NNNN)" },
  },
  async ({ provisioned, mapi, gateway }) => {
    const apiId = await provisioned.apiId();
    const ctx = await provisioned.contextPath();

    await test.step("<what the initial state proves>", async () => {
      await mapi.waitForApiMatches(
        apiId,
        { state: "STARTED", visibility: "PUBLIC", lifecycleState: "PUBLISHED" },
        { timeoutMs: 30_000 },
      );
      await gateway.assertResponds(ctx, { status: 200 });
    });

    await test.step("<what the change proves>", async () => {
      await provisioned.update({ state: "STOPPED" });
      await mapi.waitForApiStopped(apiId, { timeoutMs: 30_000 });
      await gateway.assertResponds(ctx, { status: 404 });
    });
  },
  { state: "STARTED" },
);
```

### Variant table (the shape to use instead of a second folder)

When the story is the same and only the configuration differs, keep one folder
and one assertion body:

```ts
const VARIANTS = [
  { name: "http-get", manifest: "gko/api-http-get.yaml", entrypoint: "http-get" },
  { name: "sse", manifest: "gko/api-sse.yaml", entrypoint: "sse" },
] as const;

for (const variant of VARIANTS) {
  forEachProvisioner(
    {
      title: `Publish a message API over ${variant.name}`,
      provisioners: { gko: gkoScenario({ manifests: [path.join(here, variant.manifest)], … }), … },
      xray: { gko: XRAY.MESSAGE_APIS[variant.xrayKey], … },
      tags: [TAGS.REGRESSION],
    },
    async ({ provisioned, mapi }) => { /* one shared body */ },
  );
}
```

---

## `gko/<resource>.yaml`

```yaml
<LICENCE>

# <What this manifest represents in the journey, and what re-applying a sibling
# manifest changes.>
apiVersion: gravitee.io/v1alpha1
kind: ApiV4Definition
metadata:
  name: <journey-scoped-name>
  labels:
    gravitee.io/e2e: "true"
spec:
  contextRef:
    name: "dev-ctx"
    namespace: "default"
  name: "<journey-scoped-name>"
  description: "<why this API exists in the journey>"
  version: "1.0.0"
  type: PROXY
  state: STARTED
  visibility: PUBLIC
  lifecycleState: PUBLISHED
  definitionContext:
    origin: KUBERNETES
    syncFrom: MANAGEMENT
  listeners:
    - type: HTTP
      paths:
        - path: "/<context-path>"
      entrypoints:
        - type: http-proxy
          qos: AUTO
  endpointGroups:
    - name: Default HTTP proxy group
      type: http-proxy
      endpoints:
        - name: Default HTTP proxy
          type: http-proxy
          inheritConfiguration: false
          configuration:
            target: https://api.gravitee.io/echo
  flowExecution:
    mode: DEFAULT
    matchRequired: false
  plans:
    KeyLess:
      name: "Free plan"
      description: "This plan does not require any authentication"
      security:
        type: "KEY_LESS"
```

Notes:

- `metadata.name` must be RFC 1123 DNS-safe.
- The HRID is `namespace + "-" + name`, from which APIM derives a deterministic
  UUIDv3. Use it to correlate a CR with its APIM resource.
- Never delete the shared `dev-ctx` ManagementContext in a test.

---

## `terraform/main.tf`

```hcl
# <What this configuration represents, and which variable the scenario re-applies.>
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {}

variable "environment_id" {
  type    = string
  default = "DEFAULT"
}

variable "organization_id" {
  type    = string
  default = "DEFAULT"
}

variable "state" {
  type    = string
  default = "STARTED"
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "<journey-scoped-name>-tf"
  name            = "<journey-scoped-name>-tf"
  description     = "<same intent as the GKO arm>"
  version         = "1"
  type            = "PROXY"
  state           = var.state
  lifecycle_state = "PUBLISHED"
  visibility      = "PUBLIC"

  listeners = [
    {
      http = {
        type        = "HTTP"
        paths       = [{ path = "/<context-path>-tf/" }]
        entrypoints = [{ type = "http-proxy" }]
      }
    }
  ]

  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      endpoints = [
        {
          name                  = "default-endpoint"
          type                  = "http-proxy"
          inherit_configuration = false
          configuration         = jsonencode({ target = "https://api.gravitee.io/echo" })
        }
      ]
    }
  ]

  plans = [
    {
      hrid       = "keyless"
      name       = "Free plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security   = { type = "KEY_LESS" }
    }
  ]
}

# Output contract: the provisioner resolves roles to APIM ids by these names.
output "api_id" {
  value = apim_apiv4.api.id
}

output "api_context_path" {
  value = "/<context-path>-tf"
}
```

Default role to output mapping: `api` → `api_id`, `subscription` → `sub_id`,
`application` → `app_id`, `group` → `group_id`, context path →
`api_context_path`. Override with the scenario's `outputs` /
`contextPathOutput`.

The TF arm's resource names carry a `-tf` suffix so the two arms never collide
in APIM's shared name space.

---

## `params.ts` (only when the drivers diverge structurally)

One shared param type plus the per-provisioner closures. Full example:
`api-consumer/subscribe-and-call/params.ts`.

```ts
<LICENCE>

/**
 * Parameter binding for <journey>. ONE shared param type drives every
 * provisioner; the per-provisioner closures translate it into the GKO manifest
 * and the Terraform tfvars. This is the scenario-specific seam the provisioner
 * core deliberately does not own.
 */

import type { KubectlEngine } from "../../../../../src/provisioners/index.js";

export interface MyParams {
  /** … */
}

/**
 * Per-process suffix so re-runs get fresh values. APIM enforces uniqueness on
 * some fields (api-key values per API) across active AND revoked state, and the
 * local MongoDB persists across cluster restarts, so hardcoded values collide on
 * a second run.
 */
export const RUN_ID = `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`;

export function gkoApply(/* … */) {
  return async (kubectl: KubectlEngine, params: MyParams): Promise<void> => {
    await kubectl.applyString(/* built manifest */);
  };
}

export function toVars(params: MyParams): Record<string, unknown> {
  return { /* tfvars */ };
}
```

---

## `README.md`

```markdown
# <Journey title>

**As an <persona>, I <story>.**

<Two or three sentences: what the journey stands up, what it asserts, and any
stage that exists to catch a specific regression.>

| Arm | Fixtures | Xray |
|---|---|---|
| GKO | [`gko/`](./gko/) | @GKO-NNNN |
| Terraform | [`terraform/`](./terraform/) | @GKO-MMMM |

## Running it

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-NNNN
npm --prefix test/platform-test run e2e -- --grep @GKO-MMMM --provision-with terraform
```
```
