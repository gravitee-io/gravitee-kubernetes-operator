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

import { LOG, setNoQuoteEscape, setQuoteEscape } from "../lib/index.mjs";
import * as env from "./env.mjs";

export let K3D_IMAGES_REGISTRY = `${env.K3D_IMAGES_REGISTRY_NAME}:${env.K3D_IMAGES_REGISTRY_PORT}`;

export async function initRegistry() {
    try {
        await $`k3d registry list | grep -q "${env.K3D_IMAGES_REGISTRY_NAME}"`;

        LOG.magenta(`K3d images registry ${env.K3D_IMAGES_REGISTRY_NAME} already exists, skipping.
    `);
    } catch (error) {
        LOG.magenta(`Initializing registry ${env.K3D_IMAGES_REGISTRY_NAME}
    `);

        await $`k3d registry create ${env.K3D_IMAGES_REGISTRY_NAME} --port ${env.K3D_IMAGES_REGISTRY_PORT}`;
    }

    K3D_IMAGES_REGISTRY = `k3d-${K3D_IMAGES_REGISTRY}`;
}

export async function registerImages() {
    const apimImages = new Map([
        [
            `${env.APIM_IMAGE_REGISTRY}/apim-gateway:${env.APIM_IMAGE_TAG}`,
            `${K3D_IMAGES_REGISTRY}/apim-gateway:${env.APIM_IMAGE_TAG}`,
        ],
        [
            `${env.APIM_IMAGE_REGISTRY}/apim-management-api:${env.APIM_IMAGE_TAG}`,
            `${K3D_IMAGES_REGISTRY}/apim-management-api:${env.APIM_IMAGE_TAG}`,
        ],
        [
            `${env.APIM_IMAGE_REGISTRY}/apim-management-ui:${env.APIM_IMAGE_TAG}`,
            `${K3D_IMAGES_REGISTRY}/apim-management-ui:${env.APIM_IMAGE_TAG}`,
        ],
    ]);

    const dependencyImages = [
        [
            `mongo:${env.MONGO_IMAGE_TAG}`,
            `${K3D_IMAGES_REGISTRY}/mongo:${env.MONGO_IMAGE_TAG}`,
        ],
        [
            `docker.io/bitnami/nginx-ingress-controller:${env.NGINX_CONTROLLER_IMAGE_TAG}`,
            `${K3D_IMAGES_REGISTRY}/nginx-ingress-controller:${env.NGINX_CONTROLLER_IMAGE_TAG}`,
        ],
        [
            `docker.io/bitnami/nginx:${env.NGINX_BACKEND_IMAGE_TAG}`,
            `${K3D_IMAGES_REGISTRY}/nginx:${env.NGINX_BACKEND_IMAGE_TAG}`,
        ],
    ];

    const allImages = new Map([...apimImages, ...dependencyImages]);

    const images = env.pullMode === "all" ? allImages : apimImages;

    setNoQuoteEscape();

    LOG.magenta(`Pulling docker images ...
      `);

    await Promise.all(
        Array.from(images.keys()).map(
            (image) => $`docker pull ${image} > /dev/null`
        )
    );

    LOG.magenta(`Tagging docker images ...
      `);

    await Promise.all(
        Array.from(images.entries()).map(
            ([image, tag]) => $`docker tag ${image} ${tag}`
        )
    );

    LOG.magenta(`Pushing docker images to ${K3D_IMAGES_REGISTRY} ...
      `);

    await Promise.all(
        Array.from(images.values()).map((tag) => $`docker push ${tag}  > /dev/null`)
    );

    setQuoteEscape();
}

export async function deployTestBackends() {
    await $`docker pull mccutchen/go-httpbin`;
    await $`docker tag mccutchen/go-httpbin ${K3D_IMAGES_REGISTRY}/go-httpbin`;
    await $`docker push ${K3D_IMAGES_REGISTRY}/go-httpbin`;
    await $`kubectl apply -f backends`;
}
