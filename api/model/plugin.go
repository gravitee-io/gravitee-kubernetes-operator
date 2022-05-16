package model

type PluginReference struct {
	Namespace string `json:"namespace,omitempty"`
	Resource  string `json:"resource,omitempty"`
	Name      string `json:"name,omitempty"`
}

func NewPluginReference() *PluginReference {
	return &PluginReference{
		Namespace: "default",
	}
}

type PluginRevision struct {
	PluginReference *PluginReference `json:"pluginReference,omitempty"`
	Generation      int64            `json:"generation,omitempty"`
	Plugin          *Plugin          `json:"plugin,omitempty"`
	HashCode        string           `json:"hashCode,omitempty"`
}
