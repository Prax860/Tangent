package entities

import (
	"strings"

	"github.com/prax860/tangent/internal/types"
)

func Extract(input string, intent types.IntentType) map[string]string {

	args := make(map[string]string)

	words := strings.Fields(input)

	switch intent {

	case types.IntentInstallPackage:

		// install express
		if len(words) >= 2 {
			args["package"] = words[1]
		}

	case types.IntentCreateBranch:

		// create branch feature/login
		if len(words) >= 3 {
			args["branch"] = strings.Join(words[2:], " ")
		}

	case types.IntentRunProject:

		// run main.py
		if len(words) >= 2 {
			args["file"] = words[1]
		}
	}

	return args
}
