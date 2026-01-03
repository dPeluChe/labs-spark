# ⚡ SPARK

**SPARK** (System Intelligence & Update Utility) is a professional, lightweight CLI tool designed to keep your development environment and AI tools synchronized. 

Inspired by the life-force of the Transformers (Autobots), Spark provides a clean, grouped interface to monitor and update your mission-critical CLI tools.

## 🚀 Features

- **v0.3.0 Expanded Coverage**: New IDE category (VSCode, Cursor, Zed), Oh My Zsh support, Toad CLI integration, and cleaner Homebrew output.
- **v0.2.5 Visual Clarity**: Replaces redundant version numbers with a clean "✔ Up to date" message when no actions are needed.
- **v0.2.4 Real-Time Intelligence**: Pre-fetches accurate upgrade data from Homebrew to distinguish between "Latest" and actual available versions.
- **v0.2.1 Summary Report**: Displays a clean summary of all updated tools and version changes at the end of the process.
- **v0.2.0 Surgical Precision**: Categorized updates into AI Tools, Utilities, and Critical Runtimes.
- **v0.2.0 Runtime Safety**: Adds a "Danger Zone" confirmation before updating sensitive runtimes like Node.js or Python.
- **v0.1.2 App Support**: Monitors terminal emulators like iTerm2, Ghostty, and Warp Terminal.
- **Smart Updates**: Automatically compares local vs remote versions and skips tools that are already up to date.
- **Intelligence Dashboard**: Grouped view of all system components.
- **Session Protection**: Scans for active processes/sessions before updating.

## 🛠 Supported Tools

### AI Development Tools (CODE)
- Claude CLI, Droid, Gemini, OpenCode, Codex, Crush, Toad.

### Terminal Emulators (TERM)
- iTerm2, Ghostty, Warp Terminal.

### IDEs and Code Editors (IDE)
- Visual Studio Code, Cursor, Zed Editor.

### Safe Utilities (UTILS)
- **Shell**: Oh My Zsh, Zellij, Tmux.
- **Core**: Git, Bash, SQLite, Heroku, Pre-commit, Watchman, Direnv.

### Critical Runtimes (RUNTIME) ⚠️
- **Languages**: Node.js, Python 3.13, Go, Ruby.
- **Databases**: PostgreSQL 16.
*Updates for this category require explicit confirmation.*

### System Tools (SYS)
- Homebrew, NPM Global Packages.

## 📦 Installation

1. **Clone or move** the `labs-spark` folder to your preferred projects directory.
2. **Make it executable**:
   ```bash
   chmod +x labs-spark/spark.sh
   ```
3. **Add an alias** to your `.zshrc` or `.bashrc`:
   ```bash
   # Add this line (replace /path/to/ with your actual path)
   alias spark='/path/to/labs-spark/spark.sh'
   ```
4. **Reload your shell**:
   ```bash
   source ~/.zshrc
   ```

## ⌨️ Usage

Simply type `spark` in your terminal:

```bash
spark
```

Follow the interactive menu to select your update strategy.

## 🧠 How it Works (Internal Rules)

SPARK follows a strict logic to ensure system stability:

1.  **Binary Discovery**: Uses `command -v` for CLI tools and searches `/Applications` for macOS apps.
2.  **Version Extraction**: 
    - **NPM**: Uses `npm list -g` and `npm view`.
    - **Casks (Apps)**: Uses `defaults read` for local Info.plist and `brew info --cask` for remote.
    - **Brew**: Uses `brew --version`.
3.  **Session Guard**: Runs `pgrep -f` (or `pgrep -fi` for apps) before updates. 
4.  **Smart Skipping**: If `Local Version == Remote Version`, the update is bypassed.

## 📄 License

MIT - Feel free to use and improve.