package rules

import "github.com/prax860/tangent/internal/types"

type Rule interface {
	Name() string

	Match(
		intent types.IntentType,
		workspace types.WorkspaceType,
		arguments map[string]string,
	) bool

	Generate(
		arguments map[string]string,
	) types.Command
}
