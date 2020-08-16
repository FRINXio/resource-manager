package pools

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/net-auto/resourceManager/graph/graphql/model"
)

func TestWasmerInvokeJsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	wasmer, err := NewWasmerUsingEnvVars()
	if err != nil {
		t.Fatalf("Unable create wasmer - %s", err)
	}
	script := `
function invoke() {
	log('log');
	return {result: userInput.name, respool: resourcePool.ResourcePoolName, currentRes: currentResources};
}
`
	userInput := make(map[string]interface{})
	userInput["name"] = "Alice"
	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = "testpool"

	var r0 model.ResourceInput
	r0.Properties = map[string]interface{}{"value": 1}
	var r1 model.ResourceInput
	r1.Properties = map[string]interface{}{"value": 100}
	var currentResources []*model.ResourceInput
	currentResources = append(currentResources, &r0, &r1)

	actual, stderr, err := wasmer.invokeJs(script, userInput, resourcePool, currentResources)
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}
	expectedJSON := `{"result": "Alice", "respool": "testpool",
	 "currentRes": [ {"Properties": {"value": 1}}, {"Properties": {"value": 100}}]}`
	expected := make(map[string]interface{})
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("Cannot unmarshal expected response - %s", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected evaluation response: %v, should be %v", actual, expected)
	}
	if stderr != "log\n" {
		t.Fatalf("Unexpected logging result: \"%s\"", stderr)
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
	script := `
log('log')
return {"result": userInput["name"], "respool": resourcePool["ResourcePoolName"], "currentRes": currentResources}
`
	userInput := make(map[string]interface{})
	userInput["name"] = "Alice"
	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = "testpool"

	var r0 model.ResourceInput
	r0.Properties = map[string]interface{}{"value": 1}
	var r1 model.ResourceInput
	r1.Properties = map[string]interface{}{"value": 100}
	var currentResources []*model.ResourceInput
	currentResources = append(currentResources, &r0, &r1)

	actual, stderr, err := wasmer.invokePy(script, userInput, resourcePool, currentResources)
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}
	expectedJSON := `{"result": "Alice", "respool": "testpool",
	 "currentRes": [ {"Properties": {"value": 1}}, {"Properties": {"value": 100}}]}`
	expected := make(map[string]interface{})
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("Cannot unmarshal expected response - %s", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected evaluation response: %v, should be %v", actual, expected)
	}
	if stderr != "log\n" {
		t.Fatalf("Unexpected logging result: \"%s\"", stderr)
	}
}
