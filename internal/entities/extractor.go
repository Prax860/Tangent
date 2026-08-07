package entities

import (
	"strings"

	"github.com/prax860/tangent/internal/types"
)

func Extract(input string, intent types.IntentType) map[string]string {

	args := make(map[string]string)

	switch intent {

	case types.IntentInstallPackage:

		words := strings.Fields(input)

		if len(words) >= 2 {

			args["package"] = words[len(words)-1]

		}

	case types.IntentCreateBranch:

		words := strings.Fields(input)

		if len(words) >= 3 {

			args["branch"] = words[len(words)-1]

		}

	default:
		// Handle unknown intents
	}

	return args

}