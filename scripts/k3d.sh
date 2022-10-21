#!/bin/bash
# Copyright (C) 2015 The Gravitee team (http://gravitee.io)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -e

# APIM Images
APIM_IMAGE_REGISTRY="${APIM_IMAGE_REGISTRY:-graviteeio}"
APIM_IMAGE_TAG="${APIM_IMAGE_TAG:-3.19}"

# Docker dependencies images tags
NGINX_CONTROLLER_IMAGE_TAG=1.3.0
NGINX_BACKEND_IMAGE_TAG=1.22.0
MONGO_IMAGE_TAG=4.4
ELASTIC_IMAGE_TAG=7.17.5 

# K3d cluster config
K3D_CLUSTER_NAME=graviteeio
K3D_CLUSTER_AGENTS=1
K3D_API_PORT=6950
K3D_NAMESPACE_NAME=apim-dev
K3D_LOAD_BALANCER_PORT=9000
K3D_IMAGES_REGISTRY_NAME="${K3D_CLUSTER_NAME}.docker.localhost"
K3D_IMAGES_REGISTRY_PORT=12345
K3D_IMAGES_REGISTRY="${K3D_IMAGES_REGISTRY_NAME}:${K3D_IMAGES_REGISTRY_PORT}"

echo "

    Installing the latest version of k3d (if not present) ...

    See https://k3d.io/

"

curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

echo "

    Initialising a local docker images registry for k3d images (if not present) ...

    /!\ Warning /!\  

    Assuming that host "${K3D_IMAGES_REGISTRY_NAME}" points to 127.0.0.1

    You might needd to edit your /etc/hosts file before going further.

"

if k3d registry list | grep -q "${K3D_IMAGES_REGISTRY_NAME}"; then
    echo "
         K3d images registry ${K3D_IMAGES_REGISTRY_NAME} already exists, skipping.
    "
else
    echo "
        Initialising registry ${K3D_IMAGES_REGISTRY_NAME}
    "

    k3d registry create ${K3D_IMAGES_REGISTRY_NAME} --port ${K3D_IMAGES_REGISTRY_PORT}
fi

K3D_IMAGES_REGISTRY="k3d-${K3D_IMAGES_REGISTRY}"

echo "

    Deleting K3d cluster ${K3D_CLUSTER_NAME} (if present) ...

"

k3d cluster list | grep ${K3D_CLUSTER_NAME} && k3d cluster delete ${K3D_CLUSTER_NAME}


echo "

    Creating a K3d cluster with name ${K3D_CLUSTER_NAME}

"

k3d cluster create --wait \
    --agents ${K3D_CLUSTER_AGENTS} \
    --api-port ${K3D_API_PORT} \
    -p "${K3D_LOAD_BALANCER_PORT}:80@loadbalancer" \
    --k3s-arg "--disable=traefik@server:*" \
    --registry-use=${K3D_IMAGES_REGISTRY_NAME} \
    ${K3D_CLUSTER_NAME}

echo "

    Registering docker images to ${K3D_IMAGES_REGISTRY} ...

"

# docker pull "docker.io/bitnami/mongodb:${MONGO_IMAGE_TAG}"
# docker pull "docker.elastic.co/elasticsearch/elasticsearch:${ELASTIC_IMAGE_TAG}"
# docker pull "docker.io/bitnami/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}"
# docker pull "docker.io/bitnami/nginx:${NGINX_BACKEND_IMAGE_TAG}"
# docker pull "${APIM_IMAGE_REGISTRY}/apim-gateway:${APIM_IMAGE_TAG}"
# docker pull "${APIM_IMAGE_REGISTRY}/apim-management-api:${APIM_IMAGE_TAG}"
# docker pull "${APIM_IMAGE_REGISTRY}/apim-management-ui:${APIM_IMAGE_TAG}"

# docker tag "docker.io/bitnami/mongodb:${MONGO_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/mongodb:${MONGO_IMAGE_TAG}"
# docker tag "docker.elastic.co/elasticsearch/elasticsearch:${ELASTIC_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/elasticsearch:${ELASTIC_IMAGE_TAG}"
# docker tag "docker.io/bitnami/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}"
# docker tag "docker.io/bitnami/nginx:${NGINX_BACKEND_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/nginx:${NGINX_BACKEND_IMAGE_TAG}"
# docker tag "${APIM_IMAGE_REGISTRY}/apim-gateway:${APIM_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway:${APIM_IMAGE_TAG}"
# docker tag "${APIM_IMAGE_REGISTRY}/apim-management-api:${APIM_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-api:${APIM_IMAGE_TAG}"
# docker tag "${APIM_IMAGE_REGISTRY}/apim-management-ui:${APIM_IMAGE_TAG}" "${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-ui:${APIM_IMAGE_TAG}"

# docker push "${K3D_IMAGES_REGISTRY}/mongodb:${MONGO_IMAGE_TAG}"
# docker push "${K3D_IMAGES_REGISTRY}/elasticsearch:${ELASTIC_IMAGE_TAG}"
# docker push "${K3D_IMAGES_REGISTRY}/nginx-ingress-controller:${NGINX_CONTROLLER_IMAGE_TAG}"
# docker push "${K3D_IMAGES_REGISTRY}/nginx:${NGINX_BACKEND_IMAGE_TAG}"
# docker push "${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway:${APIM_IMAGE_TAG}"
# docker push "${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-api:${APIM_IMAGE_TAG}"
# docker push "${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-ui:${APIM_IMAGE_TAG}"

echo "

    Creating Kubernetes namespace ${K3D_NAMESPACE_NAME} ...

"

kubectl create namespace ${K3D_NAMESPACE_NAME}
kubectl config set-context --current --namespace ${K3D_NAMESPACE_NAME}

echo "

    Adding Helm repositories (if not presents) ...

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
    --set resources.limits.memory=2048Mi \
    --set resources.requests.memory=2048Mi \
    --set resources.limits.cpu=2000m \
    --set resources.requests.cpu=2000m \
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

helm install \
    --namespace ${K3D_NAMESPACE_NAME} \
    --set "gateway.image.tag=${APIM_IMAGE_TAG}" \
    --set "gateway.image.repository=${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway" \
    --set "gateway.services.sync.kubernetes.enabled=true" \
    --set "gateway.ingress.hosts[0]=localhost" \
    --set "gateway.ingress.path=/gateway/?(.*)" \
    --set "gateway.ingress.tls=false" \
    --set "gateway.ingress.annotations.'nginx\.ingress\.kubernetes\.io/rewrite-target'=/\$1" \
    --set "gateway.autoscaling.enabled=false" \
    --set "gateway.resources.requests.memory=2048Mi" \
    --set "gateway.resources.limits.memory=2048Mi" \
    --set "gateway.env[0].name=GIO_MIN_MEM" \
    --set "gateway.env[0].value=1024m" \
    --set "gateway.env[1].name=GIO_MAX_MEM" \
    --set "gateway.env[1].value=1024m" \
    --set "api.image.tag=${APIM_IMAGE_TAG}" \
    --set "api.image.repository=${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-api" \
    --set "api.ingress.management.hosts[0]=localhost" \
    --set "api.ingress.management.tls=false" \
    --set "api.portal.enabled=false" \
    --set "api.resources.requests.memory=2048Mi" \
    --set "api.resources.limits.memory=2048Mi" \
    --set "api.resources.requests.cpu=1500m" \
    --set "api.resources.limits.cpu=1500m" \
    --set "api.env[0].name=GIO_MIN_MEM" \
    --set "api.env[0].value=1024m" \
    --set "api.env[1].name=GIO_MAX_MEM" \
    --set "api.env[1].value=1024m" \
    --set "api.startupProbe.initialDelaySeconds=5" \
    --set "api.startupProbe.timeoutSeconds=10" \
    --set "api.autoscaling.enabled=false" \
    --set "api.replicaCount=1" \
    --set "ui.image.tag=${APIM_IMAGE_TAG}" \
    --set "ui.image.repository=${K3D_IMAGES_REGISTRY}/graviteeio/apim-management-ui" \
    --set "ui.ingress.hosts[0]=localhost" \
    --set "ui.ingress.tls=false" \
    --set "ui.autoscaling.enabled=false" \
    --set "ui.env[0].name=CONSOLE_BASE_HREF" \
    --set "ui.env[0].value=/console/" \
    --set "ui.baseURL=http://localhost:${K3D_LOAD_BALANCER_PORT}/management/organizations/DEFAULT/environments/DEFAULT" \
    --set "elasticsearch.enabled=false" \
    --set "es.endpoints[0]=http://elasticsearch-master:9200" \
    --set "mongo.dbhost=mongodb" \
    --set "mongodb-replicaset=false" \
    --set "mongo.rsEnabled=false" \
    --set "portal.enabled=false" \
    apim graviteeio/apim3

echo "

    APIM containers are starting ...

    Version: ${APIM_IMAGE_TAG}

    Available endpoints are:

        Gateway       http://localhost:${K3D_LOAD_BALANCER_PORT}/gateway
        Management    http://localhost:${K3D_LOAD_BALANCER_PORT}/management
        Console       http://localhost:${K3D_LOAD_BALANCER_PORT}/console/#!/login

    To update APIM components (e.g. APIM Gateway) to use a new docker image run:

    > docker tag <image> "${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway:${APIM_IMAGE_TAG}"
    > docker push "${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway:${APIM_IMAGE_TAG}"
    > kubectl rollout restart deployment apim-apim3-gateway
"

echo "
    Waiting for APIM to be ready ...
    
    Press ctrl+c to exit this script without waiting ...
"

kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=apim3 --timeout=360s
