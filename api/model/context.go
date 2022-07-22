// +kubebuilder:object:generate=true
package model

type ContextRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type Context struct {
	// +kubebuilder:validation:Pattern=`^http(s?):\/\/.+$`
	BaseUrl string `json:"baseUrl"`
	// +kubebuilder:validation:Required
	OrgId string `json:"organizationId"`
	// +kubebuilder:validation:Required
	EnvId string `json:"environmentId"`
	// +kubebuilder:validation:Required
	Auth *Auth `json:"auth"`
}

type Auth struct {
	BearerToken string     `json:"bearerToken,omitempty"`
	Credentials *BasicAuth `json:"credentials,omitempty"`
}

type BasicAuth struct {
	// +kubebuilder:validation:Required
	Username string `json:"username,omitempty"`
	// +kubebuilder:validation:Required
	Password string `json:"password,omitempty"`
}
