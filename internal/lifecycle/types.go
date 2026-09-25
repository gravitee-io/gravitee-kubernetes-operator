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

// Package lifecycle is the generic reconcile and admission skeleton for new Gravitee resources.
// The framework owns fetch, finalizer, spec hash, templates, conditions, and requeue.
// ResourceLifecycle and AdmissionLifecycle are the holes a resource fills.
package lifecycle

import (
	"context"

	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/store"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// ClientFactoryFunc builds the API client from the CR (own contextRef or inherited parent).
// Error means the context cannot be used.
type ClientFactoryFunc[T core.ContextAwareObject, C core.APIClient] func(ctx context.Context, obj T) (C, error)

// ToDTOFunc maps a CR to the wire payload. Do not fetch refs or call the API.
// Identity() must be set so Delete and GetRemote can address the remote object.
type ToDTOFunc[T core.ContextAwareObject, D store.Identifiable] func(obj T) D

// UpsertFunc PUT/creates the remote resource from dto. Never sees the CR.
type UpsertFunc[C any, D store.Identifiable, R any] func(ctx context.Context, client C, dto D) (R, error)

// DeleteFunc removes the remote resource. Use dto.Identity() as the key. Never sees the CR.
type DeleteFunc[C any, D store.Identifiable] func(ctx context.Context, client C, dto D) error

// DeleteGuardFunc returns an error when the CR must not be deleted (still referenced).
type DeleteGuardFunc[T core.ContextAwareObject] func(ctx context.Context, obj T) error

// PostUpsertFunc mutates the CR after a successful Upsert (status, annotations).
type PostUpsertFunc[T core.ContextAwareObject, R any] func(ctx context.Context, obj T, response R) error

// RefResolverFunc writes resolved in-cluster refs onto obj (secrets, parent CRs). Mutates obj.
type RefResolverFunc[T any] func(ctx context.Context, obj T, namespace string) error

// AdmissionCheckFunc is extra spec/cluster validation.
// Severe errors reject; warnings admit.
type AdmissionCheckFunc[T core.ContextAwareObject] func(ctx context.Context, obj T) *gerrors.AdmissionErrors

// ImmutableFieldsFunc returns severe errors when frozen spec fields changed between old and new.
type ImmutableFieldsFunc[T core.ContextAwareObject] func(oldValue, newValue T) *gerrors.AdmissionErrors

// DryRunFunc validates dto against the remote API without persisting. Never maps the CR.
type DryRunFunc[C any, D store.Identifiable] func(ctx context.Context, client C, dto D) *gerrors.AdmissionErrors

// GetRemoteFunc returns the live remote object for dto.Identity().
// Return the client error as-is, including 404.
type GetRemoteFunc[C any, D store.Identifiable] func(ctx context.Context, client C, dto D) (D, error)

// ResourceLifecycle is the reconcile holes for one CR kind.
type ResourceLifecycle[T core.ContextAwareObject, D store.Identifiable, C core.APIClient, R core.OrgEnvIDGetter] struct {
	// Finalizer is added on every reconcile and removed only after a successful Delete.
	Finalizer string

	// ResolveRefs runs after template compile/release, before ClientFactory.
	// Nil skips resolution and the ResolvedRefs condition. Use the same func as AdmissionLifecycle.ResolveRefs.
	ResolveRefs RefResolverFunc[T]

	// ClientFactory runs after refs, before ToDTO. Required.
	ClientFactory ClientFactoryFunc[T, C]

	// ToDTO runs after ClientFactory, before Upsert or Delete. Required.
	// Use the same func as AdmissionLifecycle.ToDTO.
	ToDTO ToDTOFunc[T, D]

	// DeleteGuard runs on delete, after ToDTO, before Delete. Nil skips.
	DeleteGuard DeleteGuardFunc[T]

	// Delete runs on delete, after DeleteGuard. Required.
	Delete DeleteFunc[C, D]

	// Upsert runs on create/update, after ToDTO. Required.
	Upsert UpsertFunc[C, D, R]

	// PostUpsert runs on create/update, after a successful Upsert. Nil skips.
	PostUpsert PostUpsertFunc[T, R]
}

// AdmissionLifecycle is the admission holes for one CR kind.
type AdmissionLifecycle[T core.ContextAwareObject, D store.Identifiable, C core.APIClient] struct {
	// ResolveRefs runs on create/update after template compile, before ClientFactory.
	// Nil skips resolution. Use the same func as ResourceLifecycle.ResolveRefs.
	ResolveRefs RefResolverFunc[T]

	// ClientFactory runs on create/update after refs, before PreCheck.
	// Builds C for DryRun and GetRemote. Required when either is set.
	ClientFactory ClientFactoryFunc[T, C]

	// PreCheck runs on create/update after ClientFactory, before ImmutableFields. Nil skips.
	PreCheck AdmissionCheckFunc[T]

	// ImmutableFields runs on update after PreCheck, before ToDTO. Nil skips.
	ImmutableFields ImmutableFieldsFunc[T]

	// ToDTO runs on create/update after PreCheck and ImmutableFields, before DryRun.
	// Required for DryRun and drift. Use the same func as ResourceLifecycle.ToDTO.
	ToDTO ToDTOFunc[T, D]

	// DryRun runs on create/update after ToDTO, before PostCheck. Nil skips.
	DryRun DryRunFunc[C, D]

	// PostCheck runs on create/update after DryRun, before drift. Nil skips.
	PostCheck AdmissionCheckFunc[T]

	// GetRemote runs on update after PostCheck. Nil disables drift.
	GetRemote GetRemoteFunc[C, D]

	// DeleteGuard runs on delete review only. Nil skips.
	DeleteGuard DeleteGuardFunc[T]
}
