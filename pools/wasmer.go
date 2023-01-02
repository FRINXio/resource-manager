package pools

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	strategies "github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/generated"

	log "github.com/net-auto/resourceManager/logging"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	"github.com/pkg/errors"
)

type Wasmer struct {
	maxTimeout    time.Duration
	wasmerBinPath string
	jsBinPath     string
	pythonBinPath string
	pythonLibPath string
}

const wasmerMaxTimeoutMillisDefault = "5000"
const wasmerBinDefault = "./.wasmer/bin/wasmer"
const jsBinDefault = "./wasm/quickjs/quickjs.wasm"
const pyBinDefault = "./wasm/python/bin/python.wasm"
const pyLibDefault = "wasm/python/lib"

func loadEnvVar(key string, defaultValue string) (string, error) {
	value, found := os.LookupEnv(key)
	if !found {
		if defaultValue != "" {
			return defaultValue, nil
		}

		err := errors.Errorf("Environment variable \"%s\" not found", key)
		log.Error(nil, err, "Environment variable missing")
		return "", err
	}
	return value, nil
}

func NewWasmerUsingEnvVars() (*Wasmer, error) {

	maxTimeoutMillisStr, err := loadEnvVar("WASMER_MAX_TIMEOUT_MILLIS", wasmerMaxTimeoutMillisDefault)
	if err != nil {
		return nil, err
	}
	maxTimeoutMillis, err := strconv.Atoi(maxTimeoutMillisStr)
	if err != nil {
		err := errors.Wrapf(err, "Cannot convert \"%s\" to int", maxTimeoutMillisStr)
		log.Error(nil, err, "Conversion error")
		return nil, err
	}
	maxTimeout := time.Duration(maxTimeoutMillis) * time.Millisecond

	wasmerPath, err := loadEnvVar("WASMER_BIN", wasmerBinDefault)
	if err != nil {
		return nil, err
	}
	jsPath, err := loadEnvVar("WASMER_JS", jsBinDefault)
	if err != nil {
		return nil, err
	}
	pyPath, err := loadEnvVar("WASMER_PY", pyBinDefault)
	if err != nil {
		return nil, err
	}
	pyLibPath, err := loadEnvVar("WASMER_PY_LIB", pyLibDefault)
	if err != nil {
		return nil, err
	}

	wasmer := NewWasmer(maxTimeout, wasmerPath, jsPath, pyPath, pyLibPath)
	return &wasmer, nil
}

func NewWasmer(maxTimeout time.Duration, wasmerBinPath string, jsBinPath string, pythonBinPath string, pythonLibPath string) Wasmer {
	return Wasmer{maxTimeout, wasmerBinPath, jsBinPath, pythonBinPath, pythonLibPath}
}

type ScriptInvoker interface {
	invokeJs(strategyScript string, userInput map[string]interface{},
		resourcePool model.ResourcePoolInput,
		currentResources []*model.ResourceInput,
		poolPropertiesMaps map[string]interface{},
		functionName string,
	) (map[string]interface{}, string, error)
	invokePy(strategyScript string, userInput map[string]interface{},
		resourcePool model.ResourcePoolInput,
		currentResources []*model.ResourceInput,
		poolPropertiesMaps map[string]interface{},
		functionName string,
	) (map[string]interface{}, string, error)
}

func InvokeAllocationStrategy(
	ctx context.Context,
	invoker ScriptInvoker,
	strat *ent.AllocationStrategy,
	userInput map[string]interface{},
	resourcePool model.ResourcePoolInput,
	currentResources []*model.ResourceInput,
	poolPropertiesMaps map[string]interface{},
	functionName string,
) (map[string]interface{}, string, error) {

	switch strat.Lang {
	case allocationstrategy.LangJs:
		return invoker.invokeJs(strat.Script, userInput, resourcePool, currentResources, poolPropertiesMaps, functionName)
	case allocationstrategy.LangPy:
		return invoker.invokePy(strat.Script, userInput, resourcePool, currentResources, poolPropertiesMaps, functionName)
	case allocationstrategy.LangGo:
		return invokeGo(ctx, strat, userInput, resourcePool, currentResources, poolPropertiesMaps, functionName)
	default:
		err := errors.Errorf("Unknown language \"%s\" for strategy \"%s\"", strat.Lang, strat.Name)
		log.Error(nil, err, "Unknown strategy language")
		return nil, "", err
	}
}

func currentResourcesToArray(currentResources []*model.ResourceInput) ([]map[string]interface{}, error) {
	var mapInterface []map[string]interface{}
	for _, element := range currentResources {
		m := make(map[string]interface{})
		inrec, err := json.Marshal(element)
		if err != nil {
			err := errors.Wrap(err, "Cannot serialize currentResources into json")
			log.Error(nil, err, "Cannot serialize currentResources")
			return mapInterface, err
		}
		if err := json.Unmarshal(inrec, &m); err != nil {
			log.Error(nil, err, "Parsing error")
			return mapInterface, errors.Wrapf(err, "Parsing error")
		}

		mapInterface = append(mapInterface, m)
	}
	return mapInterface, nil
}

type GoStrategy interface {
	Invoke() (map[string]interface{}, error)
	Capacity() (map[string]interface{}, error)
}

func runStrategy(goStrategy GoStrategy, functionName string) (map[string]interface{}, error) {
	switch functionName {
	case "invoke()":
		return goStrategy.Invoke()
	case "capacity()":
		return goStrategy.Capacity()
	default:
		return nil, errors.New("Not known function")
	}
}

func invokeGo(
	ctx context.Context,
	strategy *ent.AllocationStrategy,
	userInput map[string]interface{},
	resourcePool model.ResourcePoolInput,
	currentResources []*model.ResourceInput,
	poolPropertiesMaps map[string]interface{},
	functionName string,
) (map[string]interface{}, string, error) {
	var output map[string]interface{}
	var goStrategy GoStrategy

	currentResourcesArray, err := currentResourcesToArray(currentResources)
	if err != nil {
		return nil, "", err
	}
	log.Debug(nil, "CurrentResources:\n %s\n poolProperties:\n %s", currentResourcesArray, poolPropertiesMaps)

	switch strategy.Name {
	case "vlan":
		// TODO: Pass currentResourcesArray as pointer
		vlan := strategies.NewVlan(currentResourcesArray, poolPropertiesMaps, userInput)
		goStrategy = &vlan
	case "unique_id":
		id := strategies.NewUniqueId(ctx, resourcePool.ResourcePoolID, poolPropertiesMaps, userInput)
		goStrategy = &id
	case "ipv4":
		// TODO: Pass currentResourcesArray as pointer
		id := strategies.NewIpv4(currentResourcesArray, poolPropertiesMaps, userInput)
		goStrategy = &id
	case "ipv6":
		// TODO: Pass currentResourcesArray as pointer
		id := strategies.NewIpv6(currentResourcesArray, poolPropertiesMaps, userInput)
		goStrategy = &id
	case "ipv6_prefix":
		// TODO: Pass currentResourcesArray as pointer
		id := strategies.NewIpv6Prefix(currentResourcesArray, poolPropertiesMaps, userInput)
		goStrategy = &id
	case "ipv4_prefix":
		// TODO: Pass currentResourcesArray as pointer
		id := strategies.NewIpv4Prefix(currentResourcesArray, poolPropertiesMaps, userInput)
		goStrategy = &id
	default:
		return nil, "", errors.New("Not known go strategy")
	}
	output, err = runStrategy(goStrategy, functionName)

	if err != nil {
		return nil, "", err
	}
	log.Debug(nil, "Stdout: %s", output)
	return output, "", nil
}

func serializeJsVariable(name string, data interface{}) (string, error) {
	userInputBytes, err := json.Marshal(data)
	if err != nil {
		err := errors.Wrap(err, "Cannot serialize userInput into json")
		log.Error(nil, err, "Cannot serialize userInput")
		return "", err
	}
	return "const " + name + " = " + string(userInputBytes[:]) + ";\n", nil
}

func (wasmer Wasmer) invokeJs(
	strategyScript string,
	userInput map[string]interface{},
	resourcePool model.ResourcePoolInput,
	currentResources []*model.ResourceInput,
	poolPropertiesMaps map[string]interface{},
	functionName string,
) (map[string]interface{}, string, error) {

	// Append script to invoke the function, parse inputs and serialize outputs
	header := `
console.error = function(...args) {
	std.err.puts(args.join(' '));
	std.err.puts('\n');
}
const log = console.error;
`
	if userInput == nil {
		// default in case of nil
		userInput = map[string]interface{}{}
	}
	addition, err := serializeJsVariable("userInput", userInput)
	if err != nil {
		return nil, "", err
	}
	header += addition

	addition, err = serializeJsVariable("resourcePoolProperties", poolPropertiesMaps)
	if err != nil {
		return nil, "", err
	}
	header += addition

	addition, err = serializeJsVariable("resourcePool", resourcePool)
	if err != nil {
		return nil, "", err
	}
	header += addition

	if currentResources == nil {
		// default in case of nil
		currentResources = []*model.ResourceInput{}
	}
	addition, err = serializeJsVariable("currentResources", currentResources)
	if err != nil {
		return nil, "", err
	}
	header += addition

	footer := `
let result = ` + functionName + `
if (result != null) {
	if (typeof result === 'object') {
		result = JSON.stringify(result);
	}
	std.out.puts(result);
}
`

	scriptWithInvoker := header + strategyScript + footer
	log.Debug(nil, "Executing:\n %s", scriptWithInvoker)
	return wasmer.invoke(wasmer.wasmerBinPath, wasmer.jsBinPath, "--", "--std", "-e", scriptWithInvoker)
}

func (wasmer Wasmer) invoke(name string, arg ...string) (map[string]interface{}, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), wasmer.maxTimeout)
	defer cancel()

	command := exec.CommandContext(ctx, name, arg...)
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	command.Stdout = &stdoutBuffer
	command.Stderr = &stderrBuffer

	err := command.Run()
	stdout := stdoutBuffer.Bytes()
	stderr := string(stderrBuffer.Bytes())

	if err != nil {
		err := errors.Wrapf(err,
			"Error invoking user script. Stdout: \"%s\", Stderr: \"%s\"", string(stdout), stderr)
		log.Error(nil, err, "Error invoking user script")
		return nil, "", err
	}

	log.Debug(nil, "Stdout: %s", string(stdout[:]))
	log.Debug(nil, "Stderr: %s", stderr[:])

	m := make(map[string]interface{})
	if err := json.Unmarshal(stdout, &m); err != nil {
		log.Error(nil, err, "Parsing error")
		return nil, stderr, errors.Wrapf(err,
			"Unable to parse allocation function output as flat JSON: \"%s\". "+
				"Error output: \"%s\"", string(stdout), stderr)
	}
	return m, stderr, nil
}

func prefixLines(str string, prefix string) string {
	lines := strings.Split(str, "\n")
	result := ""
	for _, line := range lines {
		result += prefix + line + "\n"
	}
	return result
}

func serializePythonVariable(name string, data interface{}) (string, error) {
	userInputBytes, err := json.Marshal(data)
	if err != nil {
		log.Error(nil, err, "Cannot serialize")
		return "", errors.Wrapf(err, "Cannot serialize %s into json", name)
	}
	// Encode twice so that the result starts and ends with a quote. All inner quotes will be escaped.
	userInputBytes, err = json.Marshal(string(userInputBytes[:]))
	if err != nil {
		log.Error(nil, err, "Cannot serialize")
		return "", errors.Wrapf(err, "Cannot serialize %s into json second time", name)
	}
	return name + "=json.loads(" + string(userInputBytes[:]) + ")\n", nil
}

func (wasmer Wasmer) invokePy(
	script string,
	userInput map[string]interface{},
	resourcePool model.ResourcePoolInput,
	currentResources []*model.ResourceInput,
	poolPropertiesMaps map[string]interface{},
	functionName string,
) (map[string]interface{}, string, error) {
	header := `
import sys,json
def log(*args, **kwargs):
  print(*args, file=sys.stderr, **kwargs)
`
	if userInput == nil {
		// default in case of nil
		userInput = map[string]interface{}{}
	}
	addition, err := serializePythonVariable("userInput", userInput)
	if err != nil {
		return nil, "", err
	}
	header += addition

	addition, err = serializePythonVariable("resourcePool", resourcePool)
	if err != nil {
		return nil, "", err
	}
	header += addition

	addition, err = serializePythonVariable("resourcePoolProperties", poolPropertiesMaps)
	if err != nil {
		return nil, "", err
	}
	header += addition

	if currentResources == nil {
		// default in case of nil
		currentResources = []*model.ResourceInput{}
	}
	addition, err = serializePythonVariable("currentResources", currentResources)
	if err != nil {
		return nil, "", err
	}
	header += addition

	header += `
def ` + functionName + `:
`
	footer := `
result = ` + functionName + `
if not result is None:
  if isinstance(result, str):
    sys.stdout.write(result)
  else:
    sys.stdout.write(json.dumps(result))
`
	script = header + prefixLines(script, "  ") + footer

	log.Debug(nil, "Executing:\n %s ", script)

	// options:
	// -q: quiet, do not print python version
	// -B: do not write .pyc files on import
	// -c script: execute passed script
	return wasmer.invoke(wasmer.wasmerBinPath,
		wasmer.pythonBinPath,
		"--mapdir=lib:"+wasmer.pythonLibPath,
		"--",
		"-B",
		"-q",
		"-c",
		script,
	)
}