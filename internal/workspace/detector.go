package workspace

import (
	"os"

	"github.com/prax860/tangent/internal/types"
)

// Detect identifies the current project's workspace
// by checking for well-known project files.
func Detect() types.WorkspaceType {

	// Node.js
	if fileExists("package.json") {
		return types.WorkspaceNode
	}

	// Python
	if fileExists("pyproject.toml") || fileExists("requirements.txt") {
		return types.WorkspacePython
	}

	// Go
	if fileExists("go.mod") {
		return types.WorkspaceGo
	}

	// Rust
	if fileExists("Cargo.toml") {
		return types.WorkspaceRust
	}

	// Java
	if fileExists("pom.xml") || fileExists("build.gradle") {
		return types.WorkspaceJava
	}

	return types.WorkspaceUnknown
}

// fileExists checks if a file exists in the current directory.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}