package apis

import (
	"encoding/base64"

	uuid "github.com/satori/go.uuid" //nolint:gomodguard // to replace with google implementation
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

const separator = "/"

// Return Spec CrossId or generate a new one from api Name & Namespace.
func RetrieveCrossId(api *gio.ApiDefinition) string {
	// If a CrossID is defined at the API level, reuse it.
	// If not, just generate a new CrossID
	if api.Spec.CrossId == "" {
		// The ID of the API will be based on the API Name and Namespace to ensure consistency
		return ToUUID(types.NamespacedName{Namespace: api.Namespace, Name: api.Name}.String())
	}

	return api.Spec.CrossId
}

// Generate UUID.
func generateId() string {
	return uuid.NewV4().String()
}

func setSpecIdsFromStatus(api *gio.ApiDefinition) {
	api.Spec.CrossId = api.Status.CrossID
	api.Spec.Id = api.Status.ID

	plans := api.Spec.Plans
	for _, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = ToUUID(api.Spec.Id + separator + plan.Name)
		}
	}
}

func ToUUID(decoded string) string {
	encoded := base64.RawStdEncoding.EncodeToString([]byte(decoded))
	return uuid.NewV3(uuid.NamespaceURL, encoded).String()
}
