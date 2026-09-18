package validation

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestPointerRequiredAndNestedPath(t *testing.T) {
	type Child struct {
		Name string `json:"name" validate:"required"`
	}
	type Input struct {
		ID       *int    `json:"id" validate:"required"`
		Children []Child `json:"children"`
	}
	input := Input{Children: []Child{{}}}
	result, err := New().Validate(input)
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if len(result) != 2 || result[0].Field != "id" || result[1].Field != "children[0].name" {
		t.Fatalf("unexpected errors: %#v", result)
	}
}

func TestEqFieldSupportsNumbers(t *testing.T) {
	type Input struct {
		A int `json:"a"`
		B int `json:"b" validate:"eqfield=A"`
	}
	result, err := New().Validate(Input{A: 1, B: 2})
	if !errors.Is(err, ErrValidation) || len(result) != 1 || result[0].Field != "b" {
		t.Fatalf("unexpected result: %#v, %v", result, err)
	}
	if _, err := New().Validate(Input{A: 1, B: 1}); err != nil {
		t.Fatalf("equal values should pass: %v", err)
	}
}

func TestMalformedRulesReturnError(t *testing.T) {
	type Input struct {
		Name string `validate:"min"`
	}
	_, err := New().Validate(Input{})
	var ruleErr *InvalidRuleError
	if !errors.As(err, &ruleErr) {
		t.Fatalf("expected InvalidRuleError, got %v", err)
	}
}

func TestLanguageIsolatedPerValidator(t *testing.T) {
	one, two := New(), New()
	one.SetLanguage(map[string]string{"required": "one"})
	type Input struct {
		Name string `validate:"required"`
	}
	result, err := two.Validate(Input{})
	if !errors.Is(err, ErrValidation) || result[0].Message[0] != Lang["required"] {
		t.Fatalf("language leaked between instances: %#v", result)
	}
}

func TestIPValidatorsRejectInvalidValues(t *testing.T) {
	type Input struct {
		V4 string `json:"v4" validate:"ipv4"`
		V6 string `json:"v6" validate:"ipv6"`
	}
	result, err := New().Validate(Input{V4: "999.1.1.1", V6: "2001:db8::1"})
	if !errors.Is(err, ErrValidation) || len(result) != 1 || result[0].Field != "v4" {
		t.Fatalf("unexpected result: %#v, %v", result, err)
	}
}

func TestFileValidateMalformedAndNil(t *testing.T) {
	if _, err := New().FileValidate(nil, "maxsize=1"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
	file, err := os.CreateTemp("", "validation-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err := New().FileValidate(file, "maxsize"); err == nil {
		t.Fatal("expected malformed file rule error")
	}
}

func FuzzValidateNeverPanics(f *testing.F) {
	f.Add("required,email")
	f.Add("min=2,max=20")
	f.Fuzz(func(t *testing.T, rules string) {
		type Input struct {
			Value string `validate:"required"`
		}
		inputType := reflect.StructOf([]reflect.StructField{{Name: "Value", Type: reflect.TypeOf(""), Tag: reflect.StructTag(`validate:"` + rules + `"`)}})
		input := reflect.New(inputType).Elem()
		input.Field(0).SetString("x")
		_, _ = New().Validate(input.Interface())
	})
}
