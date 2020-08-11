package pools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	pools "github.com/net-auto/resourceManager/pools/allocating_strategies"
	"github.com/pkg/errors"
)

type Wasmer struct {
	wasmerBinPath string
	jsBinPath     string
	pythonBinPath string
	pythonLibPath string
}

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
		return "", errors.Errorf("Environment variable \"%s\" not found", key)
	}
	return value, nil
}

func NewWasmerUsingEnvVars() (*Wasmer, error) {
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

	wasmer := NewWasmer(wasmerPath, jsPath, pyPath, pyLibPath)
	return &wasmer, nil
}

func NewWasmer(wasmerBinPath string, jsBinPath string, pythonBinPath string, pythonLibPath string) Wasmer {
	return Wasmer{wasmerBinPath, jsBinPath, pythonBinPath, pythonLibPath}
}

type ScriptInvoker interface {
	invokeJs(strategyScript string) (map[string]interface{}, error)
	invokePy(strategyScript string) (map[string]interface{}, error)
}

// TODO implement killing the wasmer process after max timeout
func (wasmer Wasmer) invokeJs(strategyScript string) (map[string]interface{}, error) {
	// TODO pass additional inputs

	// Append script to invoke the function, parse inputs and serialize outputs
	scriptWithInvoker := strategyScript + pools.STRATEGY_INVOKER
	command := exec.Command(wasmer.wasmerBinPath, wasmer.jsBinPath, "--", "--std", "-e", scriptWithInvoker)
	return invoke(command)
}

func invoke(command *exec.Cmd) (map[string]interface{}, error) {
	output, err := command.CombinedOutput()
	// TODO Add logging string(output[:])
	if err != nil {
		return nil, errors.Wrapf(err,
			"Error invoking allocation strategy. Output: \"%s\"", output)
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(output, &m); err != nil {
		return nil, errors.Wrapf(err,
			"Unable to parse allocation function output as flat JSON: \"%s\"", output)
	}
	return m, nil
}

func prefixLines(str string, prefix string) string {
	lines := strings.Split(str, "\n")
	result := ""
	for _, line := range lines {
		result += prefix + line + "\n"
	}
	return result
}

func (wasmer Wasmer) invokePy(script string) (map[string]interface{}, error) {
	head := `
import sys,json
def log(*args, **kwargs):
  print(*args, file=sys.stderr, **kwargs)

def script_fun():
`
	foot := `
result = script_fun()
if not result is None:
  if isinstance(result, str):
    sys.stdout.write(result)
  else:
    sys.stdout.write(json.dumps(result))
`
	script = head + prefixLines(script, "  ") + foot

	fmt.Println("Executing\n" + script)

	// options:
	// -q: quiet, do not print python version
	// -B: do not write .pyc files on import
	// -c script: execute passed script
	command := exec.Command(wasmer.wasmerBinPath,
		wasmer.pythonBinPath,
		"--mapdir=lib:"+wasmer.pythonLibPath,
		"--",
		"-B",
		"-q",
		"-c",
		script,
	)
	return invoke(command)
}
