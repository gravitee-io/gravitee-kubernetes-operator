# G.K.O.

The Gravitee Kubernetes Operator Helm Chart

## Installing the Chart

To install the chart with the release name `gko`

```console
$ helm repo add graviteeio https://helm.gravitee.io
$ helm install gko graviteeio/gko
```

## Upgrading the Operator

Assuming that the repository as been aliased as graviteeio and that the release name is `gko`

```console
$ helm repo update graviteeio
$ helm upgrade --install gko graviteeio/gko
```

## Requirements

Kubernetes: `>=1.16.0-0`

## Parameters

### RBAC

Required RBAC resources are created by default for all components involved in the release.

| Name                    | Description                                                                   | Value                    |
| ----------------------- | ----------------------------------------------------------------------------- | ------------------------ |
| `serviceAccount.create` | Specifies if a service account should be created for the manager pod.         | `true`                   |
| `serviceAccount.name`   | Specifies the service account name to use.                                    | `gko-controller-manager` |
| `rbac.create`           | Specifies if RBAC resources should be created.                                | `true`                   |
| `rbac.skipClusterRoles` | Specifies if cluster roles should be created when RBAC resources are created. | `false`                  |

### RBAC Proxy

Kube RBAC Proxy is deployed as a sidecar container and restricts access to the prometheus metrics endpoint.

⚠️ If this is disabled, the prometheus metrics endpoint will be exposed with no access control at all.

| Name                         | Description                                                                       | Value                            |
| ---------------------------- | --------------------------------------------------------------------------------- | -------------------------------- |
| `rbacProxy.enabled`          | Specifies if the kube-rbac-proxy sidecar should be enabled.                       | `true`                           |
| `rbacProxy.image.repository` | Specifies the docker registry and image name to use.                              | `quay.io/brancz/kube-rbac-proxy` |
| `rbacProxy.image.tag`        | Specifies the docker image tag to use.                                            | `v0.18.1`                        |
| `rbacProxy.image.pullPolicy` | Specifies the pullPolicy to use when starting a new container                     | `IfNotPresent`                   |
| `rbacProxy.image.pullSecret` | Specifies the secret holding the credentials used to pull image from the registry | `{}`                             |

### Controller Manager

This is where you can configure the deployment itself and the way the operator will interact with APIM and Custom Resources in your cluster.

| Name                                                             | Description                                                                                                                                     | Value                            |
| ---------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------- |
| `manager.image.repository`                                       | Specifies the docker registry and image name to use.                                                                                            | `graviteeio/kubernetes-operator` |
| `manager.image.tag`                                              | Specifies the docker image tag to use. If no value is set, the chart version will be used.                                                      | `""`                             |
| `manager.image.pullPolicy`                                       | Specifies the pullPolicy to use when starting a new container                                                                                   | `IfNotPresent`                   |
| `manager.image.pullSecret`                                       | Specifies the secret holding the credentials used to pull image from the registry                                                               | `{}`                             |
| `manager.logs.json`                                              | Whether to output manager logs in JSON format.                                                                                                  | `true`                           |
| `manager.configMap.name`                                         | The name of the config map used to set the manager config from this values.                                                                     | `gko-config`                     |
| `manager.resources.limits.cpu`                                   | The CPU resources limits for the GKO Manager container                                                                                          | `500m`                           |
| `manager.resources.limits.memory`                                | The memory resources limits for the GKO Manager container                                                                                       | `256Mi`                          |
| `manager.resources.requests.cpu`                                 | The requested CPU for the GKO Manager container                                                                                                 | `500m`                           |
| `manager.resources.requests.memory`                              | The requested memory for the GKO Manager container                                                                                              | `128Mi`                          |
| `manager.tolerations`                                            | Set pods tolerations. Please see https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/ about this topic.                | `[]`                             |
| `manager.scope.cluster`                                          | Use false to listen only in the release namespace.                                                                                              | `true`                           |
| `manager.scope.namespaces`                                       | Setting this list of namespaces will result in the operator listening only in this namespaces.                                                  | `[]`                             |
| `manager.applyCRDs`                                              | 👎 This feature is deprecated and will be replaced in a future release. If true, the manager will patch Custom Resource Definitions on startup. | `true`                           |
| `manager.metrics.enabled`                                        | If true, a metrics server will be created so that metrics can be scraped using prometheus.                                                      | `true`                           |
| `manager.httpClient.insecureSkipCertVerify`                      | If true, the manager HTTP client will not verify the certificate used by the Management API.                                                    | `false`                          |
| `manager.httpClient.timeoutSeconds`                              | he timeout (in seconds) used when issuing request to the Management API.                                                                        | `10`                             |
| `manager.webhook.enabled`                                        | If true, the manager will register a webhook server operating on custom resources.                                                              | `true`                           |
| `manager.webhook.service.name`                                   | The service used to expose the webhook server.                                                                                                  | `gko-webhook`                    |
| `manager.webhook.service.port`                                   | Which port the webhook server will listen to.                                                                                                   | `9443`                           |
| `manager.webhook.cert.create`                                    | If true, a secret will be created to store the webhook server certificate.                                                                      | `true`                           |
| `manager.webhook.cert.secret.name`                               | The name of the secret storing the webhook server certificate.                                                                                  | `gko-webhook-cert`               |
| `manager.webhook.admission.checkApiContextPathConflictInCluster` | check if the same API context path exists in the whole cluster.                                                                                 | `false`                          |

### ingress

Configure the behavior of the ingress controller.

When storing templates stored in config maps, the config map should contain a content key and a contentType key e.g.
```yaml
content: '{ "message": "Not Found" }'
contentType: application/json
```

| Name                              | Description                                                                      | Value            |
| --------------------------------- | -------------------------------------------------------------------------------- | ---------------- |
| `ingress.controller.enabled`      | indicates if GKO ingress controller is enabled or not                            | `true`           |
| `ingress.ingressClasses`          | list of ingress classes that the gateway will handle.                            | `["graviteeio"]` |
| `ingress.templates.404.name`      | Name of the config map storing the HTTP 404 ingress response template.           | `""`             |
| `ingress.templates.404.namespace` | Namespace of the config map storing the HTTP 404 ingress response template.      | `""`             |

### HTTP Client

👎 This section is deprecated and will be removed in version 1.0.0 The httpClient property
should now be set under the manager section instead.

| Name                                | Description                                   | Value   |
| ----------------------------------- | --------------------------------------------- | ------- |
| `httpClient.insecureSkipCertVerify` | see manager.httpClient.insecureSkipCertVerify | `false` |
