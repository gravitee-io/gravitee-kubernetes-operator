package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	model "github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	v1alpha1 "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	netV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// IngressReconciler watches and reconciles Ingress objects
type IngressReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch

// Reconcile perform reconciliation logic for Ingress resource that is managed
// by the operator.
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.NamespacedName)

	instance := &netV1.Ingress{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Look at the gravitee template annotations to refer to the template
	name, value := instance.Annotations[keys.IngressTemplateAnnotation]

	// Create a dummy ApiDefinition for Ingresses which are refering to an unknown template
	template := &v1alpha1.ApiDefinition{
		Spec: v1alpha1.ApiDefinitionSpec{
			Api: model.Api{
				Name: "default-keyless",
			},
		},
	}

	if value {
		template = &v1alpha1.ApiDefinition{}

		// Retrieve the ApiDefinition template
		err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: req.Namespace}, template)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Associate the gravitee finalizer to the ingress to keep track of deletion in the future
	operation, err := util.CreateOrUpdate(ctx, r.Client, instance, func() error {
		if !instance.DeletionTimestamp.IsZero() {
			if util.ContainsFinalizer(instance, keys.IngressFinalizer) {
				util.RemoveFinalizer(instance, keys.IngressFinalizer)
			}
			return nil
		}

		if !util.ContainsFinalizer(instance, keys.IngressFinalizer) {
			util.AddFinalizer(instance, keys.IngressFinalizer)
			return nil
		}

		return nil
	})

	if err != nil {
		log.Error(err, "An error occurs while updating the ingress", "Operation", operation)
		return ctrl.Result{}, err
	}

	err = r.mergeApiDefitinion(log, template, instance)
	if err != nil {
		log.Error(err, "Unable to merge the ingress with API Definition")
		return ctrl.Result{}, err
	}

	// TODO: templating
	// TODO: transform
	log.Info("Sync ingress DONE")

	return ctrl.Result{}, nil
}

// Transform the ingress as an API Definition as per https://kubernetes.io/docs/concepts/services-networking/ingress/#the-ingress-resource
func (r *IngressReconciler) mergeApiDefitinion(log logr.Logger, template *v1alpha1.ApiDefinition, ingress *netV1.Ingress) error {
	log.Info("Merge Ingress with API Definition")

	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {

			// Create an API Definition from the template
			api := &v1alpha1.ApiDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ingress.Name,
					Namespace: ingress.Namespace,
				},
				Enabled: true,
				Spec:    *template.Spec.DeepCopy(),
			}

			service := path.Backend.Service

			//TODO: How-to dedal with PathType ?
			api.Spec.Proxy = &model.Proxy{
				VirtualHosts: []*model.VirtualHost{
					{
						Path: path.Path,
					},
				},
				Groups: []*model.EndpointGroup{
					{
						Name: "default",
						Endpoints: []*model.HttpEndpoint{
							{
								Name:   service.Name,
								Target: fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", service.Name, ingress.Namespace, service.Port.Number),
							},
						},
					},
				},
			}
		}
	}

	return nil
}

func (r *IngressReconciler) ingressClassEventFilter() predicate.Predicate {
	isGraviteeIngress := func(o runtime.Object) bool {
		switch e := o.(type) {
		case *netV1.Ingress:
			return e.GetAnnotations()[keys.IngressClassAnnotation] == keys.IngressClassAnnotationValue
		default:
			return false
		}
	}

	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return isGraviteeIngress(e.Object)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return isGraviteeIngress(e.ObjectNew)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return isGraviteeIngress(e.Object)
		},
	}
}

// SetupWithManager initializes ingress controller manager
func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&netV1.Ingress{}).
		Owns(&v1alpha1.ApiDefinition{}).
		WithEventFilter(r.ingressClassEventFilter()).
		Complete(r)
}
