# How to test the examples

## Prerequisites

This guide assumes that the following binaries are installed on your device:

  * docker
  * helm
  * kubectl
  * kind
  * curl
  * go
  
Commands are expected to be executed from the root of this respository.

## Run a kind cluster

```sh
kind create cluster --config hack/kind/kind.conformance.yaml
```

> This will start Kind with ports `80`, `443` and `9092` bound to your host.

## Run kind cloud provider

Install the `cloud-provider-kind` binary by running the following command:

```sh
GOBIN=$(pwd)/bin go install sigs.k8s.io/cloud-provider-kind@latest 
```

Then run the following command:

```sh
sudo bin/cloud-provider-kind
```

> sudo is required to open ports on the system and to connect to the container runtime.

Once cloud-provider-kind is running, run the rest of the commands from a new shell.

## Install the Gravitee Kubernetes Operator

```sh
IMG=gko TAG=dev make docker-build \
    && kind load docker-image gko:dev --name gravitee \
    && helm upgrade --install gko helm/gko \
    --set manager.image.repository=gko \
    --set manager.image.tag=dev \
    --set manager.metrics.enabled=false \
    --set gatewayAPI.controller.enabled=true
```

> This will build the head of the repository and deploy the operator using Helm. The required resources to enable gateway-api support on your cluster will be installed by the operator.

## Install and configure cert-manager

```sh
helm repo add jetstack https://charts.jetstack.io
helm repo update jetstack
helm upgrade --install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.17.0 \
  --set crds.enabled=true \
  --set config.kind="ControllerConfiguration" \
  --set config.enableGatewayAPI=true
```

> This will enable gateway API support for cert-manager. 

Create a self-signed cluster issuer using the following command:

```sh
kubectl apply -f examples/gateway-api/cert-manager-cluster-issuer.yaml
```

## Create a Gravitee license secret

Checkout our [website](https://www.gravitee.io/try-gravitee) to get a license.

> This step is required to test the KafkaRoute ressource. Your license file MUST be named `license.key`.

```sh
kubectl create secret generic gravitee-license --from-file=license.key
```

## Apply the gravitee.io GatewayClassParameters resource

```sh
kubectl apply -f examples/gateway-api/gateway-class-parameters-with-license.yaml
```

You can check if the license secret has been resolved in the gateway class status conditions by running the following command:

```sh
kubectl get gatewayclassparameters -o jsonpath='{"LicenseFound="}{.items[*].status.conditions[?(@.type == "ResolvedRefs")].status}{"\n"}' 
```

> This should output `LicenseFound=True`.

## Apply the kubernetes.io GatewayClass resource

```sh
kubectl apply -f examples/gateway-api/gateway-class.yaml
```

You can check that the gateway class has been accepted by running the following command:

```
kubectl get gatewayclass -o jsonpath='{"Accepted="}{.items[*].status.conditions[?(@.type == "Accepted")].status}{"\n"}' 
```

> This should output `Accepted=True`.

## Apply the kubernetes.io Gateway resource

```sh
kubectl apply -f examples/gateway-api/gateway.yaml
```

You can check that all gateway listeners have been programmed by running the following command:

```sh
kubectl get gateways -o jsonpath='{"Programmed="}{.items[*].status.conditions[?(@.type == "Programmed")].status}{"\n"}' 
```

> This should output `Programmed=True`.

Wait for the gateway pod to be in a ready state:

```sh
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=gravitee-gateway --timeout=300s
```

> Depending on network conditions, this might take time as the image needs to be pulled for the gateway deployment to be ready.

You can export the IP address of your gateway to an environment variable for further usage using the following command:

```sh
export GW_ADDR=$(kubectl get gateway gravitee-gateway -o jsonpath='{.status.addresses[0].value}')
echo "$GW_ADDR"
```

> This should output the IP address assigned to the LoadBalancer service of the gateway.

You can check the connectivity with each of the gateway listener by running the following commands:

```sh
nc -w 1 -vz "$GW_ADDR" 80
nc -w 1 -vz "$GW_ADDR" 443
nc -w 1 -vz "$GW_ADDR" 9092
```

## Get the CA certificate of the HTTPS server

```sh
kubectl get secret https-server -o json | jq '.data."ca.crt"' | tr -d '"' | base64 --decode > examples/gateway-api/tmp/https.ca.crt
```

> This certificates can be used to configure your clients to trust the HTTPS gateway listener.

## Bind route hostnames to the Gateway listeners IP address in your /etc/hosts file

```sh
cp /etc/hosts examples/gateway-api/tmp/hosts
sudo -- sh -c "echo $GW_ADDR demo.apis.example.dev demo.kafka.example.dev broker-0-demo.kafka.example.dev >> /etc/hosts"
```

You should be able to issue the following call using curl:

```sh
curl -i http://demo.apis.example.dev 
```

> At this point the call should result in an HTTP 404 status because no route as been created yet.

## Deploy httpbin backends

```sh
kubectl apply -f examples/gateway-api/http-backends.yaml 
```

> This will install three instances of httpbin with httpbin-(1|2) as a service name.

## Apply the kubernetes.io HTTPRoute resource

```sh
kubectl apply -f examples/gateway-api/http-route.yaml
```

You can check that the route has been accepted by running the following command

```sh
kubectl get httproutes -o jsonpath='{"Accepted="}{.items[*].status.parents[0].conditions[?(@.type == "Accepted")].status}{"\n"}'
```

> This should output `Accepted=True`.

## Test the HTTP route

There are two rules defined in the HTTPRoute. 

The first one can be tested by issuing the following curl call:

```sh
curl -i http://demo.apis.example.dev/bin/hostname
```

> This should result in the httpbin-1 pod hostname being shown as an output with an HTTP 200 status.

The second one can be tested by issuing the following curl call:

```sh
curl -i -H "env: canary" http://demo.apis.example.dev/bin/hostname
```

> This should result in the httpbin-2 pod hostname being shown as an output with an HTTP 200 status.

## Traffic splitting between multiple backends

You can update your HTTP route to add traffic splitting between the two backends by running the following command:

```sh
kubectl apply -f examples/gateway-api/http-route-with-traffic-splitting.yaml
```

> Traffic splitting is apply on the first rule of the route.

Then you should be able to reproduce what's described in the following [guide](https://gateway-api.sigs.k8s.io/guides/traffic-splitting/) by issuing the following call:

```sh
curl -i http://demo.apis.example.dev/bin/hostname
```

## Header modifier filter

You update your HTTP route test request and response header modifiers by running the following command:

```sh
kubectl apply -f examples/gateway-api/http-route-with-header-modifiers.yml
```

Then you can check that headers are modified on the request and on the response by issuing the following call:

```sh
curl -i -H "x-tag: acme.com" -H "x-impl: acme.com" -H "x-rm: true" http://demo.apis.example.dev/bin/headers
```

## HTTP Redirect

You can update your HTTP route to test request redirects by running the following command:

```sh
kubectl apply -f examples/gateway-api/http-route-with-request-redirects.yml
```

The example test several configurations.

### Replacing path prefix with host and scheme exlicitely defined

Issuing the following call:

```sh
curl -iL -H "x-rule-match: first" http://demo.apis.example.dev/bin/headers
```

Should redirect to `https://httpbin.org/anything/headers`

### Replacing full path with host and scheme exlicitely defined

Issuing the following call:

```sh
curl -iL  -H "x-rule-match: second" http://demo.apis.example.dev/bin/headers
```

Should redirect to `https://api.gravitee.io/echo`

### Replacing full path with host and scheme from request

Issuing the following call:

```sh
curl -iL  -H "x-rule-match: third" http://demo.apis.example.dev/bin/hostname
```

Should redirect to `http://demo.apis.example.dev/bin/404`

> In that case following the redirect leads to a 404 response status.

## Add a kafka listener to the gateway

You can update the gateway to add a kafka listener by running the following command:

```sh
kubectl apply -f examples/gateway-api/gateway-kafka.yaml
```

## Get the CA certificate of the Kafka server

```sh
kubectl get secret kafka-server -o json | jq '.data."ca.crt"' | tr -d '"' | base64 --decode > examples/gateway-api/tmp/kafka.ca.crt
```

> This will be used later to configure the kafka client.

## Start a Kafka cluster with Strimzi

```sh
kubectl create -f 'https://strimzi.io/install/latest?namespace=default'
```

> This will install the Strimzi kubernetes operator.

```sh
kubectl apply -f https://strimzi.io/examples/latest/kafka/kraft/kafka-single-node.yaml
kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s
```
> This will install a single node kafka cluster that will act as a backend for the Kafka route.

## Apply the gravitee.io KafkaRoute resource

```sh
kubectl apply -f examples/gateway-api/kafka-route.yaml
```

You can check that the route has been accepted by running the following command

```sh
kubectl get kafkaroutes -o jsonpath='{"Accepted="}{.items[*].status.parents[0].conditions[?(@.type == "Accepted")].status}{"\n"}'
```

> This should output `Accepted=True`.

## Install a kafka client

You can download and install a Kafka client by running the following commands:

```sh
mkdir -p examples/gateway-api/tmp/kafka-client 
export KAFKA_CLI_DL=https://dlcdn.apache.org/kafka/4.0.0/kafka_2.13-4.0.0.tgz
curl -s -L  $KAFKA_CLI_DL | tar xvz - -C examples/gateway-api/tmp/kafka-client --strip-components=1
```

## Test the Kafka route

In one shell, run the following commands and start producing messages by writing to the stdin:

```sh
export PRODUCE=examples/gateway-api/tmp/kafka-client/bin/kafka-console-producer.sh
export PROPS=examples/gateway-api/kafka-consumer.properties

$PRODUCE --bootstrap-server demo.kafka.example.dev:9092 \
    --topic demo \
    --producer.config $PROPS
```

In another shell, run the following commands to consume the messages:

```sh
export CONSUME=examples/gateway-api/tmp/kafka-client/bin/kafka-console-consumer.sh
export PROPS=examples/gateway-api/kafka-consumer.properties

$CONSUME --bootstrap-server demo.kafka.example.dev:9092 \
    --topic demo \
    --consumer.config $PROPS \
    --from-beginning
```

## Add access controls to the Kafka route

```sh
kubectl apply -f examples/gateway-api/kafka-route-with-acl.yaml
```

> This will configure the existing Kafka route, adding an ACL filter that grants read and write operations on the `demo` topic only, which means that any other topic will be forbidden if it is not added in the access control list.

For instance, run the following command to try and produce in a `prices` topic:

```sh
export PRODUCE=examples/gateway-api/tmp/kafka-client/bin/kafka-console-producer.sh
export PROPS=examples/gateway-api/kafka-consumer.properties

$PRODUCE --bootstrap-server demo.kafka.example.dev:9092 \
    --topic prices \
    --producer.config $PROPS
```

This should result in a warning containing the following message being printed to your stdout:

```
The metadata response from the cluster reported a recoverable issue with correlation id 13 : {prices=TOPIC_AUTHORIZATION_FAILED} 
```

## Cleaning up everything

```sh
pgrep cloud-provider-kind | sudo xargs kill
kind delete cluster --name gravitee 
[ -f examples/gateway-api/tmp/hosts ] \
    && sudo -- sh -c "cat examples/gateway-api/tmp/hosts >/etc/hosts"
find examples/gateway-api/tmp -not -name ".gitignore"  -delete
```
