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

package subscription

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Upsert creates or updates the subscription in APIM.
func Upsert(ctx context.Context, client *apim.SDKClient, dto DTO) (Response, error) {
	resp, err := client.UpsertApiSubscriptionWithResponse(ctx, dto.ApiHrid, dto.Legacy.upsertParams(), dto.Subscription)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return Response{}, err
	}
	if resp.JSON200 == nil {
		return Response{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return Response{State: *resp.JSON200, Legacy: dto.Legacy}, nil
}

// Delete removes the subscription from APIM. A subscription already gone counts as deleted.
func Delete(ctx context.Context, client *apim.SDKClient, dto DTO) error {
	resp, err := client.DeleteApiSubscriptionWithResponse(ctx, dto.ApiHrid, dto.Hrid, dto.Legacy.deleteParams())
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the upsert response into the CR status. A subscription upserted by its hrid becomes
// Automation API managed; a legacy subscription, upserted by its UUID, stays as it is. A response without an id
// keeps the one the status has. The deprecated processing status is cleared: the lifecycle reports through the
// conditions, and a value left by the controller before it would never change again.
func UpdateStatus(obj *v1alpha1.Subscription, resp Response) error {
	if resp.State.StartingAt == nil {
		return fmt.Errorf("subscription [%s] has no start date in APIM", resp.State.Hrid)
	}
	startedAt, err := time.Parse(time.RFC3339, *resp.State.StartingAt)
	if err != nil {
		return fmt.Errorf("parse the start date of subscription [%s]: %w", resp.State.Hrid, err)
	}

	obj.Status.ProcessingStatus = ""
	if resp.State.Id != nil {
		obj.Status.ID = *resp.State.Id
	}
	obj.Status.StartedAt = startedAt.Format(time.RFC3339)
	obj.Status.EndingAt = ""
	if resp.State.EndingAt != nil {
		obj.Status.EndingAt = *resp.State.EndingAt
	}

	if !resp.Legacy.HridContainsUUID {
		k8s.AddAutomationAPIManagedCondition(obj)
	}
	return nil
}

// UpdateStatusAndRecount copies the upsert response into the CR status, then recounts the subscriptions of the
// API and the application the subscription references.
func UpdateStatusAndRecount(ctx context.Context, sub *v1alpha1.Subscription, resp Response) error {
	if err := UpdateStatus(sub, resp); err != nil {
		return err
	}
	return RecountParents(ctx, sub)
}
