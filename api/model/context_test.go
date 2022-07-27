package model

import (
	"net/http"
	"testing"
)

func Test_BuildUrl(t *testing.T) {
	tests := []struct {
		name     string
		ctx      Context
		path     string
		expected string
	}{
		{
			"With Context with an env and an org ID",
			Context{
				BaseUrl: "http://localhost:8083",
				EnvId:   "DEFAULT",
				OrgId:   "DEFAULT",
				Auth:    nil,
			},
			"/apis",
			"http://localhost:8083/management/organizations/DEFAULT/environments/DEFAULT/apis",
		},
		{
			"With Context with only an org ID",
			Context{
				BaseUrl: "http://localhost:8083",
				EnvId:   "",
				OrgId:   "DEFAULT",
				Auth:    nil,
			},
			"/user",
			"http://localhost:8083/management/organizations/DEFAULT/user",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			given := test.ctx.BuildUrl(test.path)
			if given != test.expected {
				t.Fail()
				t.Logf("Expected %s to be %s", given, test.expected)
			}
		})
	}
}

func Test_Authenticate(t *testing.T) {
	tests := []struct {
		name     string
		ctx      Context
		expected string
	}{
		{
			"With basic auth",
			Context{
				BaseUrl: "http://localhost:8083",
				EnvId:   "DEFAULT",
				OrgId:   "DEFAULT",
				Auth: &Auth{
					Credentials: &BasicAuth{
						Username: "admin",
						Password: "admin",
					},
				},
			},
			"Basic YWRtaW46YWRtaW4=",
		},
		{
			"With empty credentials",
			Context{
				BaseUrl: "http://localhost:8083",
				EnvId:   "DEFAULT",
				OrgId:   "DEFAULT",
				Auth: &Auth{
					Credentials: &BasicAuth{
						Username: "",
						Password: "",
					},
				},
			},
			"",
		},
		{
			"With bearer token",
			Context{
				BaseUrl: "http://localhost:8083",
				EnvId:   "DEFAULT",
				OrgId:   "DEFAULT",
				Auth: &Auth{
					BearerToken: "563c8597-7ec8-4cf2-aee4-97acb22a52c5",
				},
			},
			"Bearer 563c8597-7ec8-4cf2-aee4-97acb22a52c5",
		},
		{
			"With no auth",
			Context{
				BaseUrl: "http://localhost:8083",
				EnvId:   "DEFAULT",
				OrgId:   "DEFAULT",
				Auth:    nil,
			},
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &http.Request{Header: make(map[string][]string)}
			test.ctx.Authenticate(req)
			given := req.Header.Get("Authorization")
			if given != test.expected {
				t.Fail()
				t.Logf("Expected %s to be %s", given, test.expected)
			}
		})
	}
}
