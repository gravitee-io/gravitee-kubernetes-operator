#!/usr/bin/env zx

/**
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { K3D_IMAGES_REGISTRY } from "./docker.mjs";
import * as env from "./env.mjs";

export async function helmInstall() {
const helmInstallNginx = $`
helm install \
    --namespace ${env.K3D_NAMESPACE_NAME} \
    --set "image.registry=${K3D_IMAGES_REGISTRY}" \
    --set "image.repository=nginx-ingress-controller" \
    --set "defaultBackend.image.registry=${K3D_IMAGES_REGISTRY}" \
    --set "defaultBackend.image.repository=nginx" \
    --set "image.tag=${env.NGINX_CONTROLLER_IMAGE_TAG}" \
    --set "defaultBackend.image.tag=${env.NGINX_BACKEND_IMAGE_TAG}" \
    nginx-ingress bitnami/nginx-ingress-controller > /dev/null
`;

const helmInstallMongo = $`
helm install \
    --namespace ${env.K3D_NAMESPACE_NAME} \
    --set "image.registry=${K3D_IMAGES_REGISTRY}" \
    --set "image.repository=mongodb" \
    --set "image.tag=${env.MONGO_IMAGE_TAG}" \
    --set auth.enabled=false \
    --set readinessProbe.periodSeconds=30 \
    --set readinessProbe.timeoutSeconds=30 \
    --set livenessProbe.timeoutSeconds=30 \
    --set resources.limits.memory=2048Mi \
    --set resources.requests.memory=2048Mi \
    --set resources.limits.cpu=2000m \
    --set resources.requests.cpu=2000m \
    mongodb bitnami/mongodb > /dev/null
`;

const helmInstallApim = $`
helm install \
    --namespace ${env.K3D_NAMESPACE_NAME} \
    --set "portal.enabled=false" \
    --set "gateway.image.repository=${K3D_IMAGES_REGISTRY}/apim-gateway" \
    --set "gateway.services.sync.kubernetes.enabled=true" \
    --set "gateway.services.sync.kubernetes.namespaces=default" \
    --set "gateway.ingress.enabled=false" \
    --set "gateway.service.type=LoadBalancer" \
    --set "gateway.autoscaling.enabled=false" \
    --set "gateway.resources.requests.memory=2048Mi" \
    --set "gateway.resources.limits.memory=2048Mi" \
    --set "gateway.env[0].name=GIO_MIN_MEM" \
    --set "gateway.env[0].value=1024m" \
    --set "gateway.env[1].name=GIO_MAX_MEM" \
    --set "gateway.env[1].value=1024m" \
    --set "gateway.reporters.elasticsearch.enabled=false" \
    --set "gateway.image.tag=${env.APIM_IMAGE_TAG}" \
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
    --set "api.env[2].name=ANALYTICS_TYPE" \
    --set "api.env[2].value=none" \
    --set "api.startupProbe.initialDelaySeconds=5" \
    --set "api.startupProbe.timeoutSeconds=10" \
    --set "api.image.tag=${env.APIM_IMAGE_TAG}" \
    --set "api.image.repository=${K3D_IMAGES_REGISTRY}/apim-management-api" \
    --set "api.notifiers.smtp.enabled=false" \
    --set "ui.ingress.hosts[0]=localhost" \
    --set "ui.ingress.tls=false" \
    --set "ui.autoscaling.enabled=false" \
    --set "ui.image.repository=${K3D_IMAGES_REGISTRY}/apim-management-ui" \
    --set "ui.env[0].name=CONSOLE_BASE_HREF" \
    --set "ui.env[0].value=/console/" \
    --set "ui.image.tag=${env.APIM_IMAGE_TAG}" \
    --set "ui.baseURL=http://localhost:${env.NGINX_LOAD_BALANCER_PORT}/management/organizations/DEFAULT/environments/DEFAULT" \
    --set "elasticsearch.enabled=false" \
    --set "es.endpoints[0]=http://elasticsearch-master:9200" \
    --set "mongo.dbhost=mongodb" \
    --set "mongodb-replicaset=false" \
    --set "mongo.rsEnabled=false" \
    apim graviteeio/apim3 > /dev/null
`;

    await Promise.all([
        helmInstallApim,
        helmInstallMongo,
        helmInstallNginx,
    ]);
}

export async function addHelmRepos() {
    await $`helm repo add bitnami https://charts.bitnami.com/bitnami`;
    await $`helm repo add graviteeio https://helm.gravitee.io`;
}

export async function updateGraviteeRepo() {
    await $`helm repo update graviteeio`;
}
