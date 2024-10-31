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

package template

import (
	"bufio"
	"bytes"
	"context"
	coreerrors "errors"
	"fmt"
	"strings"
	"text/template"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// example my-configmap/key1.
const ksPropertyLength = 2

type CustomError struct{}

func (*CustomError) Error() string { return "heyo !" }

func Compile(ctx context.Context, obj runtime.Object) error {
	switch t := obj.(type) {
	case *v1alpha1.ApiDefinition, *v1alpha1.ApiV4Definition, *v1alpha1.ManagementContext,
		*v1alpha1.Application, *netv1.Ingress, *v1alpha1.ApiResource, *v1alpha1.Subscription:
		return exec(ctx, obj)
	default:
		return fmt.Errorf("unsupported object type %v", t)
	}
}

func exec(ctx context.Context, obj runtime.Object) error {
	text, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	funcMap := map[string]interface{}{
		"configmap": func(name string) (string, error) {
			return resolveConfigmap(ctx, obj, name)
		},
		"secret": func(name string) (string, error) {
			return resolveSecret(ctx, obj, name)
		},
	}

	tmpl, err := template.New("gko").Funcs(template.FuncMap(funcMap)).Delims("[[", "]]").Parse(string(text))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	if err = tmpl.Execute(writer, make(map[string]string)); err != nil {
		uErr := coreerrors.Unwrap(coreerrors.Unwrap(err))
		return uErr
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return yaml.Unmarshal(buf.Bytes(), obj)
}

func resolveConfigmap(ctx context.Context, obj runtime.Object, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty configmap name")
	}

	sp := strings.Split(name, "/")
	if len(sp) != ksPropertyLength {
		return "", fmt.Errorf("wrong configmap name. Example my-configmap/key1")
	}

	innerObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return "", err
	}

	name = sp[0]
	u := unstructured.Unstructured{Object: innerObj}
	nn := types.NamespacedName{Namespace: u.GetNamespace(), Name: name}
	cm := new(v1.ConfigMap)
	cli := k8s.GetClient()

	err = cli.Get(ctx, nn, cm)
	if errors.IsNotFound(err) {
		return "", fmt.Errorf("configmap [%s/%s] not found", u.GetNamespace(), name)
	}

	if err != nil {
		return "", err
	}

	if err = addFinalizer(ctx, cm); err != nil {
		return "", err
	}

	key := sp[1]
	v := cm.Data[key]
	if v == "" {
		return "", fmt.Errorf("key [%s] not found in configmap [%s/%s]", key, u.GetNamespace(), name)
	}
	return v, nil
}

func resolveSecret(ctx context.Context, obj runtime.Object, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty secret name")
	}

	sp := strings.Split(name, "/")
	if len(sp) != ksPropertyLength {
		return "", fmt.Errorf("wrong secret name. Example my-secret/key1")
	}

	innerObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return "", err
	}

	name = sp[0]
	u := unstructured.Unstructured{Object: innerObj}
	nn := types.NamespacedName{Namespace: u.GetNamespace(), Name: name}
	sec := new(v1.Secret)
	cli := k8s.GetClient()

	err = cli.Get(ctx, nn, sec)
	if errors.IsNotFound(err) {
		return "", fmt.Errorf("secret [%s/%s] not found", u.GetNamespace(), name)
	}

	if err != nil {
		return "", err
	}

	if err = addFinalizer(ctx, sec); err != nil {
		return "", err
	}

	key := sp[1]
	v := sec.Data[key]
	if len(v) == 0 {
		return "", fmt.Errorf("key [%s] not found in secret [%s/%s]", key, u.GetNamespace(), name)
	}
	return string(v), nil
}

func addFinalizer(ctx context.Context, obj client.Object) error {
	if !util.ContainsFinalizer(obj, core.TemplatingFinalizer) {
		var object client.Object
		switch obj.(type) {
		case *v1.ConfigMap:
			object = new(v1.ConfigMap)
		case *v1.Secret:
			object = new(v1.Secret)
		}

		nn := types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.GetName()}
		cli := k8s.GetClient()
		if err := cli.Get(ctx, nn, object); err != nil {
			return err
		}

		util.AddFinalizer(object, core.TemplatingFinalizer)

		return cli.Update(ctx, object)
	}

	return nil
}
