# Using the Gateway API with a Dedicated APIM Gateway

This guide describes how to use the Kubernetes Gateway API for ingress traffic routing while delegating API management concerns (policies, plans, subscriptions, analytics, etc.) to a dedicated Gravitee APIM gateway.

## Overview

In this architecture, two gateways work together:

- A **Gateway API gateway** deployed by the Gravitee Kubernetes Operator handles incoming
  traffic and routes it based on standard `HTTPRoute` resources.
- A **dedicated APIM gateway** deployed via Helm handles the heavy lifting: policy
  enforcement, plan-based access control, rate limiting, subscriptions, and analytics
  reporting.

```
                                        ┌─────────────────────┐
  Client ─── HTTPRoute ───▶  Gateway API Gateway  ───▶  APIM Gateway  ───▶  Backend
              (ingress)      (routing only)              (policies,         (httpbin, etc.)
                                                          plans, etc.)
                                        └─────────────────────┘
```

The HTTPRoute defines the ingress path and hostname, then forwards traffic to the APIM
gateway service. An `ApiV4Definition` with the same name defines the actual API behavior
(listeners, endpoints, plans, flows) on the APIM gateway.

## Prerequisites

- A Kubernetes cluster with the Gravitee Kubernetes Operator installed
- An APIM gateway deployed via the Gravitee Helm chart
- A `ManagementContext` resource configured to connect the operator to the APIM
  management API

## Configuration

### 1. Enable `skipAPIDefinition` in GKO Helm values

By default, the HTTPRoute reconciler creates an intermediate `ApiV4Definition` CR with the same name as the route. This conflicts with user-managed API definitions that share that name. Setting `skipAPIDefinition` to `true` makes the reconciler write a ConfigMap instead, leaving the CR name available for your own `ApiV4Definition`.

```yaml
gatewayAPI:
  controller:
    enabled: true
    skipAPIDefinition: true
```

### 2. Prevent the APIM gateway from loading Gateway API definitions

The ConfigMap created by the HTTPRoute reconciler contains an API definition tagged with
the Gateway resource name (e.g. `gravitee/gravitee-gateway`). If the APIM gateway has
Kubernetes sync enabled and no tag filtering, it will deploy this definition. Because the
API's backend points back to the APIM gateway service itself (from the HTTPRoute
`backendRef`), this creates a self-referencing loop that results in a 504 timeout.

#### Option A: Disable Kubernetes sync (recommended)

If the APIM gateway only needs to sync API definitions from the management API (which is the typical setup), disable Kubernetes sync entirely:

```yaml
# APIM Helm values
gateway:
  services:
    sync:
      kubernetes:
        enabled: false
```

This is the simplest and safest default for a dedicated APIM gateway that receives its
API definitions through a `ManagementContext`.

#### Option B: Use sharding tags

If the APIM gateway needs Kubernetes sync enabled for other purposes (e.g. other
ConfigMap-based API definitions), use sharding tags to exclude Gateway API definitions:

```yaml
# APIM Helm values
gateway:
  sharding_tags: "!gravitee/gravitee-gateway"
  services:
    sync:
      kubernetes:
        enabled: true
```

The `!` prefix tells the gateway to exclude APIs tagged with
`gravitee/gravitee-gateway`, while still allowing other ConfigMap-based definitions
through.

### 3. Deploy the Gateway API resources

Apply the `GatewayClass`, `GatewayClassParameters`, and `Gateway`:

```yaml
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1
metadata:
  name: gravitee-gateway
spec:
  controllerName: apim.gravitee.io/gateway
  parametersRef:
    kind: GatewayClassParameters
    group: gravitee.io
    name: gravitee-gateway
    namespace: gravitee
```

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gravitee-gateway
spec:
  gatewayClassName: gravitee-gateway
  listeners:
    - name: http
      port: 80
      protocol: HTTP
```

### 4. Define the HTTPRoute and ApiV4Definition together

Both resources share the same `metadata.name`. The HTTPRoute handles ingress routing to the APIM gateway, and the `ApiV4Definition` configures the API behavior on it.

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: echo
spec:
  parentRefs:
    - name: gravitee-gateway
      kind: Gateway
      group: gateway.networking.k8s.io
      namespace: gravitee
  hostnames:
    - echo.apis.example.dev
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /echo
      backendRefs:
        - kind: Service
          group: ""
          name: apim-gateway
          port: 82
---
apiVersion: gravitee.io/v1alpha1
kind: ApiV4Definition
metadata:
  name: echo
spec:
  contextRef:
    name: "dev-ctx"
  name: "echo"
  description: "Echo API managed by Gravitee Kubernetes Operator"
  version: "1.0"
  type: PROXY
  state: STARTED
  listeners:
    - type: HTTP
      paths:
        - path: "/echo/"
          host: echo.apis.example.dev
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
            target: http://httpbin-1.gravitee.svc.cluster.local:8080
          secondary: false
          sharedConfigurationOverride:
            http:
              propagateClientHost: false
  flowExecution:
    mode: DEFAULT
    matchRequired: false
  plans:
    KeyLess:
      name: "Free plan"
      description: "This plan does not require any authentication"
      security:
        type: "KEY_LESS"
  notifyMembers: false
```

Key points:

- The **HTTPRoute** `backendRef` targets the APIM gateway service (`apim-gateway:82`).
  The Gateway API gateway will forward matching requests there.
- The **ApiV4Definition** uses a `contextRef` to sync the API through the APIM management
  plane. The APIM gateway picks it up through its regular sync mechanism and handles
  policy enforcement, plan validation, and backend proxying.
- The listener `host` in the `ApiV4Definition` must match the HTTPRoute `hostname` so
  the APIM gateway can match incoming requests by virtual host.

## Request flow

1. A client sends `GET /echo/hostname` with `Host: echo.apis.example.dev`
2. The **Gateway API gateway** matches the request via the HTTPRoute and proxies it to
   `apim-gateway:82`
3. The **APIM gateway** matches the request against the `echo` API (by path and virtual
   host), applies plans and policies, then proxies to the backend
   (`httpbin-1.gravitee.svc.cluster.local:8080`)
4. The backend responds and the response flows back through both gateways to the client

