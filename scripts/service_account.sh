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


SERVICE_ACCOUNT_NAME=gravitee-sa
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

kubectl -n "${SERVICE_ACCOUNT_NAMESPACE}" create serviceaccount ${SERVICE_ACCOUNT_NAME}

kubectl create clusterrolebinding "${CLUSTER_ROLE_BINDING_NAME}" --clusterrole=cluster-admin --serviceaccount=${SERVICE_ACCOUNT_NAMESPACE}:${SERVICE_ACCOUNT_NAME}

cat << EOF | kubectl create -f -

apiVersion: v1

kind: Secret

metadata:

  name: ${SERVICE_ACCOUNT_NAME}
  namespace: ${SERVICE_ACCOUNT_NAMESPACE}
  annotations:
    kubernetes.io/service-account.name: ${SERVICE_ACCOUNT_NAME}

type: kubernetes.io/service-account-token

EOF

TOKEN=$(kubectl get secrets -n ${SERVICE_ACCOUNT_NAMESPACE} -o jsonpath="{.items[?(@.metadata.annotations['kubernetes\.io/service-account\.name']=='${SERVICE_ACCOUNT_NAME}')].data.token}"|base64 --decode)

kubectl config set-credentials ${SERVICE_ACCOUNT_NAME} --token="${TOKEN}"
kubectl config set-context --current --user="${SERVICE_ACCOUNT_NAME}"

echo "
    You are now authenticated against ${KUBE_CONTEXT} as ${SERVICE_ACCOUNT_NAME}
"

