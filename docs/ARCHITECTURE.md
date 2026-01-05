# 🏗️ SPARK - Architecture Documentation

## Overview

SPARK is a Go-based Terminal User Interface (TUI) application built using the **Bubble Tea** framework (Elm Architecture). It manages development environment updates with a focus on safety, performance, and user experience.

**Version**: 0.6.0
**Language**: Go 1.24+
**Framework**: Bubble Tea (Elm Architecture)

---

## Project Structure

```
labs-spark/
├── cmd/
│   └── spark/
│       └── main.go              (Entry point - 39 lines)
│
├── internal/
│   ├── core/                    (146 lines - Domain layer)
│   │   ├── types.go            - Type definitions & enums
│   │   └── inventory.go        - Tool catalog (71 tools)
│   │
│   ├── updater/                 (340 lines - Detection layer)
│   │   ├── detector.go         - Version detection logic
│   │   └── version.go          - Regex-based version parsing
│   │
│   └── tui/                     (1,470 lines - Presentation layer)
│       ├── model.go            - Business logic & state management
│       ├── view.go             - Main dashboard rendering
│       ├── styles.go           - Centralized theming
│       ├── summary.go          - Summary screen
│       ├── preview.go          - Dry-run preview screen
│       └── states.go           - State machine documentation
│
├── docs/                        (Documentation)
├── config/                      (Legacy bash config)
└── lib/                         (Legacy bash modules)

Total: ~1,956 lines of Go code
```

---

## Architectural Layers

### 1. **Entry Point** (`cmd/spark/main.go`)

**Responsibility**: Application bootstrap

```go
func main() {
    // 1. Setup debug logging
    f, _ := tea.LogToFile("spark_debug.log", "debug")
    defer f.Close()

    // 2. Panic recovery
    defer func() {
        if r := recover() { /* handle */ }
    }()

    // 3. Initialize Bubble Tea program
    p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
    p.Run()
}
```

**Key Features**:
- Debug logging to `spark_debug.log`
- Panic recovery with stack trace
- Alternate screen mode (preserves terminal state)

---

### 2. **Core Domain** (`internal/core/`)

**Responsibility**: Domain models and business rules

#### `types.go` - Type System

```go
// Update methods (how tools are updated)
type UpdateMethod string
const (
    MethodBrew, MethodNpmSys, MethodNpmPkg,
    MethodBrewPkg, MethodMacApp, MethodClaude,
    MethodDroid, MethodToad, MethodOpencode,
    MethodOmz, MethodManual
)

// Categories (logical grouping)
type Category string
const (
    CategoryCode, CategoryTerm, CategoryIDE,
    CategoryProd, CategoryInfra, CategoryUtils,
    CategoryRuntime, CategorySys
)

// Tool status lifecycle
type ToolStatus int
const (
    StatusChecking, StatusInstalled, StatusOutdated,
    StatusMissing, StatusUnmanaged, StatusManualCheck,
    StatusUpdating, StatusUpdated, StatusFailed
)

// Static tool metadata
type Tool struct {
    ID, Name, Binary, Package string
    Category Category
    Method   UpdateMethod
    Description string
}

// Runtime tool state
type ToolState struct {
    Tool          Tool
    Status        ToolStatus
    LocalVersion  string
    RemoteVersion string
    Message       string
}
```

**Design Principles**:
- Type-safe enums (no magic strings)
- Clear separation: static metadata vs runtime state
- Immutable tool definitions

#### `inventory.go` - Tool Catalog

```go
func GetInventory() []Tool {
    tools := []Tool{
        // AI Development (7 tools)
        {Name: "Claude CLI", Binary: "claude", ...},

        // Terminals (3 tools)
        {Name: "iTerm2", Binary: "iterm", ...},

        // ... 71 tools total across 8 categories
    }

    // Auto-assign IDs: S-01, S-02, ...
    for i := range tools {
        tools[i].ID = fmt.Sprintf("S-%02d", i+1)
    }

    return tools
}
```

**Features**:
- Single source of truth for all tools
- Automatic ID generation
- Easy to add new tools

---

### 3. **Updater Layer** (`internal/updater/`)

**Responsibility**: Version detection and comparison

#### `detector.go` - Version Detection

```go
type Detector struct {
    brewCache, brewCaskCache string
}

func (d *Detector) GetLocalVersion(t core.Tool) string {
    // Strategy pattern based on update method
    switch {
    case t.Method == core.MethodMacApp:
        return d.getMacAppVersion(t.Binary)
    case t.Binary == "omz":
        return d.getOmzVersion()  // Git hash
    case t.Binary == "antigravity":
        return d.getAntigravityVersion()  // Multiple paths
    default:
        return d.getCliToolVersion(t)
    }
}
```

**Detection Strategies**:
- **macOS Apps**: Read `Info.plist` via `defaults read`
- **CLI Tools**: Run `--version` with 2s timeout
- **Git-based**: Get commit hash (Oh My Zsh)
- **Special cases**: Claude (multiple paths), AWS CLI, Go

#### `version.go` - Regex-based Parsing

```go
var (
    semverPattern      = regexp.MustCompile(`v?(\d+\.\d+\.\d+[\w\-\+]*)`)
    majorMinorPattern  = regexp.MustCompile(`v?(\d+\.\d+)`)
    dateVersionPattern = regexp.MustCompile(`(\d{4}\.\d+\.\d+)`)
    gitHashPattern     = regexp.MustCompile(`\b([a-f0-9]{7,40})\b`)
)

func CleanVersionString(output string) string {
    // Try patterns in order of specificity
    if version := extractPattern(output, semverPattern); version != "" {
        return cleanVersion(version)
    }
    // ... fallbacks
}

func ParseToolSpecificVersion(toolBinary, output string) string {
    // Custom logic for 11 tools: aws, go, python, docker, git, etc.
}
```

**Supported Formats**:
- Semantic versioning: `1.2.3`, `v1.2.3-beta`
- Major.Minor: `20.11`
- Date-based: `2024.1.15`
- Git hashes: `abc123f`
- Tool-specific: AWS CLI, Go, Python, Docker, etc.

---

### 4. **TUI Layer** (`internal/tui/`)

**Responsibility**: User interface and interaction

#### Elm Architecture (Bubble Tea)

```
┌──────────────────────────────────────┐
│          Model (State)               │
│  - Current screen                    │
│  - Tool list                         │
│  - User selections                   │
│  - Search query                      │
│  - Progress tracking                 │
└──────────────────────────────────────┘
         ▲                    │
         │ Messages           │ Commands
         │                    ▼
┌────────┴────────┐   ┌────────────────┐
│  Update Logic   │   │  View Renderer │
│  (Event Handler)│   │  (Pure func)   │
└─────────────────┘   └────────────────┘
```

#### `model.go` - State Management

```go
type Model struct {
    state         sessionState      // Current screen
    items         []core.ToolState  // All tools
    detector      *updater.Detector // Version checker
    cursor        int               // Selected item
    checked       map[int]bool      // User selections
    width, height int               // Terminal size
    loading       int               // Tools being checked
    updating      int               // Tools being updated
    totalUpdate   int               // Total to update
    progress      progress.Model    // Progress bar
    searchQuery   string            // Search filter
    filteredItems []int             // Filtered indices
}

// Elm Architecture: Init
func (m Model) Init() tea.Cmd {
    return tea.Batch(tick(), m.checkAllVersions())
}

// Elm Architecture: Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case CheckResultMsg:
        // Update tool version info
    case UpdateResultMsg:
        // Update tool status after update
    case tea.KeyMsg:
        // Handle user input
    }
}

// Elm Architecture: View
func (m Model) View() string {
    switch m.state {
    case stateSplash:   return m.ViewSplash()
    case stateMain:     return m.ViewMain()
    case stateSearch:   return m.ViewMain()  // With search bar
    case statePreview:  return m.ViewPreview()
    case stateConfirm:  return m.overlayModal("")
    case stateUpdating: return m.ViewMain()  // With progress
    case stateSummary:  return m.ViewSummary()
    }
}
```

#### Message Types

```go
// Async version check result
type CheckResultMsg struct {
    Index         int
    LocalVersion  string
    RemoteVersion string
    Status        core.ToolStatus
    Message       string
}

// Async update result
type UpdateResultMsg struct {
    Index   int
    Success bool
    Message string
}

// Timer tick (for splash screen)
type TickMsg time.Time
```

#### `styles.go` - Centralized Theming

```go
// Color Palette
var (
    cGreen  = lipgloss.Color("#04B575")
    cBlue   = lipgloss.Color("#2E7DE1")
    cPurple = lipgloss.Color("#A78BFA")
    cGray   = lipgloss.Color("#6B7280")
    cRed    = lipgloss.Color("#EF4444")
    cYellow = lipgloss.Color("#F59E0B")
)

// Pre-rendered status indicators
var (
    statusChecking = lipgloss.NewStyle().Foreground(cYellow).Render("⟳ Checking...")
    statusUpToDate = lipgloss.NewStyle().Foreground(cGray).Render("✔ Up to date")
    statusUpdating = lipgloss.NewStyle().Foreground(cBlue).Render("➜ Updating...")
    statusSuccess  = lipgloss.NewStyle().Foreground(cGreen).Render("✔ Updated")
    statusFailed   = lipgloss.NewStyle().Foreground(cRed).Render("✘ Failed")
)
```

#### `view.go` - Modular Rendering

**14 specialized functions**:
- `ViewMain()` - Main orchestrator
- `renderHeader()` - Top bar
- `renderSearchBar()` - Search input
- `renderProgressBar()` - Update progress
- `renderGrid()` - Two-column layout
- `renderCategoryCard()` - Tool category
- `renderToolLine()` - Individual tool
- `renderItemStatus()` - Status indicator
- `renderHelpBar()` - Bottom help text
- ... and more

---

## State Machine

### States (7 total)

```
stateSplash → stateMain ←──┐
                ├─→ stateSearch ──┘
                ├─→ statePreview ─┐
                └─→ stateConfirm ─┤
                        ↓         ↓
                   stateUpdating
                        ↓
                   stateSummary
                        ↓
                      EXIT
```

### State Transitions (Validated)

```go
validTransitions := map[sessionState][]sessionState{
    stateSplash:  {stateMain},
    stateMain:    {stateSearch, statePreview, stateConfirm, stateUpdating},
    stateSearch:  {stateMain},
    statePreview: {stateMain, stateConfirm, stateUpdating},
    stateConfirm: {stateMain, stateUpdating},
    stateUpdating: {stateSummary},
    stateSummary: {}, // Terminal
}
```

See `docs/STATES.md` for detailed state machine documentation.

---

## Concurrency Model

### Parallel Version Checking

```go
func (m Model) checkAllVersions() tea.Cmd {
    var cmds []tea.Cmd
    for i := range m.items {
        cmds = append(cmds, m.checkVersion(i))  // Goroutine per tool
    }
    return tea.Batch(cmds...)  // Run all in parallel
}
```

**Benefits**:
- Checks 71 tools in ~2 seconds (vs 142s sequential)
- Non-blocking UI
- Message-based safe communication

### Command Execution with Timeout

```go
func runCmd(name string, args ...string) string {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, name, args...)
    // ... execute
}
```

---

## Design Patterns

### 1. **Elm Architecture** (Entire TUI)
- Immutable state
- Pure view functions
- Message-based updates
- Side effects as commands

### 2. **Strategy Pattern** (Updater)
```go
switch t.Method {
case MethodBrew:      // brew upgrade
case MethodMacApp:    // Read plist
case MethodNpmPkg:    // npm update -g
}
```

### 3. **Repository Pattern** (Core)
```go
func GetInventory() []Tool  // Single source of truth
```

### 4. **Command Pattern** (Bubble Tea)
```go
type tea.Cmd func() tea.Msg  // Encapsulated side effects
```

### 5. **Observer Pattern** (Async updates)
- Goroutines send `CheckResultMsg`
- Model updates state
- View automatically reflects changes

---

## Data Flow

### Version Check Flow

```
User launches SPARK
    │
    ▼
Init() → checkAllVersions()
    │
    ├──→ [Goroutine 1] Check Claude    → CheckResultMsg
    ├──→ [Goroutine 2] Check Node      → CheckResultMsg
    ├──→ [Goroutine 3] Check Python    → CheckResultMsg
    ...
    │
    ▼
Update(CheckResultMsg)
    │
    ├─→ items[i].LocalVersion = msg.LocalVersion
    ├─→ items[i].Status = msg.Status
    └─→ loading--
    │
    ▼
View() renders updated state
```

### Update Execution Flow

```
User presses ENTER
    │
    ▼
Check for dangerous runtimes?
    │
    ├─ YES → stateConfirm
    │           │
    │           ▼
    │       User confirms (Y)
    │
    └─ NO ──────┘
    │
    ▼
startUpdates()
    │
    ├──→ [Goroutine 1] Update tool 1 → UpdateResultMsg
    ├──→ [Goroutine 2] Update tool 2 → UpdateResultMsg
    ...
    │
    ▼
Update(UpdateResultMsg)
    │
    ├─→ items[i].Status = StatusUpdated/StatusFailed
    └─→ updating--
    │
    ▼
All complete? → stateSummary
```

---

## Performance Characteristics

| Metric | Value | Notes |
|--------|-------|-------|
| **Binary size** | ~4.5 MB | Optimized with `-ldflags="-s -w"` |
| **Startup time** | ~50ms | Includes splash animation |
| **Version check** | ~2s | 71 tools checked in parallel |
| **Memory usage** | ~20 MB | Lightweight TUI |
| **CPU usage** | <5% | During checks, ~0% idle |

---

## Extension Points

### Adding a New Tool

1. **Add to `inventory.go`**:
```go
{Name: "New Tool", Binary: "newtool", Package: "newtool-pkg",
 Category: CategoryProd, Method: MethodBrewPkg}
```

2. **Optional**: Add custom version parsing in `version.go`:
```go
case "newtool":
    // Custom parsing logic
    return cleanVersion(output)
```

### Adding a New Update Method

1. **Define enum in `types.go`**:
```go
MethodCustom UpdateMethod = "custom"
```

2. **Implement in `detector.go`** (future):
```go
case core.MethodCustom:
    // Custom update logic
```

### Adding a New Screen

1. **Define state in `model.go`**:
```go
stateNewScreen sessionState = iota
```

2. **Add view in separate file**:
```go
func (m Model) ViewNewScreen() string { ... }
```

3. **Update `View()` router**:
```go
case stateNewScreen: return m.ViewNewScreen()
```

---

## Testing Strategy

### Unit Tests (Recommended)
- `detector_test.go` - Version parsing
- `version_test.go` - Regex patterns
- `states_test.go` - State transitions

### Integration Tests (Future)
- Full user flows
- Mock Bubble Tea runtime

### Manual Testing
- See `docs/TESTING.md`

---

## Build System

### Standard Build
```bash
go build -o spark-tui cmd/spark/main.go
```

### Optimized Build (Production)
```bash
go build -ldflags="-s -w" -o spark-tui cmd/spark/main.go
```

**Flags**:
- `-s`: Omit symbol table
- `-w`: Omit DWARF debug info
- Result: ~40% smaller binary

---

## Dependencies

### Direct
- `github.com/charmbracelet/bubbletea` v1.3.10 - TUI framework
- `github.com/charmbracelet/lipgloss` v1.1.0 - Styling
- `github.com/charmbracelet/bubbles/progress` v0.21.0 - Progress bar

### Transitive (16 total)
- Terminal handling, color support, input handling
- No external system dependencies

---

## Future Architecture Plans

1. **Persistent State**: Save selections/preferences
2. **Plugin System**: Dynamic tool loading
3. **Remote API**: Check versions from registries
4. **Config File**: User-defined tools
5. **Update History**: Track changes over time
