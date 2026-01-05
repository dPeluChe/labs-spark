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
	// Special handling for macOS applications
	if t.Method == core.MethodMacApp {
		return d.getMacAppVersion(t.Binary)
	}

	// Special handling for Oh My Zsh (git-based)
	if t.Binary == "omz" {
		return d.getOmzVersion()
	}

	// Special handling for Antigravity (multiple paths)
	if t.Binary == "antigravity" {
		return d.getAntigravityVersion()
	}

	// Generic CLI tool detection
	return d.getCliToolVersion(t)
}

// getMacAppVersion detects version of macOS .app bundles
func (d *Detector) getMacAppVersion(binary string) string {
	appPaths := map[string]string{
		"iterm":    "/Applications/iTerm.app",
		"ghostty":  "/Applications/Ghostty.app",
		"warp":     "/Applications/Warp.app",
		"code":     "/Applications/Visual Studio Code.app",
		"cursor":   "/Applications/Cursor.app",
		"zed":      "/Applications/Zed.app",
		"windsurf": "/Applications/Windsurf.app",
		"docker":   "/Applications/Docker.app",
	}

	appPath, ok := appPaths[binary]
	if !ok {
		return "MISSING"
	}

	return getAppVersion(appPath)
}

// getOmzVersion gets Oh My Zsh git commit hash
func (d *Detector) getOmzVersion() string {
	omzPath := os.Getenv("HOME") + "/.oh-my-zsh"
	if _, err := os.Stat(omzPath); err != nil {
		return "MISSING"
	}

	cmd := exec.Command("git", "--git-dir="+omzPath+"/.git", "--work-tree="+omzPath, "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "Installed"
	}
	return strings.TrimSpace(string(out))
}

// getAntigravityVersion checks multiple possible installation paths
func (d *Detector) getAntigravityVersion() string {
	customPath := os.Getenv("HOME") + "/.antigravity/antigravity/bin/antigravity"
	if _, err := os.Stat(customPath); err == nil {
		output := runCmd(customPath, "--version")
		return ParseToolSpecificVersion("antigravity", output)
	}
	output := runCmd("antigravity", "--version")
	return ParseToolSpecificVersion("antigravity", output)
}

// getCliToolVersion detects version for standard CLI tools
func (d *Detector) getCliToolVersion(t core.Tool) string {
	// Check if binary exists in PATH
	path, err := exec.LookPath(t.Binary)
	if err != nil {
		return "MISSING"
	}
	_ = path

	// Run --version command
	output := runCmd(t.Binary, "--version")
	if output == "MISSING" {
		return "MISSING"
	}

	// Use tool-specific parser if available, otherwise generic
	return ParseToolSpecificVersion(t.Binary, output)
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
	
	// Use robust version parser
	output := strings.TrimSpace(out.String())
	return CleanVersionString(output)
}
