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

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| httpClient.insecureSkipCertVerify | bool | `false` | If true, the manager HTTP client will not verify the certificate used by the Management API. |
| ingress.templates.404.name | string | `""` | name of the config map storing the HTTP 404 ingress response template. A default template is used if this entry is empty. The config map should contain a content key and a contentType key. The default template is used if one of the key is missing. |
| ingress.templates.404.namespace | string | `""` | namespace of the config map storing the HTTP 404 ingress response template. A default template is used if this entry is empty. The config map should contain a content key and a contentType key. The default template is used if one of the key is missing.        |
| manager.applyCRDs | bool | `true` | If true, the manager will apply custom resource definitions on startup. |
| manager.configMap.name | string | `"gko-config"` | The name of the config map used to set the manager config from this values. |
| manager.image.repository | string | `"graviteeio/kubernetes-operator"` | Specifies the docker registry and image name to use. |
| manager.image.tag | string | `"latest"` | Specifies the docker image tag to use. |
| manager.logs.json | bool | `true` | Whether to output manager logs in JSON format. |
| manager.scope.cluster | bool | `true` | If true, the manager listens to resources created in the whole cluster. Use false to listen only in the release namespace. |
