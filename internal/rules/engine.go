package rules

import "github.com/prax860/tangent/internal/types"

func Generate(
	intent types.IntentType,
	workspace types.WorkspaceType,
	arguments map[string]string,
) types.Command {

	for _, rule := range Rules() {

		if rule.Match(intent, workspace) {
			return rule.Generate(arguments)
		}

	}

	return types.Command{
		Command:     "",
		Explanation: "No matching rule found.",
		Safe:        false,
	}
}
