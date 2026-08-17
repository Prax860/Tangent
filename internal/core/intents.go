package core

import (
	"strings"

	"github.com/prax860/tangent/internal/types"
)

func Resolve(input string) types.IntentType {

	switch {

	case strings.Contains(input, "virtual environment"),
		strings.Contains(input, "venv"),
		strings.Contains(input, "virtualenv"):

		return types.IntentCreateVenv

	case strings.HasPrefix(input, "create project"),
		strings.Contains(input, "create react app"),
		strings.Contains(input, "create next app"),
		strings.Contains(input, "create vue app"),
		strings.Contains(input, "create svelte app"),
		strings.Contains(input, "scaffold"),
		strings.Contains(input, "init project"),
		strings.Contains(input, "docker init"),
		strings.HasPrefix(input, "npm create"),
		strings.HasPrefix(input, "npx create"),
		(strings.HasPrefix(input, "create") && strings.Contains(input, "app")),
		(strings.HasPrefix(input, "cargo") && strings.Contains(input, "init")),
		(strings.Contains(input, "cargo init")):

		return types.IntentCreateProject

	case strings.HasPrefix(input, "install"):

		return types.IntentInstallPackage

	case strings.Contains(input, "git init"):

		return types.IntentGitInit

	case strings.Contains(input, "git status"):

		return types.IntentGitStatus

	case strings.Contains(input, "create branch"),
		strings.Contains(input, "new branch"):

		return types.IntentCreateBranch

	case strings.Contains(input, "push"):

		return types.IntentPushChanges

	default:

		return types.IntentUnknown
	}
}
