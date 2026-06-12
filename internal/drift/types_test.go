package drift_test

import (
	"reflect"
	"testing"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/zeebo/assert"
)

type SingleWithPtr struct {
	Value *string `drift:"empty-is-nil" json:"value,omitempty"`
}
type DoubleWithPtr struct {
	Name  *string `drift:"empty-is-nil" json:"name,omitempty"`
	Order *int    `drift:"empty-is-nil" json:"order,omitempty"`
}

type Nested struct {
	Value  *string `drift:"empty-is-nil" json:"value,omitempty"`
	Double DoubleWithPtr
}

// TODO do an official registry

// TODO cleanup

// TODO embedded


func EmptyIsNilString(crd any, api any) drift.Equivalence {
	if crd == nil && api != nil && api == "" {
		return drift.Equivalence{Equivalent: drift.Equivalent}
	}
	if api == nil && crd != nil && crd == "" {
		return drift.Equivalence{Equivalent: drift.Equivalent}
	}
	return drift.FromDeepEqual(crd, api)
}

func EmptyIsNilInt(crd any, api any) drift.Equivalence {
	if crd == nil && api != nil && api == 0 {
		return drift.Equivalence{Equivalent: drift.Equivalent}
	}
	if api == nil && crd != nil && crd == 0 {
		return drift.Equivalence{Equivalent: drift.Equivalent}
	}
	return drift.FromDeepEqual(crd, api)
}

func init() {
	drift.Register("empty-is-nil", reflect.String, EmptyIsNilString)
	drift.Register("empty-is-nil", reflect.Int, EmptyIsNilInt)
}

func TestNested(t *testing.T) {

	s := "foo"
	o1 := 1
	o2 := 2
	d1 := DoubleWithPtr{Name: &s, Order: &o1}
	d2 := DoubleWithPtr{Name: &s, Order: &o2}

	crd := Nested{Value: &s, Double: d1}
	api := Nested{Value: &s, Double: d2}


	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), true)
	// FIXME assert.Equal(t, r.String(), "value: foo <> bar")

}

func TestSinglePropertyPointer(t *testing.T) {

	valueCrd := "foo"
	valueApi := "bar"
	crd := SingleWithPtr{Value: &valueCrd}
	api := SingleWithPtr{Value: &valueApi}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), true)
	assert.Equal(t, r.String(), "value: foo <> bar")

}

func TestDoublePropertyPointer(t *testing.T) {

	s := "foo"
	o1 := 1
	o2 := 2
	crd := DoubleWithPtr{Name: &s, Order: &o1}
	api := DoubleWithPtr{Name: &s, Order: &o2}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), true)
	assert.Equal(t, r.String(), "order: 1 <> 2")

}

func TestDoublePropertyPointerAllInequivalent(t *testing.T) {

	s1 := "foo"
	s2 := "bar"
	o1 := 1
	o2 := 2
	crd := DoubleWithPtr{Name: &s1, Order: &o1}
	api := DoubleWithPtr{Name: &s2, Order: &o2}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), true)
	assert.Equal(t, r.String(),
		`name: foo <> bar
order: 1 <> 2`)

}

func TestDoublePropertyPointerEquivalent(t *testing.T) {

	s1 := "foo"
	s2 := "foo"
	o1 := 1
	o2 := 1
	crd := DoubleWithPtr{Name: &s1, Order: &o1}
	api := DoubleWithPtr{Name: &s2, Order: &o2}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), false)

}

// TODO complete cases
// TODO arrays

func TestSingleProperty(t *testing.T) {

	type TestNoPtr struct {
		Value string `drift:"empty-is-nil" json:"value,omitempty"`
	}

	crd := TestNoPtr{Value: "foo"}
	api := TestNoPtr{Value: "bar"}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), true)
	assert.Equal(t, r.String(), "value: foo <> bar")

}

func TestSinglePropertyNoEquivalence(t *testing.T) {

	type TestNoPtr struct {
		Value string `json:"value,omitempty"`
	}

	crd := TestNoPtr{Value: "foo"}
	api := TestNoPtr{Value: "bar"}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), true)
	assert.Equal(t, r.String(), "value: foo <> bar (default equal)")

}

func TestSinglePropertyNoEquivalenceNoJSON(t *testing.T) {

	type TestNoPtr struct {
		Value string
	}

	crd := TestNoPtr{Value: "foo"}
	api := TestNoPtr{Value: "bar"}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), true)
	assert.Equal(t, r.String(), "value: foo <> bar (default equal)")

}

func TestEmptyIsNilEquivalenceFunc(t *testing.T) {
	assert.Equal(t, EmptyIsNilString("", nil), drift.Equivalence{Equivalent: drift.Equivalent})
	assert.Equal(t, EmptyIsNilString(nil, ""), drift.Equivalence{Equivalent: drift.Equivalent})
	assert.Equal(t, EmptyIsNilString("", ""), drift.Equivalence{Equivalent: drift.Equivalent})
	assert.Equal(t, EmptyIsNilString(nil, nil), drift.Equivalence{Equivalent: drift.Equivalent})
}

func TestEmptyIsNilEquivalenceStringAgainstNil(t *testing.T) {

	type Test struct {
		Value *string `drift:"empty-is-nil"`
	}

	x := ""
	crd := Test{Value: &x}
	api := Test{Value: nil}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), false)
	assert.Equal(t, r.String(), "")

}

func TestEmptyIsNilEquivalenceNilAgainstString(t *testing.T) {

	type Test struct {
		Value *string `drift:"empty-is-nil"`
	}

	x := ""
	crd := Test{Value: nil}
	api := Test{Value: &x}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), false)
	assert.Equal(t, r.String(), "")

}

func TestEmptyIsNilEquivalenceNilAgainstNil(t *testing.T) {

	type Test struct {
		Value *string `drift:"empty-is-nil"`
	}

	crd := Test{Value: nil}
	api := Test{Value: nil}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), false)
	assert.Equal(t, r.String(), "")

}

func TestEmptyIsNilEquivalenceStringAgainstString(t *testing.T) {

	type Test struct {
		Value *string `drift:"empty-is-nil"`
	}

	s1 := ""
	s2 := ""
	crd := Test{Value: &s1}
	api := Test{Value: &s2}

	r := drift.Detect(crd, api)
	assert.Equal(t, r.DriftDetected(), false)
	assert.Equal(t, r.String(), "")

}
