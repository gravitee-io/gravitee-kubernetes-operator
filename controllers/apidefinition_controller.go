/*
Copyright 2022 DAVID BRASSELY.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apis "github.com/gravitee-io/gravitee-kubernetes-operator/delegates/apis"
	gioCtx "github.com/gravitee-io/gravitee-kubernetes-operator/delegates/context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
)

const requeueAfterTime = 5

// ApiDefinitionReconciler reconciles a ApiDefinition object.
type ApiDefinitionReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get
//+kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApiDefinition object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ApiDefinitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("namespace", req.Namespace, "name", req.Name)

	// Fetch the ApiDefinition apiDefinition
	apiDefinition := &gio.ApiDefinition{}
	requeueAfter := time.Second * requeueAfterTime

	err := r.Get(ctx, req.NamespacedName, apiDefinition)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("API Definition resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get AP IDefinition")
		return ctrl.Result{}, err
	}

	ctxDelegate := gioCtx.NewDelegate(ctx, r.Client)
	managementContext, err := ctxDelegate.Get(apiDefinition)

	if client.IgnoreNotFound(err) != nil {
		log.Info("Management context will be ignored for further operations (not found)")
	}

	apisDelegate := apis.NewDelegate(ctx, managementContext, r.Client)

	if apiDefinition.GetLabels()[keys.CrdApiDefinitionTemplate] == "true" {
		log.Info("Creating a new API Definition template", "template", apiDefinition.Name)

		requeue, importErr := apisDelegate.ImportApiDefinitionTemplate(apiDefinition, req.Namespace)
		if importErr != nil {
			log.Error(importErr, "Failed to sync template")
			return ctrl.Result{}, importErr
		}

		if requeue {
			return ctrl.Result{RequeueAfter: requeueAfter}, nil
		}
	}

	_, err = util.CreateOrUpdate(ctx, r.Client, apiDefinition, func() error {
		if !apiDefinition.ObjectMeta.DeletionTimestamp.IsZero() {
			return apisDelegate.Delete(apiDefinition)
		}
		if apiDefinition.Status.CrossID == "" {
			return apisDelegate.Create(apiDefinition)
		}
		return apisDelegate.Update(apiDefinition)
	})

	if err == nil {
		log.Info("ApiDefinition has been reconciled")
		return ctrl.Result{}, nil
	}

	// Should we keep this re-queuing strategy ?
	return ctrl.Result{RequeueAfter: requeueAfter}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApiDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gio.ApiDefinition{}).
		//		Owns(&v1.Secret{}).
		Complete(r)
}
