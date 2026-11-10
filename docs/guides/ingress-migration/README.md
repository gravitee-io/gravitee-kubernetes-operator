# Migrating from the Ingress controller to the Gateway API

The GKO Ingress controller is deprecated. It turns each Gravitee `Ingress` into a v2 `ApiDefinition`, and v2 APIs will
be removed in October 2027.

| Version | Ingress controller |
|---|---|
| Today | Enabled by default, deprecated: GKO records a `Deprecated` warning event on each Gravitee Ingress |
| 4.15 | Disabled by default. Set `ingress.controller.enabled: true` to keep it while you migrate |
| October 2027 | Removed, along with v2 APIs |

This guide maps each Ingress feature to the Gateway API (`Gateway` + `HTTPRoute`), lists what has no equivalent, and
gives a cutover procedure.

## What changes in the architecture

With the Ingress controller, **your APIM gateway serves the APIs**. GKO generates a v2 `ApiDefinition` per Ingress.
Your gateway loads it from a ConfigMap (Kubernetes sync), or from APIM when the API template has a `contextRef`.

With the Gateway API, **GKO deploys the gateway**. For each `Gateway`, GKO creates a Gravitee gateway: a Deployment,
a Service and their configuration. That gateway has no database and no Management API: it only serves the APIs
generated from the `HTTPRoute`s attached to it. Each `HTTPRoute` becomes an `ApiV4Definition` with one keyless plan.

```
Ingress:      Ingress ──▶ ApiDefinition (v2) ──▶ your APIM gateway ──▶ Service
Gateway API:  HTTPRoute ──▶ ApiV4Definition ──▶ gateway deployed by GKO for the Gateway ──▶ Service
```

## Choose your path

| Your Ingress | Path |
|---|---|
| Plain routing: no `gravitee.io/template`, or a template with only a keyless plan and no policy | [Route with the Gateway API](#route-with-the-gateway-api) |
| A template with plans (API key, JWT, OAuth2…), policies or a `contextRef` | [Keep API management on your APIM gateway](#keep-api-management-on-your-apim-gateway) |

An `HTTPRoute` cannot carry plans or policies: there is no equivalent to Ingress templates.

## Feature mapping

| Ingress | Gateway API |
|---|---|
| Ingress class `graviteeio` (`spec.ingressClassName` or `kubernetes.io/ingress.class`) | A `GatewayClass` with `controllerName: apim.gravitee.io/gateway`, a `Gateway`, and `parentRefs` on the route |
| `rules[].host` | `HTTPRoute.spec.hostnames`, intersected with the listener `hostname` |
| `pathType: Prefix` or `ImplementationSpecific` | `path.type: PathPrefix` |
| `pathType: Exact` | `path.type: Exact` |
| `backend.service.name` + `port.number` | `backendRefs` (`kind: Service`, `port`) |
| `spec.tls[].secretName` + the `kubernetes-pem-registry` ConfigMap | `Gateway` listener `tls.certificateRefs`. GKO manages the PEM registry of the gateway it deploys |
| `gravitee.io/template` (plans, policies, `contextRef`) | No equivalent, see [Keep API management on your APIM gateway](#keep-api-management-on-your-apim-gateway) |
| 404 response template (`ingress.templates.404`) | No equivalent |
| — | New: header, query parameter and method matches, weighted backends, request and response header modifiers, redirects, URL rewrites, backends in another namespace (with a `ReferenceGrant`) |

Unchanged: `spec.defaultBackend` and named Service ports are not supported by either.

## Route with the Gateway API

### 1. Enable the Gateway API controller

In the GKO Helm values:

```yaml
gatewayAPI:
  controller:
    enabled: true
```

With `gatewayAPI.applyCRDs: true` (the default), GKO applies the standard Gateway API CRDs it bundles (v1.6.1) at
every startup, with a forced server-side apply. If another Gateway API implementation already manages those CRDs in
your cluster (Istio, Cilium, Envoy Gateway…), set `gatewayAPI.applyCRDs: false` and keep the CRDs it installed, at a
compatible version.

The Gateway API controller is still flagged experimental in the chart. Check its supported features in the
[conformance report](../../../test/conformance/kubernetes.io/gateway-api/report) before you migrate.

### 2. Create the gateway

A `GatewayClassParameters` describes the gateway GKO deploys (image, replicas, Service type, license, extra
`gravitee.yml`), a `GatewayClass` points to it, and a `Gateway` declares the listeners:

```yaml
apiVersion: gravitee.io/v1alpha1
kind: GatewayClassParameters
metadata:
  name: gravitee-gateway
  namespace: default
spec:
  kubernetes:
    deployment:
      template:
        spec:
          containers:
            - name: gateway
              image: graviteeio/apim-gateway
---
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: gravitee-gateway
spec:
  controllerName: apim.gravitee.io/gateway
  parametersRef:
    group: gravitee.io
    kind: GatewayClassParameters
    name: gravitee-gateway
    namespace: default
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gravitee-gateway
  namespace: default
spec:
  gatewayClassName: gravitee-gateway
  listeners:
    - name: http
      port: 80
      protocol: HTTP
    - name: https
      port: 443
      protocol: HTTPS
      hostname: httpbin.example.com
      tls:
        certificateRefs:
          - kind: Secret
            name: httpbin.example.com
```

The `https` listener reuses the TLS Secret of the Ingress. A Secret in another namespace needs a `ReferenceGrant`.
Each TLS listener takes exactly one certificate; several TLS listeners on one port are told apart by SNI.

More examples, including `GatewayClassParameters` with a license: [examples/gateway-api](../../../examples/gateway-api).

### 3. Convert each Ingress to an HTTPRoute

Before:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: httpbin
  annotations:
    kubernetes.io/ingress.class: graviteeio
spec:
  tls:
    - hosts:
        - httpbin.example.com
      secretName: httpbin.example.com
  rules:
    - host: httpbin.example.com
      http:
        paths:
          - path: /secured
            pathType: Prefix
            backend:
              service:
                name: httpbin-1
                port:
                  number: 8080
```

After:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: httpbin
  namespace: default
spec:
  parentRefs:
    - name: gravitee-gateway
      sectionName: https
  hostnames:
    - httpbin.example.com
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /secured
      backendRefs:
        - name: httpbin-1
          port: 8080
```

TLS moves from the route to the `Gateway` listener. One Ingress rule with several paths becomes one `HTTPRoute` rule
per path, or one rule with several `matches` when they share a backend.

GKO generates an `ApiV4Definition` with the name and namespace of the `HTTPRoute`.

### 4. Check the route

```bash
kubectl get gateway gravitee-gateway
kubectl get httproute httpbin -o jsonpath='{.status.parents[*].conditions}'
kubectl get apiv4definition httpbin
```

The `Gateway` must be `Programmed`, and the route `Accepted` with `ResolvedRefs` true. The address of the gateway is in
`kubectl get gateway gravitee-gateway -o jsonpath='{.status.addresses}'`: call it with the route host before you
switch any traffic.

## Keep API management on your APIM gateway

When your Ingress template carries plans, policies or a `contextRef`, move the API itself to an `ApiV4Definition`
managed by GKO and served by your APIM gateway, as any other v4 API:

1. Write an `ApiV4Definition` with the plans, policies and `contextRef` of the template, and the paths and endpoints of
   the Ingress rules. See the [API reference](../../api/reference.md#apiv4definition).
2. Expose your APIM gateway the way you expose any service: an Ingress handled by another controller, a LoadBalancer
   Service, or the Gateway API.

To route to it with the GKO Gateway API, follow
[Using the Gateway API with a dedicated APIM gateway](../gateway-api-with-apim-gateway/README.md): an `HTTPRoute`
forwards to the APIM gateway Service, and an `ApiV4Definition` with the same name defines the API on it.

## Cut over

GKO deploys a new gateway with its own Service, so both setups can run side by side:

1. Create the `Gateway` and the `HTTPRoute`s, and check them on the new gateway address.
2. Point your DNS records to the new address.
3. Delete the Ingresses. Each one owns the v2 API generated for it, so Kubernetes deletes that API too, and GKO
   removes it from APIM when the template had a `contextRef`.
4. Delete the API templates (`gravitee.io/template: "true"`) once no Ingress uses them, and the `kubernetes-pem-registry`
   ConfigMap if nothing else reads it.
5. Disable the Ingress controller: `ingress.controller.enabled: false` (the default from 4.15).

After a restart, the operator logs no Ingress deprecation warning, and `kubectl get events` shows no new `Deprecated`
event on an Ingress.
