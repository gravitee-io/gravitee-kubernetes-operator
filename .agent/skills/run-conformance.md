# Skill: Run the Gateway API Conformance Suite Locally

## When to Use

Use this skill **before** reporting any change under `controllers/gateway-api/`,
`internal/k8s/gateway*.go` or `test/conformance/` as done, working or verified.

`go build`, `make lint-sources` and `make unit` do not execute a single reconcile. They are
necessary, they are not evidence. The conformance suite is the only thing in this repo that
runs the controllers against a real API server and sends real traffic through a real
gateway. A gateway-api change that has not been through it is unverified, and must be
described that way.

If any step below cannot be completed, **stop and say which one**. Do not substitute a
green unit suite for a conformance run in what you report.

Passing the suite once is the entry price, not the finish line. How many runs are required,
and what else is required when the target is a known flaky test, is defined by
"Verifying a fix for a flaky test" in `.agent/rules/gateway-standards.md`. Read it before
reporting anything as fixed: on a suite with a non-trivial failure rate, three passes in a
row happen by chance often enough to be worthless as evidence on their own.

## Boundaries

Everything destructive in this skill is confined to the **kind cluster**. Outside it, nothing
is removed or reconfigured:

- Do not modify Docker credentials or `~/.docker`, do not run `docker login`, do not delete or
  prune images on the host.
- Do not delete files outside the repository, and do not `git reset`, amend or push.

## Steps

### 1. Cluster: reuse it, but prove it is the right shape

The suite needs the cluster created from `hack/kind/kind.conformance.yaml`: named
`gravitee`, host ports 80/443/9092, and a bind mount of `/tmp/coverage` that the coverage
PersistentVolume depends on. A `gravitee` cluster left over from `make start-cluster` looks
identical to `kind get clusters` and will fail later in a confusing way.

```bash
mkdir -p /tmp/coverage

if kind get clusters 2>/dev/null | grep -qx gravitee; then
  # Prove it is conformance shaped before trusting it.
  docker inspect gravitee-control-plane --format '{{range .Mounts}}{{.Source}}->{{.Destination}} {{end}}' \
    | grep -q '/tmp/coverage' || echo "WRONG CLUSTER SHAPE: recreate it"
else
  make start-conformance-cluster
fi
```

If the shape is wrong, `make delete-cluster && make start-conformance-cluster`. That is the
one case where deleting is correct.

### 2. Clean the cluster, do not delete it

Cleaning means removing everything a previous run left standing, so the next run starts from
the same state CI starts from. Conformance namespaces routinely hang in `Terminating`
because an object still carries a finalizer whose operator is gone, so strip finalizers
rather than waiting.

```bash
helm uninstall gko -n default --ignore-not-found

kubectl delete gatewayclass gravitee-gateway --ignore-not-found
kubectl delete gatewayclassparameters gravitee-gateway -n default --ignore-not-found

# Strip finalizers first, or the namespace delete never returns.
for kind in httproutes gateways kafkaroutes apiv4definitions; do
  kubectl get $kind -A -o custom-columns=NS:.metadata.namespace,N:.metadata.name --no-headers 2>/dev/null \
    | while read -r ns name; do
        kubectl patch "$kind" "$name" -n "$ns" --type=merge -p '{"metadata":{"finalizers":null}}' 2>/dev/null
      done
done

kubectl delete ns gateway-conformance-infra gateway-conformance-app-backend \
  gateway-conformance-web-backend --ignore-not-found --timeout=120s

# Generated resources are cluster wide and outlive their namespaces if orphaned.
kubectl delete cm,apiv4definitions -A -l managed-by=gravitee.io --ignore-not-found
```

Verify it is actually clean before continuing — a leftover Gateway that never reaches
`Programmed` blocks `NamespacesMustBeReady` for every later test:

```bash
kubectl get gateways,httproutes -A
```

### 3. Gateway image: pull if the registry is reachable, otherwise build it

`test/conformance/gateway-class-parameters.yaml` pins `image: gateway:latest`. Whichever
path is taken, the image must end up loaded into the cluster under **exactly** that tag.

Check registry auth without guessing — `docker login` state is not directly queryable, so
ask the registry:

```bash
VERSION=$(grep -A2 '^  image:' hack/apim.yaml | awk '/version:/{print $2}')
docker manifest inspect graviteeio.azurecr.io/apim-gateway:$VERSION >/dev/null 2>&1 \
  && echo AUTHENTICATED || echo NOT_AUTHENTICATED
```

**Authenticated** — use the CI path:

```bash
make kind-load-gateway
```

**Not authenticated** — build from the APIM checkout (`../gravitee-api-management`) and load
it by hand. `task docker-backend` tags with `DOCKER_TAG`, which defaults to `local`:

```bash
cd ../gravitee-api-management
task build-quick        # maven, slow; skip if the distribution target is already built
task docker-backend
docker tag graviteeio.azurecr.io/apim-gateway:local gateway:latest
kind load docker-image gateway:latest --name gravitee
```

### 4. Load balancer: cloud-provider-kind, never MetalLB

The suite reads the address the operator writes into `Gateway.status.addresses`, which comes
from the LoadBalancer Service, and then sends HTTP from the test binary **on the host**.

- CI uses MetalLB. On macOS a MetalLB address lives only inside the Docker network and is
  **not routable from the host**, so the suite cannot work against it.
- `make cloud-lb` on Darwin runs `sudo go tool cloud-provider-kind` — it needs root, prompts
  interactively, and stays in the foreground.

**An agent cannot start it. Do not try, do not run sudo.** Check, and if it is not running,
stop and ask the user to run `make cloud-lb` in their own terminal.

```bash
pgrep -f cloud-provider-kind >/dev/null && echo RUNNING || echo "NOT RUNNING - ask the user"
kubectl get ns metallb-system --ignore-not-found   # must be empty on macOS
```

Running is not the same as working. Confirm an address is actually issued and reachable
once a Gateway exists (step 6), before trusting a failure:

```bash
kubectl get gateways -A -o custom-columns=NAME:.metadata.name,ADDR:.status.addresses[0].value
curl -sS -o /dev/null -w '%{http_code}\n' http://<address>/ --max-time 5
```

If every test fails with a connection error, suspect the load balancer, not the operator.

### 5. Bootstrap, in this order

The coverage PVC must exist before the operator, because `operator.values.yaml` mounts it
and the manager pod stays `Pending` without it.

```bash
kubectl apply -f ./test/conformance/coverage.yaml

IMG=gko TAG=latest make docker-build-cover        # Dockerfile.cover, not the normal build
kind load docker-image gko:latest --name gravitee

helm upgrade --install gko helm/gko -n default \
  --set manager.image.repository=gko \
  --set manager.image.tag=latest \
  -f ./test/conformance/operator.values.yaml

kubectl rollout status deployment/gko-controller-manager -n default --timeout=120s

kubectl apply -f ./test/conformance/gateway-class-parameters.yaml
kubectl apply -f ./test/conformance/gateway-class.yaml
kubectl wait --for=condition=Accepted gatewayclass/gravitee-gateway --timeout=120s
```

`operator.values.yaml` sets `skipAPIDefinition: true` and `matchAcrossRoutes: true`, so this
exercises the merged ConfigMap mode. The other three programming modes are **not** covered
by a conformance run; changing them needs its own reasoning.

### 6. Run

Full suite, matching what CI runs:

```bash
make conformance
```

`GATEWAY_API_MATCH_ACROSS_ROUTES=true` decides whether `HTTPRouteMatchingAcrossRoutes` runs.
`conformance_test.go` reads it from the **test binary's** environment. It is *also* set on the
operator pod via `operator.values.yaml` — same name, different process, and setting one does
nothing for the other. Forgetting it on the test binary is what produced the stale partial
reports that used to sit in `report/`: the operator supported the feature, the test never ran,
and the report said "partial" anyway.

Which value to use depends on what the run is for:

| Goal | Command |
|---|---|
| Reproduce `Run Conformance Test` (pre-merge) | `CONFORMANCE_SKIP_STARVED_TESTS=true make conformance` |
| Reproduce a report run, or produce evidence | `GATEWAY_API_MATCH_ACROSS_ROUTES=true CONFORMANCE_RERUN=0 make conformance` |

Say which one you ran. `CONFORMANCE_SKIP_STARVED_TESTS` skips `HTTPRouteWeight` and
`HTTPRouteRedirectPortAndScheme`, which starve on a small runner; the pre-merge job sets it and
`job-conformance-report` deliberately does not.

To iterate on one test, set `opts.RunTest` in `conformance_test.go` temporarily. Two tests
share the HTTPRoute name `request-header-modifier` and are the known reproducer for
delete/recreate races — run them together, in order:

```
HTTPRouteRequestHeaderModifier then HTTPRouteBackendRequestHeaderModifier
```

### 7. Report what actually happened

- Give failing test names and counts, from `/tmp/junit/reports/conformance.xml`, not a
  summary impression.
- A partial or filtered run is a partial result. Say which tests ran.
- If the suite could not be run at all, say that plainly and say why. "Builds and unit tests
  pass" is not a substitute and must not be presented as one.
- For CI results, the failing tests and the cluster dump can be pulled from the CircleCI API
  (`/api/v2/project/gh/gravitee-io/gravitee-kubernetes-operator/<build>/tests`, and the
  `cluster-dump.txt` artifact — fetch it with `curl -L`, it redirects).

## Known Traps

| Trap | Symptom |
|---|---|
| Cluster not created from `kind.conformance.yaml` | coverage PV unbound, manager pod `Pending` |
| `coverage.yaml` applied after the operator | manager pod `Pending` on a missing PVC |
| Normal `docker-build` instead of `docker-build-cover` | suite runs, coverage step fails at the end |
| Image not tagged exactly `gateway:latest` | gateway pods `ImagePullBackOff` |
| MetalLB on macOS | every request test fails with connection refused/timeout |
| Namespace stuck `Terminating` | a finalizer is still set; strip it, do not wait |
| Leftover Gateway from a previous run | `NamespacesMustBeReady` blocks for 600s per test |
