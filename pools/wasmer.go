package pools

import (
	"encoding/json"
	"os/exec"

	pools "github.com/net-auto/resourceManager/pools/allocating_strategies"
	"github.com/pkg/errors"
)

type Wasmer struct {
	wasmerBin string
	jsBin     string

	strategyInvokerScript string
}

const wasmerBinDefault = "./.wasmer/bin/wasmer"
const jsBinDefault = "./wasm/quickjs/quickjs.wasm"

func NewWasmerDefault() Wasmer {
	return NewWasmer(wasmerBinDefault, jsBinDefault)
}

func NewWasmer(wasmerBin string, jsBin string) Wasmer {
	return Wasmer{wasmerBin, jsBin, pools.STRATEGY_INVOKER}
}

type ScriptInvoker interface {
	invokeJs(strategyScript string) (map[string]interface{}, error)
}

func (wasmer Wasmer) invokeJs(strategyScript string) (map[string]interface{}, error) {
	// TODO pass additional inputs

	// Append script to invoke the function, parse inputs and serialize outputs
	scriptWithInvoker := strategyScript + wasmer.strategyInvokerScript
	command := exec.Command(wasmer.wasmerBin, wasmer.jsBin, "--", "--std", "-e", scriptWithInvoker)
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

// TODO invokePy
