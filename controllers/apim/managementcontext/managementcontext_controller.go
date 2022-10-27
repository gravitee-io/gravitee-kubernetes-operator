/*
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

package managementcontext

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// Reconciler reconciles a ManagementContext object.
type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=gravitee.io,resources=managementcontexts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gravitee.io,resources=managementcontexts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gravitee.io,resources=managementcontexts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ManagementContext object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("ManagementContext", req.NamespacedName)

	instance := &gio.ManagementContext{}

	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling ManagementContext instance")

	// Update API resources that reference this context and are in a failed state
	apis, err := r.listFailedApiDefinitionResources(ctx, instance.Name, instance.Namespace)
	if err != nil {
		log.Error(err, "unable to list API definitions resources, skipping update")
		return ctrl.Result{}, nil
	}

	for i := range apis {
		api := apis[i]

		api.Status.ProcessingStatus = gio.ProcessingStatusReconciling

		if err = r.Status().Update(ctx, &api); err != nil {
			log.Error(err, "unable to update API definition status, skipping update")
		}
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) listFailedApiDefinitionResources(
	ctx context.Context, contextName, contextNamespace string,
) ([]gio.ApiDefinition, error) {
	log := log.FromContext(ctx)
	statusFailed := string(gio.ProcessingStatusFailed)
	apiDefinitionList := &gio.ApiDefinitionList{}

	contextNameFilter := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.contextRef.name": contextName}),
	}

	contextNamespaceFilter := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.contextRef.namespace": contextNamespace}),
	}

	statusFilter := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"status.processingStatus": statusFailed}),
	}

	if err := r.Client.List(ctx, apiDefinitionList, contextNameFilter, contextNamespaceFilter, statusFilter); err != nil {
		log.Error(err, "unable to list API definitions, skipping update")
		return nil, err
	}

	return apiDefinitionList.Items, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gio.ManagementContext{}).
		Complete(r)
}
