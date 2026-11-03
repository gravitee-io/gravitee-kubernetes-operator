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

package clientsecret

import (
	"context"
	"errors"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ErrApplicationNotReady the application exists but AM has not accepted it yet (Status.Key empty). Recoverable.
var ErrApplicationNotReady = errors.New("not yet created in AM")

// ErrSecretEntryMissing the Secret exists but its entry for secretRef.key is missing or empty. Recoverable.
var ErrSecretEntryMissing = errors.New("entry not found or empty")

// getApplication reads the application in obj's namespace. A failed read names it, and a NotFound stays detectable.
func getApplication(ctx context.Context, obj *v1alpha1.AMClientSecret) (*v1alpha1.AMApplication, error) {
	name := obj.Spec.ApplicationRef.Name
	app := &v1alpha1.AMApplication{}
	if err := k8s.GetClient().Get(ctx, types.NamespacedName{Namespace: obj.Namespace, Name: name}, app); err != nil {
		return nil, readError("AMApplication", name, err)
	}
	return app, nil
}

// ResolveRefs is the reconcile resolver, strict: it fills obj.Resolved with the domain key of the application
// and, unless obj is being deleted, the value of the Secret entry.
// A missing application, or on delete a missing domain, is returned as a NotFound: on delete the lifecycle then
// releases the finalizer, AM deleted the secret with its parent (S2).
// The value is not read on delete (S1/S4). Errors name the Secret and the key, never the value (AC3).
func ResolveRefs(ctx context.Context, obj *v1alpha1.AMClientSecret, _ string) error {
	app, err := getApplication(ctx, obj)
	if err != nil {
		return err
	}
	if !obj.IsBeingDeleted() && app.Status.Key == "" {
		return fmt.Errorf("AMApplication [%s]: %w", app.Name, ErrApplicationNotReady)
	}
	obj.Resolved.DomainKey = am.DomainKey(app)

	if obj.IsBeingDeleted() {
		return am.ResolveDomain(ctx, app, true)
	}

	value, err := readSecretValue(ctx, obj)
	if err != nil {
		return err
	}
	obj.Resolved.Value = value
	return nil
}

// ResolveAdmissionRefs is the admission resolver, tolerant: a missing or not ready application, a missing
// Secret or entry leave obj.Resolved partly empty and are admitted (apply in any order); PreCheck warns
// about them. Other errors are returned (AdmissionLifecycle turns them into a severe rejection).
func ResolveAdmissionRefs(ctx context.Context, obj *v1alpha1.AMClientSecret, ns string) error {
	if err := ResolveRefs(ctx, obj, ns); !isMissingRef(err) {
		return err
	}
	return nil
}

// CreateClient builds the AM client from the AMContext of the application's domain.
func CreateClient(ctx context.Context, obj *v1alpha1.AMClientSecret) (*am.Client, error) {
	app, err := getApplication(ctx, obj)
	if err != nil {
		return nil, err
	}
	return am.CreateSubResourceClient(ctx, app)
}

// AdmissionClient builds the AM client for admission. An application or domain that is missing or not yet
// created in AM gives no client and no error: the secret is admitted and the AM calls are skipped.
func AdmissionClient(ctx context.Context, obj *v1alpha1.AMClientSecret) (*am.Client, error) {
	app, err := getApplication(ctx, obj)
	if apierrors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if app.Status.Key == "" {
		return nil, nil
	}
	return am.AdmissionSubResourceClient(ctx, app)
}

// GetOwner returns the application, as a plain owner (false, not controller) (AC2).
func GetOwner(ctx context.Context, obj *v1alpha1.AMClientSecret) (client.Object, bool, error) {
	app, err := getApplication(ctx, obj)
	if err != nil {
		return nil, false, err
	}
	return app, false, nil
}

// readSecretValue reads the entry secretRef points to, in obj's namespace.
func readSecretValue(ctx context.Context, obj *v1alpha1.AMClientSecret) (string, error) {
	ref := obj.Spec.SecretRef
	secret := &corev1.Secret{}
	err := k8s.GetClient().Get(ctx, types.NamespacedName{Namespace: obj.Namespace, Name: ref.Name}, secret)
	if err != nil {
		return "", readError("Secret", ref.Name, err)
	}
	value := secret.Data[ref.Key]
	if len(value) == 0 {
		return "", fmt.Errorf("key [%s] of Secret [%s]: %w", ref.Key, ref.Name, ErrSecretEntryMissing)
	}
	return string(value), nil
}

// readError wraps a failed read with the kind and name of the referenced object, keeping a NotFound detectable.
func readError(kind, name string, err error) error {
	if apierrors.IsNotFound(err) {
		return fmt.Errorf("%s [%s] not found: %w", kind, name, err)
	}
	return fmt.Errorf("%s [%s] could not be read: %w", kind, name, err)
}

// isMissingRef tells a dependency that does not exist yet, which apply in any order allows.
func isMissingRef(err error) bool {
	return apierrors.IsNotFound(err) || errors.Is(err, ErrApplicationNotReady) || errors.Is(err, ErrSecretEntryMissing)
}
