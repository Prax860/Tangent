package core

import (
	"strings"

	"github.com/prax860/tangent/internal/types"
)

func Extract(input string, intent types.IntentType) map[string]string {

	args := make(map[string]string)

	words := strings.Fields(input)

	switch intent {

	case types.IntentInstallPackage:

		if len(words) >= 2 {
			args["package"] = words[1]
		}

	case types.IntentCreateBranch:

		if len(words) >= 3 {
			args["branch"] = strings.Join(words[2:], " ")
		}

	case types.IntentRunProject:

		if len(words) >= 2 {
			args["file"] = words[1]
		}
	}

	return args
}
