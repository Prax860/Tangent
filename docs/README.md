# Tangent

**Tangent** is a command-line tool that translates natural language instructions into executable developer commands. A user writes a plain-English request like `"install gin"` or `"create branch feature/login"`, and Tangent detects the project's workspace (Node.js, Python, Go, Rust, Java, or generic), resolves the intent, extracts arguments, and generates the correct shell command — with a preview and confirmation step before execution.

It is written in **Go 1.25.5** and uses [Cobra](https://github.com/spf13/cobra) for its CLI scaffolding.

---

## 1. Project Layout

```
Tangent/
├── main.go                  # Entry point: calls cmd.Execute()
├── cmd/                     # Cobra command definitions
│   ├── root.go              # Root command "tangent"
│   └── run.go               # Sub-command "tangent run <instruction>"
├── internal/
│   ├── core/                # NLP pipeline (pure logic, no I/O except workspace detection)
│   │   ├── parser.go        # Text normalization
│   │   ├── intents.go       # Intent resolution from normalized text
│   │   ├── entities.go      # Argument/entity extraction per intent
│   │   ├── workspace.go     # Workspace type detection (via file probes)
│   │   └── pipeline.go      # Orchestrates: parse → resolve → extract → detect → generate
│   ├── rules/               # Rule engine for command generation
│   │   ├── rule.go          # Rule interface
│   │   ├── registry.go      # Slice-based registry & init() auto-registration
│   │   ├── engine.go        # Match loop → first-matching rule generates command
│   │   ├── git.go           # 4 Git rules: init / status / create-branch / push
│   │   ├── go.go            # 1 Go rule: install package
│   │   ├── node.go          # 2 Node rules: install package / run project
│   │   └── python.go        # 1 Python rule: create virtual env
│   ├── ui/                  # User interaction (terminal I/O + process execution)
│   │   ├── preview.go       # Formatted pretty-print of Response
│   │   ├── confirm.go       # y/N prompt with warning for unsafe commands
│   │   └── executor.go      # Spawns subprocess, captures stdout/stderr/exit/duration
│   └── types/               # Shared domain types (zero-dependency, imported by all layers)
│       └── types.go         # WorkspaceType, IntentType, Request, Response, Command, ExecutionResult
├── go.mod
├── go.sum
└── docs/
    ├── README.md            # This file
    ├── SRS.md               # Software Requirements Specification
    ├── HLD.md               # High-Level Design
    └── LLD.md               # Low-Level Design
```

### Module path
```
module github.com/prax860/tangent
```
All internal import paths use this prefix (e.g. `github.com/prax860/tangent/internal/core`).

---

## 2. Quick Start

### Prerequisites
- Go 1.25.5 or later
- Windows (the subprocess executor currently shells out to `cmd /C`; see LLD §executor for Linux/macOS extension notes)

### Build
```powershell
go build -o tangent.exe .
```

### Run
```powershell
# Example: in a Go project
.\tangent.exe run "install gin"

# Example: in a Node.js project
.\tangent.exe run "install express"

# Example: in any git repo
.\tangent.exe run "create branch feature/login"
.\tangent.exe run "git status"
.\tangent.exe run "git init"
.\tangent.exe run "push"

# Example: in a Python project
.\tangent.exe run "create virtual environment"
```

The flow for every invocation is:
1. `core.Process(input)` → produces a `types.Response`
2. `ui.Show(response)` → prints a formatted preview
3. `ui.Ask(command)` → prompts y/N (warns if `command.Safe == false`)
4. If yes → `ui.Execute(command)` → prints exit code, stdout, stderr, duration

---

## 3. How It Works — 30-Second Overview

```
user input "install gin"
        │
        ▼
  core.Parse        → "install gin" (lowercased, trimmed, deduped spaces, punctuation stripped)
        │
        ▼
  core.Resolve      → IntentInstallPackage
        │
        ▼
  core.Extract      → {"package": "gin"}
        │
        ▼
  core.Detect       → WorkspaceGo (because go.mod exists in cwd)
        │
        ▼
  rules.Generate    → finds GoInstallPackageRule (matches IntentInstallPackage + WorkspaceGo)
                     returns Command{Command:"go get gin", Safe:true, Explanation:"..."}
        │
        ▼
  build Response(Request, Command)
        │
        ├── ui.Show         → pretty-printed preview to stdout
        ├── ui.Ask          → y/N prompt
        └── ui.Execute      → cmd /C "go get gin"  → ExecutionResult
```

---

## 4. Key Design Decisions (Relevant for Anyone Extending the Code)

### 4.1 Layered Architecture, No Circular Imports
```
types (zero deps)
  ▲   ▲   ▲
  │   │   │
core  rules  ui
  │         │
  └────┬────┘
       ▼
      cmd
       ▲
       │
     main
```
- `types` is imported by everything; it imports nothing internal.
- `core` imports `rules` and `types`; it does NOT import `ui`.
- `rules` imports only `types`.
- `ui` imports only `types`.
- `cmd` imports `core` + `ui`.
- `main` imports only `cmd`.

**Rationale:** This means pipeline logic is pure and unit-testable without touching a terminal. The rule engine is self-contained. UI concerns (terminal, subprocesses) are isolated from business logic.

### 4.2 Rule Engine Uses `init()` Auto-Registration
Each rule file (git.go, go.go, node.go, python.go) has one or more concrete `Rule` structs and an `func init() { Register(SomeRule{}) }`. Because Go runs all `init()` functions of imported packages, merely importing `rules` from `core/pipeline.go` is sufficient to populate the registry.

**To add a new rule:**
1. Pick (or create) the workspace-specific file under `internal/rules/`
2. Define a struct implementing the three methods of `Rule`
3. Add an `init()` that calls `Register(YourStruct{})`
4. No other file needs to be touched.

### 4.3 First-Matching-Rule Wins
`rules.Generate` iterates in registry order and returns the first rule whose `Match()` returns true. Order is the order of `init()` calls within the rules package (file-alphabetical, but treat as "registration order"). If specificity matters, put narrower matches before broader matches in the registration.

### 4.4 Safe vs. Unsafe Commands
Every generated `Command` has a `Safe bool` flag:
- `Safe = true` → confirmation prompt says `"Execute command? (y/N): "`
- `Safe = false` → confirmation prompt prints a **⚠ WARNING** block with the explanation, then `"Continue anyway? (y/N): "`

The "no matching rule" fallback returns `Safe: false` defensively.

---

## 5. Current Capability Matrix

Supported intents × workspaces (i.e., rules that exist today):

| Intent | Node | Python | Go | Rust | Java | Workspace-agnostic |
|---|---|---|---|---|---|---|
| Install package | ✅ `npm install <pkg>` | – | ✅ `go get <module>` | – | – | – |
| Run project | ✅ `npm run <file>` | – | – | – | – | – |
| Create venv | – | ✅ `python -m venv <venv>` | – | – | – | – |
| Git init | – | – | – | – | – | ✅ `git init` |
| Git status | – | – | – | – | – | ✅ `git status` |
| Create branch | – | – | – | – | – | ✅ `git checkout -b <branch>` |
| Push | – | – | – | – | – | ✅ `git push` |

**Anything not in this matrix** (e.g. install package in a Python project) → `Explanation: "No matching rule found."`, `Safe: false`, `Command: ""`.

---

## 6. Where to Find Things in This Docs Folder

| File | Who reads it | Purpose |
|---|---|---|
| `README.md` | Everyone | Project overview, quick start, layout, key design choices, capability matrix |
| `SRS.md` | Product, stakeholders, future AI agents adding features | Requirements, use cases, functional/non-functional spec, acceptance criteria |
| `HLD.md` | Architects, AI agents planning new modules | Components, data/control flow, interfaces, module responsibilities, deployment view |
| `LLD.md` | Developers, AI agents writing/changing code | File-by-file + function-by-function reference: signatures, types, logic, invariants, call graphs, edge cases |

---

## 7. Eraser.io Diagram Prompts

(Paste these into [eraser.io](https://eraser.io) diagram-from-text to regenerate the visuals referenced in HLD/LLD.)

**Architecture (layers + dependencies):**
```
Rectangle main [main.go]
Rectangle cmd [cmd\nroot.go + run.go]
Group core (core.Parse, core.Resolve, core.Extract, core.Detect, core.Process)
Group rules (Rule interface, Registry, engine, git, go, node, python)
Group ui (ui.Show, ui.Ask, ui.Execute)
Group types (WorkspaceType, IntentType, Command, Request, Response, ExecutionResult)

main -> cmd
cmd -> core
cmd -> ui
core -> rules
core -> types
rules -> types
ui -> types
```

**Data flow for `tangent run "install gin"`:**
```
Input ["install gin"] -> Parse -> Resolve -> Extract -> Detect -> rules.Generate -> Response -> Preview -> Confirm -> Execute -> ExecutionResult
```
