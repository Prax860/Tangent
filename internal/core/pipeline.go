package core

import (
	"github.com/prax860/tangent/internal/rules"
	"github.com/prax860/tangent/internal/types"
)

func Process(input string) types.Response {

	normalized := Parse(input)

	intent := Resolve(normalized)

	arguments := Extract(normalized, intent)

	workspaceType := Detect()

	if pkg, ok := arguments["package"]; ok {

		resolved := ResolveForWorkspace(workspaceType, pkg)

		arguments["package"] = resolved.Target
	}

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