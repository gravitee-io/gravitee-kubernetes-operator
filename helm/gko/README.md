# G.K.O.

The Gravitee Kubernetes Operator Helm Chart

![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![Version: 0.8.4](https://img.shields.io/badge/Version-0.8.4-informational?style=flat-square) ![AppVersion: 0.8.4](https://img.shields.io/badge/AppVersion-0.8.4-informational?style=flat-square)

## Installing the Chart

To install the chart with the release name `graviteeio-gko`:

```console
$ helm repo add graviteeio https://helm.gravitee.io
$ helm install graviteeio-gko graviteeio/gko
```

## Requirements

Kubernetes: `>=1.14.0-0`

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| httpClient.insecureSkipCertVerify | bool | `false` | If true, the manager HTTP client will not verify the certificate used by the Management API. |
| ingress.templates.404.name | string | `""` | name of the config map storing the HTTP 404 ingress response template. A default template is used if this entry is empty. The config map should contain a content key and a contentType key. The default template is used if one of the key is missing. |
| ingress.templates.404.namespace | string | `""` | namespace of the config map storing the HTTP 404 ingress response template. A default template is used if this entry is empty. The config map should contain a content key and a contentType key. The default template is used if one of the key is missing.        |
| manager.applyCRDs | bool | `true` | If true, the manager will apply Custom Resource Definitions on startup. Please be aware that this will apply to Custom Resource Definitions  (which are the Open API model for Custom Resources such as API Definitions),  not to Custom Resources themselves. Custom Resources will be reconciled if the manager restarts whatever the value of this flag is. Because helm upgrades do not update CRDs once they have been installed, it is recommended to set this flag to true. |
| manager.configMap.name | string | `"gko-config"` | The name of the config map used to set the manager config from this values. |
| manager.image.repository | string | `"graviteeio/kubernetes-operator"` | Specifies the docker registry and image name to use. |
| manager.image.tag | string | `"latest"` | Specifies the docker image tag to use. |
| manager.logs.json | bool | `true` | Whether to output manager logs in JSON format. |
| manager.metrics.enabled | bool | `true` |  |
| manager.scope.cluster | bool | `true` | If true, the manager listens to resources created in the whole cluster. Use false to listen only in the release namespace. |
| rbac.create | bool | `true` | Specifies if RBAC resources should be created. |
| rbac.skipClusterRoles | bool | `false` | Specifies if cluster roles should be created when RBAC resources are created. |
| rbacProxy.enabled | bool | `true` | Specifies if the kube-rbac-proxy sidecar should be enabled.  Note that if this is disabled, the prometheus metrics endpoint will be exposed with no access control at all. |
| rbacProxy.image.repository | string | `"quay.io/brancz/kube-rbac-proxy"` | Specifies the docker registry and image name to use. |
| rbacProxy.image.tag | string | `"v0.14.3"` | Specifies the docker image tag to use. |
| serviceAccount.create | bool | `true` | Specifies if a service account should be created for the manager pod. |
| serviceAccount.name | string | `"gko-controller-manager"` | Specifies the service account name to use. |

