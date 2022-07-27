// +kubebuilder:object:generate=true
package model

import (
	"net/http"
	"strings"
)

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

func (ctx Context) BuildUrl(path string) string {
	orgId, envId := ctx.OrgId, ctx.EnvId
	baseUrl := strings.TrimSuffix(ctx.BaseUrl, "/")
	url := baseUrl + "/management/organizations/" + orgId
	if envId != "" {
		url = url + "/environments/" + envId
	}
	return url + path
}

func (ctx Context) Authenticate(req *http.Request) {
	if ctx.Auth == nil {
		return
	}

	bearerToken := ctx.Auth.BearerToken
	if bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bearerToken)
	} else if ctx.Auth.Credentials != nil {
		username := ctx.Auth.Credentials.Username
		password := ctx.Auth.Credentials.Password
		setBasicAuth(req, username, password)
	}
}

func setBasicAuth(request *http.Request, username, password string) {
	if username != "" {
		request.SetBasicAuth(username, password)
	}
}
