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

package drift

import (
	"log"
	"reflect"
	"strings"
	"unicode"
)

type valuePair struct {
	Value     reflect.Value
	Interface any
}

// Detect detects the drift between two structs.
// Goes recursively through the structs and detects the drift.
// It uses the drift.EquivalenceFunc to determine if the two values are equivalent.
// Each field is tagged with a drift.EquivalenceFunc name to determine which EquivalenceFunc to use.
// By default, the EquivalenceFunc is reflect.DeepEqual.
// The result is a tree of Result. It can be printed in a pseudo-yaml format.
func Detect(crd any, api any) Result {
	res := Result{Children: []*Result{}}
	if crd != nil || api != nil {
		assertRootIsStruct(crd, api)
		detectStruct(crd, api, &res, false)
	}
	return res
}

func assertRootIsStruct(crd any, api any) {
	if crd != nil {
		if reflect.TypeOf(crd).Kind() != reflect.Struct {
			log.Panicf("detect drift only supports structs, crd was '%T'.", crd)
		}
	}
	if api != nil {
		if reflect.TypeOf(api).Kind() != reflect.Struct {
			log.Panicf("detect drift only supports structs, api was '%T'.", api)
		}
	}
}

func detectStruct(crd any, api any, this *Result, ordered bool) {
	t, skip := getTypeOrSkip(crd, api)
	if skip {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		// get info to find an Equivalence func
		field := t.Field(i)
		funcName := field.Tag.Get("drift")
		// use json tag or infer the name of the field
		property := getProperty(field)

		// remove pointer type if present
		fieldType := dereferenced(field.Type)

		// get the value of the field
		crdPair := fromField(crd, i, fieldType)
		apiPair := fromField(api, i, fieldType)

		if fieldType.Kind() == reflect.Slice {
			equivalenceFunc := equivalenceRegistry.Get(funcName, fieldType.Kind())
			equivalent := equivalenceFunc(crdPair.Interface, apiPair.Interface)
			// if the equivalence on the slice does not skip items
			if !equivalent.Skip {
				// process all items
				detectItems(property, crdPair.Value, apiPair.Value, this, ordered)
			}
		}

		if fieldType.Kind() == reflect.Map {
			equivalenceFunc := equivalenceRegistry.Get(funcName, fieldType.Kind())
			equivalent := equivalenceFunc(crdPair.Interface, apiPair.Interface)
			// if the equivalence on the map does not skip items
			if !equivalent.Skip {
				detectMapItems(property, funcName, crdPair.Value, apiPair.Value, this)
			}
		}

		if fieldType.Kind() == reflect.Struct {
			handleStructField(property, funcName, field, crdPair.Interface, apiPair.Interface, this)
		} else {
			// handle the field as a simple value
			equivalenceFunc := equivalenceRegistry.Get(funcName, fieldType.Kind())
			equivalent := equivalenceFunc(crdPair.Interface, apiPair.Interface)
			this.AppendChild(&Result{
				Property:    property,
				Equivalence: equivalent,
				CRDValue:    crdPair.Interface,
				APIValue:    apiPair.Interface,
			}, ordered)
		}
	}
}

func getTypeOrSkip(crd any, api any) (reflect.Type, bool) {
	t := reflect.TypeOf(crd)
	if t == nil {
		t = reflect.TypeOf(api)
	}
	if crd == nil && api == nil {
		return nil, true
	}
	return t, false
}

func getProperty(field reflect.StructField) string {
	property := field.Tag.Get("json")
	if property == "" {
		property = unTitle(field.Name)
	} else {
		// json can have after the name (inline, omitempty)
		property = trimAfterComma(property)
	}
	return property
}

func unTitle(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func trimAfterComma(property string) string {
	return strings.Split(property, ",")[0]
}

func dereferenced(fieldType reflect.Type) reflect.Type {
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}
	return fieldType
}

func fromField(v any, i int, t reflect.Type) valuePair {
	if v == nil {
		return valuePair{
			Value: reflect.Zero(t),
		}
	}
	value := reflect.ValueOf(v).Field(i)
	return valuePair{
		Value:     value,
		Interface: asInterface(value),
	}
}

func asInterface(v reflect.Value) any {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return v.Elem().Interface()
	}
	return v.Interface()
}

func detectItems(property string, crdItems reflect.Value, apiItems reflect.Value, parent *Result, ordered bool) {
	crdSize := crdItems.Len()
	apiSize := apiItems.Len()

	if crdSize == apiSize {
		if crdSize > 0 {
			detectSymmetrical(property, crdItems, apiItems, parent, ordered)
		}
	} else {
		detectAsymmetrical(property, crdItems, apiItems, crdSize > apiSize, parent, ordered)
	}
}

func detectSymmetrical(property string, crdItems reflect.Value, apiItems reflect.Value, parent *Result, ordered bool) {
	for i := 0; i < crdItems.Len(); i++ {
		crdItem := crdItems.Index(i)
		apiItem := apiItems.Index(i)
		detectItem(property, i, crdItem, apiItem, parent, ordered)
	}
}

func detectAsymmetrical(property string, crdItems reflect.Value, apiItems reflect.Value, crdIsLarger bool, parent *Result, ordered bool) {
	var leader reflect.Value
	var follower reflect.Value
	if crdIsLarger {
		leader = crdItems
		follower = apiItems
	} else {
		leader = apiItems
		follower = crdItems
	}
	followerSize := follower.Len()
	for i := 0; i < leader.Len(); i++ {
		leaderItem := leader.Index(i)
		var followerItem reflect.Value
		if followerSize > i {
			followerItem = follower.Index(i)
		} else {
			// Create an empty item
			followerItem = reflect.Zero(leaderItem.Type())
		}
		if crdIsLarger {
			detectItem(property, i, leaderItem, followerItem, parent, ordered)
		} else {
			detectItem(property, i, followerItem, leaderItem, parent, ordered)
		}
	}
}

func detectItem(property string, i int, crdItem reflect.Value, apiItem reflect.Value, parent *Result, ordered bool) {
	if crdItem.Kind() == reflect.Struct {
		child := parent.AppendChild(&Result{Property: property, Children: []*Result{}, Index: &i}, ordered)
		detectStruct(asInterface(crdItem), asInterface(apiItem), child, false)
	} else {
		crdValue := asInterface(crdItem)
		apiValue := asInterface(apiItem)
		equivalenceFunc := equivalenceRegistry.Get("", crdItem.Kind())
		equivalence := equivalenceFunc(crdValue, apiValue)
		parent.AppendChild(&Result{
			Property:    property,
			Index:       &i,
			CRDValue:    crdValue,
			APIValue:    apiValue,
			Equivalence: equivalence,
		}, ordered)
	}
}

func detectMapItems(property string, funcName string, crdEntries reflect.Value, apiEntries reflect.Value, parent *Result) {
	crdKeyValues := crdEntries.MapKeys()
	apiKeyValues := apiEntries.MapKeys()
	if len(crdKeyValues) == 0 && len(apiKeyValues) == 0 {
		return
	}

	// create the result for the whole map
	child := parent.AppendChild(&Result{Property: property, Children: []*Result{}}, false)

	all := make(map[string]reflect.Value)
	collectKeys(crdKeyValues, all)
	collectKeys(apiKeyValues, all)

	for key, keyValue := range all {
		var typ reflect.Type

		// get the value and kind if the entry exist
		crdValue, crdInterface, ok := getValue(crdEntries, keyValue)
		if ok {
			typ = dereferenced(crdValue.Type())
		}
		apiValue, apiInterface, ok := getValue(apiEntries, keyValue)
		if ok {
			typ = dereferenced(apiValue.Type())
		}

		// both are nil, skip
		if crdInterface == nil && apiInterface == nil {
			continue
		}

		detectEntry(key, funcName, typ, valuePair{crdValue, crdInterface}, valuePair{apiValue, apiInterface}, child)
	}
}

func collectKeys(keyValues []reflect.Value, recipient map[string]reflect.Value) {
	for _, keyValue := range keyValues {
		if key, ok := asInterface(keyValue).(string); ok {
			recipient[key] = keyValue
		} else {
			log.Panicf("map key must be of type string, got %T", asInterface(keyValue))
		}
	}
}

func getValue(crdEntries reflect.Value, keyValue reflect.Value) (reflect.Value, any, bool) {
	if value := crdEntries.MapIndex(keyValue); value.IsValid() {
		return value, asInterface(value), true
	}
	return reflect.Value{}, nil, false
}

func detectEntry(key, funcName string, typ reflect.Type, crd, api valuePair, parent *Result) {
	switch {
	case typ.Kind() == reflect.Struct:
		crd.setZeroIfNilValue(typ)
		api.setZeroIfNilValue(typ)
		child := parent.AppendChild(&Result{Property: key, Children: []*Result{}}, true)
		detectStruct(crd.Interface, api.Interface, child, true)
	case typ.Kind() == reflect.Slice:
		crd.setEmptySliceIfNilValue(typ)
		api.setEmptySliceIfNilValue(typ)
		detectItems(key, crd.Value, api.Value, parent, true)
	case typ.Kind() == reflect.Map:
		crd.setEmptyMapIfNilValue(typ)
		api.setEmptyMapIfNilValue(typ)
		detectMapItems(key, funcName, crd.Value, api.Value, parent)
	default:
		equivalenceFunc := equivalenceRegistry.Get(funcName, typ.Kind())
		equivalent := equivalenceFunc(crd.Interface, api.Interface)
		parent.AppendChild(&Result{
			Property:    key,
			Equivalence: equivalent,
			CRDValue:    crd.Interface,
			APIValue:    api.Interface,
		}, true)
	}
}

func (p *valuePair) setZeroIfNilValue(typ reflect.Type) {
	if p.Interface == nil {
		zero := reflect.Zero(typ)
		p.Interface = asInterface(zero)
	}
}

func (p *valuePair) setEmptyMapIfNilValue(typ reflect.Type) {
	if p.Interface == nil {
		p.Value = reflect.MakeMap(typ)
		p.Interface = asInterface(p.Value)
	}
}
func (p *valuePair) setEmptySliceIfNilValue(typ reflect.Type) {
	if p.Interface == nil {
		p.Value = reflect.MakeSlice(typ, 0, 0)
		p.Interface = asInterface(p.Value)
	}
}

func handleStructField(property string, funcName string, field reflect.StructField, crd any, api any, this *Result) {
	if isEmbeddedStruct(field) {
		detectStruct(crd, api, this, false)
	} else {
		equivalenceFunc := equivalenceRegistry.Get(funcName, reflect.Struct)
		this.Equivalence = equivalenceFunc(crd, api)
		// if the equivalence on the struct does not skip items
		if !this.Equivalence.Skip {
			child := &Result{Property: property, Children: []*Result{}}
			this.Children = append(this.Children, child)
			detectStruct(crd, api, child, false)
		}
	}
}

func isEmbeddedStruct(field reflect.StructField) bool {
	if !field.Anonymous {
		return false
	}
	if field.Type.Kind() == reflect.Struct {
		return true
	}
	if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct {
		return true
	}
	return false
}
