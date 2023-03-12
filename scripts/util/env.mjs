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

export const pullMode = argv.pull || "all";
export const APIM_IMAGE_REGISTRY = `${
    process.env.APIM_IMAGE_REGISTRY || "graviteeio"
}`;
export const APIM_IMAGE_TAG = `${process.env.APIM_IMAGE_TAG || "latest"}`;
export const NGINX_CONTROLLER_IMAGE_TAG = "1.3.0";
export const NGINX_BACKEND_IMAGE_TAG = "1.22.0";
export const MONGO_IMAGE_TAG = "4.4";
export const ELASTIC_IMAGE_TAG = "7.17.5";

// K3d cluster config
export const K3D_CLUSTER_NAME = "graviteeio";
export const K3D_CLUSTER_AGENTS = 1;
export const K3D_API_PORT = 6950;
export const K3D_NAMESPACE_NAME = "default";
export const NGINX_LOAD_BALANCER_PORT = 9000;
export const GATEWAY_LOAD_BALANCER_PORT = 9001;
export const K3D_IMAGES_REGISTRY_NAME = `${K3D_CLUSTER_NAME}.docker.localhost`;
export const K3D_IMAGES_REGISTRY_PORT = 12345;

// APIM credentials
export const APIM_CONTEXT_SECRET_NAME = "apim-context-credentials";
export const GATEWAY_KEY_STORE_SECRET = "gw-keystore";
export const HTTPBIN_EXAMPLE_COM = "httpbin.example.com";
export const GATEWAY_KEY_STORE_CREDENTIALS_SECRET_NAME = "gw-keystore-credentials";
export const TEMPLATING_SECRET_CONFIGMAP_NAME = "graviteeio-templating";