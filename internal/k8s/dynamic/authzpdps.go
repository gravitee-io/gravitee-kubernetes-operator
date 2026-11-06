// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dynamic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

var (
	// ErrAuthzPdpNotFound is returned when the AuthzPdp a pdpRef names does not exist or is being
	// deleted. A resource referencing it waits for it.
	ErrAuthzPdpNotFound = errors.New("authz pdp not found")
	// ErrAuthzPdpNotSynced is returned when the AuthzPdp a pdpRef names has no routing scope yet:
	// the platform does not have it. A resource referencing it waits for it.
	ErrAuthzPdpNotSynced = errors.New("authz pdp not synced")
)

// ResolveAuthzPdpRoutingScope returns the routing scope of the AuthzPdp ref names, its namespace
// defaulted to the dependent's: what a schema, an entity or a policy sends as its target.
//
// A missing, deleted or not yet synced AuthzPdp is a state to wait out (see IsAuthzPdpPending):
// callers report ResolvedRefs=False and requeue, and admission warns and skips the dry run. An
// AuthzPdp of another ManagementContext is an error that no event resolves, whatever its state:
// its routing scope does not exist in the dependent's environment. It is checked first, so a
// deleted or unsynced AuthzPdp of another context is refused rather than waited for.
func ResolveAuthzPdpRoutingScope(
	ctx context.Context,
	dependent core.ContextAwareObject,
	ref refs.NamespacedName,
) (string, error) {
	if ref.Namespace == "" {
		ref.Namespace = dependent.GetNamespace()
	}

	pdp := &v1alpha1.AuthzPdp{}
	key := types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}
	if err := k8s.GetClient().Get(ctx, key, pdp); err != nil {
		if apierrors.IsNotFound(err) {
			return "", fmt.Errorf("%w: [%s]", ErrAuthzPdpNotFound, ref.String())
		}
		return "", fmt.Errorf("unable to read authz pdp [%s]: %w", ref.String(), err)
	}

	if want, got := ContextOf(dependent), ContextOf(pdp); want != got {
		return "", fmt.Errorf(
			"authz pdp [%s] belongs to management context [%s], not [%s]",
			ref.String(), got, want,
		)
	}

	if pdp.IsBeingDeleted() {
		return "", fmt.Errorf("%w: [%s] is being deleted", ErrAuthzPdpNotFound, ref.String())
	}

	if pdp.Status.RoutingScope == "" {
		return "", fmt.Errorf("%w: [%s] has no routing scope yet", ErrAuthzPdpNotSynced, ref.String())
	}

	return pdp.Status.RoutingScope, nil
}

// ResolveAuthzTarget returns what target sends on the wire: an empty string without a target, the
// scope as written, or the routing scope its pdpRef resolves to (see ResolveAuthzPdpRoutingScope,
// whose errors it returns unchanged).
func ResolveAuthzTarget(
	ctx context.Context,
	dependent core.ContextAwareObject,
	target *authzpdp.Target,
) (string, error) {
	if target == nil {
		return "", nil
	}
	if target.PdpRef == nil {
		return target.Scope, nil
	}
	return ResolveAuthzPdpRoutingScope(ctx, dependent, *target.PdpRef)
}

// ResolveAuthzTargets returns what targets send on the wire, in declaration order: each scope as
// written, each pdpRef resolved by ResolveAuthzPdpRoutingScope. A target that fails for another
// reason than a pending AuthzPdp is returned at once. The pending ones are returned together, as
// one error that IsAuthzPdpPending recognizes and whose message names each of them. No targets
// resolve to nil.
func ResolveAuthzTargets(
	ctx context.Context,
	dependent core.ContextAwareObject,
	targets []authzpdp.Target,
) ([]string, error) {
	var resolved []string
	var pending pendingTargets
	for i := range targets {
		scope, err := ResolveAuthzTarget(ctx, dependent, &targets[i])
		switch {
		case err == nil:
			resolved = append(resolved, scope)
		case IsAuthzPdpPending(err):
			pending = append(pending, err)
		default:
			return nil, err
		}
	}
	if len(pending) > 0 {
		return nil, pending
	}
	return resolved, nil
}

// pendingTargets holds the pdpRefs of a list of targets still to wait out, in declaration order.
type pendingTargets []error

func (p pendingTargets) Error() string {
	messages := make([]string, 0, len(p))
	for _, err := range p {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func (p pendingTargets) Unwrap() []error {
	return p
}

// IsAuthzPdpPending tells whether err is a pdpRef to wait out: an AuthzPdp missing, being
// deleted or not synced yet.
func IsAuthzPdpPending(err error) bool {
	return errors.Is(err, ErrAuthzPdpNotFound) || errors.Is(err, ErrAuthzPdpNotSynced)
}

// ContextOf returns the namespaced name of obj's ManagementContext, its namespace defaulted to
// obj's, or an empty string without one. Two resources sync to the same environment when their
// ContextOf are equal: a pdpRef resolves only then.
func ContextOf(obj core.ContextAwareObject) string {
	if !obj.HasContext() {
		return ""
	}
	ctxRef := obj.ContextRef()
	ns := ctxRef.GetNamespace()
	if ns == "" {
		ns = obj.GetNamespace()
	}
	return ns + "/" + ctxRef.GetName()
}
