# Gravitee

## Table of Contents

| API channel  | Implementation version                    | Mode    | Report                                                 |
|--------------|-------------------------------------------|---------|--------------------------------------------------------|
| standard     | [version-4.8.5](https://github.com/gravitee-io/gravitee-kubernetes-operator/releases/tag/4.8.5) | default | [version-4.8.5 report](./standard-4.8.5-default-report.yaml) |


## Prerequisites

The following binaries are assumed to be present on your devide
  
  - [docker](https://docs.docker.com/get-started/get-docker/)
  - [kubectl](https://kubernetes.io/docs/tasks/tools/)
  - [kind](https://github.com/kubernetes-sigs/kind)
  - [go](https://go.dev/learn/)

The reproducer as been tested on macOS and Linux only.

## Reproducer

1. Clone the Gravitee Kubernetes Operator repository

```bash
git clone --depth 1 --branch 4.8.5 https://github.com/gravitee-io/gravitee-kubernetes-operator.git
```

2. Start a kind cluster

```bash
make start-conformance-cluster
```

3. Run a local load balancer service

> The make target runs [cloud-provider-kind](https://kind.sigs.k8s.io/docs/user/loadbalancer). If you are reproducing on a macOS device, the binary requires sudo privileges and your password will be asked. For Linux devices, cloud-provider-kind will be run using Docker compose.

```bash
make cloud-lb
```

1. Run the operator

```bash
make run
```

5. Install the Gravitee GatewayClass

```bash
kubectl apply -f ./test/conformance/gateway-class-parameters.report.yaml -f ./test/conformance/gateway-class.yaml
```

6. Run the conformance tests

```bash
make conformance
```

6. Print report

```bash
cat test/conformance/kubernetes.io/gateway-api/report/standard-4.8.5-default-report.yaml
```

