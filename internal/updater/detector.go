package updater

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/dpeluche/spark/internal/core"
)

// Detector handles version checking logic
type Detector struct {
	cacheMutex    sync.RWMutex
	outdatedCache map[string]string // Package Name -> Latest Version
	hasWarmedUp   bool
}

func NewDetector() *Detector {
	return &Detector{
		outdatedCache: make(map[string]string),
	}
}

// WarmUpCache fetches brew info once to speed up subsequent checks
func (d *Detector) WarmUpCache() {
	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	if d.hasWarmedUp {
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Fetch Brew Outdated
	go func() {
		defer wg.Done()
		d.fetchBrewOutdated()
	}()

	// Fetch NPM Outdated
	go func() {
		defer wg.Done()
		d.fetchNpmOutdated()
	}()

	wg.Wait()
	d.hasWarmedUp = true
}

type brewOutdatedItem struct {
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"` // This is actually the "latest" available in brew formulae usually?
	// Brew JSON output for outdated:
	// [{"name":"fzf","installed_versions":["0.45.0"],"current_version":"0.46.0",...}]
}

func (d *Detector) fetchBrewOutdated() {
	// brew outdated --json=v2
	cmd := exec.Command("brew", "outdated", "--json=v2")
	var out bytes.Buffer
	cmd.Stdout = &out
	// Ignore errors, brew outdated returns non-zero if outdated items exist
	_ = cmd.Run()

	var data struct {
		Formulae []brewOutdatedItem `json:"formulae"`
		Casks    []brewOutdatedItem `json:"casks"`
	}

	if err := json.Unmarshal(out.Bytes(), &data); err == nil {
		for _, item := range data.Formulae {
			d.outdatedCache[item.Name] = item.CurrentVersion
		}
		for _, item := range data.Casks {
			d.outdatedCache[item.Name] = item.CurrentVersion
		}
	}
}

type npmOutdatedItem struct {
	Current  string `json:"current"`
	Wanted   string `json:"wanted"`
	Latest   string `json:"latest"`
	Location string `json:"location"`
}

func (d *Detector) fetchNpmOutdated() {
	// npm outdated -g --json
	cmd := exec.Command("npm", "outdated", "-g", "--json")
	var out bytes.Buffer
	cmd.Stdout = &out
	// Ignore errors
	_ = cmd.Run()

	var data map[string]npmOutdatedItem
	if err := json.Unmarshal(out.Bytes(), &data); err == nil {
		for pkg, info := range data {
			d.outdatedCache[pkg] = info.Latest
		}
	}
}

func (d *Detector) GetRemoteVersion(t core.Tool, localVersion string) string {
	d.cacheMutex.RLock()
	defer d.cacheMutex.RUnlock()

	// If we haven't warmed up or cache is empty, we might return "Unknown" or force a check.
	// But assuming WarmUp runs first.

	if localVersion == "MISSING" {
		return "Unknown" // We don't check for uninstalled tools yet
	}

	// If checking a package
	if latest, ok := d.outdatedCache[t.Package]; ok {
		return latest
	}

	// If not in outdated list, and we have a local version,
	// it usually means Local is Latest.
	if d.hasWarmedUp {
		return localVersion
	}

	return "Checking..."
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
	// 1. Try finding binary in PATH
	path, err := exec.LookPath(t.Binary)
	if err == nil && path != "" {
		// Try standard --version
		output := runCmd(t.Binary, "--version")
		if output != "" && output != "MISSING" && output != "Unknown" {
			return ParseToolSpecificVersion(t.Binary, output)
		}
		
		// ... (version / -v checks) ...
	}

	// 1.5 Fallback: Check ~/.local/bin explicitly (Common for Toad, Droid, Python tools)
	home := os.Getenv("HOME")
	localBin := home + "/.local/bin/" + t.Binary
	if _, err := os.Stat(localBin); err == nil {
		output := runCmd(localBin, "--version")
		if output != "" && output != "MISSING" {
			return ParseToolSpecificVersion(t.Binary, output)
		}
	}

	// 2. Fallback: Check NPM Global List (if it's an NPM tool)
	if t.Method == core.MethodNpmPkg || t.Method == core.MethodNpmSys || t.Package != "" {
		// ... existing npm logic ...
		cmd := exec.Command("npm", "list", "-g", "--depth=0", "--json", t.Package)
		out, err := cmd.Output()
		if err == nil {
			outStr := string(out)
			if strings.Contains(outStr, "\"version\":") {
				parts := strings.Split(outStr, "\"version\":")
				if len(parts) > 1 {
					ver := strings.Split(parts[1], "\"")[1]
					return CleanVersionString(ver)
				}
			}
		}
	}

	// 3. Fallback: Check Homebrew explicitly (if it's a Brew tool)
	if t.Method == core.MethodBrew || t.Method == core.MethodBrewPkg {
		// brew list --versions <package>
		// Output: "kubernetes-cli 1.28.2"
		cmd := exec.Command("brew", "list", "--versions", t.Package)
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			fields := strings.Fields(string(out))
			if len(fields) >= 2 {
				// The version is usually the second field
				return CleanVersionString(fields[len(fields)-1])
			}
		}
	}

	return "MISSING"
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
	// Increase timeout to 5 seconds for slower tools (AI runtimes)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		// If fails, return empty to signal caller to try fallback or return MISSING
		return ""
	}
	
	// Some tools print version to stderr (e.g. java sometimes, or python)
	output := strings.TrimSpace(out.String())
	if output == "" {
		output = strings.TrimSpace(stderr.String())
	}
	
	return CleanVersionString(output)
}
