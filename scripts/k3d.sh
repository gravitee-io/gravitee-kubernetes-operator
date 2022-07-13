#!/bin/bash

# Docker images tags
NGINX_CONTROLLER_IMAGE_TAG=1.3.0
NGINX_BACKEND_IMAGE_TAG=1.22.0
MONGO_IMAGE_TAG=4.4
ELASTIC_IMAGE_TAG=7.17.5 
APIM_IMAGE_TAG=3.17.3

# K3d cluster config
K3D_CLUSTER_NAME=graviteeio
K3D_CLUSTER_AGENTS=1
K3D_API_PORT=6950
K3D_NAMESPACE_NAME=apim-dev
K3D_LOAD_BALANCER_PORT=9000

echo "

    Installing the latest k3d version if not already present

"

curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

echo "

    Creating a k3d cluster with name ${K3D_CLUSTER_NAME}

"

k3d cluster create --wait \
    --agents ${K3D_CLUSTER_AGENTS} \
    --api-port ${K3D_API_PORT} \
    -p "${K3D_LOAD_BALANCER_PORT}:80@loadbalancer" \
    --k3s-arg "--disable=traefik@server:*" \
    ${K3D_CLUSTER_NAME}

echo "

    Trying to import docker images to the cluster

"

# Images preloading is not running if running on arm64 (see https://github.com/k3d-io/k3d/issues/1025)
if [[ $(uname -m) == 'arm64' ]]; then
    echo "

        Important Notice !

        You are running on an arm64, which prevents from preloading docker images in the cluster.

        Images will be pulled on container start and your components may take some time to be ready.

        See https://github.com/k3d-io/k3d/issues/1025

    "
else
    docker pull docker.io/bitnami/nginx:${NGINX_BACKEND_IMAGE_TAG}
    docker pull docker.io/bitnami/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}
    docker pull docker.io/bitnami/mongodb:${MONGO_IMAGE_TAG}
    docker pull docker.elastic.co/elasticsearch/elasticsearch:${ELASTIC_IMAGE_TAG}
    docker pull graviteeio/apim-gateway:${APIM_IMAGE_TAG}
    docker pull graviteeio/apim-management-api:${APIM_IMAGE_TAG}
    docker pull graviteeio/apim-management-ui:${APIM_IMAGE_TAG}

    k3d image import \
        -m tools-node \
        -c ${K3D_CLUSTER_NAME} \
            docker.io/bitnami/nginx:${NGINX_BACKEND_IMAGE_TAG} \
            docker.io/bitnami/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG} \
            bitnami/mongodb:${MONGO_IMAGE_TAG} \
            docker.elastic.co/elasticsearch/elasticsearch:${ELASTIC_IMAGE_TAG} \
            graviteeio/apim-gateway:${APIM_IMAGE_TAG} \
            graviteeio/apim-management-api:${APIM_IMAGE_TAG} \
            graviteeio/apim-management-ui:${APIM_IMAGE_TAG}
fi

echo "

    Adding helm repositories

"

# Add Helm repos
helm repo add elastic https://helm.elastic.co
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add graviteeio https://helm.gravitee.io

echo "

    Installing components in namespace ${K3D_NAMESPACE_NAME}

        Mongodb         ${MONGO_IMAGE_TAG}
        Elasticsearch   ${ELASTIC_IMAGE_TAG}
        Nginx ingress   ${NGINX_CONTROLLER_IMAGE_TAG}           
        Nginx backend   ${NGINX_BACKEND_IMAGE_TAG}
        Gravitee APIM   ${APIM_IMAGE_TAG}

"

# Install Helm charts
helm install \
    --create-namespace \
    --namespace ${K3D_NAMESPACE_NAME} \
    --set replicas=1 \
    --set "imageTag=${ELASTIC_IMAGE_TAG}" \
    elastic elastic/elasticsearch

helm install \
    --namespace ${K3D_NAMESPACE_NAME} \
    --set "image.tag=${MONGO_IMAGE_TAG}" \
    --set auth.enabled=false \
    --set readinessProbe.periodSeconds=30 \
    --set readinessProbe.timeoutSeconds=30 \
    --set livenessProbe.timeoutSeconds=30 \
    mongodb bitnami/mongodb

helm install \
    --namespace ${K3D_NAMESPACE_NAME} \
    --set "image.tag=${NGINX_CONTROLLER_IMAGE_TAG}" \
    --set "defaultBackend.image.tag=${NGINX_BACKEND_IMAGE_TAG}" \
    nginx-ingress bitnami/nginx-ingress-controller

helm install \
    --namespace ${K3D_NAMESPACE_NAME} \
    -f helm/apim-values.yml \
    --set "gateway.image.tag=${APIM_IMAGE_TAG}" \
    --set "api.image.tag=${APIM_IMAGE_TAG}" \
    --set "ui.image.tag=${APIM_IMAGE_TAG}" \
    --set "ui.baseURL=http://localhost:${K3D_LOAD_BALANCER_PORT}/management/organizations/DEFAULT/environments/DEFAULT/" \
    apim graviteeio/apim3

echo "
    
    Switching to namespace ${K3D_NAMESPACE_NAME}

"

kubectl config set-context --current --namespace ${K3D_NAMESPACE_NAME}

echo "

    APIM should be ready in a few minutes ...

    Version: ${APIM_IMAGE_TAG}

    Available endpoints are:

        Gateway       http://localhost:${K3D_LOAD_BALANCER_PORT}/gateway
        Management    http://localhost:${K3D_LOAD_BALANCER_PORT}/management
        Console       http://localhost:${K3D_LOAD_BALANCER_PORT}/console/#!/login
"
