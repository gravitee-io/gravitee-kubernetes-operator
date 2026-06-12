package drift

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"unicode"
)

type EquivalentStatus byte

const (
	CannotCompare EquivalentStatus = iota
	Equivalent    EquivalentStatus = 1
	Inequivalent  EquivalentStatus = 2
)

type Equivalence struct {
	Equivalent EquivalentStatus
	Reason     string
}

type Result struct {
	Equivalence
	Property string
	CRDValue any
	APIValue any
	Children []Result
}

func (r Result) String() string {
	var builder strings.Builder
	format(r, &builder, -2)
	return builder.String()
}

func (r Result) DriftDetected() bool {
	if r.Equivalent == Inequivalent {
		return true
	}
	if len(r.Children) > 0 {
		for _, child := range r.Children {
			if child.DriftDetected() {
				return true
			}
		}
	}
	return false
}

func format(r Result, b *strings.Builder, indent int) {
	if indent > 0 {
		b.WriteString(strings.Repeat(" ", indent))
	}
	if len(r.Children) == 0 {
		if r.Equivalent == Inequivalent {
			b.WriteString(fmt.Sprintf("%s: %v <> %v", r.Property, resolve(r.CRDValue), resolve(r.APIValue)))
			if r.Reason != "" {
				b.WriteString(fmt.Sprintf(" (%s)", r.Reason))
			}
			return
		}
		if r.Equivalent == CannotCompare {
			b.WriteString(fmt.Sprintf("%s: [incomparable]\n", r.Property))
		}
		return
	}

	for i, child := range r.Children {
		if r.Property != "" {
			b.WriteString(fmt.Sprintf("%s:", r.Property))
		}
		format(child, b, indent+2)
		if i < len(r.Children)-1 && child.Equivalence.Equivalent == Inequivalent {
			b.WriteString("\n")
		}
	}
}

func resolve(v any) any {
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return "<nil>"
		}
		return rv.Elem().Interface()
	}

	return v
}

type EquivalenceFunc func(crd any, api any) Equivalence

type registryEntry struct {
	kind reflect.Kind
	name string
}

type registry struct {
	registry map[registryEntry]EquivalenceFunc
}

func (r *registry) Get(name string, t reflect.Kind) EquivalenceFunc {
	if f, ok := r.registry[registryEntry{
		name: name,
		kind: t,
	}]; ok {
		return f
	}
	return defaultEquivalence
}

var equivalenceRegistry = &registry{
	registry: make(map[registryEntry]EquivalenceFunc),
}

// Register registers a named EquivalenceFunc function 'f' to compare value of kind 'k'.
func Register(name string, k reflect.Kind, f EquivalenceFunc) {
	if k == reflect.Ptr {
		panic("cannot register a pointer to a struct, use a concrete type or an interface")
	}

	equivalenceRegistry.registry[registryEntry{
		kind: k,
		name: name,
	}] = f

}

func Detect(crd any, api any) Result {
	res := Result{Children: []Result{}}
	detect(crd, api, &res)
	return res
}

func detect(crd any, api any, parent *Result) {

	t := reflect.TypeOf(crd)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		log.Panicf("expected a struct got a: %t", t)
	}

	for i := 0; i < t.NumField(); i++ {
		// get info to find an Equivalence func
		funcName := t.Field(i).Tag.Get("drift")
		property := t.Field(i).Tag.Get("json")
		if property == "" {
			property = unTitle(t.Field(i).Name)
		} else {
			property = trimAfterComma(property)
		}
		fieldType := t.Field(i).Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		crdValue := cleanValue(reflect.ValueOf(crd).Field(i))
		apiValue := cleanValue(reflect.ValueOf(api).Field(i))

		if fieldType.Kind() == reflect.Struct {
			child := &Result{Property: property, Children: []Result{}}
			parent.Children = append(parent.Children, *child)
			detect(crdValue, apiValue, child)
		}

		r := Result{
			Property:    property,
			Equivalence: equivalenceRegistry.Get(funcName, fieldType.Kind())(crdValue, apiValue),
			CRDValue:    crdValue,
			APIValue:    apiValue,
		}
		parent.Children = append(parent.Children, r)
	}
}

func defaultEquivalence(crd any, api any) Equivalence {
	e := FromDeepEqual(crd, api)
	if e.Equivalent == Inequivalent {
		e.Reason = "default equal"
	}
	return e
}

func cleanValue(v reflect.Value) any {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return v.Elem().Interface()
	}
	return v.Interface()
}

func trimAfterComma(property string) string {
	return strings.Split(property, ",")[0]
}

func unTitle(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

func FromDeepEqual(crd any, api any) Equivalence {
	eq := reflect.DeepEqual(api, crd)
	if eq {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}
