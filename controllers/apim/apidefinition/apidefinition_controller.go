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

package apidefinition

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/utils"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
)

const requeueAfterTime = 5

// Reconciler reconciles a ApiDefinition object.
type Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApiDefinition object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
//
//nolint:funlen,gocognit // We prefer to have a long function with visible intentions than several functions making it more complex to understand
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("namespace", req.Namespace, "name", req.Name)

	log.Info(fmt.Sprintf("Starting reconcile loop for %v", req.NamespacedName))
	defer log.Info(fmt.Sprintf("Finish reconcile loop for %v", req.NamespacedName))

	requeueAfter := time.Second * requeueAfterTime

	// Fetch the Api Definition apiDefinition
	apiDefinition, err := r.getApiDefinition(ctx, req.NamespacedName)
	if errors.IsNotFound(err) {
		log.Info("API Definition resource not found. Ignoring since object must be deleted")
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Failed to get API Definition")
		return ctrl.Result{}, err
	}

	log = log.WithValues("apiDefinitionName", apiDefinition.Name)

	// Ini Delegate and set management context if defined
	apisDelegate := internal.NewDelegate(ctx, r.Client, log)
	event := utils.NewEvent(r.Recorder)

	if apiDefinition.Spec.Context != nil {
		managementContext, ctxErr := apisDelegate.ResolveContext(apiDefinition.Spec.Context)

		if ctxErr != nil {
			log.Error(ctxErr, "And error has occurred while trying to retrieve ManagementContext")
			event.NormalEvent(
				apiDefinition,
				"ManagementContextUnLiked",
				fmt.Sprintf(
					"An error has occurred while trying to retrieve ManagementContext %s[%s]",
					apiDefinition.Spec.Context.Name, apiDefinition.Spec.Context.Namespace,
				),
			)
		}

		apisDelegate.SetManagementContext(managementContext)
	}

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

	if resourcesErr := apisDelegate.ResolveResources(apiDefinition); resourcesErr != nil {
		return ctrl.Result{}, resourcesErr
	}

	// Executes delegate actions
	switch {
	case !apiDefinition.HasDeletionFinalizer():
		log.Info("Add Finalizer to API definition")
		err = apisDelegate.AddDeletionFinalizer(apiDefinition)
		if err != nil {
			break
		}
		event.NormalEvent(apiDefinition, "AddedFinalizer", "Added Finalizer for the API definition")

	case apiDefinition.IsBeingDeleted():
		log.Info("Deleting API definition")
		event.NormalEvent(apiDefinition, "Deleting", "Deleting API definition")
		err = apisDelegate.Delete(apiDefinition)
		if err != nil {
			break
		}
		event.NormalEvent(apiDefinition, "Deleted", "Deleted API definition")

	case apiDefinition.IsBeingCreated():
		log.Info("Creating API definition")
		event.NormalEvent(apiDefinition, "Creating", "Creating API definition")
		err = apisDelegate.CreateOrUpdate(apiDefinition)
		if err != nil {
			break
		}
		event.NormalEvent(apiDefinition, "Created", "Created API definition")

	case apiDefinition.IsBeingUpdated():
		log.Info("Updating API definition")
		event.NormalEvent(apiDefinition, "Updating", "Updating API definition")
		err = apisDelegate.CreateOrUpdate(apiDefinition)
		if err != nil {
			break
		}
		event.NormalEvent(apiDefinition, "Updated", "Updated API definition")
	default:
		log.Info("No action to perform")
	}

	if err == nil {
		log.Info("API Definition has been reconciled")
		return ctrl.Result{}, nil
	}

	log.Error(err, "API Definition has not been reconciled")

	// Error handling
	err = apisDelegate.UpdateStatusAndReturnError(apiDefinition, err)
	event.WarningEvent(apiDefinition, string(apiDefinition.Status.ProcessingStatus), err.Error())

	if internal.IsRecoverableError(err) {
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gio.ApiDefinition{}).
		//		Owns(&v1.Secret{}).
		Complete(r)
}

func (r *Reconciler) getApiDefinition(
	ctx context.Context,
	namespacedName types.NamespacedName,
) (*gio.ApiDefinition, error) {
	apiDefinition := &gio.ApiDefinition{}
	err := r.Get(ctx, namespacedName, apiDefinition)
	if err != nil {
		return nil, err
	}
	return apiDefinition, nil
}
