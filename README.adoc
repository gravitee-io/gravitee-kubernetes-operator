= Gravitee.io Kubernetes Operator

image:https://img.shields.io/badge/License-Apache%202.0-blue.svg["License", link="https://github.com/gravitee-io/gravitee-kubernetes-operator/blob/master/LICENSE.txt"]
image:https://circleci.com/gh/gravitee-io/gravitee-kubernetes-operator/tree/alpha.svg?style=svg[link="https://app.circleci.com/pipelines/github/gravitee-io/gravitee-kubernetes-operator?branch=alpha"]
image:https://img.shields.io/badge/semantic--release-📦🚀-e10079?logo=semantic-release["semantic-release: 📦🚀", link="https://github.com/semantic-release/semantic-release"]
image:https://goreportcard.com/badge/github.com/gravitee-io/gravitee-kubernetes-operator["Go Report Card", link="https://goreportcard.com/report/github.com/gravitee-io/gravitee-kubernetes-operator"]

image:./.assets/gravitee-logo-cyan.svg["Gravitee.io",400]

== Introduction

APIM 3.19.0 has introduced the first release of the Gravitee Kubernetes Operator (GKO).

You can use the GKO to define, deploy, and publish APIs to your API Portal and API Gateway  through Custom Resource Definitions (CRDs).

== User documentation

You can find detailed information about the Gravitee Kubernetes Operator in the following sections of the Gravitee user documentation:

  * link:https://docs.gravitee.io/apim/3.x/apim_kubernetes_operator_overview.html[Overview^]
  * link:https://docs.gravitee.io/apim/3.x/apim_kubernetes_operator_architecture.html[Architecture^]
  * link:https://docs.gravitee.io/apim/3.x/apim_kubernetes_operator_quick_start.html[Quick Start^]
  * link:https://docs.gravitee.io/apim/3.x/apim_kubernetes_operator_installation.html[Installation and deployment^]
  * link:https://docs.gravitee.io/apim/3.x/apim_kubernetes_operator_user_guide.html[User Guide^]

The GKO API reference documentation is available https://github.com/gravitee-io/gravitee-kubernetes-operator/blob/master/docs/api/reference.md[here].

== Developer guide

=== Initialize your environment

* Install link:https://www.docker.com/[Docker^]
* Install link:https://kubernetes.io/docs/tasks/tools/#kubectl[kubectl^]
* Install link:https://helm.sh/docs/intro/install/[Helm^]
* Install link:https://nodejs.org/en/download/[NodeJs^]
* Install the operator-sdk: `brew install operator-sdk`

=== Install tooling for development

All the tool needed to run the make targets used during development can be installed by running the following command:

[source,shell]
----
make install-tools
----

To get more information about the available make targets, run:

[source,shell]
----
make help
----

=== Run the operator locally

To run the operator locally against an APIM-ready link:https://k3d.io/[k3d^] cluster, run the following commands:

[source,shell]
----
# Initialize a local kubernetes cluster running APIM
make k3d-init

# Install the operator CRDs into the cluster
make install

# Run the operator controller on your local machine
make run
----

=== Debug the operator and APIM

To be able to run the operator against a local instance of both an APIM Gateway and an APIM Management API, you will need to:

* Attach to a local cluster context.
* Create a local service account to authenticate the gateway against the local cluster.
* Create a Management Context pointing to your local APIM Management API.
* Run what you need to debug in debug mode.

[source,shell]
----
# Create a service account token with 'cluster-admin' role in the current context and
# use this token to authenticate against the current cluster
make k3d-admin

make run # or run using a debugger if you need to debug the operator as well

# Create the debug Management Context resource for APIM
kubectl apply -f ./config/samples/context/debug/api-context-with-credentials.yml

# Create a basic API Definition resource
kubectl apply -f ./config/samples/apim/api-with-context.yml
----

=== Run the operator as a deployment on the k3d cluster

Some features and behaviors of the operator can only be tested when running it as a deployment on the k3d cluster.

This is a case for e.g. for link:https://sdk.operatorframework.io/docs/building-operators/golang/webhook/[webhooks^] or 
when testing the operator deployed in multiple namespaces.

You can deploy the operator on your k3d cluster by running the following commands:

[source,shell]
----
make k3d-build k3d-push k3d-deploy
----

=== Working with the repo

When committing your contributions, please follow link:https://www.conventionalcommits.org/en/v1.0.0/[conventional commits^] and semantic release best practices.

=== Release process

To release a new version of the operator go to the dedicated https://github.com/gravitee-io/gravitee-kubernetes-operator/actions/workflows/trigger-release.yml[Trigger Release] workflow and run it with the default parameters.

This will:

* checkout a new `ci-prepare-release` branch from alpha
* rebase the master branch on top of it
* open a pull request from `ci-prepare-release` to `master`

Once the pull request is merged, semantic release will automatically create a new release, which includes:

* updating the changelog using conventional commits
* publishing a new docker image on docker hub with the new version tag
* publishing a new helm chart on https://helm.gravitee.io/index.yaml[helm.gravitee.io] with the new version tag

== Troubleshooting

=== Local Docker image registry

The k3d registry host used to share images between your host and your k3d cluster is defined as `k3d-graviteeio.docker.localhost`. On most linux / macos platforms, `*.localhost`` should resolve to 127.0.0.1. If this is not the case on your device, you need to add the following entry in your `/etc/hosts` file:

[source,shell]
----
127.0.0.1 k3d-graviteeio.docker.localhost
----
