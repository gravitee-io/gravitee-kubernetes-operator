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

| Name                         | Description                                                  | Value                            |
| ---------------------------- | ------------------------------------------------------------ | -------------------------------- |
| `rbacProxy.enabled`          | Specifies if the kube-rbac-proxy sidecar should be enabled.  | `true`                           |
| `rbacProxy.image.repository` | Specifies the docker registry and image name to use.         | `quay.io/brancz/kube-rbac-proxy` |
| `rbacProxy.image.tag`        | Specifies the docker image tag to use.                       | `v0.14.3`                        |

### Controller Manager

This is where you can configure the deployment itself and the way the operator will interact with APIM and Custom Resources in your cluster.

| Name                                        | Description                                                                                  | Value                            |
| ------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------- |
| `manager.image.repository`                  | Specifies the docker registry and image name to use.                                         | `graviteeio/kubernetes-operator` |
| `manager.image.tag`                         | Specifies the docker image tag to use.                                                       | `latest`                         |
| `manager.logs.json`                         | Whether to output manager logs in JSON format.                                               | `true`                           |
| `manager.configMap.name`                    | The name of the config map used to set the manager config from this values.                  | `gko-config`                     |
| `manager.scope.cluster`                     | Use false to listen only in the release namespace.                                           | `true`                           |
| `manager.applyCRDs`                         | If true, the manager will patch Custom Resource Definitions on startup.                      | `true`                           |
| `manager.metrics.enabled`                   | If true, a metrics server will be created so that metrics can be scraped using prometheus.   | `true`                           |
| `manager.httpClient.insecureSkipCertVerify` | If true, the manager HTTP client will not verify the certificate used by the Management API. | `false`                          |

### Ingress

Configure the behavior of the ingress controller.

When storing templates stored in config maps, the config map should contain a content key and a contentType key e.g.
```yaml
content: '{ "message": "Not Found" }'
contentType: application/json
```

| Name                              | Description                                                                      | Value |
| --------------------------------- | -------------------------------------------------------------------------------- | ----- |
| `ingress.templates.404.name`      | Name of the config map storing the HTTP 404 ingress response template.           | `""`  |
| `ingress.templates.404.namespace` | Namespace of the config map storing the HTTP 404 ingress response template.      | `""`  |

### Cert Manager

This section allows you to enable and configure the cert-manager dependency.
cert-manager is required to enabled webhook conversions and validation needed by the operator.

⚠️ cert-manager manages non-namespaced resources in your cluster 
and care must be taken to ensure that it is installed exactly once.
Enabling the cert-manager dependency will tie the lifecycle of cert-manager to the operator.
This property is essentially available for testing facility purposes.
When deploying in production, it is recommended that you install cert-manager as a separate component.

If enabling the dependency, please note that the namespace defined to install cert-manager
must have been created beforehand.

See https://cert-manager.io/docs/installation/helm

| Name                     | Description                                                  | Value          |
| ------------------------ | ------------------------------------------------------------ | -------------- |
| `cert-manager.enabled`   | If true, cert-manager will be installed as a dependency.     | `false`        |
| `cert-manager.namespace` | Defines the namespace where cert-manager will be installed.  | `cert-manager` |

### HTTP Client

👎 This section is deprecated and will be removed in version 1.0.0 The httpClient property
should now be set under the manager section instead.

| Name                                | Description                                   | Value   |
| ----------------------------------- | --------------------------------------------- | ------- |
| `httpClient.insecureSkipCertVerify` | see manager.httpClient.insecureSkipCertVerify | `false` |
