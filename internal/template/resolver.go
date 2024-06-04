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
	"fmt"
	"strings"
	"text/template"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// example my-configmap/key1.
const ksPropertyLength = 2

func Compile(ctx context.Context, obj runtime.Object) error {
	switch t := obj.(type) {
	case *v1alpha1.ApiDefinition, *v1alpha1.ApiV4Definition, *v1alpha1.ManagementContext,
		*v1alpha1.Application, *netv1.Ingress, *v1alpha1.ApiResource:
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
		return err
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

	u := unstructured.Unstructured{Object: innerObj}
	nn := types.NamespacedName{Namespace: u.GetNamespace(), Name: sp[0]}
	cm := new(v1.ConfigMap)
	cli := k8s.GetClient()
	if err = cli.Get(ctx, nn, cm); err != nil {
		return "", err
	}

	if err = addFinalizer(ctx, cm); err != nil {
		return "", err
	}

	return cm.Data[sp[1]], nil
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

	u := unstructured.Unstructured{Object: innerObj}
	nn := types.NamespacedName{Namespace: u.GetNamespace(), Name: sp[0]}
	sec := new(v1.Secret)
	cli := k8s.GetClient()
	if err = cli.Get(ctx, nn, sec); err != nil {
		return "", err
	}

	if err = addFinalizer(ctx, sec); err != nil {
		return "", err
	}

	return string(sec.Data[sp[1]]), nil
}

func addFinalizer(ctx context.Context, obj client.Object) error {
	if !util.ContainsFinalizer(obj, keys.TemplatingFinalizer) {
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

		util.AddFinalizer(object, keys.TemplatingFinalizer)

		return cli.Update(ctx, object)
	}

	return nil
}
