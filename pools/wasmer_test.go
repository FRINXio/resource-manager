package pools

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/net-auto/resourceManager/graph/graphql/model"
)

func TestPrefixLines(t *testing.T) {
	actual := prefixLines("aa\nbb", "  ")
	expected := "  aa\n  bb\n"
	if actual != expected {
		t.Fatalf("Got \"%s\", expected \"%s\"", actual, expected)
	}
}

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
	log(JSON.stringify({respool: resourcePool.ResourcePoolName, currentRes: currentResources}));
	return {vlan: userInput.desiredVlan};
}
	`
	userInput := make(map[string]interface{})
	userInput["desiredVlan"] = 1
	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = "testpool"
	now := time.Now()
	currentResources := createCurrentResources(now)

	actual, logString, err := wasmer.invokeJs(script, userInput, resourcePool, currentResources)
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}
	checkResult(t, now, actual, logString)
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
log(json.dumps({ "respool": resourcePool["ResourcePoolName"], "currentRes": currentResources }))
return {"vlan": userInput["desiredVlan"]}
`
	userInput := make(map[string]interface{})
	userInput["desiredVlan"] = 1
	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = "testpool"
	now := time.Now()
	currentResources := createCurrentResources(now)

	actual, logString, err := wasmer.invokePy(script, userInput, resourcePool, currentResources)
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}
	checkResult(t, now, actual, logString)
}

func createCurrentResources(now time.Time) []*model.ResourceInput {
	var r0 model.ResourceInput
	r0.Properties = map[string]interface{}{"value": 1}
	r0.Status = "claimed"
	r0.UpdatedAt = now.String()

	var r1 model.ResourceInput
	r0.Properties = map[string]interface{}{"value": 100}
	r1.Status = "claimed"
	r1.UpdatedAt = now.String()

	var currentResources []*model.ResourceInput
	currentResources = append(currentResources, &r0, &r1)
	return currentResources
}

func checkResult(t *testing.T, now time.Time, actual map[string]interface{}, logString string) {
	// check stdout
	expected := make(map[string]interface{})
	json.Unmarshal([]byte(`{"vlan":1}`), &expected)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected evaluation response: %v, should be %v", actual, expected)
	}
	// check stderr
	actualLog := make(map[string]interface{})
	json.Unmarshal([]byte(logString), &actualLog)
	expectedLog := make(map[string]interface{})
	json.Unmarshal([]byte(`{
		"respool":"testpool",
		"currentRes":[
			{"Properties":[{"Name":"value","IntVal":1,"StringVal":null,"FloatVal":null,"Type":"int","Mandatory":true}],"UpdatedAt":"`+now.String()+`","Status":"claimed"},
			{"Properties":[{"Name":"value","IntVal":100,"StringVal":null,"FloatVal":null,"Type":"int","Mandatory":true}],"UpdatedAt":"`+now.String()+`","Status":"claimed"}]}
			`), &expectedLog)
	if !reflect.DeepEqual(actualLog, expectedLog) {
		t.Fatalf("Unexpected logging result: %v, should be %v", actualLog, expectedLog)
	}
}
