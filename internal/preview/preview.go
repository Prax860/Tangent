package preview

import (
	"fmt"

	"github.com/prax860/tangent/internal/types"
)

func Show(response types.Response) {

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("              ⚡ Tangent")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	fmt.Println("📝 Request")
	fmt.Println(response.Request.RawInput)

	fmt.Println()

	fmt.Println("🧠 Intent")
	fmt.Println(response.Request.Intent)

	fmt.Println()

	fmt.Println("📂 Workspace")
	fmt.Println(response.Request.Workspace)

	fmt.Println()

	if len(response.Request.Arguments) > 0 {

		fmt.Println("📦 Arguments")

		for key, value := range response.Request.Arguments {
			fmt.Printf("%s : %s\n", key, value)
		}

		fmt.Println()
	}

	fmt.Println("⚙ Generated Command")
	fmt.Println(response.Command.Command)

	fmt.Println()

	fmt.Println("💡 Explanation")
	fmt.Println(response.Command.Explanation)

	fmt.Println()

	fmt.Println("🔒 Safe :", response.Command.Safe)

}
