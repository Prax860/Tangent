package rules

import "github.com/prax860/tangent/internal/types"

type RustScaffoldInitRule struct{}

func (r RustScaffoldInitRule) Name() string {
	return "rust.scaffold_init"
}

func (r RustScaffoldInitRule) Match(
	intent types.IntentType,
	workspace types.WorkspaceType,
	arguments map[string]string,
) bool {

	if intent != types.IntentCreateProject {
		return false
	}

	fw := arguments["framework"]
	if workspace == types.WorkspaceRust {
		return true
	}
	if fw == "rust" {
		return true
	}
	return false
}

func (r RustScaffoldInitRule) Generate(
	arguments map[string]string,
) types.Command {

	return types.Command{
		Command:     "cargo init",
		Explanation: "Initializes a new Cargo (Rust) project in the current directory. Cargo will prompt for options interactively.",
		Safe:        true,
		Interactive: true,
	}
}

func init() {
	Register(RustScaffoldInitRule{})
}
