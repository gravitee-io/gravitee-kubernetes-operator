# Scrapping operator metrics with the prometheus stack

## Start a kind  cluster


```sh
kind create cluster --config hack/kind/kind.yaml
```

## Install the prometheus stack


We will install the prometheus operator using the [
kube-prometheus-stack](https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack) provided by the prometheus community. This will install [Grafana](https://grafana.com/)  alongside the prometheus operator and allow to gather metrics exported by the operator metrics server. 


```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update prometheus-community
helm upgrade --install kube-prometheus-stack \
  --create-namespace \
  --namespace kube-prometheus-stack \
  prometheus-community/kube-prometheus-stack
```

> The prometheus stack will be installed in a dedicated `kube-prometheus-stack` namespace.  

## Install the operator

```sh
IMG=gko TAG=dev make docker-build \
    && kind load docker-image gko:dev --name gravitee \
    && helm upgrade --install gko helm/gko \
    -n gravitee --create-namespace \
    --set manager.image.repository=gko \
    --set manager.image.tag=dev \
    --set manager.scope.cluster=false \
    --set manager.metrics.enabled=true \
    --set manager.metrics.prometheus.instance.create=true \
    --set manager.metrics.prometheus.monitor.create=true \
    -n gravitee
```

> The last two values will deploy a dedicated prometheus statefulset in the operator namespace and create a service monitor to scrape metrics exported by the operator metrics server.

When metrics are enabled, the address of the prometheus endpoint is printed in the install notes.

The data is available using the following prometheus datasource URL:

`http://prometheus-operated.gravitee.svc:9090`

## Exploring metrics using Grafana

The `prometheus-kube-stack` install notes provide two command to access the Grafana UI from your host.

### Getting the Grafana admin password

```sh
kubectl --namespace kube-prometheus-stack get secrets kube-prometheus-stack-grafana -o jsonpath="{.data.admin-password}" | base64 -d ; echo
```

### Port forwarding the Grafana UI to your host

```sh
export POD_NAME=$(kubectl --namespace kube-prometheus-stack get pod -l "app.kubernetes.io/name=grafana,app.kubernetes.io/instance=kube-prometheus-stack" -oname)
kubectl --namespace kube-prometheus-stack port-forward $POD_NAME 3000
```

Once this is done, you can access the Grafana console at http://localhost:3000 and log in using the `admin` username and the password obtained using the previous commands.

### Adding the Operator prometheus instance as a datasource in Grafana.

In the left side menu, select `Data Sources` and add a new prometheus datasource using http://prometheus-operated.gravitee.svc:9090 as a Prometheus server URL.


