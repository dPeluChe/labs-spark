# ⚡ SPARK TUI

**SPARK** is a professional, cinematic Terminal User Interface (TUI) for managing development environment updates with surgical precision.

Built with **Go**, **Bubble Tea**, and **Lip Gloss**, it replaces the legacy bash script with a high-performance, concurrent, and beautiful interactive dashboard.

## 🚀 Features (v0.5.3)

### 🎨 Cinematic Experience
- **Splash Screen**: Animated startup sequence.
- **Grid Dashboard**: 2-column layout with categorized cards.
- **Visual Status**: Real-time checking, updating, and success indicators.

### ⚡ Performance
- **Parallel Scanning**: Checks 40+ tools simultaneously using Goroutines.
- **Instant Navigation**: Jump between categories with shortcuts (`C`, `T`, `I`...).
- **Smart Selection**: Toggle individual items (`Space`) or entire groups (`G`).

### 🛡️ Safety First
- **Danger Zone**: Explicit confirmation modal when updating critical runtimes (Node, Python, Postgres).
- **Update Focus**: visual isolation of active updates to reduce noise.

## 🛠 Supported Stack

| Category | Shortcuts | Tools |
|----------|-----------|-------|
| **AI Development** | `[C]` | Claude, Droid, Gemini, OpenCode, Codex, Crush, Toad |
| **Terminals** | `[T]` | iTerm2, Ghostty, Warp |
| **IDEs** | `[I]` | VS Code, Cursor, Zed, Windsurf, Antigravity |
| **Productivity** | `[P]` | JQ, FZF, Ripgrep, Bat, HTTPie, LazyGit, TLDR |
| **Infrastructure** | `[F]` | Docker, K8s, Helm, Terraform, AWS, Ngrok |
| **Utilities** | `[U]` | Git, Tmux, Zellij, Oh My Zsh, SQLite, Watchman |
| **Runtimes** | `[R]` | Node.js, Python 3.13, Go, Ruby, PostgreSQL 16 |
| **System** | `[S]` | Homebrew Core, NPM Globals |

## 📦 Installation

### From Source
```bash
# Clone repository
git clone https://github.com/dpeluche/spark.git
cd spark

# Build and Install
go build -ldflags="-s -w" -o spark-tui cmd/spark/main.go
mv spark-tui ~/.local/bin/spark

# Add to shell (if not already)
echo 'alias spark="~/.local/bin/spark"' >> ~/.zshrc
```

## ⌨️ Controls

| Key | Action |
|-----|--------|
| `↑/↓` `j/k` | Navigate items |
| `Space` | Select/Deselect item |
| `G` | Select/Deselect **entire group** |
| `Tab` | Jump to next group |
| `C, T, I...` | Jump to specific category |
| `Enter` | **Update selected** (or current if none selected) |
| `Q` / `Esc` | Quit |

## 🧠 Architecture

Spark is a pure Go application using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework (ELM architecture).

- **`cmd/spark`**: Entry point.
- **`internal/core`**: Tool definitions and inventory.
- **`internal/updater`**: Version detection logic (using `os/exec`).
- **`internal/tui`**: The UI logic, rendering, and state machine.

## 📄 License

MIT
