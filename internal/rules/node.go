package rules

import "github.com/prax860/tangent/internal/types"

type NodeInstallPackageRule struct{}

func (r NodeInstallPackageRule) Name() string {
	return "node.install_package"
}

func (r NodeInstallPackageRule) Match(intent types.IntentType, workspace types.WorkspaceType) bool {
	return intent == types.IntentInstallPackage &&
		workspace == types.WorkspaceNode
}

func (r NodeInstallPackageRule) Generate(
	arguments map[string]string,
) types.Command {
	packageName := arguments["package"]

	return types.Command{
		Command:     "npm install " + packageName,
		Explanation: "Installs Node.js package '" + packageName + "'",
		Safe:        true,
	}
}

func init() {
	Register(NodeInstallPackageRule{})
}

// ------------------------------------------

type NodeRunProjectRule struct{}

func (r NodeRunProjectRule) Name() string {
	return "node.run_project"
}

func (r NodeRunProjectRule) Match(intent types.IntentType, workspace types.WorkspaceType) bool {
	return intent == types.IntentRunProject &&
		workspace == types.WorkspaceNode
}

func (r NodeRunProjectRule) Generate(
	arguments map[string]string,
) types.Command {
	file := arguments["file"]

	return types.Command{
		Command:     "npm run " + file,
		Explanation: "Starts the Node.js development server '" + file + "'",
		Safe:        true,
	}
}

func init() {
	Register(NodeRunProjectRule{})
}
