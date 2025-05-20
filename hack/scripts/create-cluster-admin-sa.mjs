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

import { LOG, toggleVerbosity } from "./lib/index.mjs";

toggleVerbosity(argv.verbose);

const SERVICE_ACCOUNT_NAME = "gravitee-sa";
const SERVICE_ACCOUNT_NAMESPACE = "kube-system";
const CLUSTER_ROLE_BINDING_NAME = "add-on-cluster-admin";
const KUBE_CONTEXT = await $`kubectl config current-context`;

LOG.yellow(
  `⚠️ You are about to create a service account with 'cluster-admin' in kube context ${KUBE_CONTEXT}`,
);

const CONTINUE = await question("Do you want to proceed? (yes/no)");

if (CONTINUE !== "yes") {
  LOG.blue(`
Aborting ...`);
  process.exit(0);
}

LOG.blue(`
Creating service account ${SERVICE_ACCOUNT_NAME} in namespace ${SERVICE_ACCOUNT_NAMESPACE} ...`);

await $`kubectl -n "${SERVICE_ACCOUNT_NAMESPACE}" create serviceaccount ${SERVICE_ACCOUNT_NAME}`;

LOG.blue(`
Creating cluster role binding ${CLUSTER_ROLE_BINDING_NAME} ...`);

await $`kubectl create clusterrolebinding "${CLUSTER_ROLE_BINDING_NAME}" --clusterrole=cluster-admin --serviceaccount=${SERVICE_ACCOUNT_NAMESPACE}:${SERVICE_ACCOUNT_NAME}
`;

const SECRET_DEFINITION = `
apiVersion: v1
kind: Secret
metadata:
  name: ${SERVICE_ACCOUNT_NAME}
  namespace: ${SERVICE_ACCOUNT_NAMESPACE}
  annotations:
    kubernetes.io/service-account.name: ${SERVICE_ACCOUNT_NAME}
type: kubernetes.io/service-account-token
`;

await $`echo ${SECRET_DEFINITION}| kubectl apply -f -`;

const SECRET =
  await $`kubectl get secret ${SERVICE_ACCOUNT_NAME} -n ${SERVICE_ACCOUNT_NAMESPACE} -o yaml`;

const TOKEN = YAML.parse(String(SECRET)).data.token;
const DECODED_TOKEN = Buffer.from(TOKEN, "base64").toString("ascii");

LOG.blue(`
Setting kubectl context to use service account ${SERVICE_ACCOUNT_NAME} ...`);

await $`kubectl config set-credentials ${SERVICE_ACCOUNT_NAME} --token="${DECODED_TOKEN}"`;

await $`kubectl config set-context --current --user="${SERVICE_ACCOUNT_NAME}"`;

LOG.green(`
Service account ${SERVICE_ACCOUNT_NAME} created successfully!`);
