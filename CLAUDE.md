# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**SPARK** (System Intelligence & Update Utility) is a professional bash-based CLI tool for managing development environment updates. The tool categorizes updates into safe groups (AI tools, utilities) and dangerous zones (runtimes like Node.js, Python) with explicit confirmation requirements.

The project is **modular** (v0.4.2+) and organized into:
- `spark.sh`: Main entry point.
- `config/`: Configuration files (tool definitions).
- `lib/`: Functional modules (logic, UI, detection).

## Architecture

### Core Design Philosophy

SPARK follows a **modular surgical precision** approach:
- **Modular Structure**: Split into `config/tools.conf` and `lib/*.sh` for scalability.
- **Pre-fetch intelligence**: Queries Homebrew once upfront to avoid redundant network calls.
- **Spark IDs**: Tools are indexed (S-01, S-02) for targeted updates.

### Main Components

**1. Entry Point (`spark.sh`)**
- Orchestrates module loading and main execution flow.

**2. Configuration (`config/tools.conf`)**
- Central registry of all managed tools.

**3. Modules (`lib/`)**
- `common.sh`: Styling and global state.
- `detect.sh`: Local and remote version discovery logic.
- `update.sh`: Method-based update execution.
- `ui.sh`: Banner, table rendering, and interactive menu.

### Tool Categories

| Category | Risk Level | Tools |
|----------|-----------|-------|
| **CODE** | Low | Claude, Droid, Gemini, Codex, etc. |
| **TERM** | Low | iTerm2, Ghostty, Warp |
| **IDE**  | Low | VS Code, Cursor, Zed, Windsurf, Antigravity |
| **PROD** | Low | JQ, FZF, Ripgrep, Bat, LazyGit |
| **INFRA**| Medium | Docker, Kubectl, Terraform, AWS, Ngrok |
| **UTILS**| Low | Git, Tmux, Zellij, Oh My Zsh |
| **RUNTIME**| **High** | Node, Python, Go, Ruby, Postgres |
| **SYS**  | Medium | Homebrew, NPM |

## Version History

- **v0.4.2**: Modular Refactor + Infrastructure & Productivity categories + Spark IDs (S-XX).
- **v0.4.0**: IDE Expansion (Windsurf, Antigravity) + Smart status indicators.
- **v0.3.1**: AI Tools Intelligence - Accurate version detection for Claude, Droid, OpenCode, Toad.
- **v0.3.0**: Expanded coverage - New IDE category, OMZ support.

## Development Guidelines

### Adding a New Tool

Edit `config/tools.conf`. Add to the `TOOLS` array:
```bash
"CATEGORY:binary_cmd:package_name:Display Name:update_method"
```

### Testing Changes

1. Run `bash -n spark.sh` and `bash -n lib/*.sh` to check for syntax errors.
2. Test specific tool updates using Spark IDs (e.g., `./spark.sh S-01`).