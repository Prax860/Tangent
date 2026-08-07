package rules

import "github.com/prax860/tangent/internal/types"

type GoInstallPackageRule struct{}

func (r GoInstallPackageRule) Name() string {
	return "go.install_package"
}

func (r GoInstallPackageRule) Match(intent types.IntentType, workspace types.WorkspaceType) bool {
	return intent == types.IntentInstallPackage &&
		workspace == types.WorkspaceGo
}

func (r GoInstallPackageRule) Generate() types.Command {
	return types.Command{
		Command:     "go get <module>",
		Explanation: "Installs a Go module.",
		Safe:        true,
	}
}

func init() {
	Register(GoInstallPackageRule{})
}