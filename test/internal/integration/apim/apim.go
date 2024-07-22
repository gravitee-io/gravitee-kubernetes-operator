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

package apim

import (
	"context"

	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
)

type APIM struct {
	*apim.APIM

	Subscriptions *service.Subscriptions
	Pages         *service.Pages
	Env           *service.Env
<<<<<<< HEAD
=======
	Org           *service.Org
	Export        *Export
>>>>>>> f7e6318 (refactor: reduce dependency between internal and api)
}

func NewClient(ctx context.Context) *APIM {
	context := fixture.Builder().
		WithContext(constants.ContextWithCredentialsFile).
		Build().
		Context

	apim, err := apim.FromContext(ctx, context.Spec.Context)
	Expect(err).ToNot(HaveOccurred())

	subscriptions := service.NewSubscriptions(apim.APIs.Client)
	pages := service.NewPages(apim.APIs.Client)
	env := service.NewEnv(apim.APIs.Client)
<<<<<<< HEAD
=======
	org := service.NewOrg(apim.APIs.Client)
	export := NewExport(apim.APIs.Client)
>>>>>>> f7e6318 (refactor: reduce dependency between internal and api)

	return &APIM{
		APIM:          apim,
		Subscriptions: subscriptions,
		Pages:         pages,
		Env:           env,
<<<<<<< HEAD
=======
		Org:           org,
		Export:        export,
>>>>>>> f7e6318 (refactor: reduce dependency between internal and api)
	}
}
