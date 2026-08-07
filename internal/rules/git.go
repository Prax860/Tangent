package rules

import "github.com/prax860/tangent/internal/types"

type GitInitRule struct{}

func (r GitInitRule) Name() string {
	return "git.init"
}

func (r GitInitRule) Match(intent types.IntentType, workspace types.WorkspaceType) bool {
	return intent == types.IntentGitInit
}

func (r GitInitRule) Generate(
	arguments map[string]string,
) types.Command {
	return types.Command{
		Command:     "git init",
		Explanation: "Initializes a new Git repository.",
		Safe:        true,
	}
}

func init() {
	Register(GitInitRule{})
}

// ------------------------------------------

type GitStatusRule struct{}

func (r GitStatusRule) Name() string {
	return "git.status"
}

func (r GitStatusRule) Match(
	intent types.IntentType,
	workspace types.WorkspaceType,
) bool {
	return intent == types.IntentGitStatus
}

func (r GitStatusRule) Generate(
	arguments map[string]string,
) types.Command {
	return types.Command{
		Command:     "git status",
		Explanation: "Displays the current Git repository status.",
		Safe:        true,
	}
}

func init() {
	Register(GitStatusRule{})
}

// ------------------------------------------

type GitCreateBranchRule struct{}

func (r GitCreateBranchRule) Name() string {
	return "git.branch"
}

func (r GitCreateBranchRule) Match(intent types.IntentType, workspace types.WorkspaceType) bool {
	return intent == types.IntentCreateBranch
}

func (r GitCreateBranchRule) Generate(
	arguments map[string]string,
) types.Command {
	branch := arguments["branch"]

	return types.Command{
		Command:     "git checkout -b " + branch,
		Explanation: "Creates and switches to a new Git branch '" + branch + "'",
		Safe:        true,
	}
}

func init() {
	Register(GitCreateBranchRule{})
}

// ------------------------------------------

type GitPushRule struct{}

func (r GitPushRule) Name() string {
	return "git.push"
}

func (r GitPushRule) Match(intent types.IntentType, workspace types.WorkspaceType) bool {
	return intent == types.IntentPushChanges
}

func (r GitPushRule) Generate(
	arguments map[string]string,
) types.Command {
	return types.Command{
		Command:     "git push",
		Explanation: "Pushes commits to the remote repository.",
		Safe:        true,
	}
}

func init() {
	Register(GitPushRule{})
}
