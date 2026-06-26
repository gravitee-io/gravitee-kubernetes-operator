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
func Detect(crd any, remote any) Result {
	res := Result{Children: []*Result{}}
	if crd != nil || remote != nil {
		assertRootIsStruct(crd, remote)
		detectStruct(crd, remote, &res, false)
	}
	return res
}

func assertRootIsStruct(crd any, remote any) {
	if crd != nil {
		if reflect.TypeOf(crd).Kind() != reflect.Struct {
			log.Panicf("detect drift only supports structs, crd was '%T'.", crd)
		}
	}
	if remote != nil {
		if reflect.TypeOf(remote).Kind() != reflect.Struct {
			log.Panicf("detect drift only supports structs, remote was '%T'.", remote)
		}
	}
}

func detectStruct(crd any, remote any, this *Result, ordered bool) {
	t, skip := getTypeOrSkip(crd, remote)
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
		remotePair := fromField(remote, i, fieldType)

		if fieldType.Kind() == reflect.Slice {
			equivalenceFunc := equivalenceRegistry.Get(funcName, fieldType.Kind())
			equivalent := equivalenceFunc(crdPair.Interface, remotePair.Interface)
			// if the equivalence on the slice does not skip items
			if !equivalent.Skip {
				// process all items
				detectItems(property, crdPair.Value, remotePair.Value, this, ordered)
			}
		}

		if fieldType.Kind() == reflect.Map {
			equivalenceFunc := equivalenceRegistry.Get(funcName, fieldType.Kind())
			equivalent := equivalenceFunc(crdPair.Interface, remotePair.Interface)
			// if the equivalence on the map does not skip items
			if !equivalent.Skip {
				detectMapItems(property, funcName, crdPair.Value, remotePair.Value, this)
			}
		}

		if fieldType.Kind() == reflect.Struct {
			handleStructField(property, funcName, field, crdPair.Interface, remotePair.Interface, this)
		} else {
			// handle the field as a simple value
			equivalenceFunc := equivalenceRegistry.Get(funcName, fieldType.Kind())
			equivalent := equivalenceFunc(crdPair.Interface, remotePair.Interface)
			this.AppendChild(&Result{
				Property:    property,
				Equivalence: equivalent,
				CRDValue:    crdPair.Interface,
				APIValue:    remotePair.Interface,
			}, ordered)
		}
	}
}

func getTypeOrSkip(crd any, remote any) (reflect.Type, bool) {
	t := reflect.TypeOf(crd)
	if t == nil {
		t = reflect.TypeOf(remote)
	}
	if crd == nil && remote == nil {
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

func detectItems(property string, crdItems reflect.Value, remoteItems reflect.Value, parent *Result, ordered bool) {
	crdSize := crdItems.Len()
	remoteSize := remoteItems.Len()

	if crdSize == remoteSize {
		if crdSize > 0 {
			detectSymmetrical(property, crdItems, remoteItems, parent, ordered)
		}
	} else {
		detectAsymmetrical(property, crdItems, remoteItems, crdSize > remoteSize, parent, ordered)
	}
}

func detectSymmetrical(property string, crdItems reflect.Value, remoteItems reflect.Value, parent *Result, ordered bool) {
	for i := 0; i < crdItems.Len(); i++ {
		crdItem := crdItems.Index(i)
		remoteItem := remoteItems.Index(i)
		detectItem(property, i, crdItem, remoteItem, parent, ordered)
	}
}

func detectAsymmetrical(property string, crdItems reflect.Value, remoteItems reflect.Value,
	crdIsLarger bool, parent *Result, ordered bool) {
	var leader reflect.Value
	var follower reflect.Value
	if crdIsLarger {
		leader = crdItems
		follower = remoteItems
	} else {
		leader = remoteItems
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

func detectItem(property string, i int, crdItem reflect.Value, remoteItem reflect.Value, parent *Result, ordered bool) {
	if crdItem.Kind() == reflect.Struct {
		child := parent.AppendChild(&Result{Property: property, Children: []*Result{}, Index: &i}, ordered)
		detectStruct(asInterface(crdItem), asInterface(remoteItem), child, false)
	} else if crdItem.Kind() == reflect.Map {
		detectIndexedMapItems(property, &i, "", crdItem, remoteItem, parent)
	} else {
		runEquivalence := true
		if crdItem.Kind() == reflect.Interface {
			runEquivalence = detectAny(property, "",
				valuePair{crdItem, asInterface(crdItem)},
				valuePair{remoteItem, asInterface(remoteItem)},
				parent)
		}
		if runEquivalence {
			crdValue := asInterface(crdItem)
			remoteValue := asInterface(remoteItem)
			equivalenceFunc := equivalenceRegistry.Get("", crdItem.Kind())
			equivalence := equivalenceFunc(crdValue, remoteValue)
			parent.AppendChild(&Result{
				Property:    property,
				Index:       &i,
				CRDValue:    crdValue,
				APIValue:    remoteValue,
				Equivalence: equivalence,
			}, ordered)
		}
	}
}

func detectIndexedMapItems(
	property string,
	i *int,
	funcName string,
	crdEntries reflect.Value,
	remoteEntries reflect.Value,
	parent *Result) {
	crdKeyValues := crdEntries.MapKeys()
	remoteKeyValues := remoteEntries.MapKeys()
	if len(crdKeyValues) == 0 && len(remoteKeyValues) == 0 {
		return
	}

	// create the result for the whole map
	child := parent.AppendChild(&Result{Property: property, Index: i, Children: []*Result{}}, false)

	all := make(map[string]reflect.Value)
	collectKeys(crdKeyValues, all)
	collectKeys(remoteKeyValues, all)

	for key, keyValue := range all {
		var typ reflect.Type

		// get the value and kind if the entry exist
		crdValue, crdInterface, ok := getValue(crdEntries, keyValue)
		if ok {
			typ = dereferenced(crdValue.Type())
		}
		remoteValue, remoteInterface, ok := getValue(remoteEntries, keyValue)
		if ok {
			typ = dereferenced(remoteValue.Type())
		}

		// both are nil, skip
		if crdInterface == nil && remoteInterface == nil {
			continue
		}

		detectEntry(key, funcName, typ, valuePair{crdValue, crdInterface}, valuePair{remoteValue, remoteInterface}, child)
	}
}

func detectMapItems(property string, funcName string, crdEntries reflect.Value, remoteEntries reflect.Value, parent *Result) {
	detectIndexedMapItems(property, nil, funcName, crdEntries, remoteEntries, parent)
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

func detectEntry(key, funcName string, typ reflect.Type, crd, remote valuePair, parent *Result) {
	switch {
	case typ.Kind() == reflect.Struct:
		crd.setZeroIfNilValue(typ)
		remote.setZeroIfNilValue(typ)
		child := parent.AppendChild(&Result{Property: key, Children: []*Result{}}, true)
		detectStruct(crd.Interface, remote.Interface, child, true)
	case typ.Kind() == reflect.Slice:
		crd.setEmptySliceIfNilValue(typ)
		remote.setEmptySliceIfNilValue(typ)
		detectItems(key, crd.Value, remote.Value, parent, true)
	case typ.Kind() == reflect.Map:
		crd.setEmptyMapIfNilValue(typ)
		remote.setEmptyMapIfNilValue(typ)
		detectMapItems(key, funcName, crd.Value, remote.Value, parent)
	default:
		runEquivalence := true
		if typ.Kind() == reflect.Interface {
			runEquivalence = detectAny(key, funcName, crd, remote, parent)
		}
		if runEquivalence {
			equivalenceFunc := equivalenceRegistry.Get(funcName, typ.Kind())
			equivalent := equivalenceFunc(crd.Interface, remote.Interface)
			parent.AppendChild(&Result{
				Property:    key,
				Equivalence: equivalent,
				CRDValue:    crd.Interface,
				APIValue:    remote.Interface,
			}, true)
		}
	}
}

func detectAny(key string, funcName string, crd valuePair, remote valuePair, parent *Result) bool {
	crdElem := reflect.ValueOf(crd.Interface)
	remoteElem := reflect.ValueOf(remote.Interface)
	if crdElem.Kind() == reflect.Struct || remoteElem.Kind() == reflect.Struct {
		// create a zero struct with the same type so it can be introspected
		if remoteElem.Kind() == reflect.Invalid {
			remote.setZeroIfNilValue(crdElem.Type())
		} else if crdElem.Kind() == reflect.Invalid {
			crd.setZeroIfNilValue(remoteElem.Type())
		}
		child := parent.AppendChild(&Result{Property: key, Children: []*Result{}}, true)
		detectStruct(crd.Interface, remote.Interface, child, true)
		return false
	}
	if crdElem.Kind() == reflect.Map || remoteElem.Kind() == reflect.Map {
		// create a map with the same type so it can be introspected
		if remoteElem.Kind() == reflect.Invalid {
			remoteElem = reflect.MakeMap(crdElem.Type())
		} else if crdElem.Kind() == reflect.Invalid {
			crdElem = reflect.MakeMap(remoteElem.Type())
		}
		detectMapItems(key, funcName, crdElem, remoteElem, parent)
		return false
	}
	if crdElem.Kind() == reflect.Slice || remoteElem.Kind() == reflect.Slice {
		// create a slice with the same type so it can be introspected
		if remoteElem.Kind() == reflect.Invalid {
			remoteElem = reflect.MakeSlice(crdElem.Type(), 0, 0)
		} else if crdElem.Kind() == reflect.Invalid {
			crdElem = reflect.MakeSlice(remoteElem.Type(), 0, 0)
		}
		detectItems(key, crdElem, remoteElem, parent, true)
		return false
	}
	return true
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

func handleStructField(property string, funcName string, field reflect.StructField, crd any, remote any, this *Result) {
	if isEmbeddedStruct(field) {
		detectStruct(crd, remote, this, false)
	} else {
		equivalenceFunc := equivalenceRegistry.Get(funcName, reflect.Struct)
		this.Equivalence = equivalenceFunc(crd, remote)
		// if the equivalence on the struct does not skip items
		if !this.Equivalence.Skip {
			child := &Result{Property: property, Children: []*Result{}}
			this.Children = append(this.Children, child)
			detectStruct(crd, remote, child, false)
			if this.PostFunc != nil {
				this.PostFunc(child)
			}
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
