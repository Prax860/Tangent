package types

import "time"

// =======================
// Workspace Types
// =======================

type WorkspaceType string

const (
	WorkspaceUnknown WorkspaceType = "unknown"
	WorkspacePython  WorkspaceType = "python"
	WorkspaceNode    WorkspaceType = "node"
	WorkspaceGo      WorkspaceType = "go"
	WorkspaceRust    WorkspaceType = "rust"
	WorkspaceJava    WorkspaceType = "java"
)

// =======================
// Intent Types
// =======================

type IntentType string

const (
	IntentUnknown IntentType = "unknown"

	IntentCreateVenv     IntentType = "create_venv"
	IntentInstallPackage IntentType = "install_package"
	IntentRunProject     IntentType = "run_project"

	IntentGitInit      IntentType = "git_init"
	IntentGitStatus    IntentType = "git_status"
	IntentCreateBranch IntentType = "create_branch"
	IntentPushChanges  IntentType = "push_changes"
)

// =======================
// Command Model
// =======================

type Command struct {
	Command     string
	Explanation string
	Safe        bool
}

// =======================
// Request Model
// =======================

type Request struct {
	RawInput   string
	Normalized string

	Intent    IntentType
	Workspace WorkspaceType

	Arguments map[string]string
}

// =======================
// Response Model
// =======================

type Response struct {
	Request Request
	Command Command
}

// =======================
// Execution Result
// =======================

type ExecutionResult struct {
	Command  string
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}
