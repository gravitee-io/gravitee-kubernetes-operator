# Deploy Gravitee Kubernetes Operator Resources using Helm

The purpose of this guide is to give you an overview of how to use [Helm](https://helm.sh/docs/) templates capabilities to reuse properties across you Gravitee Kubernetes Operator resources.

## Disclaimer

The [sources](./src) coming with the guide are provided as an example and are not intended to fit with any 
requirement. 

## Context

In this example, we suppose that we have a `dev` and `staging` environment, each deployed in a Kubernetes namespace, with a set of APIs sharing a given Gravitee APIM [environment](https://docs.gravitee.io/apim/3.x/apim_adminguide_organizations_and_environments.html) and a set of policies enforced for this environment.

For some reason, we want any API deployed in the `dev` namespace to be accessible using a [keyless plan](https://docs.gravitee.io/apim/3.x/apim_publisherguide_plan_security.html#keyless_plans) while any API deployed in the `staging` namespace will be accessible using an [API key](https://docs.gravitee.io/apim/3.x/apim_publisherguide_plan_security.html#api_key_plans) plan.

For the same unknown reason, any API deployed in the `dev` namespace will enforce caching with a TTL of 1 minute, while 
APIs deployed in the `staging` namespace will rate limit requests to 100 requests per seconds.

We will use the [Bored API](https://www.boredapi.com/) and the [CatFact API](https://catfact.ninja/) as targets for our API proxy definitions.

## Project structure

```bash
.
├── templates
└── values
    ├── dev
    └── staging
```

The `templates` directory contains the templates used to reuse the values defined in the `values` directory.

The values directory is split into two sub directories:
  - The `dev` directory will be used to deploy our APIs in the dev namespace.
  - The `staging` directory will be used to deploy our APIs in the staging namespace.

In each of those directories, a `common.yaml` files defines values to be reused across each APIs deployed with the release, and the `apis.yaml` file contains a simple set of specific properties for each API (`name`, `version`, and `proxy`).

Here is, for instance, the structure of the `dev` values directory.

```bash
values/dev
├── apis.yaml
└── common.yaml
```

## Prerequisites

The Gravitee Kubernetes Operator resources allow you to simply source sensitive values from an existing secret (more on that later).

For this reason, we need to create one secret for each of our namespace containing this token.

```bash
kubectl create namespace dev
kubectl create secret generic gravitee-secrets \
  --from-literal=token=${DEV_TOKEN} \
  --namespace dev

kubectl create namespace staging
kubectl create secret generic gravitee-secrets \
  --from-literal=token=${STAGING_TOKEN} \
  --namespace staging
```

## The API definitions values

The `apis.yaml` file for our `dev` environment contains the following values definition.

```yaml
apis:
- name: bored
  version: "1.0.0"
  proxy:
    groups:
      - endpoints:
          - target: "https://www.boredapi.com/api/activity"
- name: catfact
  version: "1.0.0"
  proxy:
    groups:
      - endpoints:
          - target: "https://catfact.ninja/fact"
```

## Common values

Here is a rough overview of the `common.yaml` file defined for our `dev` environment.

```yaml
api:
  plans:
  - # A set of plans to be reused across our APIs
  flows:
  - # A set of flows to be reused across our APIs

# The management context referring to our dev APIM instance / environment.
context:
  name: dev
  baseUrl: http://localhost:9000
  environmentId: DEFAULT
  organizationId: DEFAULT
  token: "[[ secret `gravitee-secrets/token` ]]"

resources:
  - # A set of resources to be reused across our APIs.
```

💡 The context token will be sourced from the `gravitee-secrets` we created earlier in the `dev` namespace using the `token` key. Because our secret name contains an hyphen, we *must* wrap our secret reference with backticks in order to avoid any error in the Operator when interpreting the secret reference.

## The API definition template

Without getting into too much details, here is an overview of how we can build an API definition resource that will play nicely with set of values.

```yaml
{{- $common := .Values }}
{{- $ns := .Release.Namespace}}
{{- range $base := $common.apis }}
{{- $api := deepCopy $common.api | merge $base}}
---
apiVersion: "gravitee.io/v1alpha1"
kind: "ApiDefinition"
metadata:
    name: {{ regexReplaceAll "\\W+" $api.name "-" | lower }}
spec:
    name: {{ $api.name }}
    version: "{{ $api.version }}"
    # [...]
{{- end}}
```

The whole API definition is wrapped inside a [range](https://helm.sh/docs/chart_template_guide/control_structures/#looping-with-the-range-action) loop. Inside this loop, we can merge the API definition coming from our `apis.yaml` file with a set of common property defined in the `common.yaml` file defined for our environment (plans, flows, resources ...), giving precedence to what has been defined in the `apis` file.

💡 Because of Helm internals, we need to capture the global scope of our template values in order to be able to reuse those references inside the range loop. This is done using the following instructions.

```yaml
{{- $common := .Values }}
{{- $ns := .Release.Namespace}}
```

## Deploying all your namespace in a single command

With this setup, we can update a whole environment each time we need using this simple helm command.

```sh
NS=dev helm upgrade --install "$NS" . -f "values/$NS/common.yaml" -f "values/$NS/apis.yaml" -n "$NS"
```
