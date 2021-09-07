// +build ignore

package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

const START_TAG = "// STRATEGY_START\n"
const END_TAG = "// STRATEGY_END\n"
const FILENAME = "./builtin_strategies.go"
const JS_SUFFIX = "_strategy.js"
const GO_SUFFIX = "_strategy.go"
const STRATEGY_DIR = "./strategies/generated/"

func main() {
	var files []string

	var content = `// Code generated, DO NOT EDIT.

package pools

`
	files, err := filepath.Glob(STRATEGY_DIR + "*" + GO_SUFFIX)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		strategyB, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		strategy := string(strategyB)
		var startPos = 0
		var endPos = 0
		if (strings.Contains(strategy, START_TAG)) {
			startPos = strings.Index(strategy, START_TAG)
		}
		if (strings.Contains(strategy, END_TAG)) {
			endPos = strings.Index(strategy, END_TAG)
		} else {
			endPos = len(strategy) - 1
		}
		strategy = strategy[startPos:endPos]

		content += strategy
	}

	files, err = filepath.Glob(STRATEGY_DIR + "*" + JS_SUFFIX)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		strategyB, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		content += "const "
		content += strings.ToUpper(strings.TrimSuffix(filepath.Base(file), JS_SUFFIX))
		content += " = `\n"

		strategy := string(strategyB)

		// remove trailing test code
		var startPos = 0
		var endPos = 0
		if (strings.Contains(strategy, END_TAG)) {
			endPos = strings.Index(strategy, END_TAG)
		} else {
			endPos = len(strategy) - 1
		}
		strategy = strategy[startPos:endPos]

		// escape backticks
		content += strings.ReplaceAll(strategy, "`", "` + \"`\" + `")

		content += "\n`\n\n"
	}

	ioutil.WriteFile(FILENAME, []byte(content), 0644)
}
