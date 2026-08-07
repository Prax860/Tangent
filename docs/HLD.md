# High-Level Design — Tangent

**Version:** 1.0
**Date:** 2026-08-07

The HLD answers: *What are the pieces? How do they talk to each other? What data flows between them?* It is intentionally file-detail light — see the LLD for the function-by-function reference.

---

## 1. System Context

Tangent is a **single-process, single-binary, terminal-first tool** built in Go. It accepts free-form natural language from the command line, locally transforms it into a deterministic shell command tailored to the current working directory, shows a preview, asks for confirmation, and optionally executes it while capturing the child-process result.

```
                                 ┌──────────────────────────────────────────────────────┐
                                 │                  Tangent (one process)               │
                                 │                                                      │
   CLI argv ─────────────────────▶  cmd (Cobra)                                         │
                                 │    │                                                 │
                                 │    ▼                                                 │
                                 │  core pipeline ──────▶  rules engine                │
                                 │    │        ▲                                        │
                                 │    │        │ types (shared domain models)           │
                                 │    ▼        │                                        │
                                 │  ui layer ◀─┘                                        │
                                 │    ├─ preview (stdout)                               │
                                 │    ├─ confirm (stdin/stdout)                         │
                                 │    └─ executor ────▶ child process via cmd.exe /C   │
                                 └──────────────────────────────────────────────────────┘
```

There are **no network calls, no databases, no background threads, no event loops**. Everything happens synchronously in the main goroutine from process start to process exit.

---

## 2. Layered Architecture (the Dependency DAG)

The system is organized as 5 source-code packages + 1 entrypoint package. The import graph is a strict DAG. Arrow direction = "imports" (i.e. depends on).

```
                    ┌──────────┐
                    │   main   │  (package main)
                    └────┬─────┘
                         │ imports
                         ▼
                   ┌───────────┐
                   │    cmd    │  (cobra CLI wiring)
                   └──┬───┬────┘
            imports   │   │   imports
        ┌─────────────┘   └──────────────┐
        ▼                                  ▼
   ┌─────────┐                       ┌─────────┐
   │  core   │                       │   ui    │
   └────┬────┘                       └────┬────┘
        │  imports                         │
        ▼                                  │ imports
   ┌─────────┐                             │
   │  rules  │                             │
   └────┬────┘                             │
        │                                  │
        │         imports                  │
        └──────────────┬───────────────────┘
                       ▼
                  ┌─────────┐
                  │  types  │   (zero internal imports)
                  └─────────┘
```

**Package Responsibilities:**

| Package | What it owns | What it must NOT do |
|---|---|---|
| `main` | Entry point. Only `cmd.Execute()`. | Any business logic. |
| `cmd` | Cobra command registration, argument joining, orchestrating the ui.Show → ui.Ask → ui.Execute sequence. | Must NOT contain parsing, intent resolution, or rule logic. |
| `core` | The pure NLP-ish pipeline: parse, resolve intent, extract entities, detect workspace, stitch together via `Process`. Imports `rules` to get a `Command`. | Must NOT touch stdin/stdout/stderr (workspace detection uses `os.Stat` only). |
| `rules` | Rule interface, registry, match-first engine, and all concrete rules (git, go, node, python). | Must NOT do terminal I/O or process execution. |
| `ui` | All terminal I/O: preview printing, confirmation prompt, subprocess spawning + result capture. | Must NOT know about intents, parsers, workspaces. Only knows `Command`, `Response`, `ExecutionResult`. |
| `types` | All shared domain enums + structs. A leaf package. | Must NOT import any other tangent package. |

---

## 3. Component Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                 cmd                                      │
│                                                                          │
│   rootCmd (cobra root) ─────── addCommand ─────── runCmd                 │
│                                                  │                       │
│                                                  ▼                       │
│                              strings.Join(args) → input                  │
│                                                  │                       │
│    ┌─────────────────────────────────────────────┤                       │
│    │  ui layer                                   ▼                       │
│    │                                                                    │
│    │  ┌──────────┐   ┌──────────┐   ┌──────────────┐                    │
│    │  │ ui.Show  │   │  ui.Ask  │   │ ui.Execute   │                    │
│    │  │(preview) │   │(confirm) │   │(child proc)  │                    │
│    │  └────▲─────┘   └────▲─────┘   └──────▲───────┘                    │
│    └───────┼──────────────┼────────────────┼────────────────────────────┘
│            │Response      │Command         │Command                      │
│    ┌───────┴──────────────┴────────────────┴───────┐                     │
│    │              core pipeline                     │                     │
│    │                                                │                     │
│    │   core.Process(input) → types.Response        │                     │
│    │        │                                       │                     │
│    │        ▼                                       │                     │
│    │   parser.Parse ──────────► normalized          │                     │
│    │        │                                       │                     │
│    │        ▼                                       │                     │
│    │   intents.Resolve ───────► intent              │                     │
│    │        │                                       │                     │
│    │        ▼                                       │                     │
│    │   entities.Extract ──────► arguments           │                     │
│    │                                                │                     │
│    │   workspace.Detect ───────► workspaceType      │                     │
│    │        │                                       │                     │
│    │        ▼                                       │                     │
│    │   rules.Generate(intent, ws, args) ──► Command │                     │
│    │        ▲                                       │                     │
│    └────────┼───────────────────────────────────────┘                     │
│             │                                                             │
│    ┌────────┴───────────────────────────────────┐                        │
│    │             rules engine                    │                        │
│    │                                             │                        │
│    │   Rule interface { Name, Match, Generate }  │                        │
│    │        ▲                                    │                        │
│    │        │ registry (slice)                   │                        │
│    │        │                                    │                        │
│    │   GitInitRule  GitStatusRule  GitCreateBranchRule  GitPushRule      │
│    │   GoInstallPackageRule                                                  │
│    │   NodeInstallPackageRule  NodeRunProjectRule                         │
│    │   CreateVenvRule                                                        │
│    └─────────────────────────────────────────────┘                        │
│                                                                           │
│    ┌──────────────────────────────────────────────────────────────┐       │
│    │                     types (domain model)                     │       │
│    │                                                              │       │
│    │  WorkspaceType   IntentType   Command   Request              │       │
│    │                    Response    ExecutionResult               │       │
│    └──────────────────────────────────────────────────────────────┘       │
└───────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Data Flow (end-to-end, stepwise)

### 4.1 Control + Data Flow

1. **argv → cmd.runCmd.Run**
   - Input: `["run", "install", "gin"]`
   - Transform: `strings.Join` → `"install gin"`

2. **input → core.Process**
   Returns a `types.Response`, built internally via 5 sub-steps:

   | Step | Func | Input | Output |
   |---|---|---|---|
   | 2a | `core.Parse` | `string` input | `string` normalized |
   | 2b | `core.Resolve` | normalized | `types.IntentType` |
   | 2c | `core.Extract` | normalized + intent | `map[string]string` |
   | 2d | `core.Detect` | (none; reads filesystem) | `types.WorkspaceType` |
   | 2e | `rules.Generate` | intent + workspace + arguments | `types.Command` |
   | 2f | Assemble | all above | `types.Request` embedded in `types.Response` |

3. **Response → ui.Show**
   - Pure terminal side-effect. No return value. Prints headers, request, intent, workspace, args, generated command, explanation, safety flag.

4. **Response.Command → ui.Ask**
   - Returns `bool` = user said "y" / "yes"?
   - If `!Command.Safe`: prepends a ⚠ WARNING block.

5. **(If confirmed) Response.Command → ui.Execute**
   - Spawns: `cmd.exe /C "<Command.Command>"`.
   - Captures: stdout bytes, stderr bytes, exit code, wall-clock duration.
   - Returns: `types.ExecutionResult`.

6. **ExecutionResult → cmd.runCmd.Run (print)**
   - Prints Execution Result header block + Exit Code + STDOUT section (if non-empty) + STDERR section (if non-empty) + Duration.

### 4.2 Data-Object Lifecycle (what carries what)

```
                ┌─────────────────────┐
                │ string input        │  raw user words
                └────────┬────────────┘
                         │
                ┌────────▼────────────┐
                │ string normalized   │  lowercase, trimmed, punctuation stripped
                └────────┬────────────┘
                         │
          ┌──────────────┼──────────────────┐
          ▼              ▼                  ▼
   IntentType      WorkspaceType    map[string]string
   (intent)        (workspace)      (arguments)
          │              │                  │
          └──────────────┼──────────────────┘
                         ▼
               ┌────────────────────┐
               │ types.Command      │   Command string, Explanation, Safe bool
               └────────┬───────────┘
                        │
         ┌──────────────┴──────────────────┐
         ▼                                 ▼
  ┌──────────────┐              ┌────────────────────────┐
  │types.Request │              │types.ExecutionResult   │
  │RawInput      │              │Command string          │
  │Normalized    │              │Stdout / Stderr strings │
  │Intent        │              │ExitCode int            │
  │Workspace     │              │Duration Duration       │
  │Arguments     │              └────────────────────────┘
  └──────┬───────┘
         ▼
  ┌──────────────┐
  │types.Response│  Request + Command
  └──────────────┘
```

---

## 5. Key Interfaces and Contracts

### 5.1 `Rule` Interface (the plug-in boundary)

```
Rule:
  Name() string                                          → unique stable id, "domain.action"
  Match(intent IntentType, workspace WorkspaceType) bool → selection predicate
  Generate(arguments map[string]string) Command          → builds the shell command
```

**Invariants (hard contracts):**
- `Match` MUST be pure. It may not mutate state, touch the filesystem, read environment, or be non-deterministic.
- `Generate` MAY perform string interpolation of `arguments`, but MUST NOT execute shell commands or perform I/O.
- Registration of a rule is done via `rules.Register(r)` inside `func init()` in the same file that defines the concrete struct.

### 5.2 Public API per Package

Treat the following exported symbols as the *only* public surface. AI agents adding code MUST NOT import or rely on any unexported names outside their package.

| Package | Public (exported) API |
|---|---|
| `core` | `func Parse(string) string`<br>`func Resolve(string) IntentType`<br>`func Extract(string, IntentType) map[string]string`<br>`func Detect() WorkspaceType`<br>`func Process(string) Response` |
| `rules` | `type Rule interface{...}`<br>`func Register(Rule)`<br>`func Rules() []Rule`<br>`func Generate(IntentType, WorkspaceType, map[string]string) Command` |
| `ui` | `func Show(Response)`<br>`func Ask(Command) bool`<br>`func Execute(Command) ExecutionResult` |
| `types` | ALL types/consts in `types.go` are public. |
| `cmd` | Only `func Execute()`. The `rootCmd` and `runCmd` are unexported. |

---

## 6. Rule Engine Pattern

Tangent uses a **registry + strategy + first-match** pattern.

```
registry (slice of Rule)
   │
   ├─ GitInitRule{}
   ├─ GitStatusRule{}
   ├─ GitCreateBranchRule{}
   ├─ GitPushRule{}
   ├─ GoInstallPackageRule{}
   ├─ NodeInstallPackageRule{}
   ├─ NodeRunProjectRule{}
   └─ CreateVenvRule{}

rules.Generate(intent, workspace, arguments):
    for each rule in registry:
        if rule.Match(intent, workspace):
            return rule.Generate(arguments)
    return fallback(empty, no-match, unsafe)
```

Why first-match? It gives us predictable ordering and lets us write higher-specificity rules (e.g. workspace-specific package install) that simply beat more generic ones, *provided* the registration order places them first. Today the registry order is the lexical file order of the `init()` functions inside the `rules` package. See LLD §rules for exact ordering.

---

## 7. Deployment View (as-built)

Tangent has no separate deployment phase; the build IS the deployable.

```
go build -o tangent.exe .
  │
  ├─ statically linked single PE (Windows .exe)
  │    - embeds go stdlib + cobra + pflag + mousetrap + all tangent packages
  │
  ▼
run on any Windows machine with a compatible CPU:
    tangent.exe run "install gin"
```

**Files touched at runtime (by Tangent itself):**
- Reads via `os.Stat` for workspace detection.
- Reads `/dev/stdin` (actually Windows stdin handle) for the confirmation prompt line.
- The child process spawned by the executor may read/write arbitrary files in the CWD — that's the point of the tool and is outside Tangent's direct control.

**Network access:** None by Tangent. Child processes (e.g. `go get`, `npm install`, `git push`) WILL access the network — that is their own behaviour.

---

## 8. Extension Points — How to Add Things (the 80% cases)

### 8.1 Add a new rule for an existing intent + workspace combo
Example: Python install package (pip install \<pkg\>).
- Add a `PythonInstallPackageRule` struct in `internal/rules/python.go`.
- Implement `Name`, `Match(IntentInstallPackage, WorkspacePython)`, `Generate(arguments["package"] → "pip install " + p)`.
- Add `func init() { Register(PythonInstallPackageRule{}) }`.
- **That's it.** No other file touches.

### 8.2 Add a new intent
Example: IntentUninstallPackage.
1. Append `IntentUninstallPackage IntentType = "uninstall_package"` to `types.go`.
2. Add a `case` in `core/intents.go` resolver (e.g. `strings.HasPrefix(input, "uninstall")`).
3. If intent needs entities, add an `IntentUninstallPackage` case to `core/entities.go`.
4. Add rules that `Match(IntentUninstallPackage, *)` in their respective workspace files.
5. Update SRS + LLD.

### 8.3 Add a new workspace
Example: WorkspaceRuby.
1. Append `WorkspaceRuby = "ruby"` to `types.go`.
2. Add probe `if fileExists("Gemfile") { return WorkspaceRuby }` in the precedence list `core/workspace.go` BEFORE `WorkspaceUnknown`.
3. Add `Match(*, WorkspaceRuby)` rules.

### 8.4 Add a new UI step (e.g. dry-run flag, logging to file)
- Add the Cobra flag to `cmd/run.go`.
- Conditionally call a new public function in `ui/` (do NOT inline terminal logic in `cmd`).

### 8.5 Port executor to Linux / macOS
- Keep one file `executor.go` with build tags (`//go:build windows` and `//go:build !windows`), OR
- Switch `exec.Command("cmd", "/C", ...)` to `exec.Command("sh", "-c", ...)` via runtime.GOOS switch.
- The rest of `core`, `rules`, `types`, `cmd` is fully portable as-is.

---

## 9. Failure Modes and Mitigations

| Failure | Where it surfaces | Current behaviour | Future improvement |
|---|---|---|---|
| No rule matched | rules.Generate | Fallback returned with Safe=false, ⚠ prompt shown. | Surface a more helpful "you could add a rule for X" hint. |
| Child process fails | ui.Execute | Exit code captured, STDERR printed. | Optionally colorize STDERR. |
| Confirmation read error / EOF | ui.Ask | Returns false → cancelled. | Explicit "cancelled due to input error" msg. |
| Empty command string (buggy rule) | cmd/run.go → ui.Execute | `cmd /C ""` runs and generally exits 0 on Windows. | Add guard `if command.Command == "" { skip execution, print error }`. |
| Wrong workspace detected | workspace.Detect | Wrong rule matched. | Allow `--workspace=X` CLI override. |
| Intent mis-resolved | intents.Resolve | Wrong rule matched. | Prompt user disambiguation UI for ambiguous cases. |

---

## 10. Eraser.io Diagram Prompts

**Component/architecture diagram:**
```
Rectangle main [main.go]
Rectangle cmd [cmd\nroot.go + run.go]
Group core [
  parser.Parse
  intents.Resolve
  entities.Extract
  workspace.Detect
  pipeline.Process
]
Group rules [
  Rule interface
  Registry
  engine.Generate
  git rules (x4)
  go rules (x1)
  node rules (x2)
  python rules (x1)
]
Group ui [
  ui.Show
  ui.Ask
  ui.Execute
]
Group types [
  WorkspaceType
  IntentType
  Command
  Request
  Response
  ExecutionResult
]

main -> cmd
cmd -> core
cmd -> ui
core -> rules
core -> types
rules -> types
ui -> types
```

**Pipeline dataflow (end-to-end):**
```
Input["install gin"]
   -> Parse["parse: lowercase/trim/collapse/strip"]
   -> Resolve["resolve: IntentInstallPackage"]
   -> Extract["extract: {package:gin}"]
   -> Detect["detect: WorkspaceGo (reads go.mod)"]
   -> Generate["rules.Generate: match GoInstallPackageRule -> go get gin"]
   -> Response
   -> Preview["ui.Show prints"]
   -> Confirm["ui.Ask y/N"]
   -> Execute["spawn cmd /C -> ExecutionResult"]
```
