# 🔧 Adding New Tools to SPARK

This guide explains how to extend SPARK with new tools and libraries.

---

## Quick Start

Adding a new tool requires modifying **only 1 file** in most cases:

1. Open `internal/core/inventory.go`
2. Add your tool to the list
3. Rebuild SPARK

That's it! SPARK will automatically handle version detection for most standard CLI tools.

---

## Step-by-Step Guide

### 1. **Identify Tool Metadata**

Gather the following information about your tool:

| Field | Description | Example |
|-------|-------------|---------|
| **Name** | Display name | `"Prettier"` |
| **Binary** | Command name | `"prettier"` |
| **Package** | Install package name | `"prettier"` |
| **Category** | Logical grouping | `CategoryProd` |
| **Method** | How to update it | `MethodNpmPkg` |

### 2. **Choose the Right Category**

```go
CategoryCode     // AI Development tools (Claude, Droid, etc.)
CategoryTerm     // Terminal emulators (iTerm, Ghostty, etc.)
CategoryIDE      // Code editors (VS Code, Cursor, etc.)
CategoryProd     // Productivity CLI tools (jq, fzf, ripgrep, etc.)
CategoryInfra    // Infrastructure tools (Docker, Kubernetes, etc.)
CategoryUtils    // System utilities (Git, Tmux, Oh My Zsh, etc.)
CategoryRuntime  // Programming runtimes (Node, Python, Go, etc.)
CategorySys      // Package managers (Homebrew, NPM, etc.)
```

**Note**: `CategoryRuntime` triggers a safety confirmation modal before updates.

### 3. **Choose the Right Update Method**

```go
MethodBrew      // Homebrew formula: brew install <package>
MethodBrewPkg   // Homebrew package: brew upgrade <package>
MethodNpmSys    // System npm: npm update -g
MethodNpmPkg    // npm package: npm update -g <package>
MethodMacApp    // macOS application bundle (.app)
MethodClaude    // Custom: Claude CLI specific
MethodDroid     // Custom: Droid CLI specific
MethodToad      // Custom: Toad CLI specific
MethodOpencode  // Custom: OpenCode specific
MethodOmz       // Custom: Oh My Zsh (git-based)
MethodManual    // Requires manual intervention
```

### 4. **Add to Inventory**

Edit `internal/core/inventory.go`:

```go
func GetInventory() []Tool {
    tools := []Tool{
        // ... existing tools ...

        // 🆕 Your new tool here
        {
            Name:     "Prettier",                  // Display name
            Binary:   "prettier",                  // Command to run
            Package:  "prettier",                  // Package name
            Category: CategoryProd,                // Productivity tool
            Method:   MethodNpmPkg,                // npm global package
            Description: "Code formatter",         // Optional
        },

        // ... more tools ...
    }
    // ... rest of function
}
```

**That's it!** SPARK will auto-assign an ID (`S-XX`) and handle version detection.

---

## Examples by Type

### Example 1: npm Global Package

```go
{
    Name:     "TypeScript",
    Binary:   "tsc",
    Package:  "typescript",
    Category: CategoryProd,
    Method:   MethodNpmPkg,
}
```

**How it works**:
- SPARK runs: `tsc --version`
- Parses output with regex
- Updates with: `npm update -g typescript` (future)

---

### Example 2: Homebrew CLI Tool

```go
{
    Name:     "htop",
    Binary:   "htop",
    Package:  "htop",
    Category: CategoryUtils,
    Method:   MethodBrewPkg,
}
```

**How it works**:
- SPARK runs: `htop --version`
- Updates with: `brew upgrade htop` (future)

---

### Example 3: macOS Application

```go
{
    Name:     "Obsidian",
    Binary:   "obsidian",  // Not used for apps
    Package:  "obsidian",  // Not used for apps
    Category: CategoryProd,
    Method:   MethodMacApp,
}
```

**How it works**:
- SPARK reads: `/Applications/Obsidian.app/Contents/Info.plist`
- Extracts: `CFBundleShortVersionString`
- Updates: Manual (App Store or direct download)

**Important**: For macOS apps, you need to add the app path to `detector.go` (see Advanced Customization below).

---

### Example 4: Programming Runtime

```go
{
    Name:     "Rust",
    Binary:   "rustc",
    Package:  "rust",
    Category: CategoryRuntime,  // 🚨 Triggers safety modal
    Method:   MethodBrewPkg,
}
```

**Special behavior**:
- Selecting any `CategoryRuntime` tool triggers "DANGER ZONE" confirmation
- Prevents accidental breaking of development environments

---

## Advanced Customization

### Custom Version Detection

If the tool has non-standard version output, add custom parsing to `internal/updater/version.go`:

```go
func ParseToolSpecificVersion(toolBinary string, output string) string {
    switch toolBinary {
    // ... existing cases ...

    case "yourtool":
        // Custom parsing logic
        parts := strings.Split(output, " - ")
        if len(parts) > 1 {
            return cleanVersion(parts[1])
        }
    }

    // Default: use generic cleaner
    return CleanVersionString(output)
}
```

**Example**: If your tool outputs `YourTool - Version 1.2.3 (build 456)`, you can extract `1.2.3`.

---

### macOS App Detection

For macOS `.app` bundles, add the path to `internal/updater/detector.go`:

```go
func (d *Detector) getMacAppVersion(binary string) string {
    appPaths := map[string]string{
        // ... existing paths ...
        "yourtool": "/Applications/YourTool.app",  // 🆕 Add here
    }

    // ... rest of function
}
```

---

### Multiple Installation Paths

If your tool can be installed in multiple locations:

```go
func (d *Detector) getYourToolVersion() string {
    // Try custom path first
    customPath := os.Getenv("HOME") + "/.yourtool/bin/yourtool"
    if _, err := os.Stat(customPath); err == nil {
        output := runCmd(customPath, "--version")
        return ParseToolSpecificVersion("yourtool", output)
    }

    // Fallback to PATH
    output := runCmd("yourtool", "--version")
    return ParseToolSpecificVersion("yourtool", output)
}

// Then use in GetLocalVersion():
case "yourtool":
    return d.getYourToolVersion()
```

---

### Custom Update Method

If the tool requires a custom update mechanism:

1. **Define new method in `internal/core/types.go`**:

```go
const (
    // ... existing methods ...
    MethodYourTool UpdateMethod = "yourtool"
)
```

2. **Use in inventory**:

```go
{
    Name:     "Your Tool",
    Binary:   "yourtool",
    Package:  "yourtool-pkg",
    Category: CategoryProd,
    Method:   MethodYourTool,  // 🆕 Custom method
}
```

3. **Implement update logic** (future - in executor.go):

```go
case core.MethodYourTool:
    // Custom update command
    return executeCommand("yourtool", "self-update")
```

---

## Version Detection Strategies

SPARK tries multiple strategies to detect versions:

### Strategy 1: Standard `--version` Flag

Most CLI tools support this:

```bash
yourtool --version
# Output: yourtool 1.2.3
# SPARK extracts: 1.2.3
```

### Strategy 2: Regex Pattern Matching

SPARK uses 5 regex patterns (in order):

1. **Semantic versioning**: `v1.2.3`, `1.2.3-beta`
2. **Major.Minor**: `20.11`, `v16.0`
3. **Date-based**: `2024.1.15`
4. **Git hash**: `abc123f` (7+ hex chars)
5. **Simple number**: `123`

### Strategy 3: Tool-Specific Parsing

For tools with unique formats (see `version.go`):

```
aws-cli/2.22.35 Python/3.11.9  → 2.22.35
go version go1.23.4 darwin      → 1.23.4
Docker version 24.0.7, build... → 24.0.7
```

### Strategy 4: Plist Reading (macOS Apps)

```bash
defaults read /Applications/YourApp.app/Contents/Info.plist CFBundleShortVersionString
# Output: 1.2.3
```

---

## Testing Your Addition

### 1. **Rebuild SPARK**

```bash
cd /path/to/labs-spark
go build -ldflags="-s -w" -o spark-tui cmd/spark/main.go
cp spark-tui ~/.local/bin/spark
```

### 2. **Launch and Check**

```bash
spark
```

- Navigate to your tool's category
- Verify it appears in the list
- Check if the version is detected correctly

### 3. **Debug Version Detection**

If version shows as `MISSING` or incorrect:

```bash
# Enable debug logging
DEBUG=1 spark

# Check logs after running
cat spark_debug.log | grep "your-tool"
```

### 4. **Manual Version Check**

Test the detection command directly:

```bash
yourtool --version
```

If the output is non-standard, you'll need custom parsing (see Advanced Customization).

---

## Common Issues & Solutions

### Issue: Tool shows as "MISSING" but is installed

**Cause**: Binary not in `$PATH` or has different name

**Solution 1**: Check actual binary name

```bash
which yourtool
```

**Solution 2**: Update inventory with correct binary name

```go
{
    Binary: "actual-binary-name",  // Not "yourtool"
}
```

---

### Issue: Version shows as `"Detected"` instead of number

**Cause**: Version parsing failed

**Solution**: Add custom parsing to `version.go`

```go
case "yourtool":
    // Extract version from output
    if strings.Contains(output, "Version:") {
        parts := strings.Split(output, "Version:")
        return strings.TrimSpace(parts[1])
    }
```

---

### Issue: macOS App shows as "MISSING"

**Cause**: App path not registered

**Solution**: Add to `getMacAppVersion()` in `detector.go`

```go
appPaths := map[string]string{
    "yourapp": "/Applications/YourApp.app",
}
```

---

### Issue: Tool has multiple installation methods

**Example**: Python can be `python`, `python3`, `python3.13`

**Solution**: Choose the most common one, or add multiple entries

```go
{Name: "Python 3.13", Binary: "python3", ...},
{Name: "Python 3.12", Binary: "python3.12", ...},
```

---

## Best Practices

### 1. **Use Standard Naming**

- Binary name = lowercase package name when possible
- Example: Package `prettier` → Binary `prettier`

### 2. **Test Before Committing**

- Ensure version detection works
- Check that the tool appears in the correct category

### 3. **Document Special Cases**

If your tool requires special handling, document it:

```go
{
    Name:     "Special Tool",
    Binary:   "special",
    Package:  "special-tool",
    Category: CategoryProd,
    Method:   MethodManual,  // 🔔 Requires manual update
    Description: "Note: Must update via their website",
}
```

### 4. **Consider Impact on Users**

- Don't add tools most users won't have
- Focus on common development tools
- Use `CategoryRuntime` sparingly (triggers warnings)

### 5. **Alphabetical Order (Optional)**

Keep tools within each category alphabetically sorted for easy scanning.

---

## Real-World Examples

### Adding Bun (JavaScript Runtime)

```go
{
    Name:     "Bun",
    Binary:   "bun",
    Package:  "bun",
    Category: CategoryRuntime,  // It's a runtime
    Method:   MethodBrewPkg,
},
```

### Adding Neovim

```go
{
    Name:     "Neovim",
    Binary:   "nvim",
    Package:  "neovim",
    Category: CategoryIDE,  // It's an editor
    Method:   MethodBrewPkg,
},
```

### Adding Raycast (macOS App)

```go
{
    Name:     "Raycast",
    Binary:   "raycast",  // Placeholder
    Package:  "raycast",  // Placeholder
    Category: CategoryProd,
    Method:   MethodMacApp,
},
```

Then in `detector.go`:

```go
appPaths := map[string]string{
    // ...
    "raycast": "/Applications/Raycast.app",
}
```

---

## Contributing Your Additions

If you add a tool that others might find useful:

1. Test it thoroughly
2. Ensure version detection works
3. Submit a pull request with:
   - Updated `inventory.go`
   - Any custom detection logic
   - Documentation updates

---

## Next Steps

- See `docs/ARCHITECTURE.md` for detailed code structure
- See `docs/WORKFLOWS.md` for user interaction flows
- See `internal/tui/states.go` for state machine documentation
