# Migrating a v2 API managed by GKO to an ApiV4Definition

v2 APIs (`ApiDefinition`) will be removed in October 2027. This guide moves a v2 API that GKO manages to an
`ApiV4Definition`, **without losing its subscriptions**: the API keeps its id, each plan keeps its id, and API keys
keep working.

The migration itself happens in APIM. GKO then takes over the migrated API through an `ApiV4Definition` that names
the API and each of its plans by id.

```
ApiDefinition (v2) ──▶ detached from GKO ──▶ migrated to v4 in APIM ──▶ ApiV4Definition with the same ids
```

## Before you start

- **GKO 4.12.29, 4.13.6, 4.14.0 or later.** Earlier 4.12 and 4.13 versions ignore the plan ids of an
  `ApiV4Definition` that sets `spec.id`: the first apply replaces every plan and closes its subscriptions.
- **APIM 4.11.33, 4.12.27, 4.13.4 or later.** These versions refuse a plan id that belongs to another API.
- Subscriptions made in APIM (Console, Developer Portal, Management API) are kept. Subscriptions managed with GKO
  `Subscription` resources are not covered by this guide: GKO refuses to delete an `ApiDefinition` that a
  `Subscription` resource references.

### Why every id matters

APIM matches plans **by id**. When GKO applies an `ApiV4Definition`, for each plan of the API in APIM:

| The manifest has… | APIM… |
|---|---|
| a plan with this plan's id | updates the plan in place. Subscriptions are kept |
| no plan with this id (no id, or an id APIM does not know) | creates a new plan, then deletes the existing one after closing its subscriptions. API keys stop working |

So every plan you keep must appear in the manifest with **its id in that environment**, like `spec.id`. Plan ids
differ between environments: a plan migrated in staging and the same plan migrated in production have different ids.

The migration keeps every id, so you can read them from the v2 API, before migrating.

**The first apply decides.** Once GKO has applied an `ApiV4Definition`, it stores the plan ids in `status.plans` and
keeps sending them. Changing `plans.<key>.id` later has no effect, and a plan that was replaced does not come back.

### Commands used in this guide

The commands use `kubectl`, `curl`, `jq` and `yq`, and these variables. Set them for each environment:

```bash
export NS=apis                     # namespace of your API resources
export V2_NAME=my-api              # name of the ApiDefinition resource
export MAPI="https://apim.example.com/management/v2/organizations/DEFAULT/environments/DEFAULT"
export TOKEN=<personal access token of a user who can update the API>
```

## Choose your path

| After the migration, the v4 API… | Path |
|---|---|
| Needs no change | [Path A: migrate each environment](#path-a-migrate-each-environment) |
| Needs a change, made in a non-production environment and then promoted | [Path B: fix in non-production, then promote](#path-b-fix-in-non-production-then-promote) |

Both paths use the four steps below, in each environment.

## Step 1: note the ids of the API

Read them while GKO still manages the v2 API. The migration does not change them.

```bash
export API_ID=$(kubectl get apidefinitions.gravitee.io "$V2_NAME" -n "$NS" -o jsonpath='{.status.id}')
export CROSS_ID=$(kubectl get apidefinitions.gravitee.io "$V2_NAME" -n "$NS" -o jsonpath='{.status.crossId}')
echo "$API_ID $CROSS_ID"
curl -s -H "Authorization: Bearer $TOKEN" \
  "$MAPI/apis/$API_ID/plans?statuses=STAGING,PUBLISHED,DEPRECATED" | jq -r '.data[] | "\(.id)  \(.status)  \(.name)"'
```

```
4fb7eef4-0c35-4b51-b7ee-f40c35cb51a2  PUBLISHED  API key plan
9a1c7d02-5e4f-4a7b-9c7d-025e4f6a7b11  PUBLISHED  Keyless
```

Keep this list: it is what the `ApiV4Definition` must carry in this environment.

## Step 2: check that the API can be migrated

A dry run changes nothing. Run it before detaching anything, so an API that cannot be migrated stays as it is.

**In the Console:** open the API and click **Migrate to v4** on its main page. The dialog shows the check first;
click **Cancel** to stop there.

**With the Management API:**

```bash
curl -s -X POST -H "Authorization: Bearer $TOKEN" "$MAPI/apis/$API_ID/_migrate?mode=DRY_RUN" | jq '{state, issues}'
```

| `state` | Next |
|---|---|
| `MIGRATABLE` | Go on |
| `CAN_BE_FORCED` | Read the `issues`. If you accept them, go on and migrate with `?mode=FORCE` in step 4 |
| `IMPOSSIBLE` | Fix what the `issues` list in the v2 API, then run the dry run again |

## Step 3: detach the v2 API from GKO

Deleting an `ApiDefinition` normally deletes the API in APIM. Remove its finalizer first: GKO then lets the resource go
and leaves the API in APIM, deployed, with its plans and subscriptions. Traffic is not affected.

```bash
kubectl patch apidefinitions.gravitee.io "$V2_NAME" -n "$NS" --type=merge -p '{"metadata":{"finalizers":null}}'
kubectl delete apidefinitions.gravitee.io "$V2_NAME" -n "$NS"
```

Check that the API is still in APIM:

```bash
curl -s -H "Authorization: Bearer $TOKEN" "$MAPI/apis/$API_ID" | jq '{id, name, definitionVersion}'
```

Detach before migrating: a v2 resource still managed by GKO could push its v2 definition back over the migrated API.

> [!WARNING]
> If a GitOps tool (Argo CD, Flux) applies your manifests, remove the finalizer **before** removing the
> `ApiDefinition` from your repository, and make sure the tool does not create it again. If the tool deletes the
> resource while the finalizer is still there, GKO deletes the API from APIM, with its plans and subscriptions.

## Step 4: migrate the API in APIM

The migration turns the v2 API into a v4 API. The API keeps its id, each plan keeps its id, and subscriptions stay
accepted.

**In the Console:** open the API, click **Migrate to v4**, then **Continue** and **Start Migration**.

**With the Management API** (add `?mode=FORCE` if the dry run answered `CAN_BE_FORCED`):

```bash
curl -s -X POST -H "Authorization: Bearer $TOKEN" "$MAPI/apis/$API_ID/_migrate" | jq '{state, issues}'
```

The migration answers `state: MIGRATED`.

## Path A: migrate each environment

When the migrated API needs no change, each environment exports its own manifest, which already carries that
environment's ids. Do this in every environment, starting with non-production:

1. Run [steps 1 to 4](#step-1-note-the-ids-of-the-api).

2. Export the migrated API **from the same environment**. In the Console: **Export**, tab **CRD API Definition**,
   **Export**. Or:

   ```bash
   curl -s -H "Authorization: Bearer $TOKEN" "$MAPI/apis/$API_ID/_export/crd" > my-api.yaml
   ```

   The export carries `spec.id`, `spec.crossId` and every `plans.<key>.id` of this environment.

3. Name the resource and point it at the environment's `ManagementContext` (the export has no `contextRef`). Use the
   context the v2 API used:

   ```bash
   yq -i '.metadata.name = "my-api" | .metadata.namespace = env(NS)
          | .spec.contextRef = {"name": "apim-context", "namespace": env(NS)}' my-api.yaml
   ```

4. Apply it, and [check the result](#check-the-result):

   ```bash
   kubectl apply -f my-api.yaml
   ```

If one manifest is shared between environments (a Helm chart, a Kustomize base), the exports differ by their ids:
make `spec.id`, `spec.crossId` and every `plans.<key>.id` per-environment values, as in
[Path B](#path-b-fix-in-non-production-then-promote), production step 2.

## Path B: fix in non-production, then promote

When the migrated API needs a change, make it once in non-production, through GKO, then promote the fixed manifest.

The migrated API keeps the Kubernetes origin of the v2 API, so the Console shows it read-only: make the change in the
manifest.

### In non-production

1. Run [steps 1 to 4](#step-1-note-the-ids-of-the-api).
2. Export it, name it and set its `contextRef`, then apply it, as in [Path A](#path-a-migrate-each-environment),
   steps 2 to 4. GKO now manages the v4 API, with the non-production ids.
3. Fix the API in the manifest and apply it again. Test it.

### In production

Between the migration and the apply, production runs the migrated API **without the fix**. Prepare everything first,
then migrate and apply back to back.

**Prepare, while the v2 API still runs unchanged:**

1. Run [step 1](#step-1-note-the-ids-of-the-api) and [step 2](#step-2-check-that-the-api-can-be-migrated) **in
   production**: note the production ids, and check that the API can be migrated.

2. Copy the fixed non-production manifest. Change what differs between environments as for any promotion
   (`contextRef`, endpoints, hosts…), then replace every id with the production one:

   - `spec.id`: the production API id (`$API_ID`).
   - `spec.crossId`: the production cross id (`$CROSS_ID`).
   - `plans.<key>.id`: for each plan, the id of the production plan it stands for. Match the plans by name: the keys
     come from the non-production export.

   ```bash
   yq -i '.spec.id = env(API_ID) | .spec.crossId = env(CROSS_ID)
          | .spec.plans["api-key-plan"].id = "<production id of the API key plan>"
          | .spec.plans["keyless"].id = "<production id of the Keyless plan>"' my-api.prod.yaml
   ```

   Every production plan from step 1 that you want to keep must be in the manifest. A production plan missing from
   it is deleted, and its subscriptions are closed.

   Check that no non-production id is left:

   ```bash
   yq '{"api": .spec.id, "crossId": .spec.crossId, "plans": (.spec.plans | map_values(.id))}' my-api.prod.yaml
   ```

**Switch, back to back:**

3. [Detach the v2 API](#step-3-detach-the-v2-api-from-gko). Traffic is not affected.
4. [Migrate it](#step-4-migrate-the-api-in-apim). From now on, production runs the migrated API without the fix.
5. Apply the prepared manifest right away, then [check the result](#check-the-result):

   ```bash
   kubectl apply -f my-api.prod.yaml
   ```

## Check the result

The plans GKO sent are the ones you noted in [step 1](#step-1-note-the-ids-of-the-api):

```bash
kubectl get apiv4definitions.gravitee.io my-api -n "$NS" -o jsonpath='{.status.plans}'
curl -s -H "Authorization: Bearer $TOKEN" \
  "$MAPI/apis/$API_ID/plans?statuses=STAGING,PUBLISHED,DEPRECATED" | jq -r '.data[] | "\(.id)  \(.status)  \(.name)"'
```

Subscriptions are still accepted:

```bash
curl -s -H "Authorization: Bearer $TOKEN" "$MAPI/apis/$API_ID/subscriptions?statuses=ACCEPTED" | jq '.pagination.totalCount'
```

## When an id is wrong

| The manifest carries… | Same APIM installation | Separate APIM installation |
|---|---|---|
| The right ids for this environment | Plans updated in place, subscriptions kept | Same |
| No plan id | New plan, existing plan deleted, subscriptions closed | Same |
| A plan id from another environment | Refused at apply time, nothing written | New plan with that id, existing plan deleted, subscriptions closed |
| An API id from another environment | Refused at apply time, nothing written | A second API is created next to the migrated one, unless APIM refuses it because its path is already used |
| No entry for an existing plan | That plan deleted, subscriptions closed | Same |

A refusal comes from the APIM dry run that GKO runs when you apply the resource. Once the API is migrated, you can
trigger it without applying:

```bash
kubectl apply --dry-run=server -f my-api.prod.yaml
```

A dry run that passes does not prove the ids are right: in a separate installation, ids from another environment are
unknown, not refused. Compare the manifest with [step 1](#step-1-note-the-ids-of-the-api) before applying.
