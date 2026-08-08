package ui

import (
	"fmt"

	"github.com/prax860/tangent/internal/types"
)

func ShowResult(result types.ExecutionResult) {

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Execution Result")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Printf("Exit Code : %d\n", result.ExitCode)

	if result.Stdout != "" {
		fmt.Println()
		fmt.Println("STDOUT")
		fmt.Println(result.Stdout)
	}

	if result.Stderr != "" {
		fmt.Println()
		fmt.Println("STDERR")
		fmt.Println(result.Stderr)
	}

	fmt.Println()
	fmt.Printf("Duration : %v\n", result.Duration)
	fmt.Println()
}