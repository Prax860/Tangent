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

	case types.IntentCreateProject:

		switch {

		case strings.Contains(input, "react"):
			args["framework"] = "react"

		case strings.Contains(input, "vue"):
			args["framework"] = "vue"

		case strings.Contains(input, "svelte"):
			args["framework"] = "svelte"

		case strings.Contains(input, "next"):
			args["framework"] = "next"

		case strings.Contains(input, "docker"):
			args["framework"] = "docker"

		case strings.Contains(input, "cargo") || strings.Contains(input, "rust"):
			args["framework"] = "rust"

		case strings.Contains(input, "go") || strings.Contains(input, "golang"):
			args["framework"] = "go"
		}

		args["rawInput"] = input
	}

	return args
}
