package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prax860/tangent/internal/types"
)

func Ask(command types.Command) bool {

	if !command.Safe {

		fmt.Println()
		fmt.Println("⚠ WARNING")
		fmt.Println(command.Explanation)
		fmt.Println()

		fmt.Print("Continue anyway? (y/N): ")

	} else {

		fmt.Println()
		fmt.Print("Execute command? (y/N): ")

	}

	reader := bufio.NewReader(os.Stdin)

	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}
