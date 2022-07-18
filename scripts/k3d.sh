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
K3D_IMAGES_REGISTRY_NAME="${K3D_CLUSTER_NAME}.docker.localhost"
K3D_IMAGES_REGISTRY_PORT=12345
K3D_IMAGES_REGISTRY="${K3D_IMAGES_REGISTRY_NAME}:${K3D_IMAGES_REGISTRY_PORT}"
K3D_IMAGES_REGISTRY_LOCALHOST="localhost:${K3D_IMAGES_REGISTRY_PORT}"

echo "

    Installing the latest version of k3d (if not already present)

    See https://k3d.io/

"

curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

if [[ $(k3d registry list | grep $K3D_IMAGES_REGISTRY_NAME) ]]; then
    echo "

         K3d images registry already exists.
    "
else
    echo "

        Initialising a local docker images registry for k3d images (if not present)

    "

    k3d registry create ${K3D_IMAGES_REGISTRY_NAME} --port ${K3D_IMAGES_REGISTRY_PORT}
fi

echo "

    Creating a k3d cluster with name ${K3D_CLUSTER_NAME}

"

k3d cluster create --wait \
    --agents ${K3D_CLUSTER_AGENTS} \
    --api-port ${K3D_API_PORT} \
    -p "${K3D_LOAD_BALANCER_PORT}:80@loadbalancer" \
    --k3s-arg "--disable=traefik@server:*" \
    --registry-use=${K3D_IMAGES_REGISTRY_NAME} \
    ${K3D_CLUSTER_NAME}

K3D_IMAGES_REGISTRY="k3d-${K3D_IMAGES_REGISTRY}"

echo "

    Registering docker images to ${K3D_IMAGES_REGISTRY}

"

docker pull docker.io/bitnami/mongodb:${MONGO_IMAGE_TAG}
docker pull docker.elastic.co/elasticsearch/elasticsearch:${ELASTIC_IMAGE_TAG}
docker pull docker.io/bitnami/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}
docker pull docker.io/bitnami/nginx:${NGINX_BACKEND_IMAGE_TAG}
docker pull graviteeio/apim-gateway:${APIM_IMAGE_TAG}
docker pull graviteeio/apim-management-api:${APIM_IMAGE_TAG}
docker pull graviteeio/apim-management-ui:${APIM_IMAGE_TAG}

docker tag "docker.io/bitnami/mongodb:${MONGO_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/mongodb:${MONGO_IMAGE_TAG}"
docker tag "docker.elastic.co/elasticsearch/elasticsearch:${ELASTIC_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/elasticsearch:${ELASTIC_IMAGE_TAG}"
docker tag "docker.io/bitnami/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}"
docker tag "docker.io/bitnami/nginx:${NGINX_BACKEND_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/nginx:${NGINX_BACKEND_IMAGE_TAG}"
docker tag "graviteeio/apim-gateway:${APIM_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway:${APIM_IMAGE_TAG}"
docker tag "graviteeio/apim-management-api:${APIM_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-api:${APIM_IMAGE_TAG}"
docker tag "graviteeio/apim-management-ui:${APIM_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-ui:${APIM_IMAGE_TAG}"

docker push "${K3D_IMAGES_REGISTRY_LOCALHOST}/mongodb:${MONGO_IMAGE_TAG}"
docker push "${K3D_IMAGES_REGISTRY_LOCALHOST}/elasticsearch:${ELASTIC_IMAGE_TAG}"
docker push "${K3D_IMAGES_REGISTRY_LOCALHOST}/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}"
docker push "${K3D_IMAGES_REGISTRY_LOCALHOST}/nginx:${NGINX_BACKEND_IMAGE_TAG}"
docker push "${K3D_IMAGES_REGISTRY_LOCALHOST}/graviteeio/apim-gateway:${APIM_IMAGE_TAG}"
docker push "${K3D_IMAGES_REGISTRY_LOCALHOST}/graviteeio/apim-management-api:${APIM_IMAGE_TAG}"
docker push "${K3D_IMAGES_REGISTRY_LOCALHOST}/graviteeio/apim-management-ui:${APIM_IMAGE_TAG}"

echo "

    Creating Kubernetes namespace ${K3D_NAMESPACE_NAME}

"

kubectl create namespace ${K3D_NAMESPACE_NAME}
kubectl config set-context --current --namespace ${K3D_NAMESPACE_NAME}

echo "

    Adding Helm repositories (if not presents)

"

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

helm install \
    --namespace ${K3D_NAMESPACE_NAME} \
    --set "image.registry=${K3D_IMAGES_REGISTRY}" \
    --set "image.repository=mongodb" \
    --set "image.tag=${MONGO_IMAGE_TAG}" \
    --set auth.enabled=false \
    --set readinessProbe.periodSeconds=30 \
    --set readinessProbe.timeoutSeconds=30 \
    --set livenessProbe.timeoutSeconds=30 \
    mongodb bitnami/mongodb

helm install \
    --namespace ${K3D_NAMESPACE_NAME} \
    --set replicas=1 \
    --set "image=${K3D_IMAGES_REGISTRY}/elasticsearch" \
    --set "imageTag=${ELASTIC_IMAGE_TAG}" \
    elastic elastic/elasticsearch

helm install \
    --namespace ${K3D_NAMESPACE_NAME} \
    --set "image.registry=${K3D_IMAGES_REGISTRY}" \
    --set "image.repository=nginx-ingress-controller" \
    --set "defaultBackend.image.registry=${K3D_IMAGES_REGISTRY}" \
    --set "defaultBackend.image.repository=nginx" \
    --set "image.tag=${NGINX_CONTROLLER_IMAGE_TAG}" \
    --set "defaultBackend.image.tag=${NGINX_BACKEND_IMAGE_TAG}" \
    nginx-ingress bitnami/nginx-ingress-controller


BASEDIR="$( cd "$( dirname "$0" )" && pwd )"
helm install \
    --namespace ${K3D_NAMESPACE_NAME} \
    -f "$BASEDIR/helm/apim-values.yml" \
    --set "gateway.image.repository=${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway" \
    --set "api.image.repository=${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-api" \
    --set "ui.image.repository=${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-ui" \
    --set "gateway.image.tag=${APIM_IMAGE_TAG}" \
    --set "api.image.tag=${APIM_IMAGE_TAG}" \
    --set "ui.image.tag=${APIM_IMAGE_TAG}" \
    --set "ui.baseURL=http://localhost:${K3D_LOAD_BALANCER_PORT}/management/organizations/DEFAULT/environments/DEFAULT/" \
    apim graviteeio/apim3

echo "

    APIM should be ready in a few minutes ...

    Version: ${APIM_IMAGE_TAG}

    Available endpoints are:

        Gateway       http://localhost:${K3D_LOAD_BALANCER_PORT}/gateway
        Management    http://localhost:${K3D_LOAD_BALANCER_PORT}/management
        Console       http://localhost:${K3D_LOAD_BALANCER_PORT}/console/#!/login

    To update APIM components (e.g. APIM Gateway) to use a new docker image run:

    > docker tag <image> "${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway:${APIM_IMAGE_TAG}"
    > docker push "${K3D_IMAGES_REGISTRY_LOCALHOST}/graviteeio/apim-gateway:${APIM_IMAGE_TAG}"
    > kubectl rollout restart deployment apim-apim3-gateway
"
