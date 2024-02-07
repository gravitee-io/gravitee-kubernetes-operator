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

import {LOG, setNoQuoteEscape, time, toggleVerbosity} from "./lib/index.mjs";
import {K3D_IMAGES_REGISTRY, deployTestBackends, initRegistry, registerImages} from "./util/docker.mjs";
import {
    initCluster,
    installK3d,
    createSecret,
    createHttpBinSecret,
    createTemplatingSecret,
    createTemplatingConfigmap
} from "./util/k3d.mjs";
import {addHelmRepos, helmInstall, updateGraviteeRepo} from "./util/helm.mjs";
import * as env from "./util/env.mjs";

toggleVerbosity(argv.verbose);

LOG.green(`
Starting k3d cluster with APIM dependencies...`);

LOG.blue(`
  ☸ Installing the latest version of k3d (if not present) ...

  See https://k3d.io/
`);

await time(installK3d);

LOG.blue(`
  🐳 Initializing a local docker images registry for k3d images (if not present) ...
`);

await time(initRegistry);

LOG.yellow(`
  ⚠️ WARNING ⚠️ 

  Assuming that host "${env.K3D_IMAGES_REGISTRY_NAME}" points to 127.0.0.1

  You might need to edit your /etc/hosts file before going further.
`);

await time(initCluster);

if (env.pullMode !== "none") {
    LOG.blue(`
  🐳 Registering docker images to ${K3D_IMAGES_REGISTRY} ...
`);

    await time(registerImages);
}

LOG.blue(`
  ☸ Storing APIM context credentials as a secret ...

  The following declaration can be used in your API context to reference this secret:

  secretRef:
      name: ${env.APIM_CONTEXT_SECRET_NAME}
      namespace: ${env.K3D_NAMESPACE_NAME}
`);

await time(createSecret);

LOG.blue(`
  ☸ Storing Templating secret ...
`);

await time(createTemplatingSecret);

LOG.blue(`
  ☸ Storing Templating configmap ...
`);

await time(createTemplatingConfigmap);

LOG.blue(`
  ☸ Storing httpbin.example.com keypair as a secret ...

  The following declaration can be used in your ingress to reference this secret:

  tls:
    - hosts:
        - httpbin.example.com
      secretName: httpbin.example.com
`);

await time(createHttpBinSecret);

LOG.blue(`
⎈ Adding Helm repositories (if not presents) ...
`);

await time(addHelmRepos);

LOG.blue(`
⎈ Ensuring that Graviteeio repo is up to date ...
`);

await time(updateGraviteeRepo);

LOG.blue(`
⎈ Installing components in namespace ${env.K3D_NAMESPACE_NAME} ...

      Mongodb         ${env.MONGO_IMAGE_TAG}
      Elasticsearch   ${env.ELASTIC_IMAGE_TAG}
      Nginx ingress   ${env.NGINX_CONTROLLER_IMAGE_TAG}           
      Nginx backend   ${env.NGINX_BACKEND_IMAGE_TAG}
      Gravitee APIM   ${env.APIM_IMAGE_TAG}
`);

setNoQuoteEscape();

await time(helmInstall);

LOG.blue(`
⎈ Deploying test backends ...
`);

await time(deployTestBackends);

LOG.magenta(`
    APIM containers are starting ...

    Version: ${env.APIM_IMAGE_TAG}

    Available endpoints are:

        Gateway       http://localhost:${env.GATEWAY_LOAD_BALANCER_PORT}
        Management    http://localhost:${env.NGINX_LOAD_BALANCER_PORT}/management
        Console       http://localhost:${env.NGINX_LOAD_BALANCER_PORT}/console/#!/login

    To update APIM components (e.g. APIM Gateway) to use a new docker image run:

    > docker tag <image> "${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway:${env.APIM_IMAGE_TAG}"
    > docker push "${K3D_IMAGES_REGISTRY}/graviteeio/apim-gateway:${env.APIM_IMAGE_TAG}"
    > kubectl rollout restart deployment apim-apim3-gateway
`);

LOG.blue(`Waiting for APIM to be ready ...
    
    Press ctrl+c to exit this script without waiting ...
`);

await time(waitForApim);

async function waitForApim() {
    await $`kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=apim3 --timeout=360s`;
}

LOG.green(`
🚀 Cluster is ready!
`);
