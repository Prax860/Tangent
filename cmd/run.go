package cmd

import (
	"fmt"
	"strings"

	"github.com/prax860/tangent/internal/core"
	"github.com/prax860/tangent/internal/ui"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a natural language command",
	Long:  "Accepts a natural language instruction and processes it through the Tangent pipeline.",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("❌ Please provide a command.")
			fmt.Println(`Example: tangent run "install gin"`)
			return
		}

		input := strings.Join(args, " ")

		response := core.Process(input)

		ui.Show(response)

		if !ui.Ask(response.Command) {
			fmt.Println()
			fmt.Println("❌ Command cancelled.")
			return
		}

		if response.Command.Interactive {
			fmt.Println()
			fmt.Println("⚡ Running...")
			fmt.Println()
		}

		result := ui.Execute(response.Command)

		if response.Command.Interactive {
			fmt.Println()
			if result.ExitCode == 0 {
				fmt.Println("✓ Command completed successfully.")
			} else {
				fmt.Printf("✗ Command failed with exit code %d.\n", result.ExitCode)
			}
			fmt.Printf("Duration : %v\n", result.Duration)
			return
		}

		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("Execution Result")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		fmt.Printf("Exit Code : %d\n", result.ExitCode)

		if result.Stdout != "" {
			fmt.Println("\nSTDOUT")
			fmt.Println(result.Stdout)
		}

		if result.Stderr != "" {
			fmt.Println("\nSTDERR")
			fmt.Println(result.Stderr)
		}

		fmt.Printf("\nDuration : %v\n", result.Duration)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
