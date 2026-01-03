# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**SPARK** has evolved from a bash script into a **Go-based TUI application**.
It manages system updates for developers, focusing on AI tools, IDEs, and Infrastructure.

## Architecture (Go + Bubble Tea)

The project follows the standard Go project layout:

- **`cmd/spark/main.go`**: Entry point. initializes the Tea program.
- **`internal/core/`**:
    - `types.go`: Struct definitions (`Tool`, `ToolState`).
    - `inventory.go`: The static list of all supported tools.
- **`internal/updater/`**:
    - `detector.go`: Logic to run shell commands (`brew`, `npm`, etc.) to find local/remote versions.
- **`internal/tui/`**:
    - `model.go`: The heart of the application. Contains the Bubble Tea `Model`, `Update`, and `View` functions. Handles all UI logic, rendering, and keybindings.

## Key Features & Logic

### State Machine (`sessionState`)
1. **`stateSplash`**: Intro animation.
2. **`stateMain`**: The main dashboard grid.
3. **`stateConfirm`**: "Danger Zone" modal for critical runtimes.
4. **`stateUpdating`**: Execution phase (currently simulated).
5. **`stateSummary`**: Final report.

### Concurrency
- Uses **Goroutines** via `tea.Cmd` to check tool versions in parallel at startup (`Init`).
- Uses `tea.Batch` to manage multiple simultaneous checks.

### UX Patterns
- **Mnemonic Navigation**: Jump to categories via first letter (`C`, `T`, `R`...).
- **Safety Lock**: Prevents accidental updates of Runtimes (Node/Python) via a modal.
- **Visual Focus**: Dims non-active items during the update phase.

## Development Commands

### Running Locally
```bash
go run cmd/spark/main.go
```

### Building for Release
```bash
go build -ldflags="-s -w" -o spark-tui cmd/spark/main.go
```

### Adding a New Tool
1. Open `internal/core/inventory.go`.
2. Add a new `Tool` struct to the list.
3. If it requires custom version detection, update `internal/updater/detector.go` (`GetLocalVersion`).

## Version History

- **v0.5.3**: UX Polish - Visual focus during updates, Smart Enter, Ctrl+C fix.
- **v0.5.2**: Safety System - Danger Zone Modal, Group Selection (G), Parallel Checks.
- **v0.5.0**: **The Great Rewrite**. Migrated to Go/Bubble Tea. Grid Layout, Splash Screen.
- **v0.4.x**: Legacy Bash Script era (Archived).
