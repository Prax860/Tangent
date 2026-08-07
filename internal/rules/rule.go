package rules

import "github.com/prax860/tangent/internal/types"

type Rule interface {

	Name() string

	Match(
		intent types.IntentType,
		workspace types.WorkspaceType,
	) bool

	Generate() types.Command
}