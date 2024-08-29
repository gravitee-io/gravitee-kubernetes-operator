// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mctx

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const usCloudGateURL = "https://us.cloudgate.gravitee.io/apim/rest"

var _ = Describe("cloud token defaults", func() {
	It("without cloud token", func() {
		orgID := newUUID()
		envID := newUUID()
		underTest := &v1alpha1.ManagementContext{
			Spec: v1alpha1.ManagementContextSpec{
				Context: &management.Context{
					BaseUrl: "https://localhost:8080",
					OrgID:   orgID,
					EnvID:   envID,
					Auth: &management.Auth{
						Credentials: &management.BasicAuth{
							Username: "admin",
							Password: "password",
						},
					}},
			},
		}

		err := mctx.SetDefaults(context.Background(), underTest)
		Expect(err).To(BeNil())

		// nothing should have changed
		Expect(underTest.Spec.BaseUrl).To(Equal("https://localhost:8080"))
		Expect(underTest.Spec.OrgID).To(Equal(orgID))
		Expect(underTest.Spec.EnvID).To(Equal(envID))
		Expect(underTest.HasAuthentication()).To(BeTrue())
		Expect(underTest.GetAuth().HasCredentials()).To(BeTrue())
		Expect(underTest.GetAuth().GetCredentials().GetUsername()).To(Equal("admin"))
		Expect(underTest.GetAuth().GetCredentials().GetPassword()).To(Equal("password"))
	})

	DescribeTable("with cloud token", func(given *management.Context, expectedUrl string) {
		orgID := newUUID()
		envID := newUUID()

		token := forgeToken(mctx.CloudTokenClaimsData{
			Org:       orgID,
			Envs:      []string{envID},
			Geography: "us",
		})

		given.Cloud = &management.Cloud{Token: token}

		underTest := &v1alpha1.ManagementContext{
			Spec: v1alpha1.ManagementContextSpec{
				Context: given,
			},
		}
		err := mctx.SetDefaults(context.Background(), underTest)

		Expect(err).To(BeNil())
		Expect(underTest.Spec.BaseUrl).To(Equal(expectedUrl))
		Expect(underTest.Spec.OrgID).To(Equal(orgID))
		Expect(underTest.Spec.EnvID).To(Equal(envID))
		Expect(underTest.Spec.HasAuthentication()).To(BeTrue())
		Expect(underTest.Spec.Auth.BearerToken).To(Equal(token))
		Expect(underTest.Spec.Auth.SecretRef).To(BeNil())
		Expect(underTest.Spec.Auth.HasCredentials()).To(BeFalse())
		Expect(underTest.HasCloud()).To(BeTrue())
		Expect(underTest.GetCloud().GetToken()).To(Equal(token))

	},
		Entry("default", &management.Context{}, usCloudGateURL),
		Entry("keep url", &management.Context{BaseUrl: "https://locahost"}, "https://locahost"),
		Entry("override orgID", &management.Context{OrgID: "foo"}, usCloudGateURL),
		Entry("override secret ref", &management.Context{Auth: &management.Auth{
			SecretRef: &refs.NamespacedName{
				Name:      "foo",
				Namespace: "bar",
			}}}, usCloudGateURL),
		Entry("override bearerToken", &management.Context{
			Auth: &management.Auth{
				BearerToken: "foo",
			}}, usCloudGateURL),
		Entry("remove basic auth", &management.Context{
			Auth: &management.Auth{
				Credentials: &management.BasicAuth{
					Username: "admin",
					Password: "password",
				},
			}}, usCloudGateURL))

	DescribeTable("errors", func(givenContext *management.Context, givenToken string, partialError string) {
		givenContext.Cloud = &management.Cloud{Token: givenToken}

		underTest := &v1alpha1.ManagementContext{
			Spec: v1alpha1.ManagementContextSpec{
				Context: givenContext,
			},
		}
		err := mctx.SetDefaults(context.Background(), underTest)
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(ContainSubstring(partialError))
	},
		Entry("no env", &management.Context{}, forgeToken(mctx.CloudTokenClaimsData{
			Org:       "foo",
			Envs:      []string{},
			Geography: "us",
		}), "required claims"),
		Entry("no org", &management.Context{}, forgeToken(mctx.CloudTokenClaimsData{
			Org:       "",
			Envs:      []string{"bar"},
			Geography: "us",
		}), "required claims"),
		Entry("no region", &management.Context{}, forgeToken(mctx.CloudTokenClaimsData{
			Org:       "foo",
			Envs:      []string{"bar"},
			Geography: "",
		}), "required claims"),
		Entry("three env", &management.Context{}, forgeToken(mctx.CloudTokenClaimsData{
			Org:       "foo",
			Envs:      []string{"123", "456", "789"},
			Geography: "us",
		}), "more than one environment (3)"),
		Entry("env mismatch", &management.Context{EnvID: "777"}, forgeToken(mctx.CloudTokenClaimsData{
			Org:       "foo",
			Envs:      []string{"123", "456", "789"},
			Geography: "us",
		}), "[777], it must be one of: [123 456 789]"),
		Entry("invalid token", &management.Context{}, "invalid token", "cannot parse cloud token"),
	)

})

func forgeToken(tokenData mctx.CloudTokenClaimsData) string {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, mctx.CloudTokenClaims{
		RegisteredClaims:     jwt.RegisteredClaims{},
		CloudTokenClaimsData: tokenData,
	})
	s, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		panic(err.Error())
	}
	return s
}

func newUUID() string {
	return uuid.New().String()
}
