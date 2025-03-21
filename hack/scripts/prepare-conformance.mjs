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
  setNoQuoteEscape,
  setQuoteEscape,
  time,
} from "./lib/index.mjs";

const KIND_CONFIG = path.join(__dirname, "..", "kind");

const REDIRECT = argv.verbose ? "" : "> /dev/null";

async function createKindCluster() {
  setNoQuoteEscape();
  await $`kind create cluster --config ${KIND_CONFIG}/kind.conformance.yaml ${REDIRECT}`;
  setQuoteEscape();
}

async function runCloudProvider() {
  setNoQuoteEscape();
  await $`NET_MODE=kind docker compose -f kind/cloud-provider/compose.yaml up -d`;
  setQuoteEscape();
}

async function installCertManager() {
    await $`
        helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager \
            --set config.apiVersion="controller.config.cert-manager.io/v1alpha1" \
            --set config.kind="ControllerConfiguration" \
            --set config.enableGatewayAPI=true \
             --set crds.enabled=true \
             --create-namespace
    `
}

LOG.blue(`
  ☸ Initializing kind cluster
`);

await time(createKindCluster);

LOG.blue(`
  ☸ Running cloud provider
`);

await time(runCloudProvider);

LOG.blue(`
  ☸ Installing cert-manager
`);

await time(installCertManager);

LOG.blue(`
    ☸ Run the operator
`);

