/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"
	"strings"

	"github.com/prax860/tangent/internal/pipelines"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a natural language command",
	Long:  "Accepts a natural language instruction and processes it through the Tangent pipeline.",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("❌ Please provide a command.")
			fmt.Println(`Example: tangent run "create virtual environment"`)
			return
		}

		// Join CLI arguments into a single sentence
		input := strings.Join(args, " ")

		// Process the request through the Tangent pipeline
		response := pipeline.Process(input)	

		// Display results
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("              ⚡ Tangent")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		fmt.Printf("📝 Input       : %s\n", response.Request.RawInput)
		fmt.Printf("🧹 Normalized  : %s\n", response.Request.Normalized)
		fmt.Printf("🧠 Intent      : %s\n", response.Request.Intent)
		fmt.Printf("📂 Workspace   : %s\n", response.Request.Workspace)
		fmt.Println()
		fmt.Printf("⚙️  Command     : %s\n", response.Command.Command)
		fmt.Printf("🔒 Safe        : %t\n", response.Command.Safe)
		fmt.Printf("💡 Explanation : %s\n", response.Command.Explanation)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}