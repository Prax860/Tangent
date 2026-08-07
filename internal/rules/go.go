package rules

import "github.com/prax860/tangent/internal/types"

type GoInstallPackageRule struct{}

func (r GoInstallPackageRule) Name() string {
	return "go.install_package"
}

func (r GoInstallPackageRule) Match(
	intent types.IntentType,
	workspace types.WorkspaceType,
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
