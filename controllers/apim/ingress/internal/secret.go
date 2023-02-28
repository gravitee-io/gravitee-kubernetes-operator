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

package internal

import (
	"errors"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) createUpdateTLSSecret(ingress *v1.Ingress) error {
	if ingress.Spec.TLS == nil || len(ingress.Spec.TLS) == 0 {
		d.log.Info("no TLS will be configured")
		return nil
	}

	/*
		// We can also check if the TLS is already has a reference inside the rules.host section of the ingress
		for _, rule := range ingress.Spec.Rules
			for _, tls := range ingress.Spec.TLS
				for _, host := range tls.Hosts
					if rule.Host == host
	*/
	for _, tls := range ingress.Spec.TLS {
		secret := &core.Secret{}
		key := types.NamespacedName{Namespace: ingress.Namespace, Name: tls.SecretName}
		if err := d.k8s.Get(d.ctx, key, secret); err != nil {
			return err
		}

		// If finalizer not present, add it; This is a new object
		if !util.ContainsFinalizer(secret, keys.KeyPairFinalizer) {
			d.log.Info("adding finalizer to the tls secret")

			secret.ObjectMeta.Finalizers = append(secret.ObjectMeta.Finalizers, keys.KeyPairFinalizer)
			if err := d.k8s.Update(d.ctx, secret); err != nil {
				return client.IgnoreNotFound(err)
			}
		}

		// Secret has been deleted while it still has reference to an ingress, not allowed
		if !secret.DeletionTimestamp.IsZero() {
			return fmt.Errorf("secret can't be deleted because it has reference to an exsintg ingress [%s]", ingress.Name)
		}

		d.log.Info("Update GW keystore with new key pairs")
		if err := d.updateKeyInKeystore(secret); err != nil {
			return err
		}
	}

	d.log.Info("gateway keystore has been successfully update.")
	return nil
}

func (d *Delegate) deleteTLSSecret(ingress *v1.Ingress) error {
	for _, tls := range ingress.Spec.TLS {
		secret := &core.Secret{}
		key := types.NamespacedName{Namespace: ingress.Namespace, Name: tls.SecretName}
		if err := d.k8s.Get(d.ctx, key, secret); err != nil {
			return err
		}

		// It is possible that the same secret has been used in another ingress
		// We will not remove the finalizer and we will not remove the keypair
		// from keystore but we also don't throw any error to let the current ingress
		// be deleted
		hasReferenceToOtherIngress, err := d.secretHasReference(ingress, secret)
		if err != nil {
			return err
		}

		if hasReferenceToOtherIngress {
			d.log.Error(
				errors.New("secret has reference"),
				"secret is used by another ingress, it will not be deleted from the keystore")
		} else {
			// no reference to this secret, we can remove it from the keystore
			if err = d.removeKeyFromKeystore(secret); err != nil {
				return err
			}

			d.log.Info("removing finalizer from secret", "secret", secret.Name)
			util.RemoveFinalizer(secret, keys.KeyPairFinalizer)

			if err = d.k8s.Update(d.ctx, secret); err != nil {
				return err
			}
		}
	}

	d.log.Info("gateway keystore has been successfully updated.")
	return nil
}

func (d *Delegate) secretHasReference(ing *v1.Ingress, secret *core.Secret) (bool, error) {
	il, err := d.retrieveIngressListWithTLS(d.ctx, ing.Namespace)
	if err != nil {
		return false, err
	}

	for i := range il.Items {
		for _, tls := range il.Items[i].Spec.TLS {
			if tls.SecretName == secret.Name && il.Items[i].DeletionTimestamp.IsZero() {
				d.log.Info("the secret is already used inside an ingress resource", "resource", il.Items[i].Name)
				return true, nil
			}
		}
	}

	return false, nil
}
