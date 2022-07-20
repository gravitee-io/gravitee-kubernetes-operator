// +kubebuilder:object:generate=true
package model

type ContextRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type Context struct {
	BaseUrl string `json:"baseUrl"`
	EnvId   string `json:"environmentId,omitempty"`
	OrgId   string `json:"organizationId,omitempty"`
	Auth    *Auth  `json:"auth"`
}

type Auth struct {
	BearerToken string     `json:"bearerToken,omitempty"`
	Credentials *BasicAuth `json:"credentials,omitempty"`
}

type BasicAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
