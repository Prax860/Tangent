# Software Requirements Specification — Tangent

**Version:** 1.0
**Date:** 2026-08-07
**Module path:** `github.com/prax860/tangent`

---

## 1. Introduction

### 1.1 Purpose
Tangent is a developer CLI that converts **natural language instructions** into concrete shell commands tailored to the current project's workspace (Node, Python, Go, Rust, Java). It shields the developer from having to remember every tool's exact syntax, while still requiring explicit user confirmation before executing anything (a deliberate safety net against accidental destructive operations.

This SRS defines all functional and non-functional requirements, use cases, and acceptance criteria for the system. It is intended for:
- Human developers extending Tangent
- AI coding agents (e.g. Claude) auto-generating Tangent features, who must treat this document as the source of truth for behaviour.

### 1.2 Scope — In Scope
- Single binary built via `go build`.
- One CLI sub-command: `tangent run <free-form natural language instruction>`.
- Text-normalization pipeline (parse → resolve intent → extract entities → detect workspace → generate command).
- Rule-based command generator with pluggable rules.
- Terminal preview, yes/no confirmation, and subprocess execution with stdout/stderr/exit-code capture.
- Support for the intents and workspaces listed in §3.

### 1.3 Scope — Out of Scope (v1)
- Persistent config / config files / `--config`.
- Authentication, multi-tenancy, remote execution.
- Large language model calls. All resolution is local and rule-based.
- Interactive / multi-turn REPL. The tool runs one instruction and exits.
- Linux / macOS shells (executor hardcodes `cmd /C` on Windows; cross-platform behaviour is an extension, not v1.

### 1.4 Definitions
| Term | Definition |
|---|---|
| Instruction | The raw string the user types after `tangent run`. |
| Normalized input | The instruction after lowercasing, trimming, punctuation stripping, and space collapse. |
| Intent | An enum-like string identifying the user's goal (see `types.IntentType`). |
| Entity / Argument | A named string value extracted from the input (e.g. `"package": "gin"`). |
| Workspace | The project type of the current working directory (see `types.WorkspaceType`). |
| Rule | A stateless struct implementing the `Rule` interface: `Name`, `Match`, `Generate`. |
| Safe command | `Command.Safe == true`; shown with a standard confirmation prompt. |
| Unsafe command | `Command.Safe == false`; shown with a ⚠ WARNING block. |
| Response | The immutable result of `core.Process` containing `Request + Command`. |

---

## 2. Overall Description

### 2.1 Product Perspective
Tangent is a **standalone console application**. It has no network dependencies at runtime. Its only third-party runtime dependency is the Cobra library for CLI parsing. It runs within the user's current shell, in the current working directory.

```
┌──────────────┐   argv[1..]   ┌──────────────────┐   stdin/stdout   ┌──────────┐
│  User shell  │ ───────────▶ │  tangent binary   │ ───────────────▶ │ Terminal │
│ (powershell) │              │  (Go exe)         │ ◀───────────────│ (y/N)    │
└──────────────┘               └────────┬─────────┘                  └──────────┘
                                        │ cmd /C <generated>
                                        ▼
                               ┌──────────────────┐
                               │ Child process    │
                               │ (go/npm/git/etc)  │
                               └──────────────────┘
```

### 2.2 User Characteristics
- Primary users are software developers comfortable with a terminal.
- Secondary users include AI coding agents invoking Tangent programmatically from shells (they must parse `y`, `yes`, or nothing from the confirmation prompt accordingly; note: non-interactive/`--yes` flag is not yet in scope).
- Users are expected to understand the semantics of shell commands Tangent generates; Tangent only *translates*, it does not reason about correctness beyond intent matching.

### 2.3 Operating Environment
| Aspect | Requirement |
|---|---|
| OS | Windows 10+ (initial). POSIX-compatible layer is optional future work. |
| Go runtime | 1.25.5 minimum as specified in `go.mod`. |
| Working directory | Any; workspace detection probes relative to `os.Getwd()`. |
| Shell for child exec | `cmd.exe` via `os/exec.Command("cmd", "/C", ...)`. |

### 2.4 Constraints
1. **No user-interruption beyond the explicit confirmation prompt.** The tool MUST NOT call `Read-Host`, `pause`, `Get-Credential`, or otherwise block on stdin in any place other than `ui.Ask`.
2. **No secrets in code.** No API keys, tokens, or passwords.
3. **Single-pass pipeline.** `core.Process` is pure and deterministic given identical inputs and identical working directory contents. This is a testability requirement.
4. **Zero mutable global state outside of `rules.registry`.** The registry is append-only via `Register` during `init()`.

### 2.5 Assumptions and Dependencies
- Git is installed and on `PATH` for git commands to actually execute successfully (Tangent does not check; it simply spawns the process and forwards exit codes).
- `go`, `npm`, `python` are on `PATH` for workspace-command execution to succeed in the respective workspace types.
- File existence check `os.Stat` is an acceptable proxy for "this directory contains a project of type X". This can have false positives (a stray `go.mod` in a monorepo root). That is accepted behaviour.

---

## 3. Functional Requirements

Functional requirements are numbered `FR-<intent>-<N>` where intent codes are: `GEN` (general), `PAR` (parser), `RES` (resolver), `ENT` (entity extractor), `WS` (workspace detector), `RUL` (rules engine), `UI` (preview/confirm/executor).

### 3.1 General (GEN)

**FR-GEN-1: Required argument check.** If `tangent run` is invoked with zero positional arguments, the tool MUST print an error line starting with `❌ Please provide a command.` and an example line `Example: tangent run "install gin"`, then exit with code 0 (Cobra's default for a handler returning nil).

**FR-GEN-2: Pipeline call.** On non-empty args, the tool MUST `strings.Join(args, " ")` into a single input string and call `core.Process(input)` exactly once.

**FR-GEN-3: Preview-then-confirm-then-execute ordering.** The tool MUST, in this order: (a) call `ui.Show(response)`; (b) call `ui.Ask(response.Command)`; (c) if and only if `ui.Ask` returns `true`, call `ui.Execute(response.Command)`. Any other ordering is non-conformant.

**FR-GEN-4: Execution result display.** After a successful execution, the tool MUST print a section titled "Execution Result" between two `━━━` header lines, then print `Exit Code : N`, then non-empty STDOUT section, non-empty STDERR section, and `Duration : <Duration`.

### 3.2 Parser (PAR)

**FR-PAR-1: Lowercase.** The parser MUST map every ASCII uppercase letters to lowercase.

**FR-PAR-2: Trim.** The parser MUST trim leading and trailing Unicode whitespace.

**FR-PAR-3: Punctuation stripping.** The parser MUST strip every character that does not match `[^\w\s./-]`. Characters in regexp.

**FR-PAR-4: Space collapse.** Any run of whitespace characters (U+0020, tabs, newlines) MUST collapse into exactly one space (U+0020).

**FR-PAR-5: Determinism.** Two byte-for-byte identical inputs MUST produce identical outputs.

### 3.3 Intent Resolver (RES)

Resolver takes the result of PAR.

**FR-RES-1: IntentCreateVenv.** Normalized text containing ANY of `"virtual environment"`, `"venv"`, `"virtualenv"` → `IntentCreateVenv`.

**FR-RES-2: IntentInstallPackage.** Normalized text whose prefix `"install"` → IntentInstallPackage (prefix match, exact).

**FR-RES-3: IntentGitInit.** Contains exact substring `"git init"` → `IntentGitInit`.

**FR-RES-4: IntentGitStatus.** Contains exact substring `"git status"` → `IntentGitStatus`.

**FR-RES-5: IntentCreateBranch.** Contains either exact substring `"create branch"` OR `"new branch"` → `IntentCreateBranch`.

**FR-RES-6: IntentPushChanges.** Contains exact substring `"push"` → `IntentPushChanges`.

**FR-RES-7: Default/Unknown.** Any input not matching FR-RES-1..6 MUST resolve to `IntentUnknown`.

**FR-RES-8: Evaluation order is switch-case order as written in §7 Intent list.** Earlier conditions win.

### 3.4 Entity Extractor (ENT)

**FR-ENT-1: General signature. Extract takes normalized input and an intent type and returns `map[string]string`. Unknown intents return empty map.

**FR-ENT-2: IntentInstallPackage → `package`.** Word-tokenize input via `strings.Fields`. If `len(words) >= 2`, args `"package"` `words[1]`.

**FR-ENT-3: IntentCreateBranch → `branch`.** If `len(words) >= 3`, args `"branch"` `words[2:]`, space-joined.

**FR-ENT-4: IntentRunProject → `file`.** If `len(words) >= 2`, args `"file"` words[1]`.

**FR-ENT-5: Unknown keys absent.** The map MUST NOT contain keys other than those defined for that intent.

### 3.5 Workspace Detector (WS)

**FR-WS-1: Detection precedence (highest to lowest):**
1. `package.json` exists in cwd → `WorkspaceNode`
2. `pyproject.toml` OR `requirements.txt` exists → `WorkspacePython`
3. `go.mod` exists → `WorkspaceGo`
4. `Cargo.toml` exists → `WorkspaceRust`
5. `pom.xml` OR `build.gradle` exists → `WorkspaceJava`
6. None of the above → `WorkspaceUnknown`

**FR-WS-2: Order-dependent.** FR-WS-1 precedence MUST be implemented as written. If a project has both `package.json` and `go.mod`, the result is `WorkspaceNode`.

**FR-WS-3: File existence means `os.Stat(filename)` returning `nil`. No is nil`.

### 3.6 Rule Engine (RUL)

**FR-RUL-1: Rule interface.** Every rule MUST expose three zero-argument Name() string  Match(IntentType, WorkspaceType) bool, Generate(arguments map[string]string) types.Command.**

**FR-RUL-2: Registry semantics.** `Register(r Rule) appends. `Rules()` returns the slice in the slice order.

**FR-RUL-3: Auto-populated by `init()`.** Each concrete rule MUST be automatically registered through an `func init()` in the same Go file that defines the rule.

**FR-RUL-4: Generate first-match `rules.Generate(intent, workspace, arguments iterates Rules() and returns the first rule R where `R.Match(intent, workspace) is true. If no rule matches, MUST return fallback Command{Command:"", Explanation:"No matching rule found.", Safe:false}. This fallback is non-optional.

### 3.6.1 Git Rules

**FR-RUL-GIT-1 GitInitRule.** Match: `IntentGitInit` (any workspace). Generate: `Command{Command:"git init", Explanation:"Initializes a new Git repository.", Safe:true}`.

**FR-RUL-GIT-2 GitStatusRule.** Match: `IntentGitStatus` (any workspace). Generate: `Command:"git status", Explanation:"Displays the current Git repository status.", Safe:true`.

**FR-RUL-GIT-3 GitCreateBranchRule.** Match: `IntentCreateBranch` (any workspace). Generate: branch `= arguments["branch"]`. Command: `"git checkout -b " + branch. Explanation: "Creates and switches to a new Git branch '<branch>'". Safe: true.

**FR-RUL-GIT-4 GitPushRule.** Match: `IntentPushChanges`. Generate: Command:"git push", Explanation:"Pushes commits to the remote repository.", Safe:true.

### 3.6.2 Go Rules

**FR-RUL-GO-1 GoInstallPackageRule.** Match: `IntentInstallPackage` AND `WorkspaceGo`. Generate module := arguments["package"]. Generate: Command "go get " + module. Explanation "Installs Go module '<module>'" Safe true.

### 3.6.3 Node Rules

**FR-RUL-NODE-1 NodeInstallPackageRule.** Match IntentInstallPackage AND WorkspaceNode. packageName := arguments["package"]. Command: "npm install " + packageName. Explanation: "Installs Node.js package '<packageName>'". Safe: true.

**FR-RUL-NODE-2 NodeRunProjectRule.** Match: IntentRunProject AND WorkspaceNode. file := arguments["file"]. Command: "npm run " + file. Explanation: "Starts the Node.js development server '<file>'". Safe: true.

### 3.6.4 Python Rules

**FR-RUL-PY-1 CreateVenvRule.** Match IntentCreateVenv AND WorkspacePython. venv := arguments["venv"]. Command: "python -m venv " + venv. Explanation: "Creates Python virtual environment '<venv>'". Safe: true.

### 3.7 UI Sub-module (UI)

**FR-UI-1 Preview.** ui.Show MUST print in this order:
- box header "━━━" lines, title "⚡ Tangent"
- Section "📝 Request" → raw input
- "🧠 Intent" → Request.Intent
- "📂 Workspace" → Request.Workspace
- "📦 Arguments" block if len(Arguments) > 0, each k/v pair as "k : v"
- "⚙ Generated Command" → Command.Command
- "💡 Explanation" → Command.Explanation
- "🔒 Safe :" followed by true/false

**FR-UI-2 Confirm (Safe path).** If Command.Safe == true → print blank line → "Execute command? (y/N): " → read stdin line → trim + lowercase → true iff input is exactly "y" or "yes". Any other input or EOF returns false.

**FR-UI-3 Confirm (Unsafe path).** If Command.Safe == false → print blank line → "⚠ WARNING" header → print Command.Explanation → blank line → "Continue anyway? (y/N): " → same acceptance logic as UI-2.

**FR-UI-4 Executor runs via cmd /C.** `exec.Command("cmd", "/C", command.Command)`.

**FR-UI-5 Captures buffers stdout and stderr into separate buffers and reads Stdout and Stderr fields.**

**FR-UI-6 Exit code.** If cmd.Run returns nil, ExitCode=0. If *exec.ExitError, ExitCode = e.ExitCode(). Any other error, ExitCode = -1.

**FR-UI-7 Duration wall-clock duration between time.Now() at start, time.Since(start) at return.

---

## 4. Non-Functional Requirements

### 4.1 Performance (PERF)
| ID | Requirement |
|---|---|
| NFR-PERF-1 | End-to-end latency for `tangent run "git status"` excluding subprocess execution MUST be < 100 ms on commodity Windows hardware. |
| NFR-PERF-2 | No heap allocations in hot path (parse/resolve/extract) beyond what is required for the returned strings and map. |
| NFR-PERF-3 | Rules iteration is O(R) where R = number of registered rules. There MUST NOT be nested loops over rules × rules. |

### 4.2 Reliability (REL)
| ID | Requirement |
|---|---|
| NFR-REL-1 | If any step panics, the panic is allowed to propagate (no generic recover at top of `main`); reproducibility of stack traces takes priority over graceful recovery. |
| NFR-REL-2 | If subprocess cannot start, `ui.Execute` MUST still return a valid ExecutionResult with ExitCode = -1 and Stderr/Stdout as empty. |
| NFR-REL-3 | Confirmation prompt read error returns false. |

### 4.3 Usability (USA)
| ID | Requirement |
|---|---|
| NFR-USA-1 | Every user-visible output line that is a header MUST be visually distinct. Box-drawing characters (`━`) are the standard header delimiter. |
| NFR-USA-2 | Every argument map entry printed via Show MUST be on its own "key : value" line. |
| NFR-USA-3 | Empty STDOUT/STDERR sections MUST be omitted from execution output to reduce noise. |

### 4.4 Maintainability (MAI)
| ID | Requirement |
|---|---|
| NFR-MAI-1 | Every new rule MUST be addable by modifying exactly 1 file (the file containing the rule struct + its init()). Changes to `engine.go`, `registry.go`, `cmd/run.go` MUST NOT be required. |
| NFR-MAI-2 | All type definitions live only in `internal/types/types.go`. No file may redefine WorkspaceType/IntentType/Command/Request/Response/ExecutionResult aliases. |
| NFR-MAI-3 | Compile-time interface compliance: `go vet ./...` MUST pass on every commit. |
| NFR-MAI-4 | Imports form a DAG (no cycles). This is enforced by Go; human reviewers must not introduce `core → cmd` back-edges. |

### 4.5 Portability (POR)
| ID | Requirement |
|---|---|
| NFR-POR-1 | Windows-only is acceptable for v1. HOWEVER: code that will need to change for POSIX MUST be localized to `internal/ui/executor.go` (the only file that knows about `cmd /C`). Parse/resolve/extract/detect/rules are fully portable as-written. |

### 4.6 Security (SEC)
| ID | Requirement |
|---|---|
| NFR-SEC-1 | No command executes without user confirmation. The default action on any non-"y"/non-"yes" answer is: print "❌ Command cancelled." and exit. |
| NFR-SEC-2 | The "no rule matched" case MUST return Safe:false. This ensures the user sees ⚠ WARNING before running anything generated from a future edge-case branch that mistakenly leaves Command populated but Explanation unclear. |
| NFR-SEC-3 | Tangent never shells out without quoting. Responsibility for safely concatenating strings into command strings lives in each Generate() rule. Current rules concatenate with simple spaces; rules that take arguments with spaces MUST add quoting explicitly in the implementation and update this SRS accordingly. |

---

## 5. Use Cases (UC)

### UC-1: Install Go Package in a Go Project
| Field | Value |
|---|---|
| Actor | Developer (human) |
| Precondition | Current working directory contains `go.mod`. |
| Trigger | `tangent run "install gin"` |
| Main success scenario | 1. Parse → `"install gin"` 2. Resolve → IntentInstallPackage 3. Extract → {"package":"gin"} 4. Detect → WorkspaceGo 5. Generate → `go get gin` (Safe:true) 6. Preview prints. 7. User types "y". 8. Execute runs `cmd /C "go get gin"`. 9. Execution result printed. Exit code from go get is shown as Exit Code. |
| Alternatives | 7a. User types "n". → "❌ Command cancelled." Process exits. |
| Post-conditions | Child process side effects (e.g. `gin` added to `go.mod/go.sum`) are whatever the child performed. Tangent writes nothing beyond that itself. |

### UC-2: Create Git Branch
| Field | Value |
|---|---|
| Precondition | CWD is a git repo (or not; detection of git checkout -b is workspace-agnostic. |
| Trigger | `tangent run "create branch feature/login"` |
| Main scenario | 1. Parse → "create branch feature/login" 2. Resolve → IntentCreateBranch 3. Extract → {"branch":"feature/login"} 4. Detect → any workspace (git rules are workspace-agnostic). 5. GitCreateBranchRule matches. 6. Preview shows `git checkout -b feature/login`. Safe:true. 7. "y" → executes, exit code 0. |

### UC-3: No Matching Rule
| Field | Value |
|---|---|
| Trigger | `tangent run "compile the whales"` (in any workspace) |
| Main scenario | 1. Resolve → IntentUnknown. 2. No rule matches IntentUnknown. 3. Fallback returned: Command="", Explanation="No matching rule found.", Safe=false. 4. Preview shows empty command, explanation, Safe:false. 5. ui.Ask shows ⚠ WARNING + explanation. |
| Alternative | If user still types "y" → cmd /C "" is executed (empty command → typically exit 0 on Windows). This is technically allowed by the executor; TODO: add explicit empty-command guard. |

### UC-4: Unsafe / Unknown Workspace Pairing
| Field | Value |
|---|---|
| Trigger | `tangent run "install flask"` inside a Python-only project (pyproject.toml present, no package.json, no go.mod). |
| Scenario | 1. Intent = InstallPackage. 2. Workspace = WorkspacePython. 3. No existing rule matches (Python-install rule absent). 4. Returns fallback Safe:false. 5. User warned, must explicitly accept. |

---

## 6. Acceptance Criteria (AC)

A build is considered compliant if ALL of the following hold:

1. **Compile.** `go build -o tangent.exe .` exits 0.
2. **Vet.** `go vet ./...` exits 0.
3. **Structural tests (pseudo-code; tests not part of v1 but the assertions are review-checklist):**
   - `core.Parse("  Install!!! ") == "install"`
   - `core.Resolve("create venv project") == IntentCreateVenv`
   - `core.Extract("install express", IntentInstallPackage)["package"] == "express"`
   - `core.Extract("create branch feature/a b", IntentCreateBranch)["branch"] == "feature/a b"`
   - With `go.mod` in cwd: `core.Detect() == WorkspaceGo`
   - `rules.Generate(IntentInstallPackage, WorkspaceGo, {"package":"gin"}).Command == "go get gin"`
   - `rules.Generate(IntentUnknown, WorkspaceUnknown, {}).Safe == false` AND `.Explanation == "No matching rule found."`
4. **End-to-end.** `tangent run "git status"` produces:
   - Request section "git status", Intent "git_status", Workspace detected
   - Command "git status", Safe true
   - Prompt "Execute command? (y/N): "
   - On "n" prints cancelled. On "y" runs git status and prints exit code 0 + stdout.

---

## 7. Intent Type Registry (Canonical List)

```go
IntentUnknown         = "unknown"
IntentCreateVenv      = "create_venv"
IntentInstallPackage  = "install_package"
IntentRunProject      = "run_project"
IntentGitInit         = "git_init"
IntentGitStatus       = "git_status"
IntentCreateBranch   = "create_branch"
IntentPushChanges   = "push_changes"
```

To add a new intent:
1. Append a const to `types.go` + resolver case in `intents.go` + at least 1 rule. Update SRS §3.3 and §7. Otherwise the intent is dead code.

---

## 8. Workspace Type Registry (Canonical List)

```go
WorkspaceUnknown = "unknown"
WorkspacePython  = "python"
WorkspaceNode    = "node"
WorkspaceGo      = "go"
WorkspaceRust    = "rust"
WorkspaceJava    = "java"
```

To add a new workspace (e.g. Ruby):
1. Add `WorkspaceRuby` const.
2. Add probe in `workspace.go` precedence list BEFORE `WorkspaceUnknown`.
3. Add rules that match the workspace.
4. Update SRS §3.5 and §8.

---

## Appendix A: Traceability Matrix (Requirements → Files)

| Requirement | Primary file(s) |
|---|---|
| FR-GEN-1..4 | cmd/run.go |
| FR-PAR-1..5 | internal/core/parser.go |
| FR-RES-1..8 | internal/core/intents.go |
| FR-ENT-1..5 | internal/core/entities.go |
| FR-WS-1..3 | internal/core/workspace.go |
| FR-RUL-1..4 | internal/rules/{rule,registry,engine}.go |
| FR-RUL-GIT-1..4 | internal/rules/git.go |
| FR-RUL-GO-1 | internal/rules/go.go |
| FR-RUL-NODE-1,2 | internal/rules/node.go |
| FR-RUL-PY-1 | internal/rules/python.go |
| FR-UI-1 | internal/ui/preview.go |
| FR-UI-2,3 | internal/ui/confirm.go |
| FR-UI-4..7 | internal/ui/executor.go |
| All type invariants | internal/types/types.go |
| NFR-MAI-4 (DAG) | go.mod + whole codebase |
