package core

import (
	"os"

	"github.com/prax860/tangent/internal/types"
)

func Detect() types.WorkspaceType {

	if fileExists("package.json") {
		return types.WorkspaceNode
	}

	if fileExists("pyproject.toml") || fileExists("requirements.txt") {
		return types.WorkspacePython
	}

	if fileExists("go.mod") {
		return types.WorkspaceGo
	}

	if fileExists("Cargo.toml") {
		return types.WorkspaceRust
	}

	if fileExists("pom.xml") || fileExists("build.gradle") {
		return types.WorkspaceJava
	}

	return types.WorkspaceUnknown
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
