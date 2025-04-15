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
<<<<<<< HEAD:internal/template/resolver.go
=======
	"encoding/json"
	"errors"
>>>>>>> cc9c9c0 (fix: update templated resources on source changes):internal/template/compiler.go
	"fmt"
	"strconv"
	"strings"
	"text/template"

<<<<<<< HEAD:internal/template/resolver.go
=======
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
>>>>>>> cc9c9c0 (fix: update templated resources on source changes):internal/template/compiler.go
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

<<<<<<< HEAD:internal/template/resolver.go
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
=======
type ctxKey string

const objectIDCtxKey = ctxKey("gravitee.io/templating/objectId")
const objectAnnotationKey = ctxKey("gravitee.io/templating/annotationKey")
const totalReferenceKey = "gravitee.io/references"

func Compile(ctx context.Context, obj runtime.Object, updateObjectMetadata bool) error {
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}
	return doCompile(
		context.WithValue(
			context.
				WithValue(ctx, objectIDCtxKey, getUnstructuredObjectID(u)),
			objectAnnotationKey, getObjectAnnotationName(u),
		),
		obj,
		updateObjectMetadata,
	)
}

func doCompile(ctx context.Context, obj runtime.Object, updateObjectMetadata bool) error {
	c, err := traverse(ctx, obj, updateObjectMetadata)
>>>>>>> cc9c9c0 (fix: update templated resources on source changes):internal/template/compiler.go
	if err != nil {
		return err
	}

<<<<<<< HEAD:internal/template/resolver.go
	funcMap := map[string]interface{}{
		"configmap": func(name string) (string, error) {
			return resolveConfigmap(ctx, obj, name)
		},
		"secret": func(name string) (string, error) {
			return resolveSecret(ctx, obj, name)
=======
	return runtime.DefaultUnstructuredConverter.FromUnstructured(objData, obj)
}

func exec(ctx context.Context, text, ns string, parentResourceDeleted, updateObjectMetadata bool) (string, error) {
	funcMap := map[string]interface{}{
		"configmap": func(name string) (string, error) {
			return resolveConfigmap(ctx, ns, name, parentResourceDeleted, updateObjectMetadata)
		},
		"secret": func(name string) (string, error) {
			return resolveSecret(ctx, ns, name, parentResourceDeleted, updateObjectMetadata)
>>>>>>> cc9c9c0 (fix: update templated resources on source changes):internal/template/compiler.go
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

<<<<<<< HEAD:internal/template/resolver.go
func resolveConfigmap(ctx context.Context, obj runtime.Object, name string) (string, error) {
=======
func resolveConfigmap(ctx context.Context, ns, name string, parentResourceDeleted,
	updateObjectMetadata bool) (string, error) {
>>>>>>> cc9c9c0 (fix: update templated resources on source changes):internal/template/compiler.go
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

<<<<<<< HEAD:internal/template/resolver.go
	if err = addFinalizer(ctx, cm); err != nil {
		return "", err
	}

	return cm.Data[sp[1]], nil
}

func resolveSecret(ctx context.Context, obj runtime.Object, name string) (string, error) {
=======
	key := sp[1]
	v := cm.Data[key]
	if v == "" {
		return "", fmt.Errorf("key [%s] not found in configmap [%s/%s]", key, ns, name)
	}
	if !updateObjectMetadata {
		return v, nil
	}

	totalReference, updated, err := updateAnnotation(ctx, cm, parentResourceDeleted)
	if err != nil {
		return "", err
	}

	if updated {
		updateFinalizer(ctx, cm, parentResourceDeleted && totalReference == 0)
		if err := k8s.UpdateSafely(ctx, cm); err != nil {
			return "", err
		}
	}

	return v, nil
}

func resolveSecret(ctx context.Context, ns, name string, parentResourceDeleted,
	updateObjectMetadata bool) (string, error) {
>>>>>>> cc9c9c0 (fix: update templated resources on source changes):internal/template/compiler.go
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

<<<<<<< HEAD:internal/template/resolver.go
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
=======
	key := sp[1]
	v := sec.Data[key]
	if len(v) == 0 {
		return "", fmt.Errorf("key [%s] not found in secret [%s/%s]", key, ns, name)
	}

	if !updateObjectMetadata {
		return string(v), nil
	}

	totalReference, updated, err := updateAnnotation(ctx, sec, parentResourceDeleted)
	if err != nil {
		return "", err
	}

	if updated {
		updateFinalizer(ctx, sec, parentResourceDeleted && totalReference == 0)
		if err := k8s.UpdateSafely(ctx, sec); err != nil {
			return "", err
		}
	}

	return string(v), nil
}

func updateFinalizer(_ context.Context, obj client.Object, unreferenced bool) {
	if unreferenced && util.ContainsFinalizer(obj, core.TemplatingFinalizer) {
		util.RemoveFinalizer(obj, core.TemplatingFinalizer)
		return
>>>>>>> cc9c9c0 (fix: update templated resources on source changes):internal/template/compiler.go
	}

	if !util.ContainsFinalizer(obj, core.TemplatingFinalizer) {
		util.AddFinalizer(obj, core.TemplatingFinalizer)
	}
}

func updateAnnotation(ctx context.Context, obj client.Object, parentResourceDeleted bool) (int, bool, error) {
	annotationKey, _ := ctx.Value(objectAnnotationKey).(string)
	objID, _ := ctx.Value(objectIDCtxKey).(string)
	annotationValue, ok := obj.GetAnnotations()[annotationKey]
	if !ok {
		annotationValue = "[]"
	}
	totalReferenceString, ok := obj.GetAnnotations()[totalReferenceKey]
	if !ok {
		totalReferenceString = "0"
	}
	totalReference, err := strconv.Atoi(totalReferenceString)
	if err != nil {
		return 0, false, err
	}
	values := make([]string, 0)
	if err := json.Unmarshal([]byte(annotationValue), &values); err != nil {
		return totalReference, false, err
	}

	updated := false
	valueSet := sets.New(values...)
	if parentResourceDeleted {
		updated = true
		totalReference--
		valueSet.Delete(objID)
	} else if !valueSet.Has(objID) {
		updated = true
		totalReference++
		valueSet.Insert(objID)
	}

	if updated {
		b, err := json.Marshal(sets.List(valueSet))
		if err != nil {
			return totalReference, false, err
		}
		if obj.GetAnnotations() == nil {
			obj.SetAnnotations(make(map[string]string))
		}

		obj.GetAnnotations()[annotationKey] = string(b)
		obj.GetAnnotations()[totalReferenceKey] = fmt.Sprintf("%v", totalReference)
		return totalReference, true, nil
	}

	return totalReference, false, nil
}

func getUnstructuredObjectID(unstructured map[string]interface{}) string {
	metadata, _ := unstructured["metadata"].(map[string]interface{})
	ns := metadata["namespace"]
	name := metadata["name"]
	return fmt.Sprintf("%v/%v", ns, name)
}

func getObjectAnnotationName(unstructured map[string]interface{}) string {
	kind := unstructured["kind"]
	return "gravitee.io/" + strings.ToLower(fmt.Sprintf("%v", kind)) + "s"
}
