package apis

import (
	"encoding/base64"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"

	uuid "github.com/satori/go.uuid" //nolint:gomodguard // to replace with google implementation
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

const separator = "/"

// This function is used to generate all the IDs needed for communicating with the Management API
// It doesn't override IDs if these one have been defined.
func generateIds(api *gio.ApiDefinition) {
	// If a CrossID is defined at the API level, reuse it.
	// If not, just generate a new CrossID
	if api.Spec.CrossId == "" {
		// The ID of the API will be based on the API Name and Namespace to ensure consistency
		api.Spec.CrossId = toUUID(getNamespacedName(api))
	}

	if api.Spec.Id == "" {
		api.Spec.Id = uuid.NewV4().String()
	}

	plans := api.Spec.Plans

	for _, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = toUUID(api.Spec.Id + separator + plan.Name)
		}
		plan.Status = "PUBLISHED"
	}

	//TODO: manage metadata
}

func setIds(api *gio.ApiDefinition) {
	api.Spec.CrossId = api.Status.CrossID
	api.Spec.Id = api.Status.ID

	plans := api.Spec.Plans
	for _, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = toUUID(api.Spec.Id + separator + plan.Name)
		}
	}
}

func setOrigin(api *gio.ApiDefinition) {
	api.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: origin,
		Mode:   mode,
	}
}

func getNamespacedName(api *gio.ApiDefinition) string {
	return types.NamespacedName{Namespace: api.Namespace, Name: api.Name}.String()
}

func toUUID(decoded string) string {
	encoded := base64.RawStdEncoding.EncodeToString([]byte(decoded))
	return uuid.NewV3(uuid.NamespaceURL, encoded).String()
}
