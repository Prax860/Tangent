package rules

import "github.com/prax860/tangent/internal/types"

type GoInstallPackageRule struct{}

func (r GoInstallPackageRule) Name() string {
	return "go.install_package"
}

func (r GoInstallPackageRule) Match(
	intent types.IntentType,
	workspace types.WorkspaceType,
	_ map[string]string,
) bool {

	return intent == types.IntentInstallPackage &&
		workspace == types.WorkspaceGo
}

func (r GoInstallPackageRule) Generate(
	arguments map[string]string,
) types.Command {

	module := arguments["package"]

	return types.Command{
		Command:     "go get " + module,
		Explanation: "Installs Go module '" + module + "'",
		Safe:        true,
	}
}

func init() {
	Register(GoInstallPackageRule{})
}

// ------------------------------------------

type GoScaffoldInitRule struct{}

func (r GoScaffoldInitRule) Name() string {
	return "go.scaffold_init"
}

func (r GoScaffoldInitRule) Match(
	intent types.IntentType,
	workspace types.WorkspaceType,
	arguments map[string]string,
) bool {

	if intent != types.IntentCreateProject {
		return false
	}

	fw := arguments["framework"]
	if workspace == types.WorkspaceGo {
		return true
	}
	if fw == "go" {
		return true
	}
	return false
}

func (r GoScaffoldInitRule) Generate(
	arguments map[string]string,
) types.Command {

	return types.Command{
		Command:     "go mod init",
		Explanation: "Initializes a new Go module in the current directory.",
		Safe:        true,
		Interactive: false,
	}
}

func init() {
	Register(GoScaffoldInitRule{})
}
