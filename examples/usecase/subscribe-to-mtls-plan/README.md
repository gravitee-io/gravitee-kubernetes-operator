# Subscribing to an MTLS plan

This guide shows how to configure and create an MTLS plan for a v4 API using the operator, then how to subscribe an application to that plan using the subscription resource.

The guide assumes that you have an operator and an APIM instance up and running.

You can read more about MTLS plans in the APIM [documentation](https://documentation.gravitee.io/apim/using-the-product/managing-your-apis/preparing-apis-for-subscribers/plans/mtls).

## Important consideration regarding APIM configuration

To be able to use MTLS subscriptions, your APIM instance has to be configured to enable TLS and client authentication.

Here are, for example, the Helm values we used to write this guide:

```yaml
servers:
    - type: http
      port: 8082
      service:
        type: NodePort
        nodePort: 30082
        externalPort: 82
    - type: http
      port: 8084
      service:
        type: NodePort
        nodePort: 30084
        externalPort: 84
      ssl:
        keystore:
          type: pem
          secret: secret://kubernetes/tls-server
        clientAuth: request
  service:
    type: NodePort
```

Here, both HTTP and HTTPS are enabled respectively on node ports 30082 and 30084 and TLS is enabled with a self signed keystore and on demand client authentication.

Server keystore is sourced from a kubernetes secret that you can create like that:

```sh
kubectl create secret tls tls-server --cert=pki//server.crt --key=pki/server.key
```

Note that for this to work you need to enabled kubernetes secrets at the root of your values file:

```yaml
secrets:
  kubernetes:
    enabled: true
```

## Creating a TLS secret

The client key and certificate we will use can be found in the [pki](pki/) directory. Let's create a secret to store them so that we can configure our application later on using the client certificate:

```sh
kubectl create secret tls tls-client --cert=pki/client.crt --key=pki/client.key --dry-run=client -o yaml>resources/tls-client.yml
```

## Configuring a MTLS plan

The API definition can be found [here](resources/api.yml).

Here is extracted the plan configuration.

```yaml
plans:
    MTLS:
      name: "mtls"
      security:
        type: "MTLS"
```

## The application resource

The application resource can be found [here](resources/application.yml)

The application makes use of string templating to inject the certificate we previously stored in a secret. The certificate *must* be in a valid PEM format.

## The subscription resource

A valid subscription must:
  - reference a valid API by its name and an optional namespace
  - reference a valid plan key defined in the API
  - reference a valid application by its name and an optional namespace

```yaml
apiVersion: gravitee.io/v1alpha1
kind: Subscription
metadata:
  name: echo-client-subscription
spec:
  api:
    name: mtls-demo
  application: 
    name: echo-client
  plan: MTLS
```

> If the `api` property points to a v2 API, this must be explicitly stated by adding a kind
> property with the `ApiDefinition` value to your api reference.

## Applying the resources

Only resources holding a management context ref are supported at the moment, so let's create this first and then the rest of the resources we described in order.

> The management context must be configured according to your setup, using your management API URL and credentials.

```sh
kubectl apply -f resources/management-context.yml
kubectl apply -f resources/tls-client.yml
kubectl apply -f resources/api.yml
kubectl apply -f resources/application.yml
kubectl apply -f resources/subscription.yml
```

## Calling your API

Now that all the resources are created, calling the API on the HTTPS endpoint of the gateway without client authentication should result in a 401 response.

```sh
❯ export GW_URL=https://localhost:30084
```

```sh
❯ curl -ksi $GW_URL/mtls-demo| head -1                                               
HTTP/1.1 401 Unauthorized
```

Let's now use the key and certificate we used to configure the echo-client application:

```sh
❯ curl -ksi --cert pki/client.crt --key pki/client.key $GW_URL/mtls-demo| head -1    
HTTP/1.1 200 OK
```


