package rules

import "github.com/prax860/tangent/internal/types"

type NodeInstallPackageRule struct{}

func (r NodeInstallPackageRule) Name() string {
	return "node.install_package"
}

func (r NodeInstallPackageRule) Match(intent types.IntentType, workspace types.WorkspaceType, _ map[string]string) bool {
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

func (r NodeRunProjectRule) Match(intent types.IntentType, workspace types.WorkspaceType, _ map[string]string) bool {
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

// ------------------------------------------

type NodeScaffoldNextRule struct{}

func (r NodeScaffoldNextRule) Name() string {
	return "node.scaffold_next"
}

func (r NodeScaffoldNextRule) Match(intent types.IntentType, workspace types.WorkspaceType, arguments map[string]string) bool {
	return intent == types.IntentCreateProject &&
		arguments["framework"] == "next"
}

func (r NodeScaffoldNextRule) Generate(
	arguments map[string]string,
) types.Command {
	return types.Command{
		Command:     "npx create-next-app@latest .",
		Explanation: "Scaffolds a new Next.js project in the current directory. All prompts (TypeScript, ESLint, Tailwind, etc.) are handled interactively by create-next-app.",
		Safe:        true,
		Interactive: true,
	}
}

func init() {
	Register(NodeScaffoldNextRule{})
}

// ------------------------------------------

type NodeScaffoldViteRule struct{}

func (r NodeScaffoldViteRule) Name() string {
	return "node.scaffold_vite"
}

func (r NodeScaffoldViteRule) Match(intent types.IntentType, workspace types.WorkspaceType, arguments map[string]string) bool {
	if intent != types.IntentCreateProject {
		return false
	}

	fw := arguments["framework"]
	switch fw {
	case "react", "vue", "svelte":
		return true
	}

	if workspace == types.WorkspaceNode || workspace == types.WorkspaceUnknown {
		return fw == ""
	}

	return false
}

func (r NodeScaffoldViteRule) Generate(
	arguments map[string]string,
) types.Command {
	return types.Command{
		Command:     "npm create vite@latest .",
		Explanation: "Scaffolds a new frontend project in the current directory using Vite. Framework/variant prompts are handled interactively by Vite.",
		Safe:        true,
		Interactive: true,
	}
}

func init() {
	Register(NodeScaffoldViteRule{})
}
