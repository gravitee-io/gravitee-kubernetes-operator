package internal

import (
	"encoding/base64"
	"fmt"

	uuid "github.com/satori/go.uuid" // nolint:gomodguard // to replace with google implementation
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// This function is used to generate all the IDs needed for communicating with the Management API
// It doesn't override IDs if these one have been defined.
func GenerateIds(api *gio.ApiDefinition) {
	// If a CrossID is defined at the API level, reuse it.
	// If not, just generate a new CrossID
	if api.Spec.CrossId == "" {
		// The ID of the API will be based on the API Name and Namespace to ensure consistency
		api.Spec.CrossId = toUUID(getNamespacedName(api))
	}

	plans := api.Spec.Plans

	for i, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = toUUID(api.Spec.CrossId + fmt.Sprint(i))
		}
		plan.Status = "PUBLISHED"
	}

	//TODO: manage metadata
}

func getNamespacedName(api *gio.ApiDefinition) string {
	return types.NamespacedName{Namespace: api.Namespace, Name: api.Name}.String()
}

func toUUID(decoded string) string {
	encoded := base64.RawStdEncoding.EncodeToString([]byte(decoded))
	return uuid.NewV3(uuid.NamespaceURL, encoded).String()
}
