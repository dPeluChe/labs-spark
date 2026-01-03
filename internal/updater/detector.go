package updater

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dpeluche/spark/internal/core"
)

// Detector handles version checking logic
type Detector struct {
	brewCache     string
	brewCaskCache string
}

func NewDetector() *Detector {
	return &Detector{}
}

// WarmUpCache fetches brew info once to speed up subsequent checks
func (d *Detector) WarmUpCache() {
	// In a real implementation, we would run 'brew outdated --verbose' here and store it.
	// For simplicity in this step, we'll keep it basic or implement later.
}

func (d *Detector) GetLocalVersion(t core.Tool) string {
	switch t.Binary {
	case "iterm":
		return getAppVersion("/Applications/iTerm.app")
	case "ghostty":
		return getAppVersion("/Applications/Ghostty.app")
	case "warp":
		return getAppVersion("/Applications/Warp.app")
	case "code":
		return getAppVersion("/Applications/Visual Studio Code.app")
	case "cursor":
		return getAppVersion("/Applications/Cursor.app")
	case "zed":
		return getAppVersion("/Applications/Zed.app")
	case "windsurf":
		return getAppVersion("/Applications/Windsurf.app")
	case "antigravity":
		// Check custom path first
		if _, err := os.Stat(os.Getenv("HOME") + "/.antigravity/antigravity/bin/antigravity"); err == nil {
			return runCmd(os.Getenv("HOME")+"/.antigravity/antigravity/bin/antigravity", "--version")
		}
		return runCmd("antigravity", "--version")
	case "omz":
		omzPath := os.Getenv("HOME") + "/.oh-my-zsh"
		if _, err := os.Stat(omzPath); err == nil {
			// git --git-dir=... rev-parse --short HEAD
			cmd := exec.Command("git", "--git-dir="+omzPath+"/.git", "--work-tree="+omzPath, "rev-parse", "--short", "HEAD")
			out, err := cmd.Output()
			if err != nil {
				return "Installed"
			}
			return strings.TrimSpace(string(out))
		}
		return "MISSING"
	case "aws":
		out := runCmd("aws", "--version") // aws-cli/2.22.35 ...
		parts := strings.Fields(out)
		if len(parts) > 0 {
			// Extract 2.22.35 from aws-cli/2.22.35
			return strings.Split(parts[0], "/")[1]
		}
		return out
	case "go":
		out := runCmd("go", "version") // go version go1.23.4 darwin/arm64
		parts := strings.Fields(out)
		if len(parts) >= 3 {
			return strings.Replace(parts[2], "go", "", 1)
		}
		return out
	case "claude":
		// Try curl installation path first
		localPath := os.Getenv("HOME") + "/.local/bin/claude"
		if _, err := os.Stat(localPath); err == nil {
			return runCmd(localPath, "--version")
		}
		return runCmd("claude", "--version")
	default:
		// Generic --version check
		path, err := exec.LookPath(t.Binary)
		if err != nil {
			return "MISSING"
		}
		_ = path
		return runCmd(t.Binary, "--version")
	}
}

// Helper functions

func getAppVersion(appPath string) string {
	plistPath := appPath + "/Contents/Info.plist"
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return "MISSING"
	}
	cmd := exec.Command("defaults", "read", plistPath, "CFBundleShortVersionString")
	out, err := cmd.Output()
	if err != nil {
		return "Detected"
	}
	return strings.TrimSpace(string(out))
}

func runCmd(name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "MISSING"
	}
	
	// Basic cleanup of output (first line, first word usually)
	output := strings.TrimSpace(out.String())
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		// Custom logic per tool might be needed, but generic approach:
		// Return the last word of the first line often works, or the whole first line
		// For now return raw cleaned string, refinement needed per tool
		return cleanVersionString(lines[0])
	}
	return "Detected"
}

func cleanVersionString(s string) string {
	// Simple heuristic: take the version-looking part
	// This is a placeholder for more robust regex parsing
	fields := strings.Fields(s)
	if len(fields) > 1 {
		// Try to find the one starting with a number
		for _, f := range fields {
			if len(f) > 0 && (f[0] >= '0' && f[0] <= '9') || f[0] == 'v' {
				return f
			}
		}
		return fields[len(fields)-1]
	}
	return s
}
