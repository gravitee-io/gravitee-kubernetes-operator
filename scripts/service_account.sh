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


SERVICE_ACCOUNT_NAME=kubeconfig-sa
SERVICE_ACCOUNT_NAMESPACE=kube-system
CLUSTER_ROLE_BINDING_NAME=add-on-cluster-admin

KUBE_CONTEXT=$(kubectl config current-context)

echo "

    You are about to create a service account with 'cluster-admin' in kube context ${KUBE_CONTEXT}
    
"
read -r -p "
    Do you want to continue (type yes if you want to proceed) ? " CONTINUE

if [ "$CONTINUE" != "yes" ]; then
    echo "

    Exiting ...

    "
    exit 0
else
    echo "
    
    Proceeding ...

    "
fi

kubectl -n "${SERVICE_ACCOUNT_NAMESPACE}" create serviceaccount ${SERVICE_ACCOUNT_NAME} > /dev/null 2>&1 

kubectl create clusterrolebinding "${CLUSTER_ROLE_BINDING_NAME}" --clusterrole=cluster-admin --serviceaccount=${SERVICE_ACCOUNT_NAMESPACE}:${SERVICE_ACCOUNT_NAME} > /dev/null 2>&1

TOKEN_NAME=$(kubectl -n ${SERVICE_ACCOUNT_NAMESPACE} get serviceaccount/${SERVICE_ACCOUNT_NAME} -o jsonpath='{.secrets[0].name}')

TOKEN=$(kubectl -n ${SERVICE_ACCOUNT_NAMESPACE} get secret ${TOKEN_NAME} -o jsonpath='{.data.token}'| base64 --decode)

kubectl config set-credentials ${SERVICE_ACCOUNT_NAME} --token="${TOKEN}"
kubectl config set-context --current --user=${SERVICE_ACCOUNT_NAME}

echo "
    You are now authenticated against ${KUBE_CONTEXT} as ${SERVICE_ACCOUNT_NAME}
"
