# User Guide

## Overview

> TODO Why ? How ? What ?

## Quick Start

> TODO Install -> Deploy Gateway -> Create API -> Call

## Creating an API using the API Definition Resource

> TODO Create an API without management context and call + explanation of the How ?

## Synchronizing your API CRDs with an existing management API

The following example assumes that a management API has been deployed in a namespace called `apim-blue`, and connects to the management API using the default in memory credentials. 

### Creating a Management Context

```shell
echo """
apiVersion: gravitee.io/v1alpha1
kind: ManagementContext
metadata:
  name: apim-blue-context
  namespace: apim-blue
spec:
  baseUrl: http://apim-blue-apim3-api.apim-blue.svc:83
  environmentId: DEFAULT
  organizationId: DEFAULT
  auth:
    credentials:
      username: admin
      password: admin
""" | kubectl apply -f -
```

### Create an API referencing your context

```shell
echo """
apiVersion: gravitee.io/v1alpha1
kind: ApiDefinition
metadata:
  name: basic-api-example
spec:
  name: gko-example
  contextRef: 
    name: apim-blue-context
    namespace: apim-blue
  version: "1.0.0"
  description: "Basic api managed by Gravitee Kubernetes Operator"
  proxy:
    virtual_hosts:
      - path: "/k8s-basic"
    groups:
      - endpoints:
          - name: "Default"
            target: "https://api.gravitee.io/echo"
""" | kubectl apply -f -
```

### Update your API

> TODO (not implemented)

### Delete your API

> TODO (not implemented)

## Installation

> TODO Detailed installation guide

## Reference

### API Reference

See [api/reference.md](api/reference.md)

### The Management Context Resource

To be able to synchronize CRDs with a remote [management API](https://docs.gravitee.io/apim/3.x/apim_overview_architecture.html), you need to create a Management Context refering to an existing [organization and environment](https://docs.gravitee.io/am/current/am_adminguide_organizations_and_environments.html).

You can create as much management contexts as you want, each one targeting a specific environment, defined in a specific organization of a management API instance

Management Context can use either basic authentication or a bearer token to authenticate to your management API instance.

> Note: If both credentials and bearerToken are defined in your custom resource, the basic auth credentials will take precedence

#### Example

The following custom resource refers to a management API instance exposed at `https://gravitee-api.acme.com` and targets the `dev` environment of the `acme` organization, with the `admin` account, using basic auth.

```yaml
apiVersion: gravitee.io/v1alpha1
kind: ManagementContext
metadata:
  name: apim-blue-context
spec:
  baseUrl: https://gravitee-api.acme.com
  environmentId: dev 
  organizationId: acme
  auth:
    credentials:
      username: admin
      password: 406185a0-9adb-4097-b0bd-eb5cf13a7c6e
```

If you want to target another environment on the same API instance, just add another 
management context targeting this environment (e.g. staging)

This time, we use a bearerToken to authenticate the requests (the token must have been generated beforehand for the admin account)

```yaml
apiVersion: gravitee.io/v1alpha1
kind: ManagementContext
metadata:
  name: apim-blue-context
spec:
  baseUrl: https://gravitee-api.acme.com
  environmentId: staging 
  organizationId: acme
  auth:
    bearerToken: d70db517-a7fc-4fd5-924c-74ce1bfcf253
```

### The API Definition Resource

The APIDefinition CRD is (more or less) the yaml representation of an [API Definition](https://docs.gravitee.io/apim/3.x/apim_publisherguide_create_apis.html#import_an_existing_api_definition) in json format.

Here is a minimal example of an API Definition resource:

```yaml
apiVersion: gravitee.io/v1alpha1
kind: ApiDefinition
metadata:
  name: basic-api-example
spec:
  name: "GKO Basic"
  version: "1.1"
  description: "Basic api managed by Gravitee Kubernetes Operator"
  proxy:
    virtual_hosts:
      - path: "/k8s-basic"
    groups:
      - endpoints:
          - name: "Default"
            target: "https://api.gravitee.io/echo"
```
