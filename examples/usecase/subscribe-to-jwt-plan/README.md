# Subscribing to a JWT plan

This guide shows how to configure and create a JWT plan for a v4 API using the operator, then how to subscribe an application to that plan using the subscription resource.

The guide assumes that you have an operator and an APIM instance up and running.

## Generating the key pair

The plan will be configured using a hardcoded public key.

The key pair has been generated using openssl as follows:

```sh
ssh-keygen -t rsa -b 4096 -m PEM -f pki/private.key
openssl rsa -in jwt-demo.key -pubout -outform PEM -out pki/public.key
```

## Storing the public key in a secret

We will make use of the string templating capabilities of the operator to source the public key from a secret, so let's create the secret file first:

```sh
kubectl create secret generic jwt --from-file=pki/public.key --dry-run=client -o yaml| grep -v creationTimestamp > resources/jwt-key.yml
```

## Configuring the JWT plan

The API definition can be found [here](resources/api.yml)

Below is the extracted plan configuration, with the public key copied from the [pki](pki/) directory.

```yaml
plans:
  JWT:
    name: "jwt"
    security:
      type: "JWT"
      configuration:
        signature: "RSA_RS256"
        publicKeyResolver: "GIVEN_KEY"
        resolverParameter: "[[ secret `jwt/public.key` ]]"
        userClaim: "sub"
        clientIdClaim: "client_id"
    status: "PUBLISHED"
```

## The application resource

The application resource can be found [here](resources/application.yml)

> The client_id field defined in the application settings is mandatory to consume a JWT plan and must match the client_id claim of the token.

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
    name: jwt-demo
  application:
    name: echo-client
  plan: JWT
```

> If your API reference points to a v2 API, this must be explicitly stated by adding the `kind`
> property with the `ApiDefinition` value to your API reference.

## Applying the resources

Currently, only resources with a management context reference are supported. Create this first, followed by the other resources in the described order.

> The management context must be configured according to your setup, using your management API URL and credentials.

```sh
kubectl apply -f resources/management-context.yml
kubectl apply -f resources/jwt-key.yml
kubectl apply -f resources/api.yml
kubectl apply -f resources/application.yml
kubectl apply -f resources/subscription.yml
```

## Getting a token

You can forge a token using the [jwt.io](https://jwt.io) debugger.

Set the algorithm to `RS256` and sign your token with the provided keys and the following claims:

```json
{
  "sub": "echo-client",
  "client_id": "echo-client",
  "iat": 1516239022
}
```

If you are following this guide on macOS or Linux, you can get a token by running the [get_token.sh](pki/get_token.sh) bash script located in the pki directory.

```sh
export TOKEN=$(bash pki/get_token.sh)
```

## Calling your API

You can now use your token to call your API

```sh
GW_URL=http://localhost:30082 # replace with your gateway URL
curl -H "Authorization: Bearer $TOKEN" "$GW_URL/jwt-demo"
```

## Closing the subscription

Deleting the subscription resource results in the subscription being closed. Which means the client id associated with your token will be rejected with a 401 status on subsequent calls to the gateway.

```sh
kubectl delete -f resources/subscription.yml
```
