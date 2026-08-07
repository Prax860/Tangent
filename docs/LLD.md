# Low-Level Design — Tangent

**Version:** 1.0
**Date:** 2026-08-07

The LLD is the authoritative reference for *every source file*, *every exported and unexported function*, *every type*, *every package import*, *every control-flow branch*, and *every concrete rule*. When Claude is asked to modify Tangent code, it MUST treat the exact signatures and behaviours here as constraints. Deviations require an explicit update to this LLD.

---

## 0. File Index (exact, as of this version)

```
Tangent/
├── go.mod                          module github.com/prax860/tangent ; go 1.25.5
├── go.sum
├── main.go                         package main
├── cmd/
│   ├── root.go                     package cmd
│   └── run.go                      package cmd
└── internal/
    ├── types/
    │   └── types.go                package types
    ├── core/
    │   ├── parser.go               package core
    │   ├── intents.go              package core
    │   ├── entities.go             package core
    │   ├── workspace.go            package core
    │   └── pipeline.go             package core
    ├── rules/
    │   ├── rule.go                 package rules
    │   ├── registry.go             package rules
    │   ├── engine.go               package rules
    │   ├── git.go                  package rules
    │   ├── go.go                   package rules
    │   ├── node.go                 package rules
    │   └── python.go               package rules
    └── ui/
        ├── preview.go              package ui
        ├── confirm.go              package ui
        └── executor.go             package ui
```

**File count:** 18 `.go` source files + `go.mod`, `go.sum` (20).

---

## 1. go.mod

```
module github.com/prax860/tangent
go 1.25.5
require github.com/spf13/cobra v1.10.2
require (
    github.com/inconshreveable/mousetrap v1.1.0 // indirect
    github.com/spf13/pflag v1.0.9      // indirect
)
```

**Internal import base path:** Every `internal/*` import MUST use the prefix `github.com/prax860/tangent/internal/<package>`.

---

## 2. main.go

- **File:** [main.go](file:///d:/Tangent/main.go)
- **Package:** `package main`
- **Imports:**
  - `"github.com/prax860/tangent/cmd"` (only)

### 2.1 Function: `func main()`

```
Behaviour:
  Call cmd.Execute() exactly once.
  No recover(), no defer(), no flags, no other logic.
```

---

## 3. cmd package

### 3.1 cmd/root.go

- **File:** [cmd/root.go](file:///d:/Tangent/cmd/root.go)
- **Package:** `package cmd`
- **Imports:**
  - `"os"`
  - `"github.com/spf13/cobra"`

#### Unexported variable: `rootCmd *cobra.Command`

```go
var rootCmd = &cobra.Command{
    Use:   "tangent",
    Short: "A brief description of your application",
    Long:  `...(long description text)...`,
    // Run: nil (cobra default → print help)
}
```

#### Exported function: `func Execute()`

```
Behaviour:
  err := rootCmd.Execute()
  if err != nil → os.Exit(1)
```

#### Unexported function: `func init()`

```
Behaviour:
  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
  (This flag is currently unused by any handler; kept for future Cobra convention.)
```

### 3.2 cmd/run.go

- **File:** [cmd/run.go](file:///d:/Tangent/cmd/run.go)
- **Package:** `package cmd`
- **Imports:**
  - `"fmt"`
  - `"strings"`
  - `"github.com/prax860/tangent/internal/core"`
  - `"github.com/prax860/tangent/internal/ui"`
  - `"github.com/spf13/cobra"`

#### Unexported variable: `runCmd *cobra.Command`

```go
var runCmd = &cobra.Command{
    Use:   "run",
    Short: "Run a natural language command",
    Long:  "Accepts a natural language instruction and processes it through the Tangent pipeline.",
    Run:   <anonymous closure>,
}
```

#### Anonymous closure: `func(cmd *cobra.Command, args []string)`

```
Exact sequential behaviour:

  1. if len(args) == 0:
       fmt.Println("❌ Please provide a command.")
       fmt.Println(`Example: tangent run "install gin"`)
       return (exit handler normally; cobra returns nil → OS exit 0)

  2. input := strings.Join(args, " ")   // preserves argument order; joins with U+0020

  3. response := core.Process(input)    // pipeline; returns types.Response

  4. ui.Show(response)                  // side-effect: prints preview to stdout

  5. if !ui.Ask(response.Command):      // side-effect: prints prompt; reads stdin
       fmt.Println()
       fmt.Println("❌ Command cancelled.")
       return

  6. result := ui.Execute(response.Command)  // side-effect: spawns child process

  7. fmt.Println()
     fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
     fmt.Println("Execution Result")
     fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

  8. fmt.Printf("Exit Code : %d\n", result.ExitCode)

  9. if result.Stdout != "":
        fmt.Println("\nSTDOUT")
        fmt.Println(result.Stdout)

  10. if result.Stderr != "":
        fmt.Println("\nSTDERR")
        fmt.Println(result.Stderr)

  11. fmt.Printf("\nDuration : %v\n", result.Duration)
```

#### Unexported function: `func init()`

```
Behaviour:
  rootCmd.AddCommand(runCmd)
  (Attaches run subcommand to root. Run exactly once per process via Go's init.)
```

---

## 4. types package (leaf, zero internal imports)

### 4.1 internal/types/types.go

- **File:** [types.go](file:///d:/Tangent/internal/types/types.go)
- **Package:** `package types`
- **Imports:**
  - `"time"` (only; used by `ExecutionResult.Duration`)

#### 4.1.1 Type: `WorkspaceType string` (named string)

Constants — the canonical set. Do NOT add elsewhere:

```go
const (
    WorkspaceUnknown WorkspaceType = "unknown"
    WorkspacePython  WorkspaceType = "python"
    WorkspaceNode    WorkspaceType = "node"
    WorkspaceGo      WorkspaceType = "go"
    WorkspaceRust    WorkspaceType = "rust"
    WorkspaceJava    WorkspaceType = "java"
)
```

Precedence order in `workspace.Detect()` is: Node → Python → Go → Rust → Java → Unknown. (Not the same as const declaration order. See §5.4.)

#### 4.1.2 Type: `IntentType string` (named string)

```go
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
```

Switch/case order in the resolver (matters for conflict resolution):
`IntentCreateVenv → IntentInstallPackage → IntentGitInit → IntentGitStatus → IntentCreateBranch → IntentPushChanges → default: IntentUnknown`

#### 4.1.3 Struct: `Command`

```go
type Command struct {
    Command     string   // shell command string; passed to cmd /C as whole argument
    Explanation string   // human-readable why this command; shown in preview + unsafe warning
    Safe        bool     // determines confirmation prompt style
}
```

Invariants:
- If `Safe == true`, a preview-only accidental-execution mitigation is assumed by the UI.
- If `Safe == false`, the UI MUST print the ⚠ WARNING block.
- The "no matching rule" fallback returns `Command{Command:"", Explanation:"No matching rule found.", Safe:false}` — note empty `Command`; a future guard should prevent executing empty strings.

#### 4.1.4 Struct: `Request`

```go
type Request struct {
    RawInput   string            // exact user-joined input before Parse
    Normalized string            // after Parse
    Intent     IntentType
    Workspace  WorkspaceType
    Arguments  map[string]string // entity extraction output
}
```

#### 4.1.5 Struct: `Response`

```go
type Response struct {
    Request Request
    Command Command
}
```

This is the immutable aggregate returned by `core.Process`. It is the only data structure that flows from `core` into `ui`.

#### 4.1.6 Struct: `ExecutionResult`

```go
type ExecutionResult struct {
    Command  string        // echoed command string
    Stdout   string
    Stderr   string
    ExitCode int
    Duration time.Duration
}
```

Exit code semantics:
- `0` → child ran, exited 0 (or on Windows, `cmd /C` reported 0)
- positive int → child's exit code via `exec.ExitError.ExitCode()`
- `-1` → error happened that was NOT an `*exec.ExitError` (e.g. process failed to start)

---

## 5. core package

All files share `package core`. They live in the same directory and see each other's symbols without imports.

### 5.1 internal/core/parser.go

- **File:** [parser.go](file:///d:/Tangent/internal/core/parser.go)
- **Package:** `package core`
- **Imports:** `"regexp"`, `"strings"`

#### Exported function: `func Parse(input string) string`

```
Ordered transformations (strictly in this sequence):

  1. input = strings.ToLower(input)
       - ASCII uppercase A-Z → a-z; locale-independent
  2. input = strings.TrimSpace(input)
       - trims \t\n\v\f\r U+0020 U+0085 U+00A0 (Unicode whitespace) both ends
  3. re1 = regexp.MustCompile(`[^\w\s./-]`)
     input = re1.ReplaceAllString(input, "")
       - Whitelist character class:
         \w  → [0-9A-Za-z_] (after toLower, effectively digits + lowercase + underscore)
         \s  → any Unicode whitespace category
         .   → literal dot
         /   → literal forward slash
         -   → literal hyphen (safe because it is last in class)
       - Everything else is DELETED (replaced with empty string)
  4. re2 = regexp.MustCompile(`\s+`)
     input = re2.ReplaceAllString(input, " ")
       - Collapse any run (≥1) of any whitespace → exactly one space char (U+0020)
  5. return input
```

NOTES:
- The regexps are compiled on EVERY invocation. For higher throughput, move them to package-level `var` literals compiled at `init()`. Today it is kept simple because input volume is CLI-scale (~one per process).
- Behaviour for empty string: returns empty string.

### 5.2 internal/core/intents.go

- **File:** [intents.go](file:///d:/Tangent/internal/core/intents.go)
- **Package:** `package core`
- **Imports:**
  - `"strings"`
  - `"github.com/prax860/tangent/internal/types"`

#### Exported function: `func Resolve(input string) types.IntentType`

```
Switch/if-else chain (EVALUATION ORDER IS SIGNIFICANT):

  if Contains "virtual environment"  OR  Contains "venv"  OR  Contains "virtualenv":
      return types.IntentCreateVenv

  else if HasPrefix "install":
      return types.IntentInstallPackage

  else if Contains "git init":
      return types.IntentGitInit

  else if Contains "git status":
      return types.IntentGitStatus

  else if Contains "create branch" OR Contains "new branch":
      return types.IntentCreateBranch

  else if Contains "push":
      return types.IntentPushChanges

  else:
      return types.IntentUnknown
```

Match semantics:
- `strings.Contains` = exact substring search (case-sensitive; input is already lowercase so this is effectively case-insensitive for original ASCII input).
- `strings.HasPrefix` = exact prefix at byte 0.
- A pathological input like `"install virtual environment package"` matches `IntentCreateVenv` (CreateVenv test is checked BEFORE HasPrefix install). This is intentional, not a bug.

### 5.3 internal/core/entities.go

- **File:** [entities.go](file:///d:/Tangent/internal/core/entities.go)
- **Package:** `package core`
- **Imports:**
  - `"strings"`
  - `"github.com/prax860/tangent/internal/types"`

#### Exported function: `func Extract(input string, intent types.IntentType) map[string]string`

```
Behaviour:

  1. args := make(map[string]string)     // empty map, not nil
  2. words := strings.Fields(input)      // splits on any Unicode whitespace, drops empty tokens
  3. switch intent:

     case IntentInstallPackage:
         if len(words) >= 2:
             args["package"] = words[1]

     case IntentCreateBranch:
         if len(words) >= 3:
             args["branch"] = strings.Join(words[2:], " ")   // multi-word branch name

     case IntentRunProject:
         if len(words) >= 2:
             args["file"] = words[1]

     default:
         (nothing is set; args remains empty)

  4. return args
```

Notes:
- `map[string]string` is NEVER nil; callers may safely iterate.
- Intent-unknown branch key set is kept EMPTY intentionally; no "raw" catch-all key.
- If len(words) is too small for the intent, args returns partial (not an error).

### 5.4 internal/core/workspace.go

- **File:** [workspace.go](file:///d:/Tangent/internal/core/workspace.go)
- **Package:** `package core`
- **Imports:**
  - `"os"`
  - `"github.com/prax860/tangent/internal/types"`

#### Exported function: `func Detect() types.WorkspaceType`

```
Ordered precedence (highest to lowest):

  if fileExists("package.json"): return WorkspaceNode
  if fileExists("pyproject.toml") || fileExists("requirements.txt"): return WorkspacePython
  if fileExists("go.mod"):         return WorkspaceGo
  if fileExists("Cargo.toml"):     return WorkspaceRust
  if fileExists("pom.xml") || fileExists("build.gradle"): return WorkspaceJava
  return WorkspaceUnknown
```

Significant: If a monorepo has both `package.json` and `go.mod`, it is classified as **Node**. This is by-design-first-match-wins precedence.

#### Unexported helper: `func fileExists(filename string) bool`

```
Implementation:
  _, err := os.Stat(filename)
  return err == nil
```

Uses `os.Stat` (NOT `os.IsNotExist` negation). If Stat returns *any* non-nil error (incl. permission denied, invalid name), the function treats the file as non-existent. Callers (`Detect`) receive false and fall through.

### 5.5 internal/core/pipeline.go

- **File:** [pipeline.go](file:///d:/Tangent/internal/core/pipeline.go)
- **Package:** `package core`
- **Imports:**
  - `"github.com/prax860/tangent/internal/rules"`
  - `"github.com/prax860/tangent/internal/types"`

#### Exported function: `func Process(input string) types.Response`

This is the single entry point into the pipeline. Exact control flow, line-by-line equivalent:

```
1. normalized  := Parse(input)                // local call within package
2. intent      := Resolve(normalized)         // local call
3. arguments   := Extract(normalized, intent) // local call
4. workspaceType := Detect()                  // local call; reads filesystem
5. command     := rules.Generate(intent, workspaceType, arguments)
6. request := types.Request{
       RawInput:   input,            // original, UN-normalized
       Normalized: normalized,
       Intent:     intent,
       Workspace:  workspaceType,
       Arguments:  arguments,
   }
7. return types.Response{ Request: request, Command: command }
```

Note `Request.RawInput` = the caller's original `input` (before Parse). This is intentional so the preview can echo the user's exact wording.

---

## 6. rules package

All files share `package rules`. Live in same dir → inter-file unexported symbols shared (registry slice, etc).

### 6.1 internal/rules/rule.go — Rule interface

- **File:** [rule.go](file:///d:/Tangent/internal/rules/rule.go)
- **Package:** `package rules`
- **Imports:** `"github.com/prax860/tangent/internal/types"`

```go
type Rule interface {
    Name() string
    Match(
        intent    types.IntentType,
        workspace types.WorkspaceType,
    ) bool
    Generate(
        arguments map[string]string,
    ) types.Command
}
```

Contract for concrete rules:
- `Name()` MUST be stable across code versions; it is the human + machine identifier of the rule. Recommended format `"domain.action"`, e.g. `"go.install_package"`.
- `Match` MUST be pure.
- `Generate` MUST be pure; the returned `Command.Safe` MUST be `true` for read-only / low-risk operations. Mark `false` for destructive operations (write, delete, force-push, install). Currently v1 marks everything (even install) `true`; this is an explicit design call. Marking install as unsafe would change user experience significantly.

### 6.2 internal/rules/registry.go

- **File:** [registry.go](file:///d:/Tangent/internal/rules/registry.go)
- **Package:** `package rules`
- **Imports:** (none)

```go
var registry []Rule

func Register(rule Rule) {
    registry = append(registry, rule)
}

func Rules() []Rule {
    return registry
}
```

Invariants:
- `Register` is append-only. No removal, no de-duplication.
- `Rules()` returns the slice header directly (a copy is NOT made). Callers MUST NOT mutate it. Current callers only iterate read-only.
- Order of entries is the order in which `Register` was called during `init()` execution for the `rules` package. Because init order across files in a package is the *source file name lexical order* (per Go spec), the effective current order is:
  1. `engine.go` → (no inits)
  2. `git.go` → inits: GitInitRule, GitStatusRule, GitCreateBranchRule, GitPushRule
  3. `go.go` → GoInstallPackageRule
  4. `node.go` → NodeInstallPackageRule, NodeRunProjectRule
  5. `python.go` → CreateVenvRule
  6. `registry.go` → (no inits)
  7. `rule.go` → (no inits)

Wait — Go's lexical `init()` order within a package is the **filename order**, not declaration order inside a file. Exact current order (all file names alphabetical within `internal/rules/`):
1. `engine.go` — no init
2. `git.go` — inits at top of file → GitInitRule first, then GitStatusRule, then GitCreateBranchRule, then GitPushRule (per their physical init() positions inside `git.go`, which Go runs top-to-bottom within one file)
3. `go.go` — GoInstallPackageRule
4. `node.go` — NodeInstallPackageRule, then NodeRunProjectRule
5. `python.go` — CreateVenvRule
6. `registry.go` — no init
7. `rule.go` — no init

So `Rules()` returns: `[GitInit, GitStatus, GitCreateBranch, GitPush, GoInstall, NodeInstall, NodeRun, PythonCreateVenv]`.

### 6.3 internal/rules/engine.go — Generate match loop

- **File:** [engine.go](file:///d:/Tangent/internal/rules/engine.go)
- **Package:** `package rules`
- **Imports:** `"github.com/prax860/tangent/internal/types"`

#### Exported function: `func Generate(intent types.IntentType, workspace types.WorkspaceType, arguments map[string]string) types.Command`

```
for _, rule := range Rules():
    if rule.Match(intent, workspace):
        return rule.Generate(arguments)

// No match → fallback (non-optional):
return types.Command{
    Command:     "",
    Explanation: "No matching rule found.",
    Safe:        false,
}
```

### 6.4 internal/rules/git.go — 4 rules

- **File:** [git.go](file:///d:/Tangent/internal/rules/git.go)
- **Package:** `package rules`
- **Imports:** `"github.com/prax860/tangent/internal/types"`

All 4 git rules have `Match(intent, _)` = intent-specific, **workspace-agnostic**.

#### `type GitInitRule struct{}`

| Method | Implementation |
|---|---|
| `Name()` → `"git.init"` | |
| `Match(intent, ws)` → `intent == types.IntentGitInit` | ws ignored |
| `Generate(args)` → `Command{"git init", "Initializes a new Git repository.", true}` | args unused |

`init()` → `Register(GitInitRule{})`

#### `type GitStatusRule struct{}`

| Method | Implementation |
|---|---|
| `Name()` → `"git.status"` | |
| `Match(i, w)` → `i == IntentGitStatus` | |
| `Generate(args)` → `Command{"git status", "Displays the current Git repository status.", true}` | |

`init()` → `Register(GitStatusRule{})`

#### `type GitCreateBranchRule struct{}`

| Method | Implementation |
|---|---|
| `Name()` → `"git.branch"` | |
| `Match(i, w)` → `i == IntentCreateBranch` | |
| `Generate(args)` | `branch := args["branch"]`<br>Command:<br>- `Command: "git checkout -b " + branch`<br>- `Explanation: "Creates and switches to a new Git branch '" + branch + "'"`<br>- `Safe: true` |

Note: `branch` may be empty string if entity extractor did not set it (input short). This produces `git checkout -b ` with a trailing space. Behaviour is left as-is for v1. No quoting, no validation.

`init()` → `Register(GitCreateBranchRule{})`

#### `type GitPushRule struct{}`

| Method | Implementation |
|---|---|
| `Name()` → `"git.push"` | |
| `Match(i, w)` → `i == IntentPushChanges` | |
| `Generate(args)` → `Command{"git push", "Pushes commits to the remote repository.", true}` | |

`init()` → `Register(GitPushRule{})`

### 6.5 internal/rules/go.go — 1 rule

- **File:** [go.go](file:///d:/Tangent/internal/rules/go.go)
- **Package:** `package rules`
- **Imports:** `"github.com/prax860/tangent/internal/types"`

#### `type GoInstallPackageRule struct{}`

| Method | Implementation |
|---|---|
| `Name()` → `"go.install_package"` | |
| `Match(i, w)` → `i == IntentInstallPackage && w == WorkspaceGo` | both required |
| `Generate(args)` | `module := args["package"]`<br>Command:<br>- `Command: "go get " + module`<br>- `Explanation: "Installs Go module '" + module + "'"`<br>- `Safe: true` |

`init()` → `Register(GoInstallPackageRule{})`

### 6.6 internal/rules/node.go — 2 rules

- **File:** [node.go](file:///d:/Tangent/internal/rules/node.go)
- **Package:** `package rules`
- **Imports:** `"github.com/prax860/tangent/internal/types"`

#### `type NodeInstallPackageRule struct{}`

| Method | Implementation |
|---|---|
| `Name()` → `"node.install_package"` | |
| `Match(i, w)` → `i == IntentInstallPackage && w == WorkspaceNode` | both required |
| `Generate(args)` | `packageName := args["package"]`<br>Command:<br>- `Command: "npm install " + packageName`<br>- `Explanation: "Installs Node.js package '" + packageName + "'"`<br>- `Safe: true` |

`init()` → `Register(NodeInstallPackageRule{})`

#### `type NodeRunProjectRule struct{}`

| Method | Implementation |
|---|---|
| `Name()` → `"node.run_project"` | |
| `Match(i, w)` → `i == IntentRunProject && w == WorkspaceNode` | |
| `Generate(args)` | `file := args["file"]`<br>Command:<br>- `Command: "npm run " + file`<br>- `Explanation: "Starts the Node.js development server '" + file + "'"`<br>- `Safe: true` |

`init()` → `Register(NodeRunProjectRule{})`

### 6.7 internal/rules/python.go — 1 rule

- **File:** [python.go](file:///d:/Tangent/internal/rules/python.go)
- **Package:** `package rules`
- **Imports:** `"github.com/prax860/tangent/internal/types"`

#### `type CreateVenvRule struct{}`

| Method | Implementation |
|---|---|
| `Name()` → `"python.create_venv"` | |
| `Match(i, w)` → `i == IntentCreateVenv && w == WorkspacePython` | |
| `Generate(args)` | `venv := args["venv"]`<br>Command:<br>- `Command: "python -m venv " + venv`<br>- `Explanation: "Creates Python virtual environment '" + venv + "'"`<br>- `Safe: true` |

Note: Currently no entity extraction case sets `args["venv"]` in `entities.go`. (Entity extractor only sets `package`, `branch`, `file`.) So `venv` is **always empty string** at runtime for this rule, producing `python -m venv ` with trailing empty name. This is a **known issue** (as of v1.0). Fix requires adding an `IntentCreateVenv` case to `entities.go` that extracts e.g. the word after "venv" / "virtual environment" as the venv name.

`init()` → `Register(CreateVenvRule{})`

---

## 7. ui package

All files share `package ui`. Same directory.

### 7.1 internal/ui/preview.go

- **File:** [preview.go](file:///d:/Tangent/internal/ui/preview.go)
- **Package:** `package ui`
- **Imports:**
  - `"fmt"`
  - `"github.com/prax860/tangent/internal/types"`

#### Exported function: `func Show(response types.Response)`

Side-effect only. Returns nothing. Writes to `os.Stdout` via fmt.Println / fmt.Printf.

```
Ordered output:

1. ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  (41 U+2501 chars)
2.               ⚡ Tangent
3. ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
4. blank line
5. "📝 Request"             (label on its own line)
6. response.Request.RawInput  (echo verbatim; not normalized)
7. blank line
8. "🧠 Intent"
9. response.Request.Intent  (the type's string form, e.g. "install_package")
10. blank line
11. "📂 Workspace"
12. response.Request.Workspace  (e.g. "go", "node", "unknown")
13. blank line
14. if len(response.Request.Arguments) > 0:
       "📦 Arguments"
       for key, value := range response.Request.Arguments:
           fmt.Printf("%s : %s\n", key, value)
       blank line
15. "⚙ Generated Command"
16. response.Command.Command
17. blank line
18. "💡 Explanation"
19. response.Command.Explanation
20. blank line
21. "🔒 Safe :", response.Command.Safe   (prints "true" or "false" verbatim)
```

Important: `Arguments` iteration order is Go's randomized map iteration order. Output order between argument keys is intentionally non-deterministic. Do NOT rely on it for testing.

### 7.2 internal/ui/confirm.go

- **File:** [confirm.go](file:///d:/Tangent/internal/ui/confirm.go)
- **Package:** `package ui`
- **Imports:**
  - `"bufio"`
  - `"fmt"`
  - `"os"`
  - `"strings"`
  - `"github.com/prax860/tangent/internal/types"`

#### Exported function: `func Ask(command types.Command) bool`

```
Branch A — !command.Safe (unsafe):
  blank line
  "⚠ WARNING"
  command.Explanation  (printed on its own line)
  blank line
  prompt: "Continue anyway? (y/N): "   (no newline after colon)

Branch B — command.Safe (normal):
  blank line
  prompt: "Execute command? (y/N): "

Then, common code:
  reader := bufio.NewReader(os.Stdin)
  input, _ := reader.ReadString('\n')       // error discarded on EOF/read fail → input = ""
  input = strings.TrimSpace(strings.ToLower(input))
  return (input == "y") || (input == "yes")
```

Semantics:
- The prompt suffix `(y/N)` uses capital N meaning "default is No" convention. Any non-"y"/non-"yes" input, including empty input (just Enter) → false (cancelled).
- Read failure (e.g. closed stdin in piped invocation) → returns false safely.

### 7.3 internal/ui/executor.go

- **File:** [executor.go](file:///d:/Tangent/internal/ui/executor.go)
- **Package:** `package ui`
- **Imports:**
  - `"bytes"`
  - `"os/exec"`
  - `"time"`
  - `"github.com/prax860/tangent/internal/types"`

#### Exported function: `func Execute(command types.Command) types.ExecutionResult`

```
Ordered behaviour:

  1. start := time.Now()

  2. cmd := exec.Command("cmd", "/C", command.Command)
       - Uses cmd.exe /C on Windows. The whole command.Command is a single argument to /C.
       - cmd.exe will parse and execute command.Command with its own quoting rules.

  3. var stdout bytes.Buffer ; var stderr bytes.Buffer
     cmd.Stdout = &stdout ; cmd.Stderr = &stderr
       - Child's inherits env/cwd from tangent process, not stdin (stdin not wired, so child gets /dev/null-equivalent handle on Windows). This is intentional: we do NOT want the child to take interactive input.

  4. err := cmd.Run()        // blocks until child terminates

  5. exitCode := 0
     if err != nil:
         if e, ok := err.(*exec.ExitError); ok:
             exitCode = e.ExitCode()   // real exit code from process
         else:
             exitCode = -1              // non-exit failure (e.g. cmd.exe couldn't start)

  6. return types.ExecutionResult{
         Command:  command.Command,
         Stdout:   stdout.String(),
         Stderr:   stderr.String(),
         ExitCode: exitCode,
         Duration: time.Since(start),
     }
```

CROSS-PLATFORM NOTE: `exec.Command("cmd", "/C", command.Command)` is Windows-only. To port:
- POSIX: `exec.Command("sh", "-c", command.Command)`
- Or use build-tag file split: `executor_windows.go` with `//go:build windows` and `executor_unix.go` with `//go:build !windows`.

---

## 8. Call Graph (all cross-package edges)

```
main.main
  └─ cmd.Execute
       └─ cobra → runCmd.Run closure
            ├─ core.Process
            │    ├─ core.Parse
            │    ├─ core.Resolve
            │    ├─ core.Extract
            │    ├─ core.Detect
            │    │    └─ fileExists (local)
            │    └─ rules.Generate
            │         └─ rules.Rules
            │              └─ iterate Rule.Match / Rule.Generate over registry (contains:)
            │                   ├─ GitInitRule
            │                   ├─ GitStatusRule
            │                   ├─ GitCreateBranchRule
            │                   ├─ GitPushRule
            │                   ├─ GoInstallPackageRule
            │                   ├─ NodeInstallPackageRule
            │                   ├─ NodeRunProjectRule
            │                   └─ CreateVenvRule
            ├─ ui.Show
            ├─ ui.Ask
            └─ ui.Execute
```

---

## 9. Concrete Examples (Execution Traces)

### Example A: `tangent run "install gin"` in a Go workspace (cwd contains `go.mod`)

| Stage | Value |
|---|---|
| input | `"install gin"` |
| normalized | `"install gin"` |
| intent | `IntentInstallPackage` |
| words | `["install","gin"]` → Extract → `{"package":"gin"}` |
| workspace | `WorkspaceGo` |
| rule match | `GoInstallPackageRule` (Match: IntentInstallPackage + WorkspaceGo) |
| generated | `Command{"go get gin", "Installs Go module 'gin'", true}` |
| Preview → Ask (Safe:true) → Execute → child `cmd /C "go get gin"` runs |

### Example B: `tangent run "git status"` in any folder

| Stage | Value |
|---|---|
| normalized | `"git status"` |
| intent | `IntentGitStatus` |
| workspace | whatever (ignored by rule) |
| rule | `GitStatusRule` |
| cmd | `Command{"git status", "Displays the current Git repository status.", true}` |

### Example C: `tangent run "install flask"` in Python-only workspace (pyproject.toml exists)

| Stage | Value |
|---|---|
| intent | `IntentInstallPackage` |
| workspace | `WorkspacePython` |
| rule search order | GitInit(no), GitStatus(no), GitCreateBranch(no), GitPush(no), GoInstall(no), NodeInstall(no), NodeRun(no), CreateVenv(no) |
| result | **fallback** — `Command{"", "No matching rule found.", false}` |
| UI effect | Preview shows empty command, Explanation shows no-match msg, Safe=false → ⚠ WARNING prompt |

### Example D: `tangent run "create virtual environment"` in Python workspace (known issue)

| Stage | Value |
|---|---|
| normalized | `"create virtual environment"` |
| intent | `IntentCreateVenv` (matched via "virtual environment") |
| Extract(IntentCreateVenv) | **no case exists** in entities.go → `{}` empty |
| rule match | `CreateVenvRule` matches |
| Generate | `venv = args["venv"]` = empty string → `Command:"python -m venv "` trailing empty name |
| Current runtime behaviour on Windows: `cmd /C "python -m venv "` — error output from python about missing venv name. |

**Defect / known gap:** To fix, add this case to `Extract()` in entities.go under a new case for IntentCreateVenv.

---

## 10. When Adding Code — Checklist (used by human or AI)

1. Does the new feature require a new type? Add to `types/types.go` ONLY.
2. Does the feature require new input processing logic? Put it in `core/*` — NEVER in `cmd/` or `ui/`.
3. Does the feature introduce a new shell-command transformation? Add a rule struct + `init()` register in the workspace-appropriate file under `internal/rules/<workspace>.go`.
4. Does the feature introduce new terminal output / prompt? Add a new exported function (or method) in `internal/ui/*`. NEVER print to fmt.* directly inside `core/*` or `rules/*`.
5. New intent? 4 files MUST change:
   - `types/types.go` + const
   - `core/intents.go` + resolver case
   - `core/entities.go` + extract case (if intent needs args)
   - at least ONE rule file with a matching rule + init register
   - docs (SRS/LLD) updated
6. New workspace? 3+ files MUST change:
   - `types/types.go` + const
   - `core/workspace.go` + probe in precedence list (above unknown)
   - at least ONE rule file using the workspace
   - docs updated
7. Build + vet: `go build . ; go vet ./...` → MUST exit 0 before commit.

---

## 11. Eraser.io Diagram Prompts

**Call graph (LLD level):**
```
main.main -> cmd.Execute -> runCmd.Run_closure
runCmd.Run_closure -> core.Process -> [core.Parse, core.Resolve, core.Extract, core.Detect]
core.Process -> rules.Generate
rules.Generate -> rules.Rules -> for_each -> rule.Match -> rule.Generate
runCmd.Run_closure -> ui.Show
runCmd.Run_closure -> ui.Ask
runCmd.Run_closure -> ui.Execute -> os/exec.Command child
```

**Package dependency graph (DAG):**
```
main -> cmd
cmd -> core
cmd -> ui
core -> rules
core -> types
rules -> types
ui -> types
```

---

## 12. Full Import List Per File (for reference)

Concrete `import (...)` blocks used right now. If your code change diverges, verify this list is updated.

| File | Imports |
|---|---|
| main.go | `cmd` |
| cmd/root.go | `os`, `cobra` |
| cmd/run.go | `fmt`, `strings`, `internal/core`, `internal/ui`, `cobra` |
| types/types.go | `time` |
| core/parser.go | `regexp`, `strings` |
| core/intents.go | `strings`, `types` |
| core/entities.go | `strings`, `types` |
| core/workspace.go | `os`, `types` |
| core/pipeline.go | `rules`, `types` |
| rules/rule.go | `types` |
| rules/registry.go | (none) |
| rules/engine.go | `types` |
| rules/git.go | `types` |
| rules/go.go | `types` |
| rules/node.go | `types` |
| rules/python.go | `types` |
| ui/preview.go | `fmt`, `types` |
| ui/confirm.go | `bufio`, `fmt`, `os`, `strings`, `types` |
| ui/executor.go | `bytes`, `os/exec`, `time`, `types` |
