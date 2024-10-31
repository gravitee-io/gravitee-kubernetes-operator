# Subscribing to a JWT plan

This guide shows how to configure and create a JWT plan for a v4 API using the operator, then how to subscribe an application to that plan using the subscription resource.

The guide assumes that you have an operator and an APIM instance up and running.

## Generating the key pair

The plan will be configured using a hard coded public key.

The key pair has been generated using openssl as follow:

```sh
ssh-keygen -t rsa -b 4096 -m PEM -f pki/private.key
openssl rsa -in jwt-demo.key -pubout -outform PEM -out pki/public.key
```

## Configuring the JWT plan

The API definition can be found [here](resources/api.yml)

Here is extracted the plan configuration, with the public key copied from the [pki](pki/) directory.

```yaml
plans:
    JWT:
      name: "jwt"
      security:
        type: "JWT"
        configuration:
          signature: "RSA_RS256"
          publicKeyResolver: "GIVEN_KEY"
          resolverParameter: |
            -----BEGIN PUBLIC KEY-----
            MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAv8YGSPoQEl7lXnp8OHkb
            AOPYZ81rzXkmO83d0P8G78qWzi3gPnODm6Qxi2NbgcWXqQlZXxPkDTS3Xck1V3WY
            E9voqQE7UEwpFBolqtUHQqL4w2vr/eUtZv9t3DdtoCcIj4xLmJUw7PS7jAb9quq0
            XiVN692d6LI62T+9LyN+kcWHTpUyMBB8oxfQ9ekkGHskTc6LgYovKK+9lKoJv6gg
            0ge8YAFbpjJBZbVX3jV8qeszgw9Xdhs3w/S8QnvWa3Cv0+c47oxZjXwpAa8ARzfn
            D/5oK4CWRRy+t3QUndSR0cBR+bU0YFks3mmbl514/ywOXRf/sZmXaJkNejfNHQVa
            hJgj/Z3W3F8GKksuFF14+BK2KX30bsQL3e4SeN0Wv6DF1UloG0T396yDd/o7L3ZC
            DBlRB44OZ8sO3h8iSW7wVX0sGj/OKc4smo5dgP0r4+Fm2EVmVFU5YvEkFcy0Xoth
            QmLwq0lJc7BdRMpAfRZLbW5WSlb2jgvxA/VI/ScLTRWZI7DGbzHRBS6J8Rnt3Inq
            jo7mUV1juBs3RhpxdOmg1LpGLAtQdcSSnX3IyyEVbzTVb22Px0EGAlKzMs6bnTJf
            3TbZd/C0iqd6QOyaTh7D4Nr7ClfWAaYGZBA/FsHWA88fOsIQCtovWjp9A8i1+VQ5
            HEy1rpaHPGHt1DFt2hu+d3MCAwEAAQ==
            -----END PUBLIC KEY-----
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

> If you API reference point to a v2 API, this must be explicitly stated by adding the kind
> property with the `ApiDefinition` value to your api reference.

## Applying the resources

Only resources holding a management context ref are supported at the moment, so let's create this first and then the rest of the resources we described in order.

> The management context must be configured accordingly to your setup, using your management API URL and credentials.

```sh
kubectl apply -f resources/management-context.yml
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

If you are following this guide on macOs or linux, you can get a token by running the [get_token.sh](pki/get_token.sh) bash script located in the pki directory.

```sh
export TOKEN=$(bash pki/get_token.sh)
```

## Calling your API

You can now use your token to call your API

```sh
GW_URL=http://localhost:30082 # replace by your gateway URL
curl -H "Authorization: Bearer $TOKEN" "$GW_URL/jwt-demo"
```

