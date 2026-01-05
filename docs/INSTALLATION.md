# 📦 SPARK - Installation Guide

## Quick Start

### Prerequisites
- **Go 1.24+** (for building from source)
- **macOS** (currently supported)
- Terminal with 256 color support

### Installation from Source

```bash
# 1. Clone the repository
cd /path/to/labs-spark

# 2. Build the binary
go build -ldflags="-s -w" -o spark-tui cmd/spark/main.go

# 3. Install to local bin
mkdir -p ~/.local/bin
cp spark-tui ~/.local/bin/spark
chmod +x ~/.local/bin/spark

# 4. Verify installation
spark
```

### Shell Configuration

The alias should already be configured in your `~/.zshrc`:

```bash
alias spark='~/.local/bin/spark'
```

If not, add it manually:

```bash
echo "alias spark='~/.local/bin/spark'" >> ~/.zshrc
source ~/.zshrc
```

---

## Running SPARK

### Basic Usage

Simply run:

```bash
spark
```

This will launch the interactive TUI dashboard.

### Debug Mode

To enable debug logging:

```bash
DEBUG=1 spark
```

Debug logs are written to `spark_debug.log` in the current directory.

---

## Updating SPARK

### Rebuild After Changes

```bash
cd /path/to/labs-spark
go build -ldflags="-s -w" -o spark-tui cmd/spark/main.go
cp spark-tui ~/.local/bin/spark
```

### Quick Update Script

Create an alias for easy updates:

```bash
# Add to ~/.zshrc
alias spark-rebuild='cd ~/path/to/labs-spark && go build -ldflags="-s -w" -o spark-tui cmd/spark/main.go && cp spark-tui ~/.local/bin/spark && echo "✓ SPARK updated"'
```

---

## Troubleshooting

### Issue: "command not found: spark"

**Solution**: Ensure `~/.local/bin` is in your PATH:

```bash
# Add to ~/.zshrc if not present
export PATH="$HOME/.local/bin:$PATH"
source ~/.zshrc
```

### Issue: "could not open a new TTY"

**Cause**: SPARK requires an interactive terminal

**Solution**: Don't pipe SPARK or run it in non-interactive contexts

### Issue: Binary won't execute

**Solution**: Check permissions:

```bash
chmod +x ~/.local/bin/spark
```

### Issue: Version check shows old version

**Solution**: Rebuild and reinstall:

```bash
cd /path/to/labs-spark
go build -ldflags="-s -w" -o spark-tui cmd/spark/main.go
cp spark-tui ~/.local/bin/spark
```

---

## Uninstallation

```bash
# Remove binary
rm ~/.local/bin/spark

# Remove alias from ~/.zshrc (manual edit)
# Remove this line: alias spark='~/.local/bin/spark'
```

---

## Advanced Installation

### Custom Install Location

```bash
# Install to custom location
mkdir -p ~/bin
cp spark-tui ~/bin/spark

# Update alias
alias spark='~/bin/spark'
```

### System-wide Installation (requires sudo)

```bash
sudo cp spark-tui /usr/local/bin/spark
sudo chmod +x /usr/local/bin/spark
```

---

## Verifying Installation

### Check Version

Currently, SPARK doesn't have a `--version` flag, but you can verify it runs:

```bash
spark  # Should launch the TUI
```

Press `Q` to quit immediately.

### Check Binary Info

```bash
file ~/.local/bin/spark
# Should show: Mach-O 64-bit executable arm64

ls -lh ~/.local/bin/spark
# Should show executable permissions and size (~4-5 MB)
```

---

## Dependencies

SPARK has **no runtime dependencies**. It's a single static binary.

Build dependencies (only needed for compilation):
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/lipgloss`
- `github.com/charmbracelet/bubbles/progress`

These are automatically fetched during `go build`.
