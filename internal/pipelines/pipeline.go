package pipeline

import (
	"github.com/prax860/tangent/internal/entities"
	"github.com/prax860/tangent/internal/intents"
	"github.com/prax860/tangent/internal/parser"
	"github.com/prax860/tangent/internal/rules"
	"github.com/prax860/tangent/internal/types"
	"github.com/prax860/tangent/internal/workspace"
)

func Process(input string) types.Response {

	normalized := parser.Parse(input)

	intent := intents.Resolve(normalized)

	arguments := entities.Extract(normalized, intent)

	workspaceType := workspace.Detect()

	command := rules.Generate(
		intent,
		workspaceType,
		arguments,
	)

	request := types.Request{
		RawInput:   input,
		Normalized: normalized,
		Intent:     intent,
		Workspace:  workspaceType,
		Arguments:  arguments,
	}

	return types.Response{
		Request: request,
		Command: command,
	}
}
