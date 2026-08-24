# Skill: Add a New Controller

## When to Use

Use this skill when adding a reconciler for a CRD that syncs a resource to APIM.

Read [`.agent/rules/operator-standards.md`](../rules/operator-standards.md) first — it owns the
reconcile, status and codegen invariants. This skill is the ordered procedure and the shapes to
copy.

## The shape of the work

Every APIM controller is the same four layers, and they must be built bottom-up. Most of the real
work is in layers 1 and 2, not in the controller itself:

```
internal/apim/model/<resource>.go        1. the Automation API payload + its mapper
internal/apim/service/<resource>.go      2. the HTTP calls
internal/apim/apim.go                    3. one field on the APIM facade
controllers/apim/<resource>/             4. the reconciler (thin)
```

Pick one existing resource and read it end to end before you start. `sharedpolicygroups` is the
cleanest reference: no legacy Management API path, no bespoke wire shape, no controller-specific
ref resolution.

---

## 1. The DTO — `internal/apim/model/<resource>.go`

The DTO **is the Automation API request and response body**. It is not a drift helper; drift
merely reuses it. `json` tags must match the Automation API schema exactly.

```go
package model

type MyResourceDTO struct {
	HRID        string    `json:"hrid,omitempty" drift:"ignore"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Steps       []StepDTO `json:"steps,omitempty" drift:"empty-is-nil"`
}

// MyResourceState is what a GET returns: the payload plus server-side state.
type MyResourceState struct {
	MyResourceDTO `json:",inline"`
	Errors        status.Errors `json:"errors,omitempty"`
}

func ToMyResourceDTO(crd myresource.Type, hrid string) MyResourceDTO {
	dto := mapViaJSON[MyResourceDTO](crd)
	dto.HRID = hrid
	return dto
}
```

- The `(crd, hrid)` mapper signature is what every post-4.12 resource uses
  (`ToDictionaryDTO`, `ToPortalDTO`, `ToPortalLinkDTO`) — the HRID is derived by the caller from
  the object's namespace and name and stamped onto the DTO, never read from the spec.
- `mapViaJSON[T]` (`internal/apim/model/json_map.go`) is a JSON round-trip. Use it whenever the
  CRD model and the wire shape agree, which is the normal case.
- Write explicit mapping only where they genuinely diverge — `ToAPIV4DTO` flattens the spec's plan
  and page maps into sorted slices; `SubscriptionDTO.ToAutomation()` renames `apiId` → `apiHrid`.
- Add `drift:` tags now, while you know which fields APIM owns. `ignore` identity and
  server-managed fields (`hrid`, `crossId`, `id`, `status`); tag the rest per the table in
  [AGENTS.md](../../AGENTS.md#drift-tags-on-the-dto). They do not affect serialization.
- If APIM returns extra state alongside the payload, add a `MyResourceState` that embeds the DTO
  inline (see `GroupState`).

## 2. The service — `internal/apim/service/<resource>.go`

One file per resource, embedding `*client.Client` so the target builders are in scope.

```go
package service

const myResourcesPath = "/my-resources"

type MyResources struct {
	*client.Client
}

func NewMyResources(client *client.Client) *MyResources {
	return &MyResources{Client: client}
}

func (svc *MyResources) CreateOrUpdate(obj *v1alpha1.MyResource) (*myresource.Status, error) {
	return svc.createOrUpdate(obj, false)
}

func (svc *MyResources) DryRunCreateOrUpdate(obj *v1alpha1.MyResource) (*myresource.Status, error) {
	return svc.createOrUpdate(obj, true)
}

func (svc *MyResources) createOrUpdate(obj *v1alpha1.MyResource, dryRun bool) (myresource.Status, error) {
	url := svc.AutomationTarget(myResourcesPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := model.ToMyResourceDTO(obj.Spec.MyResource, refs.NewNamespacedNameFromObject(obj).HRID())
	status := &myresource.Status{}

	if err := svc.HTTP.Put(url.String(), dto, status); err != nil {
		return *status, err
	}

	k8s.AddAutomationAPIManagedCondition(obj)

	return *status, nil
}

func (svc *MyResources) Delete(obj *v1alpha1.MyResource) error {
	hrid := refs.NewNamespacedNameFromObject(obj).HRID()
	url := svc.AutomationTarget(myResourcesPath).WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID For test purposes only.
func (svc *MyResources) GetByHRID(hrid string) (*model.MyResourceState, error) {
	url := svc.AutomationTarget(myResourcesPath).WithPath(hrid)
	state := new(model.MyResourceState)
	if err := svc.HTTP.Get(url.String(), state); err != nil {
		return nil, err
	}
	return state, nil
}
```

That is the whole service for a resource introduced after 4.12 — `dictionaries.go`, `portals.go`,
`portallinks.go` and `portallistings.go` are all this short. Note what is *absent*: no `GetByID`
and no UUID branch. The HRID is derived from the object on every call, so there is no identity to
resolve beforehand.

**Only if the resource must also support CRs created before 4.12** do you need the migration
branch, as `applications.go` and `sharedpolicygroups.go` do:

```go
	setHridWithUUID := obj.GetID() != "" && !k8s.IsAutomationAPIManaged(obj)
	if setHridWithUUID {
		dto.HRID = obj.GetID()
		url = url.WithQueryParam("hridContainsUUID", strconv.FormatBool(true))
	}
	// ... Put ...
	if !setHridWithUUID {
		k8s.AddAutomationAPIManagedCondition(obj)
	}
```

with the matching `Delete` helper and a `GetByID` against `EnvV1Target`/`EnvV2Target` for the
UUID-addressed read. A genuinely new resource needs none of it — do not carry it over by reflex.

Rules that are easy to get wrong:

- **`AutomationTarget` only.** The four services added since 4.12 do not reference `EnvV1Target` or
  `EnvV2Target` at all; those appear only in `apis.go`, `applications.go`, `subscriptions.go` and
  `sharedpolicygroups.go`, for v2 APIs and legacy UUID reads.
- **Always ship the dry-run twin.** Every resource service has one, because the admission webhook
  needs it. It must be the same code path with `dryRun=true`, not a separate implementation.
- **Return the APIM error unwrapped**, including 404. `internal/errors.IsNotFound` and the drift
  policies depend on the `ServerError.StatusCode` surviving. Wrapping happens in the controller,
  with `gerrors.NewControlPlaneError`.
- **The `AutomationAPIManaged` condition is set by the service, not the controller.**
- `GetByHRID` exists for drift and for tests. Mark it as such, as the other services do.

## 3. Register on the facade — `internal/apim/apim.go`

```go
type APIM struct {
	// ...
	MyResources *service.MyResources
}

// in FromContext:
	MyResources: service.NewMyResources(c),
```

## 4. The CRD side

The type needs `ContextRef()` and a status carrying whatever the Automation API returns (`id`,
`orgId`, `envId`).

`PopulateIDs(mCtx, automationAPIManaged)` is part of the interface, so you must declare it — but
for a post-4.12 resource it is an **empty no-op**, exactly as on `Dictionary` and `Portal`:

```go
func (d *MyResource) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}
```

Only the types that predate the Automation API — `Application`, `ApiV4Definition`,
`SharedPolicyGroup`, `Group` — give it a real body, because those services must choose between an
HRID and a stored UUID. See `Application.PopulateIDs` in `api/v1alpha1/application_types.go` if you
are extending one of them.

---

## 5. The controller — `controllers/apim/<resource>/<resource>_controller.go`

```go
type Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
}

// +kubebuilder:rbac:groups=gravitee.io,resources=myresources,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=myresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=myresources/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	obj := &v1alpha1.MyResource{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)
	k8s.ResetConditionsExceptAutomationAPI(obj)

	// The APIM sync runs against this copy; the live object only carries finalizers and
	// annotations through util.CreateOrUpdate. Status is copied back afterwards.
	dc := obj.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, r.Client, obj, func() error {
		util.AddFinalizer(obj, core.MyResourceFinalizer)
		k8s.AddAnnotation(obj, core.LastSpecHashAnnotation, hash.Calculate(&obj.Spec))

		if obj.IsBeingDeleted() {
			if err := template.ReleaseReferences(ctx, obj); err != nil {
				return err
			}
		} else if err := template.Compile(ctx, dc, true); err != nil {
			obj.Status.ProcessingStatus = core.ProcessingStatusFailed
			return err
		}

		var err error
		if obj.IsBeingDeleted() {
			err = events.Record(event.Delete, obj, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(obj, core.MyResourceFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, obj, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(obj); err != nil {
		return ctrl.Result{}, err
	}
	obj.SetFinalizers(dc.GetFinalizers())

	if err == nil {
		log.InfoEndReconcile(ctx, obj)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, obj)
	}

	if err := internal.UpdateStatusFailure(ctx, obj, err); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, obj)
		return ctrl.Result{RequeueAfter: requeueAfterTime}, err
	}

	log.ErrorAbortingReconcile(ctx, err, obj)
	return ctrl.Result{}, nil
}
```

There is no shared reconcile wrapper — every controller inlines this. Copy it rather than
inventing a variation. The only sanctioned deviation is `controllers/apim/apidefinition/`, where
v2 and v4 share a package-level `Reconcile()` because they implement a common interface.

If the resource cannot be synced without a `ManagementContext`, bail out before doing any work,
the way `application` does:

```go
if obj.Spec.Context == nil {
	log.ErrorAbortingReconcile(ctx, fmt.Errorf(
		"no context is provided, no attempt will be made to sync with APIM"), obj)
	return ctrl.Result{}, nil
}
```

### `SetupWithManager`

```go
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MyResource{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.MyResourceContextField))

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("myresources")).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("myresources"))
	}

	return newController.Complete(r)
}
```

Add one `Watches` per kind the spec can reference. The `watch.Interface`
(`internal/watch/watch.go`) already exposes `WatchContexts`, `WatchResources`, `WatchGroups`,
`WatchNotifications`, `WatchSharedPolicyGroups`, `WatchApis`, `WatchPortals`, `WatchTLSSecret` and
`WatchTemplatingSource`. Each takes a `search.IndexField`, so a new reference kind means a new
indexer in `internal/search/indexer.go` — without it the watch silently enqueues nothing.

## 6. The internal package — `controllers/apim/<resource>/internal/`

Three files, always the same signatures.

**`update.go`** — get a client, resolve what the payload needs inlined, call, copy status back:

```go
func CreateOrUpdate(ctx context.Context, obj *v1alpha1.MyResource) error {
	apimClient, err := apim.FromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
	if err != nil {
		return err
	}

	// Only if the spec carries refs that must be inlined before the payload is built.
	if err := ResolveRefs(ctx, obj); err != nil {
		return err
	}

	status, mgmtErr := apimClient.MyResources.CreateOrUpdate(obj)
	if mgmtErr != nil {
		return gerrors.NewControlPlaneError(mgmtErr)
	}

	// Field by field where the status holds operator-owned values too, or
	// status.DeepCopyTo(&obj.Status.Status) where the whole struct comes from APIM.
	obj.Status.ID = status.ID
	obj.Status.OrgID = status.OrgID
	obj.Status.EnvID = status.EnvID

	return nil
}
```

Controllers whose `PopulateIDs` does real work call it immediately after obtaining the client and
before any ref resolution:
`obj.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(obj))`.

**`delete.go`** — guard on the finalizer, tolerate a missing remote:

```go
func Delete(ctx context.Context, obj *v1alpha1.MyResource) error {
	if !util.ContainsFinalizer(obj, core.MyResourceFinalizer) {
		return nil
	}
	apimClient, err := apim.FromContextRef(ctx, obj.Spec.Context, obj.GetNamespace())
	if err != nil {
		return err
	}
	return errors.IgnoreNotFound(apimClient.MyResources.Delete(obj))
}
```

**`status.go`** — copy verbatim; the condition builders stamp `ObservedGeneration` for you:

```go
func UpdateStatusSuccess(ctx context.Context, obj *v1alpha1.MyResource) error {
	if obj.IsBeingDeleted() {
		return nil
	}
	k8s.AddSuccessfulConditions(obj)
	obj.Status.ProcessingStatus = core.ProcessingStatusCompleted
	return k8s.GetClient().Status().Update(ctx, obj)
}

func UpdateStatusFailure(ctx context.Context, obj *v1alpha1.MyResource, err error) error {
	k8s.ErrorToCondition(obj, err)
	obj.Status.ProcessingStatus = core.ProcessingStatusFailed
	return k8s.GetClient().Status().Update(ctx, obj)
}
```

## 7. Reference resolution

Controllers all resolve references the same way, in the same order. Do not invent a fourth
mechanism.

| What you are resolving | Use | Where it runs |
|---|---|---|
| ``[[ secret `name/key` ]]`` / ``[[ configmap `name/key` ]]`` values | `template.Compile(ctx, obj, true)` | The reconcile mutate fn, before anything else |
| Another CR (`ManagementContext`, API, Application, Group, Portal, Secret) | `dynamic.Resolve*` in `internal/k8s/dynamic/` | `internal/update.go` |
| Refs that must be **inlined into the payload** | a resolver in `internal/apim/<domain>/` | `internal/update.go`, before the service call |

`dynamic.resolveRef` compiles templates on the resolved object too, so a referenced CR arrives
fully materialised:

```go
app, err := dynamic.ResolveApplication(ctx, &spec.App, ns)
api, err := dynamic.ResolveAPI(ctx, &spec.API, ns)
mCtx, err := dynamic.ResolveContext(ctx, ref, ns)
```

Inlining resolvers live next to the APIM client, not in the controller, because admission needs
them too: `internal/apim/apidefinition/{resources,groups,shared_policy_groups,notification}.go`
and `internal/apim/application/resolve.go`. `PrepareV4SpecForAutomation` orchestrates the whole
set for v4 APIs. Put yours there and have both `internal/update.go` and the webhook's drift
`RefResolver` call it — resolving in only one of the two is how CRD and remote end up
incomparable.

When a referenced object is optional and missing, degrade rather than fail: mark the
`k8s.ConditionResolvedRefs` condition false and continue, as `ResolveGroupRefs` does.

## 8. Wire it up in `main.go`

Add the reconciler to `registerAutomationAPIControllers`:

```go
if err := (&myresource.Reconciler{
	Client:   k8s.GetClient(),
	Scheme:   mgr.GetScheme(),
	Recorder: mgr.GetEventRecorderFor("myresource-controller"),
	Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.MyResourceList{}),
}).SetupWithManager(mgr); err != nil {
	setupLog.Error(err, "Unable to create controller", "controller", "MyResource")
	os.Exit(1)
}
```

Gateway API controllers go in `registerGatewayAPIsControllers` behind `env.Config.EnableGatewayAPI`;
the ingress controller behind `env.Config.EnableIngress`. Indexers are registered by
`search.InitCache` earlier in `main`, so add yours there too.

Declare the finalizer in `internal/core/keys.go` alongside the others:

```go
MyResourceFinalizer = "finalizers.gravitee.io/myresources"
```

## 9. Tests

**Unit tests only in this repo.** Put them in `test/unit/<area>/` (Ginkgo v2) and cover the parts
that are pure logic — the `To*DTO` mapping, drift tags (`test/unit/drift/apim/`), any predicate or
helper you added. Never place `_test.go` under `controllers/**/internal`.

**Everything that needs a cluster or a live APIM goes to
[`gravitee-io/gravitee-platform-e2e`](https://github.com/gravitee-io/gravitee-platform-e2e)**:
reconciliation, `.status`, conditions, finalizer-driven deletion, templating, watch-triggered
re-reconciles. Operator-specific coverage lives in `apim/tests/gko/<area>/` with fixtures in
`apim/fixtures/<area>/`; behaviour a customer could also reach through Terraform belongs in a user
journey under `apim/tests/user-journeys/<persona>/<journey>/`. That repo has its own `AGENTS.md`
and a `write-e2e-test` skill — follow them there.

Do not add anything to `test/integration/`.

## 10. Final checks

```bash
make generate manifests reference
make build
make -j4 lint-sources
make unit
```

`make manifests` regenerates the RBAC from your markers and the CRD YAML under
`helm/gko/crds/gravitee.io/`. CI fails on a dirty tree, so commit the generated files with the
source change.
