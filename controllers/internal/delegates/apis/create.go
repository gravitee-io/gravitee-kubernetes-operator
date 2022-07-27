package apis

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	gioCtx "github.com/gravitee-io/gravitee-kubernetes-operator/controllers/internal/delegates/context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	requestTimeoutSeconds = 5
	defaultPlanSecurity   = "KEY_LESS"
	defaultPlanStatus     = "PUBLISHED"
	defaultPlanName       = "G.K.O. Default"
	origin                = "kubernetes"
	mode                  = "fully_managed"
)

func (d *Delegate) Create(
	api *gio.ApiDefinition,
) error {
	// Plan is not required from the CRD, but is expected by the Gateway, so we must create at least one
	d.addPlan(api)

	ctxDelegate := gioCtx.NewDelegate(d.ctx, d.cli)
	apimCtx, err := ctxDelegate.Get(api)
	if client.IgnoreNotFound(err) != nil {
		d.log.Error(err, "Management context will be discarded in further operations")
	}

	// Ensure that IDs have been generated
	generateIds(apimCtx, api)
	setDeployedAt(api)

	api.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: origin,
		Mode:   mode,
	}

	apiJson, err := json.Marshal(api.Spec)
	if err != nil {
		d.log.Error(err, "Unable to marshall API definition as JSON")
		return err
	}

	updated, err := d.updateConfigMap(api, apimCtx, apiJson)
	if err != nil {
		d.log.Error(err, "Unable to create or update ConfigMap from API definition")
		return err
	}

	if updated {
		err = d.importToManagementApi(api, apimCtx, apiJson)
		if err != nil {
			d.log.Error(err, "Unable to import to the Management API")
			return err
		}
	}

	api.Status.ApiID = api.Spec.CrossId

	err = d.cli.Status().Update(d.ctx, api)
	if err != nil {
		d.log.Error(err, "Unexpected error while updating status")
		return err
	}

	return nil
}

// Add a default keyless plan to the api definition if no plan is defined.
func (d *Delegate) addPlan(api *gio.ApiDefinition) {
	plans := api.Spec.Plans

	if len(plans) == 0 {
		d.log.Info("Define default plan for API")
		api.Spec.Plans = []*model.Plan{
			{
				Name:     defaultPlanName,
				Security: defaultPlanSecurity,
				Status:   defaultPlanStatus,
			},
		}
	}
}

func (d *Delegate) updateConfigMap(
	api *gio.ApiDefinition,
	apimCtx *gio.ManagementContext,
	apiJson []byte,
) (bool, error) {
	// Create configmap with some specific metadata that will be used to check changes across 'Update' events.
	cm := &v1.ConfigMap{}

	cm.Namespace = api.Namespace
	cm.Name = api.Name
	cm.CreationTimestamp = metav1.Now()
	cm.Labels = map[string]string{
		"managed-by": keys.CrdGroup,
		"gio-type":   keys.CrdApiDefinitionResource + "." + keys.CrdGroup,
	}

	cm.Data = map[string]string{
		"definition":        string(apiJson),
		"definitionVersion": api.ResourceVersion,
	}

	if apimCtx != nil {
		cm.Data["organizationId"] = apimCtx.Spec.OrgId
		cm.Data["environmentId"] = apimCtx.Spec.EnvId
	}

	currentapiDefinition := &v1.ConfigMap{}
	err := d.cli.Get(d.ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, currentapiDefinition)

	if err == nil {
		if currentapiDefinition.Data["definitionVersion"] != api.ResourceVersion {
			d.log.Info("Updating ConfigMap", "id", api.Spec.Id)
			// Only update the confimap if resource version has changed (means api definition has changed).
			err = d.cli.Update(d.ctx, cm)
		} else {
			d.log.Info("No change detected on api. Skipped.", "id", api.Spec.Id)
			return false, nil
		}
	} else {
		d.log.Info("Creating configmap for api.", "id", api.Spec.Id, "name", api.Name)
		err = d.cli.Create(d.ctx, cm)
	}
	return true, err
}
