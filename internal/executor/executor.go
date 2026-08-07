package executor

import (
	"bytes"
	"os/exec"
	"time"

	"github.com/prax860/tangent/internal/types"
)

func Execute(command types.Command) types.ExecutionResult {

	start := time.Now()

	cmd := exec.Command("cmd", "/C", command.Command)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0

	if err != nil {

		if e, ok := err.(*exec.ExitError); ok {
			exitCode = e.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return types.ExecutionResult{

		Command: command.Command,

		Stdout: stdout.String(),

		Stderr: stderr.String(),

		ExitCode: exitCode,

		Duration: time.Since(start),
	}
}
