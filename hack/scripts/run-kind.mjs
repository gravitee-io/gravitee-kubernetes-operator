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

const APIM_REGISTRY = `${process.env.APIM_IMAGE_REGISTRY || "graviteeio"}`;
const APIM_TAG = `${process.env.APIM_IMAGE_TAG || "latest"}`;
const APIM_VALUES = `${process.env.APIM_VALUES || "values.yaml"}`;

const IMAGES = new Map([
  [`${APIM_REGISTRY}/apim-gateway:${APIM_TAG}`, `gravitee-apim-gateway:dev`],
  [
    `${APIM_REGISTRY}/apim-management-api:${APIM_TAG}`,
    `gravitee-apim-management-api:dev`,
  ],
  [
    `${APIM_REGISTRY}/apim-management-ui:${APIM_TAG}`,
    `gravitee-apim-management-ui:dev`,
  ],
  [`mongo:6.0.15-jammy`, `mongo:6.0.15-jammy`],
  [`mccutchen/go-httpbin:latest`, `go-httpbin:dev`],
]);

const REDIRECT = argv.verbose ? "" : "> /dev/null";

toggleVerbosity(argv.verbose);

async function createKindCluster() {
  setNoQuoteEscape();
  await $`kind create cluster --config ${KIND_CONFIG}/kind.yaml ${REDIRECT}`;
  setQuoteEscape();
}

async function loadImages() {
  setNoQuoteEscape();

  await Promise.all(
    Array.from(IMAGES.keys()).map(
      (image) => $`docker pull ${image} ${REDIRECT}`,
    ),
  );

  await Promise.all(
    Array.from(IMAGES.entries()).map(
      ([image, tag]) => $`docker tag ${image} ${tag} ${REDIRECT}`,
    ),
  );

  await Promise.all(
    Array.from(IMAGES.values()).map(
      (tag) => $`kind load docker-image ${tag} --name gravitee ${REDIRECT}`,
    ),
  );

  setQuoteEscape();
}

async function createGraviteeNamespace() {
  await $`kubectl create ns gravitee`;
}

async function createTLSSecret() {
  await $`kubectl create secret tls tls-server --cert=${PKI}/server.crt --key=${PKI}/server.key`;
}

async function helmInstallAPIM() {
  await $`helm repo add graviteeio https://helm.gravitee.io`;
  await $`helm repo update graviteeio`;
  await $`helm install apim graviteeio/apim3 -f ${KIND_CONFIG}/apim/${APIM_VALUES}`;
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

    Version: ${APIM_TAG}

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
