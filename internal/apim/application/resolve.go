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

package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ResolveRefs inlines the client certificates the application references, for the lifecycle. A deleted application
// needs none: Delete addresses it by its hrid, so a certificate Secret deleted first does not leave it in APIM.
func ResolveRefs(ctx context.Context, obj *v1alpha1.Application, namespace string) error {
	if obj.Spec.Settings == nil || obj.IsBeingDeleted() {
		return nil
	}
	return ResolveClientCertificates(ctx, obj.Spec.Settings, namespace, obj.Spec.Name)
}

func ResolveClientCertificates(ctx context.Context, settings *application.Setting, namespace, appName string) error {
	tls := settings.TLS
	if tls == nil || len(tls.ClientCertificates) == 0 {
		return nil
	}

	for i := range tls.ClientCertificates {
		cert := &tls.ClientCertificates[i]

		if err := resolveRef(ctx, namespace, cert); err != nil {
			return err
		}

		if err := decodeIfEncoded(cert); err != nil {
			return gerrors.NewResolveRefError(
				fmt.Errorf("failed to base64-decode certificate [%d]: %w", i, err),
			)
		}

		if cert.Name == "" {
			cert.Name = appName + "-" + strconv.Itoa(i)
		}
	}

	return nil
}

// resolveRef leaves a certificate that sets both content and ref as it is, for admission to reject.
func resolveRef(ctx context.Context, namespace string, cert *application.ClientCertificate) error {
	if cert.Ref == nil || cert.Content != "" {
		return nil
	}

	ref := cert.Ref
	ns := ref.Namespace
	if ns == "" {
		ns = namespace
	}

	key := types.NamespacedName{Namespace: ns, Name: ref.Name}

	content, err := fetchRefContent(ctx, ref.Kind, key, ref.Key)
	if err != nil {
		return gerrors.NewResolveRefError(
			fmt.Errorf("failed to resolve certificate ref [%s/%s] key [%s]: %w",
				key.Namespace, key.Name, ref.Key, err),
		)
	}

	cert.Content = content
	cert.Ref = nil

	return nil
}

func fetchRefContent(ctx context.Context, kind string, key types.NamespacedName, dataKey string) (string, error) {
	cli := k8s.GetClient()

	switch kind {
	case "secrets", "":
		secret := &corev1.Secret{}
		if err := cli.Get(ctx, key, secret); err != nil {
			return "", err
		}
		v, ok := secret.Data[dataKey]
		if !ok {
			return "", fmt.Errorf("key [%s] not found in secret [%s]", dataKey, key)
		}
		return string(v), nil
	case "configmaps":
		cm := &corev1.ConfigMap{}
		if err := cli.Get(ctx, key, cm); err != nil {
			return "", err
		}
		v, ok := cm.Data[dataKey]
		if !ok {
			return "", fmt.Errorf("key [%s] not found in configmap [%s]", dataKey, key)
		}
		return v, nil
	default:
		return "", fmt.Errorf("unsupported ref kind [%s]", kind)
	}
}

func decodeIfEncoded(cert *application.ClientCertificate) error {
	if !cert.Encoded || cert.Content == "" {
		return nil
	}

	decoded, err := base64.StdEncoding.DecodeString(cert.Content)
	if err != nil {
		return err
	}
	cert.Content = string(decoded)
	cert.Encoded = false

	return nil
}
