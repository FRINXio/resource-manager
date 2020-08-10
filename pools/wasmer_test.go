package pools

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestWasmerInvokeJsIntegration(t *testing.T) {
	wasmerPath, found := os.LookupEnv("WASMER_BIN")
	if !found {
		t.Fatalf("environment variable WASMER_BIN not found")
	}
	jsPath, found := os.LookupEnv("WASMER_JS")
	if !found {
		t.Fatalf("environment variable WASMER_JS not found")
	}
	wasmer := NewWasmer(wasmerPath, jsPath)
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
