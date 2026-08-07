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

func (r NodeInstallPackageRule) Generate() types.Command {
	return types.Command{
		Command:     "npm install <package>",
		Explanation: "Installs a package using npm.",
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

func (r NodeRunProjectRule) Generate() types.Command {
	return types.Command{
		Command:     "npm run dev",
		Explanation: "Starts the Node.js development server.",
		Safe:        true,
	}
}

func init() {
	Register(NodeRunProjectRule{})
}