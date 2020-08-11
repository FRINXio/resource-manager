package pools

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestWasmerInvokeJsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	wasmer, err := NewWasmerUsingEnvVars()
	if err != nil {
		t.Fatalf("Unable create wasmer - %s", err)
	}
	actual, err := wasmer.invokeJs("const invoke = function() {return {result: 1};}")
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}
	expected := make(map[string]interface{})
	if err := json.Unmarshal([]byte("{\"result\":1}"), &expected); err != nil {
		t.Fatalf("Cannot unmarshal expected response - %s", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected evaluation response: %v, should be %v", actual, expected)
	}
}

func TestPrefixLines(t *testing.T) {
	actual := prefixLines("aa\nbb", "  ")
	expected := "  aa\n  bb\n"
	if actual != expected {
		t.Fatalf("Got \"%s\", expected \"%s\"", actual, expected)
	}
}

func TestWasmerInvokePyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	wasmer, err := NewWasmerUsingEnvVars()
	if err != nil {
		t.Fatalf("Unable create wasmer - %s", err)
	}
	actual, err := wasmer.invokePy("return {'result': 1}")
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}
	expected := make(map[string]interface{})
	if err := json.Unmarshal([]byte("{\"result\":1}"), &expected); err != nil {
		t.Fatalf("Cannot unmarshal expected response - %s", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected evaluation response: %v, should be %v", actual, expected)
	}
}
