package rules

import "github.com/prax860/tangent/internal/types"

type CreateVenvRule struct{}

func (r CreateVenvRule) Name() string {
	return "python.create_venv"
}

func (r CreateVenvRule) Match(
	intent types.IntentType,
	workspace types.WorkspaceType,
) bool {

	return intent == types.IntentCreateVenv &&
		workspace == types.WorkspacePython
}

func (r CreateVenvRule) Generate(
	arguments map[string]string,
) types.Command {

	venv := arguments["venv"]

	return types.Command{
		Command:     "python -m venv " + venv,
		Explanation: "Creates Python virtual environment '" + venv + "'",
		Safe:        true,
	}
}

func init() {
	Register(CreateVenvRule{})
}
