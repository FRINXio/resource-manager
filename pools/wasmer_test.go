package pools

import (
	"encoding/json"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	"reflect"
	"testing"
	"time"
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

	actual, logString, err := wasmer.invokeJs(script, userInput, resourcePool, currentResources, nil, "invoke()")
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

	actual, logString, err := wasmer.invokePy(script, userInput, resourcePool, currentResources, nil, "script_fun()")
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}
	checkResult(t, now, actual, logString)
}

func TestVlanInvokeGo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := getContext()
	client := openDb(ctx)
	strat, _ := client.AllocationStrategy.Create().
		SetName("vlan").
		SetLang(allocationstrategy.LangGo).
		SetScript("Hello World!").
		Save(ctx)

	userInput := make(map[string]interface{})
	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = "svlan"
	resourcePool.ResourcePoolID = 1
	now := time.Now()
	var r0 model.ResourceInput
	r0.Properties = map[string]interface{}{"vlan": 0}
	r0.Status = "claimed"
	r0.UpdatedAt = now.String()

	var currentResources []*model.ResourceInput
	currentResources = append(currentResources, &r0)

	var resourcePoolResources = map[string]interface{}{"from": 0, "to": 4095}

	actual, logString, err := invokeGo(ctx, strat, userInput, resourcePool, currentResources, resourcePoolResources, "invoke()")
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}

	expected := make(map[string]interface{})
	expected["vlan"] = float64(1)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected evaluation response: %v, should be %v", actual, expected)
	}
	// check stderr
	expectedLogString := ""
	if !reflect.DeepEqual(logString, expectedLogString) {
		t.Fatalf("Unexpected logging result: %v, should be %v", logString, expectedLogString)
	}
}

func TestUniqueIdInvokeGo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := getContext()
	client := openDb(ctx)
	strat, _ := client.AllocationStrategy.Create().
		SetName("unique_id").
		SetLang(allocationstrategy.LangGo).
		SetScript("Hello World!").
		Save(ctx)

	userInput := make(map[string]interface{})
	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = "svlan"
	resourcePool.ResourcePoolID = 1

	now := time.Now()
	var r0 model.ResourceInput
	r0.Properties = map[string]interface{}{"text": "1000", "counter": 1000}
	r0.Status = "claimed"
	r0.UpdatedAt = now.String()

	var currentResources []*model.ResourceInput
	currentResources = append(currentResources, &r0)

	var resourcePoolResources = map[string]interface{}{"from": 1000, "idFormat": "{counter}"}

	actual, logString, err := invokeGo(ctx, strat, userInput, resourcePool, currentResources, resourcePoolResources, "invoke()")
	if err != nil {
		t.Fatalf("Unable run - %s", err)
	}

	expected := make(map[string]interface{})
	expected["text"] = "1001"
	expected["counter"] = float64(1001)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected evaluation response: %v, should be %v", actual, expected)
	}
	// check stderr
	expectedLogString := ""
	if !reflect.DeepEqual(logString, expectedLogString) {
		t.Fatalf("Unexpected logging result: %v, should be %v", logString, expectedLogString)
	}
}

func createCurrentResources(now time.Time) []*model.ResourceInput {
	var r0 model.ResourceInput
	r0.Properties = map[string]interface{}{"value": 1}
	r0.Status = "claimed"
	r0.UpdatedAt = now.String()

	var r1 model.ResourceInput
	r1.Properties = map[string]interface{}{"value": 100}
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
	expectedLogString := `{
		"respool":"testpool",
		"currentRes":[
			{"Properties":{"value":1},"UpdatedAt":"` + now.String() + `","Status":"claimed"},
			{"Properties":{"value":100},"UpdatedAt":"` + now.String() + `","Status":"claimed"}]}
			`
	expectedLog := make(map[string]interface{})
	json.Unmarshal([]byte(expectedLogString), &expectedLog)
	if !reflect.DeepEqual(actualLog, expectedLog) {
		t.Fatalf("Unexpected logging result: %v, should be %v", logString, expectedLogString)
	}
}
