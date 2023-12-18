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

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// example my-configmap/key1.
const ksPropertyLength = 2

type Resolver struct {
	ctx    context.Context
	client client.Client
	obj    runtime.Object
}

func NewResolver(ctx context.Context, c client.Client, obj runtime.Object) *Resolver {
	return &Resolver{ctx: ctx, client: c, obj: obj}
}

func (r *Resolver) Resolve() error {
	switch t := r.obj.(type) {
	case *v1beta1.ApiDefinition, *v1beta1.ManagementContext, *v1beta1.Application, *netv1.Ingress, *v1beta1.ApiResource:
		return r.exec()
	default:
		return fmt.Errorf("unsupported object type %v", t)
	}
}

func (r *Resolver) exec() error {
	text, err := yaml.Marshal(r.obj)
	if err != nil {
		return err
	}

	funcMap := map[string]interface{}{
		"configmap": r.resolveConfigmap,
		"secret":    r.resolveSecret,
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

	return yaml.Unmarshal(buf.Bytes(), r.obj)
}

func (r *Resolver) resolveConfigmap(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty configmap name")
	}

	sp := strings.Split(name, "/")
	if len(sp) != ksPropertyLength {
		return "", fmt.Errorf("wrong configmap name. Example my-configmap/key1")
	}

	innerObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r.obj)
	if err != nil {
		return "", err
	}

	u := unstructured.Unstructured{Object: innerObj}
	nn := types.NamespacedName{Namespace: u.GetNamespace(), Name: sp[0]}
	cm := new(v1.ConfigMap)
	if err = r.client.Get(r.ctx, nn, cm); err != nil {
		return "", err
	}

	if err = r.addFinalizer(cm); err != nil {
		return "", err
	}

	return cm.Data[sp[1]], nil
}

func (r *Resolver) resolveSecret(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty secret name")
	}

	sp := strings.Split(name, "/")
	if len(sp) != ksPropertyLength {
		return "", fmt.Errorf("wrong secret name. Example my-secret/key1")
	}

	innerObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r.obj)
	if err != nil {
		return "", err
	}

	u := unstructured.Unstructured{Object: innerObj}
	nn := types.NamespacedName{Namespace: u.GetNamespace(), Name: sp[0]}
	sec := new(v1.Secret)
	if err = r.client.Get(r.ctx, nn, sec); err != nil {
		return "", err
	}

	if err = r.addFinalizer(sec); err != nil {
		return "", err
	}

	return string(sec.Data[sp[1]]), nil
}

func (r *Resolver) addFinalizer(obj client.Object) error {
	if !util.ContainsFinalizer(obj, keys.TemplatingFinalizer) {
		var object client.Object
		switch obj.(type) {
		case *v1.ConfigMap:
			object = new(v1.ConfigMap)
		case *v1.Secret:
			object = new(v1.Secret)
		}

		nn := types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.GetName()}
		if err := r.client.Get(r.ctx, nn, object); err != nil {
			return err
		}

		util.AddFinalizer(object, keys.TemplatingFinalizer)

		return r.client.Update(r.ctx, object)
	}

	return nil
}
