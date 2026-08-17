package rules

import "github.com/prax860/tangent/internal/types"

type DockerInitRule struct{}

func (r DockerInitRule) Name() string {
	return "docker.init"
}

func (r DockerInitRule) Match(
	intent types.IntentType,
	workspace types.WorkspaceType,
	arguments map[string]string,
) bool {

	if intent != types.IntentCreateProject {
		return false
	}

	if arguments["framework"] == "docker" {
		return true
	}
	return false
}

func (r DockerInitRule) Generate(
	arguments map[string]string,
) types.Command {

	return types.Command{
		Command:     "docker init",
		Explanation: "Launches Docker's interactive project initializer. Docker will ask questions about the stack interactively.",
		Safe:        true,
		Interactive: true,
	}
}

func init() {
	Register(DockerInitRule{})
}
