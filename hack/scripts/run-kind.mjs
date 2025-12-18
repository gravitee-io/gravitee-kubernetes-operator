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

import { APIM } from "./lib/apim.mjs";
import { Mongo } from "./lib/mongo.mjs";

import {
  LOG,
  PROJECT_DIR,
  setNoQuoteEscape,
  setQuoteEscape,
  time,
  toggleVerbosity,
} from "./lib/index.mjs";

const KIND_CONFIG = path.join(PROJECT_DIR, "hack", "kind");
const PKI = path.join(
  PROJECT_DIR,
  "examples",
  "usecase",
  "subscribe-to-mtls-plan",
  "pki",
);

const APIM_IMAGE_REGISTRY = await getAPIMImageRegistry();
const APIM_IMAGE_TAG = await getAPIMImageTag();
const APIM_CHART_REGISTRY = await getAPIMChartRegistry();
const APIM_CHART_VERSION = await getAPIMChartVersion();
const MONGO_IMAGE_TAG = await Mongo.getImageTag();

const APIM_VALUES = `${$.env.APIM_VALUES || "values.yaml"}`;

const IMAGES = new Map([
  [
    `${APIM_IMAGE_REGISTRY}/apim-gateway:${APIM_IMAGE_TAG}`,
    `gravitee-apim-gateway:dev`,
  ],
  [
    `${APIM_IMAGE_REGISTRY}/apim-management-api:${APIM_IMAGE_TAG}`,
    `gravitee-apim-management-api:dev`,
  ],
  [
    `${APIM_IMAGE_REGISTRY}/apim-management-ui:${APIM_IMAGE_TAG}`,
    `gravitee-apim-management-ui:dev`,
  ],
  [`mongo:${MONGO_IMAGE_TAG}`, `mongo:${MONGO_IMAGE_TAG}`],
  [`mccutchen/go-httpbin:latest`, `go-httpbin:dev`],
]);

if (APIM_VALUES.includes("dbless")) {
  IMAGES.delete(`mongo:${MONGO_IMAGE_TAG}`);
  IMAGES.delete(`${APIM_IMAGE_REGISTRY}/apim-management-api:${APIM_IMAGE_TAG}`);
  IMAGES.delete(`${APIM_IMAGE_REGISTRY}/apim-management-ui:${APIM_IMAGE_TAG}`);
}

const REDIRECT = argv.verbose ? "" : "> /dev/null";

toggleVerbosity(argv.verbose);

async function getAPIMImageRegistry() {
  if ($.env.APIM_IMAGE_REGISTRY) {
    return $.env.APIM_IMAGE_REGISTRY;
  }
  return await APIM.getImageRegistry();
}

async function getAPIMImageTag() {
  if ($.env.APIM_IMAGE_TAG) {
    return $.env.APIM_IMAGE_TAG;
  }
  return await APIM.getImageTag();
}

async function getAPIMChartRegistry() {
  if ($.env.APIM_CHART_REGISTRY) {
    return $.env.APIM_CHART_REGISTRY;
  }
  return await APIM.getChartRegistry();
}

async function getAPIMChartVersion() {
  if ($.env.APIM_CHART_VERSION) {
    return $.env.APIM_CHART_VERSION;
  }
  return await APIM.getChartVersion();
}

async function createKindCluster() {
  setNoQuoteEscape();
  const clusters = await $`kind get clusters`.quiet();
  if (clusters.toString().split("\n").includes("gravitee")) {
    LOG.blue(`Kind cluster 'gravitee' already exists. Skipping creation...`);
    setQuoteEscape();
    return;
  }
  await $`kind create cluster --config ${KIND_CONFIG}/kind.yaml ${REDIRECT}`;
  setQuoteEscape();
}

async function loadImages() {
  setNoQuoteEscape();

  const promisesToLoad = Array.from(IMAGES.entries()).map(pullAndTag);

  const tagsToLoad = await Promise.all(promisesToLoad);

  LOG.blue(`All images pulled. Starting Kind load...`);

  for (const tag of tagsToLoad) {
    LOG.blue(`loading image tag ${tag}`);
    await $`kind load docker-image ${tag} --name gravitee`;
  }

  setQuoteEscape();
}

async function pullAndTag([image, tag]) {
  if (!image.includes("local")) {
    LOG.blue(`pulling image ${image}`);
    await $`docker pull ${image}`;
  }
  LOG.blue(`tagging image ${image} with ${tag}`);
  await $`docker tag ${image} ${tag}`;
  return tag;
}

async function createGraviteeNamespace() {
  await $`kubectl create ns gravitee --dry-run=client -o yaml | kubectl apply -f -`;
}

async function createTLSSecret() {
  await $`kubectl create secret tls tls-server --cert=${PKI}/server.crt --key=${PKI}/server.key --dry-run=client -o yaml | kubectl apply -f -`;
}

async function helmInstallAPIM() {
  await $`helm repo add graviteeio https://helm.gravitee.io`;
  await $`helm repo update graviteeio`;
  await $`helm upgrade --install apim ${APIM_CHART_REGISTRY} -f ${KIND_CONFIG}/apim/${APIM_VALUES} --version ${APIM_CHART_VERSION}`;
}

async function deployHTTPBin() {
  await $`kubectl apply -f ${KIND_CONFIG}/httpbin`;
}

async function waitForApim() {
  await $`kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=apim3 --timeout=360s`;
}

LOG.blue(`
  ☸ Initializing kind cluster
`);

await time(createKindCluster);

LOG.blue(`
  🐳 Loading docker images
`);

await time(loadImages);

LOG.blue(`
  ☸ Creating gravitee namespace
`);

await time(createGraviteeNamespace);

LOG.blue(`
  ☸ Creating APIM gateway TLS secret
`);

await time(createTLSSecret);

LOG.blue(`
  ☸ Installing APIM
`);

await time(helmInstallAPIM);

LOG.blue(`
  ☸ Deploying httpbin
`);

await time(deployHTTPBin);

LOG.magenta(`
    APIM containers are starting ...

    Version: ${APIM_IMAGE_TAG}

    Available endpoints are:
        Gateway             http://localhost:30082
        Gateway with mTLS   https://localhost:30084
        Management API      http://localhost:30083/management/organizations/DEFAULT
        Console             http://localhost:30080
`);

LOG.blue(`Waiting for services to be ready ...
    
    Press ctrl+c to exit this script without waiting ...
`);

await time(waitForApim);
