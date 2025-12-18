# Runbook: GKO Upgrade Lifecycle Test

Validates the full lifecycle of GKO-managed resources across an APIM upgrade from 4.10 to latest.

## Prerequisites

- Docker running
- `kind`, `kubectl`, `helm` installed
- Access to `graviteeio.azurecr.io` registry (or local images pre-loaded)
- GKO operator image built and available

## Variables

Set these before starting. Adjust versions and registries as needed.

```bash
# APIM
export APIM_IMAGE_REGISTRY=graviteeio.azurecr.io
export APIM_410_TAG=4.10.x-latest            # initial APIM image tag
export APIM_410_CHART=4.10.*          # initial Helm chart version
export APIM_LATEST_TAG=master-latest  # upgrade target image tag
export APIM_LATEST_CHART=4.11.*       # upgrade target chart version

# GKO
export GKO_410_IMAGE=graviteeio/kubernetes-operator:4.10.9
export GKO_LATEST_IMAGE=graviteeio/kubernetes-operator:latest

# Paths (from project root)
export PROJECT_ROOT=$(pwd)
export APIM_VALUES=hack/kind/apim/values.yaml
export KIND_CONFIG=hack/kind/kind.yaml
export PKI_JWT=examples/usecase/subscribe-to-jwt-plan/pki
export PKI_MTLS=examples/usecase/subscribe-to-mtls-plan/pki
export RUNBOOK=runbooks/upgrade-lifecycle

# Branch
export GIT_CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
export GIT_410_BRANCH=4.10.x
```

---

## Phase 1 — Create a 4.10 Cluster

### 1.1 Create the kind cluster

```bash
kind create cluster --config $KIND_CONFIG
```

### 1.2 Load and tag APIM 4.10 images

```bash
for component in apim-gateway apim-management-api apim-management-ui; do
  docker pull ${APIM_IMAGE_REGISTRY}/${component}:${APIM_410_TAG}
  docker tag  ${APIM_IMAGE_REGISTRY}/${component}:${APIM_410_TAG} gravitee-${component}:dev
  kind load docker-image gravitee-${component}:dev --name gravitee
done

# MongoDB
docker pull mongo:7.0.30-jammy
kind load docker-image mongo:7.0.30-jammy --name gravitee

```

### 1.3 Create TLS secrets and deploy APIM

```bash
# TLS secret for gateway TLS listener
kubectl create secret tls tls-server \
  --cert=${PKI_MTLS}/server.crt --key=${PKI_MTLS}/server.key \
  --dry-run=client -o yaml | kubectl apply -f -

# Install APIM 4.10 with explicit image tags
helm repo add graviteeio https://helm.gravitee.io
helm repo update graviteeio

helm upgrade --install apim oci://graviteeio.azurecr.io/helm/apim3 \
  -f $APIM_VALUES \
  --version "$APIM_410_CHART" \
  --set gateway.image.repository=gravitee-apim-gateway \
  --set gateway.image.tag=dev \
  --set api.image.repository=gravitee-apim-management-api \
  --set api.image.tag=dev \
  --set ui.image.repository=gravitee-apim-management-ui \
  --set ui.image.tag=dev

# Wait for APIM to be ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=apim3 --timeout=360s
```

### 1.4 Install GKO 4.10 in the cluster

```bash
# Load GKO image
docker pull $GKO_410_IMAGE
kind load docker-image $GKO_410_IMAGE --name gravitee

# Install GKO via Helm
helm upgrade --install gko helm/gko \
  --set manager.image.repository=$(echo $GKO_410_IMAGE | cut -d: -f1) \
  --set manager.image.tag=$(echo $GKO_410_IMAGE | cut -d: -f2) \
  --set manager.metrics.enabled=false

kubectl wait --for=condition=ready pod -l control-plane=controller-manager --timeout=120s
```

### 1.5 Create the Management Context and JWT secret

```bash
kubectl apply -f $RUNBOOK/crds/phase-1/context.yaml
kubectl apply -f $RUNBOOK/crds/phase-1/jwt-secret.yaml
```

---

## Phase 2 — Create Initial Resources

### 2.1 Create API, Application, and Subscription

```bash
kubectl apply -f $RUNBOOK/crds/phase-2/api-legacy.yaml
kubectl apply -f $RUNBOOK/crds/phase-2/app-legacy.yaml

kubectl wait --for=condition=Accepted apiv4definitions/legacy-api --timeout=30s
kubectl wait --for=condition=Accepted applications.gravitee.io/legacy-app --timeout=30s

kubectl apply -f $RUNBOOK/crds/phase-2/sub-legacy-legacy-jwt.yaml
kubectl wait --for=condition=Accepted subscriptions.gravitee.io/legacy-legacy-jwt --timeout=30s
```

### 2.2 Verify the API works

```bash
# Generate a JWT token
export TOKEN=$($PKI_JWT/get_token.sh legacy-client)

# Call the gateway with the token
curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:30082/legacy
# Expected: 200

# Without token — should be rejected
curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" http://localhost:30082/legacy
# Expected: 401
```

---

## Phase 3 — Upgrade APIM and GKO

### 3.1 Load new APIM images

```bash
# TODO remove this once APIM images are published to the public registry
export APIM_LATEST_TAG=local 
export APIM_IMAGE_REGISTRY=graviteeio
for component in apim-gateway apim-management-api apim-management-ui; do
  # TODO uncomment when APIM images are published to the public registry
  # docker pull ${APIM_IMAGE_REGISTRY}/${component}:${APIM_LATEST_TAG}
  docker tag  ${APIM_IMAGE_REGISTRY}/${component}:${APIM_LATEST_TAG} gravitee-${component}:dev
  kind load docker-image gravitee-${component}:dev --name gravitee
done
```

### 3.2 Upgrade APIM — database must not restart

The MongoDB StatefulSet must **not** be affected by the Helm upgrade. We achieve this
by pinning the MongoDB subchart values so Helm sees no diff on that resource.

```bash
helm upgrade apim oci://graviteeio.azurecr.io/helm/apim3 \
  -f $APIM_VALUES \
  --version "$APIM_LATEST_CHART" \
  --set gateway.image.repository=gravitee-apim-gateway \
  --set gateway.image.tag=dev \
  --set api.image.repository=gravitee-apim-management-api \
  --set api.image.tag=dev \
  --set ui.image.repository=gravitee-apim-management-ui \
  --set ui.image.tag=dev \
  --set mongodb.enabled=true \
  --set mongodb.architecture=standalone \
  --set mongodb.auth.enabled=false \
  --reuse-values
```

> `--reuse-values` prevents Helm from recomputing the MongoDB subchart template,
> which would trigger a StatefulSet rollout even if nothing changed.

Wait for APIM pods (not MongoDB) to restart:

```bash
kubectl rollout status deployment/apim-apim3-gateway --timeout=120s
kubectl rollout status deployment/apim-apim3-api --timeout=120s
```

Verify MongoDB pod did **not** restart:

```bash
kubectl get pod -l app.kubernetes.io/component=mongodb -o jsonpath='{.items[0].status.containerStatuses[0].restartCount}'
# Expected: 0
```

### 3.3 Verify existing resources survived the upgrade

```bash
TOKEN=$($PKI_JWT/get_token.sh legacy-client)
curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:30082/legacy
# Expected: 200
```

### 3.4 Upgrade GKO

```bash
#TODO remove this once GKO is published to the public registry
export GKO_LATEST_IMAGE=gko:dev
# TODO uncomment when GKO is published to the public registry
# docker pull $GKO_LATEST_IMAGE
kind load docker-image $GKO_LATEST_IMAGE --name gravitee

helm upgrade gko helm/gko \
  --set manager.image.repository=$(echo $GKO_LATEST_IMAGE | cut -d: -f1) \
  --set manager.image.tag=$(echo $GKO_LATEST_IMAGE | cut -d: -f2)

kubectl rollout status deployment -l control-plane=controller-manager --timeout=120s
```

---

## Phase 4 — Add mTLS Plan and Subscribe

### 4.1 Create the client certificate secret

```bash
kubectl create secret tls tls-client \
  --cert=${PKI_MTLS}/client.crt --key=${PKI_MTLS}/client.key \
  --dry-run=client -o yaml | kubectl apply -f -
```

### 4.2 Update the API with an mTLS plan (alongside JWT)

```bash
kubectl apply -f $RUNBOOK/crds/phase-4/api-legacy-with-mtls.yaml
kubectl wait --for=condition=Accepted apiv4definitions/legacy-api --timeout=30s
```

### 4.3 Update the Application with a client certificate

```bash
kubectl apply -f $RUNBOOK/crds/phase-4/app-legacy-with-cert.yaml
kubectl wait --for=condition=Accepted applications.gravitee.io/legacy-app --timeout=30s
```

### 4.4 Update the existing JWT subscription with an end date (tomorrow)

```bash
TOMORROW=$(date -u -v+1d '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || date -u -d '+1 day' '+%Y-%m-%dT%H:%M:%SZ')
sed "s/REPLACE_WITH_TOMORROW/$TOMORROW/" \
  $RUNBOOK/crds/phase-4/sub-legacy-legacy-jwt-update.yaml \
  | kubectl apply -f -
kubectl wait --for=condition=Accepted subscriptions.gravitee.io/legacy-legacy-jwt --timeout=30s
```

### 4.5 Subscribe to the mTLS plan

```bash
kubectl apply -f $RUNBOOK/crds/phase-4/sub-legacy-legacy-mtls.yaml
kubectl wait --for=condition=Accepted subscriptions.gravitee.io/legacy-legacy-mtls --timeout=30s
```

### 4.6 Verify mTLS works via the HTTPS gateway port

```bash
# Without client cert — rejected
curl -sk -o /dev/null -w "---\nstatus %{http_code}\n---\n" https://localhost:30084/legacy
# Expected: 401

# With client cert — accepted
curl -sk -o /dev/null -w "---\nstatus %{http_code}\n---\n" \
  --cert ${PKI_MTLS}/client.crt \
  --key ${PKI_MTLS}/client.key \
  https://localhost:30084/legacy
# Expected: 200
```

JWT on the HTTP port should still work:

```bash
TOKEN=$($PKI_JWT/get_token.sh legacy-client)
curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:30082/legacy
# Expected: 200
```

---

## Phase 5 — Cross-Subscriptions

### 5.1 Create a new API and a new Application

```bash
kubectl apply -f $RUNBOOK/crds/phase-5/api-weather.yaml
kubectl apply -f $RUNBOOK/crds/phase-5/app-mobile.yaml

kubectl wait --for=condition=Accepted apiv4definitions/weather-api --timeout=30s
kubectl wait --for=condition=Accepted applications.gravitee.io/mobile-app --timeout=30s
```

### 5.2 Subscribe legacy-app to weather-api

```bash
kubectl apply -f $RUNBOOK/crds/phase-5/sub-legacy-weather.yaml
kubectl wait --for=condition=Accepted subscriptions.gravitee.io/legacy-weather-jwt --timeout=30s
```

### 5.3 Subscribe mobile-app to legacy-api

```bash
kubectl apply -f $RUNBOOK/crds/phase-5/sub-mobile-legacy.yaml
kubectl wait --for=condition=Accepted subscriptions.gravitee.io/mobile-legacy-jwt --timeout=30s
```

### 5.4 Subscribe mobile-app to weather-api

```bash
kubectl apply -f $RUNBOOK/crds/phase-5/sub-mobile-weather.yaml
kubectl wait --for=condition=Accepted subscriptions.gravitee.io/mobile-weather-jwt --timeout=30s
```


### 5.5 Verify both APIs work

```bash
TOKEN=$($PKI_JWT/get_token.sh legacy-client)

curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:30082/legacy
# Expected: 200

TOKEN=$($PKI_JWT/get_token.sh legacy-client)

curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:30082/weather
# Expected: 200

TOKEN=$($PKI_JWT/get_token.sh mobile-app)

curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:30082/weather
# Expected: 200
```

---

## Phase 6 — Cleanup

Delete in reverse dependency order: subscriptions, then apps, then APIs.

```bash
# Subscriptions
kubectl delete subscriptions.gravitee.io/mobile-legacy-jwt
kubectl delete subscriptions.gravitee.io/mobile-weather-jwt
kubectl delete subscriptions.gravitee.io/legacy-weather-jwt
kubectl delete subscriptions.gravitee.io/legacy-legacy-mtls
kubectl delete subscriptions.gravitee.io/legacy-legacy-jwt


# Applications
kubectl delete applications.gravitee.io/mobile-app
kubectl delete applications.gravitee.io/legacy-app

# APIs
kubectl delete apiv4definitions/weather-api
kubectl delete apiv4definitions/legacy-api

# Context and secrets
kubectl delete -f $RUNBOOK/crds/phase-1/context.yaml
kubectl delete -f $RUNBOOK/crds/phase-1/jwt-secret.yaml
kubectl delete secret tls-client tls-server
```

Verify everything is gone:

```bash
kubectl get apiv4definitions,applications.gravitee.io,subscriptions.gravitee.io
# Expected: No resources found

curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" http://localhost:30082/legacy
# Expected: 404

curl -s -o /dev/null -w "---\nstatus %{http_code}\n---\n" http://localhost:30082/weather
# Expected: 404
```

---

## Resource Summary

| Phase | Resource | File |
|-------|----------|------|
| 1 | JWT public key secret | `crds/phase-1/jwt-secret.yaml` |
| 1 | ManagementContext | `crds/phase-1/context.yaml` |
| 2 | Legacy API (JWT) | `crds/phase-2/api-legacy.yaml` |
| 2 | Legacy App | `crds/phase-2/app-legacy.yaml` |
| 2 | JWT Subscription | `crds/phase-2/sub-legacy-legacy-jwt.yaml` |
| 4 | Legacy API (JWT + mTLS) | `crds/phase-4/api-legacy-with-mtls.yaml` |
| 4 | Legacy App (with cert) | `crds/phase-4/app-legacy-with-cert.yaml` |
| 4 | JWT Sub update (endingAt) | `crds/phase-4/sub-legacy-legacy-jwt-update.yaml` |
| 4 | mTLS Subscription | `crds/phase-4/sub-legacy-legacy-mtls.yaml` |
| 5 | Weather API (JWT) | `crds/phase-5/api-weather.yaml` |
| 5 | Mobile App | `crds/phase-5/app-mobile.yaml` |
| 5 | Legacy → Weather JWT sub | `crds/phase-5/sub-legacy-weather.yaml` |
| 5 | Mobile → Legacy JWT sub | `crds/phase-5/sub-mobile-legacy.yaml` |
| 5 | Mobile → Weather JWT sub | `crds/phase-5/sub-mobile-weather.yaml` |