package pools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	invokeJs(strategyScript string) (map[string]interface{}, string, error)
	invokePy(strategyScript string) (map[string]interface{}, string, error)
}

// TODO implement killing the wasmer process after max timeout
func (wasmer Wasmer) invokeJs(strategyScript string) (map[string]interface{}, string, error) {
	// TODO pass additional inputs

	// Append script to invoke the function, parse inputs and serialize outputs
	header := `
console.error = function(...args) {
	std.err.puts(args.join(' '));
	std.err.puts('\n');
}
log = console.error;
`
	footer := `
let result = invoke()
if (result != null) {
	if (typeof result === 'object') {
		result = JSON.stringify(result);
	}
	std.out.puts(result);
}
`
	scriptWithInvoker := header + strategyScript + footer
	command := exec.Command(wasmer.wasmerBinPath, wasmer.jsBinPath, "--", "--std", "-e", scriptWithInvoker)
	return invoke(command)
}

func invoke(command *exec.Cmd) (map[string]interface{}, string, error) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	command.Stdout = &stdoutBuffer
	command.Stderr = &stderrBuffer

	err := command.Run()
	stdout := stdoutBuffer.Bytes()
	stderr := string(stderrBuffer.Bytes())

	// TODO logging
	fmt.Println("Stdout:" + string(stdout[:]))
	fmt.Println("Stderr:" + stderr[:])

	if err != nil {
		return nil, stderr, errors.Wrapf(err,
			"Error invoking allocation strategy. Stdout: \"%s\", Stderr: \"%s\"", string(stdout), stderr)
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(stdout, &m); err != nil {
		return nil, stderr, errors.Wrapf(err,
			"Unable to parse allocation function output as flat JSON: \"%s\"", string(stdout))
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

func (wasmer Wasmer) invokePy(script string) (map[string]interface{}, string, error) {
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
