# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**SPARK** (System Intelligence & Update Utility) is a professional bash-based CLI tool for managing development environment updates. The tool categorizes updates into safe groups (AI tools, utilities) and dangerous zones (runtimes like Node.js, Python) with explicit confirmation requirements.

This is a **single-file bash application** (`spark.sh`) with no build process or external dependencies beyond the tools it manages.

## Installation & Setup

The tool is designed to be aliased in the user's shell:

```bash
# Make executable
chmod +x /path/to/labs-spark/spark.sh

# Add to ~/.zshrc or ~/.bashrc
alias spark='/path/to/labs-spark/spark.sh'

# Reload shell
source ~/.zshrc
```

## Running the Tool

```bash
# Execute directly
./spark.sh

# Or via alias
spark
```

The tool presents an interactive menu with 5 options:
1. AI & Terminals (CODE + TERM categories)
2. Utilities (UTILS category)
3. Runtimes with safety confirmation (RUNTIME category)
4. Full System (all categories)
5. Exit

## Architecture

### Core Design Philosophy

SPARK follows a **surgical precision** approach:
- **Pre-fetch intelligence**: Queries Homebrew once upfront to avoid redundant network calls
- **Smart skipping**: Only updates tools where local version ≠ remote version
- **Session protection**: Warns if tools are actively running before updates
- **Category-based safety**: Critical runtimes require explicit "yes" confirmation
- **No rollback**: Updates are irreversible; relies on preventive measures

### Main Components

**1. Configuration Section (Lines 7-63)**
- ANSI color codes for terminal styling
- `TOOLS` array: Central registry of all managed tools
- Tool format: `"CATEGORY:BinaryName:PackageName:DisplayName:UpdateMethod"`
- Global counters and cache variables

**2. Version Detection (`get_local_version`, Lines 80-124)**
- **macOS Apps**: Uses `defaults read` on Info.plist for iTerm2, Ghostty, Warp
- **NPM packages**: Parses `npm list -g` output
- **CLI tools**: Executes `--version` and extracts version string
- **Fallback**: Returns "Detected" or "Installed" if version parsing fails

**3. Remote Version Lookup (`get_remote_version`, Lines 126-143)**
- **NPM**: Uses `npm view package version`
- **Homebrew**: Queries pre-fetched `BREW_CACHE` (from `brew outdated --verbose`)
- **Others**: Returns "Latest" (requires manual verification)

**4. Session Guard (`check_active_sessions`, Lines 145-172)**
- Uses `pgrep -fi` to detect running processes
- Maps binary names to process names (e.g., `iterm` → `iTerm2`)
- Warns user but does not block updates

**5. System Analysis (`analyze_system`, Lines 174-249)**
- **Brew cache population**: Single `brew outdated --verbose` call
- **Grouped display**: Iterates tools by category (CODE, TERM, UTILS, RUNTIME, SYS)
- **Status icons**:
  - `●` = Installed and current
  - `↑` = Update available
  - `○` = Not installed
- **Update counting**: Tracks pending updates per category

**6. Update Execution (`perform_update`, Lines 251-290)**
- **Method dispatch**: Routes to correct package manager based on `UpdateMethod`
- **Success tracking**: Records updated tools with version changes for summary
- **Error handling**: Reports failures but continues execution

**7. Main Execution Flow (Lines 304-359)**
- Banner display
- System analysis with grouped output
- Session check
- Interactive mode selection
- Safety confirmation for RUNTIME/ALL modes
- Sequential update execution
- Summary report

### Tool Categories

| Category | Risk Level | Tools | Update Methods |
|----------|-----------|-------|----------------|
| **CODE** | Low | Claude CLI, Droid, Gemini, OpenCode, Codex, Crush, Toad | npm_pkg, droid, toad, opencode, brew_pkg |
| **TERM** | Low | iTerm2, Ghostty, Warp | mac_app (Homebrew casks) |
| **IDE** | Low | Visual Studio Code, Cursor, Zed Editor | mac_app (Homebrew casks) |
| **UTILS** | Low | Oh My Zsh, Zellij, Tmux, Git, Bash, SQLite, Watchman, Direnv, Heroku, Pre-commit | omz, brew_pkg |
| **RUNTIME** | **High** | Node.js, Python 3.13, Go, Ruby, PostgreSQL 16 | brew_pkg (requires confirmation) |
| **SYS** | Medium | Homebrew, NPM globals | brew, npm_sys |

### Update Methods

- **brew**: Updates Homebrew itself (`brew update && brew upgrade && brew cleanup`)
- **npm_sys**: Updates all global npm packages (`npm update -g`)
- **npm_pkg**: Installs specific npm package to latest (`npm install -g package@latest`)
- **brew_pkg**: Updates specific Homebrew formula (`brew upgrade package`)
- **mac_app**: Updates Homebrew cask (`brew upgrade --cask package`)
- **droid**: Curl-based installer (`curl -fsSL https://app.factory.ai/cli | sh`)
- **toad**: Curl-based installer (`curl -fsSL https://batrachian.ai/install | sh`)
- **opencode**: Custom upgrade command with fallback (`opencode upgrade || curl install`)
- **omz**: Git pull in Oh My Zsh directory (`cd ~/.oh-my-zsh && git pull`)

## Modifying the Tool

### Adding a New Tool

Edit the `TOOLS` array (lines 19-54):

```bash
"CATEGORY:binary_cmd:package_name:Display Name:update_method"
```

**Example**: Adding a new AI tool managed by npm:
```bash
"CODE:newtool:@company/newtool-cli:NewTool CLI:npm_pkg"
```

**Category options**: CODE, TERM, UTILS, RUNTIME, SYS

### Adding a New Update Method

1. Add case in `perform_update` function (line 266)
2. Implement version detection in `get_local_version` if needed
3. Add remote version logic in `get_remote_version`

### Modifying Categories or Safety Behavior

- **Category display order**: Controlled by `print_group` calls in `analyze_system` (lines 239-247)
- **Runtime safety**: Confirmation logic at lines 329-337
- **Session detection**: Tool-to-process mapping in `check_active_sessions` (lines 152-156)

## Version History

- **v0.3.0**: Expanded coverage - New IDE category (VSCode, Cursor, Zed), Oh My Zsh support, Toad CLI, cleaner Homebrew output, reorganized update modes
- **v0.2.5**: Visual clarity - "✔ Up to date" message instead of duplicate version numbers
- **v0.2.4**: Real-time Homebrew intelligence via pre-fetching
- **v0.2.1**: Post-execution summary report
- **v0.2.0**: Category-based surgical updates + runtime safety confirmation
- **v0.1.2**: macOS app support (iTerm2, Ghostty, Warp)

## Development Guidelines

### Testing Changes

Since this is a system update tool, **test in a safe environment**:

```bash
# Dry-run approach: Comment out actual update commands
# Replace line 267-282 with echo statements to simulate

# Test version detection
get_local_version "node"
get_remote_version "brew_pkg" "node" "20.0.0"

# Test with minimal tools
# Temporarily edit TOOLS array to include only 1-2 safe utilities
```

### Bash Best Practices

- Use `command -v` for binary existence checks (not `which`)
- Quote all variables: `"$variable"` to prevent word splitting
- Use `&>` for combined stdout/stderr redirection
- Prefer `[[ ]]` over `[ ]` for conditionals (bash-specific, more robust)
- Use `local` for function variables to avoid scope pollution

### Color Codes

Defined at lines 7-15:
- **BOLD**: Important headers
- **DIM**: Secondary/metadata information
- **GREEN**: Success states, installed tools
- **YELLOW**: Warnings, updates available
- **RED**: Errors, danger zone
- **MAGENTA**: Target versions
- **CYAN**: Branding, action indicators
- **RESET**: Always reset after colored output

### Common Patterns

**Tool iteration**:
```bash
for tool_entry in "${TOOLS[@]}"; do
    IFS=':' read -r category binary pkg display method <<< "$tool_entry"
    # Process tool
done
```

**Conditional formatting**:
```bash
printf "% -13b ${color}%-18s %-15s %b${RESET}\n" "$icon" "$display" "$current" "$target"
```

The `%b` flag interprets escape sequences in variables (for colors).

## Codebase Patterns

### State Management

- **Global cache**: `BREW_CACHE` populated once at startup (line 177)
- **Update counters**: Increment during analysis, used for display only
- **Update log**: `UPDATED_TOOLS` array accumulates successful updates for summary

### Error Handling

- **Missing binaries**: Detected as "MISSING" or "Not Installed" (lines 93-94, 203-206)
- **Version parsing failures**: Fallback to "Detected" or "Installed" (lines 109-122)
- **Update failures**: Reported but don't halt execution (lines 284-289)
- **User abort**: Clean exit on "no" confirmation (lines 333-336)

### Platform Assumptions

- **macOS-specific**: Uses `defaults read` for .app bundles, `/Applications` directory
- **Homebrew dependency**: Core package manager for most tools
- **NPM availability**: Required for npm-managed tools
- **Bash 3.2+**: Uses bash-specific features (`[[]]`, `read -p`, arrays)

## Troubleshooting

### Tool Not Detected

Check binary name matches actual command:
```bash
command -v tool_name
```

For macOS apps, verify path:
```bash
ls -la /Applications/ | grep -i AppName
```

### Version Parsing Issues

Add custom parsing in `get_local_version`:
```bash
elif [[ "$binary" == "toolname" ]]; then
    ver=$(toolname --version | awk '{print $2}')
```

### Homebrew Cache Misses

Tool won't show update if not in `brew outdated` output. This is intentional - means tool is already latest or not managed by brew.

### Session Detection False Positives

Process names are case-insensitive (`pgrep -fi`). Adjust mappings at lines 152-156 if needed.
