# ⚡ SPARK

**SPARK** (System Intelligence & Update Utility) is a professional, lightweight CLI tool designed to keep your development environment and AI tools synchronized. 

Inspired by the life-force of the Transformers (Autobots), Spark provides a clean, grouped interface to monitor and update your mission-critical CLI tools.

## 🚀 Features

- **v0.1.2 App Support**: Now monitors terminal emulators like iTerm2, Ghostty, and Warp Terminal.
- **v0.1.1 Smart Updates**: Automatically compares local vs remote versions and skips tools that are already up to date.
- **Intelligence Dashboard**: Grouped view of AI Development Tools, Terminals, and System Tools.
- **Session Protection**: Scans for active processes/sessions before updating to prevent data loss.
- **Visual Feedback**: Professional TUI-style output with color-coded status icons:
  - ● (Green): Installed & Up to date.
  - ↑ (Yellow): Update available.
  - ○ (Dim): Not installed.

## 🛠 Supported Tools

### AI Development Tools
- **Claude CLI** (@anthropic-ai/claude-code)
- **Droid CLI** (Factory AI)
- **Gemini CLI** (Google)
- **OpenCode** (OpenCode AI)
- **Codex CLI** (OpenAI)
- **Crush CLI** (Development)

### Terminal Emulators
- **iTerm2**
- **Ghostty**
- **Warp Terminal**

### System Tools
- **Homebrew**
- **NPM Global Packages**

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


MIT - Feel free to use and improve.