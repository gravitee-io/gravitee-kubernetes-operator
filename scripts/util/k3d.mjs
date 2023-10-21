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

import { LOG } from "../lib/index.mjs";
import * as env from "./env.mjs";

export async function installK3d() {
    try {
        await $`k3d version`;
    } catch (e) {
        await $`curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash`;
    }
}

export async function initCluster() {
    LOG.blue(`☸ Deleting K3d cluster ${env.K3D_CLUSTER_NAME} (if present) ...
`);

    try {
        await $`k3d cluster list| grep -q "${env.K3D_CLUSTER_NAME}"`;
        await $`k3d cluster delete ${env.K3D_CLUSTER_NAME}`;
    } catch (error) {
        LOG.magenta(`No K3d cluster with name ${env.K3D_CLUSTER_NAME}
    `);
    }

    LOG.blue(`☸ Creating a K3d cluster with name ${env.K3D_CLUSTER_NAME} ...
  `);

    await $`k3d cluster create --wait \
    --agents ${env.K3D_CLUSTER_AGENTS} \
    --api-port ${env.K3D_API_PORT} \
    -p "${env.NGINX_LOAD_BALANCER_PORT}:80@loadbalancer" \
    -p "${env.GATEWAY_LOAD_BALANCER_PORT}:82@loadbalancer" \
    --k3s-arg "--disable=traefik@server:*" \
    --registry-use=${env.K3D_IMAGES_REGISTRY_NAME} \
    --k3s-arg '--kubelet-arg=eviction-hard=imagefs.available<1%,nodefs.available<1%@agent:*' \
    --k3s-arg '--kubelet-arg=eviction-minimum-reclaim=imagefs.available=1%,nodefs.available=1%@agent:*' \
    ${env.K3D_CLUSTER_NAME}
  `;
}

export async function createSecret() {
    await $`kubectl create secret generic ${env.APIM_CONTEXT_SECRET_NAME} \
    -n ${env.K3D_NAMESPACE_NAME} \
    --from-literal=username=admin \
    --from-literal=password=admin
  `;
}

export async function createTemplatingSecret() {
    await $`kubectl create secret generic ${env.TEMPLATING_SECRET_CONFIGMAP_NAME} \
    -n ${env.K3D_NAMESPACE_NAME} \
    --from-literal=security=KEY_LESS
  `;
}

export async function createTemplatingConfigmap() {
    await $`kubectl create configmap ${env.TEMPLATING_SECRET_CONFIGMAP_NAME} \
    -n ${env.K3D_NAMESPACE_NAME} \
    --from-literal=target=https://api.gravitee.io/echo
  `;
}

export async function createHttpBinSecret() {
    await $`kubectl create secret tls ${env.HTTPBIN_EXAMPLE_COM} \
    -n ${env.K3D_NAMESPACE_NAME} \
    --key ${__dirname}/resources/httpbin.example.com.key \
    --cert ${__dirname}/resources/httpbin.example.com.crt
  `;
}
