package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prax860/tangent/internal/core"
	"github.com/prax860/tangent/internal/ui"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start the interactive Tangent shell",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("               ⚡ Tangent Shell")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Println("Type 'help' for commands.")
		fmt.Println("Type 'exit' to quit.")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)

		for {

			fmt.Print("❯ ")

			input, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println(err)
				continue
			}

			input = strings.TrimSpace(input)

			if input == "" {
				continue
			}

			switch strings.ToLower(input) {

			case "exit", "quit":
				fmt.Println()
				fmt.Println("👋 Goodbye!")
				return

			case "clear":
				fmt.Print("\033[H\033[2J")
				continue

			case "help":
				printHelp()
				continue
			}

			response := core.Process(input)

			ui.Show(response)

			if !ui.Ask(response.Command) {
				fmt.Println("❌ Cancelled.")
				fmt.Println()
				continue
			}

			if response.Command.Interactive {
				fmt.Println()
				fmt.Println("⚡ Running...")
				fmt.Println()
			}

			result := ui.Execute(response.Command)

			fmt.Println()

			if response.Command.Interactive {
				if result.ExitCode == 0 {
					fmt.Println("✓ Command completed successfully.")
				} else {
					fmt.Printf("✗ Command failed with exit code %d.\n", result.ExitCode)
				}
				fmt.Printf("Duration : %v\n", result.Duration)
				fmt.Println()
				continue
			}

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
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

func printHelp() {

	fmt.Println()

	fmt.Println("Available commands")
	fmt.Println("------------------")
	fmt.Println("help   - Show help")
	fmt.Println("clear  - Clear screen")
	fmt.Println("exit   - Exit Tangent")
	fmt.Println()

	fmt.Println("Examples")
	fmt.Println("------------------")
	fmt.Println("install gin")
	fmt.Println("git status")
	fmt.Println("create virtual environment")
	fmt.Println("create react app")
	fmt.Println("docker init")
	fmt.Println()

}
