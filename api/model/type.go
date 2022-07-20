package model

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type GenericStringMap struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	unstructured.Unstructured `json:",inline"`
}

func (in *GenericStringMap) DeepCopyInto(out *GenericStringMap) {
	// controller-gen cannot handle the interface{} type of an aliased Unstructured,
	// thus we write our own DeepCopyInto function.
	if out != nil {
		casted := in.Unstructured
		deepCopy := casted.DeepCopy()
		out.Object = deepCopy.Object
	}
}

func (in *GenericStringMap) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	in.Object = m

	return nil
}

func (in *GenericStringMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(in.Object)
}
