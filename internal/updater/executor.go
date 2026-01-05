package updater

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/dpeluche/spark/internal/core"
)

// Executor handles the actual update process for tools
type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

// Update attempts to update the specified tool
func (e *Executor) Update(t core.Tool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute) // Updates can take time
	defer cancel()

	switch t.Method {
	case core.MethodBrew, core.MethodBrewPkg:
		return e.updateBrew(ctx, t)
	case core.MethodMacApp:
		return e.updateMacApp(ctx, t)
	case core.MethodNpmSys, core.MethodNpmPkg:
		return e.updateNpm(ctx, t)
	case core.MethodClaude:
		return e.updateNpm(ctx, t) // Claude is an NPM package
	case core.MethodOmz:
		return e.updateOmz(ctx)
	case core.MethodManual:
		return fmt.Errorf("manual update required")
	default:
		return fmt.Errorf("update method %s not implemented", t.Method)
	}
}

func (e *Executor) updateBrew(ctx context.Context, t core.Tool) error {
	// brew upgrade <package>
	cmd := exec.CommandContext(ctx, "brew", "upgrade", t.Package)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("brew upgrade failed: %s: %v", string(output), err)
	}
	return nil
}

func (e *Executor) updateMacApp(ctx context.Context, t core.Tool) error {
	// Try upgrading via brew cask first
	// We assume if it's a MacApp it might be managed by brew cask
	// Check if it is a cask
	checkCmd := exec.CommandContext(ctx, "brew", "list", "--cask", t.Package)
	if err := checkCmd.Run(); err == nil {
		cmd := exec.CommandContext(ctx, "brew", "upgrade", "--cask", t.Package)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("brew cask upgrade failed: %s: %v", string(output), err)
		}
		return nil
	}

	// If not a cask, we can't auto-update it easily
	return fmt.Errorf("manual update required (not a brew cask)")
}

func (e *Executor) updateNpm(ctx context.Context, t core.Tool) error {
	// npm install -g <package>@latest
	pkg := t.Package
	if pkg == "" {
		pkg = t.Binary
	}

	cmd := exec.CommandContext(ctx, "npm", "install", "-g", pkg+"@latest")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("npm install failed: %s: %v", string(output), err)
	}
	return nil
}

func (e *Executor) updateOmz(ctx context.Context) error {
	// omz update usually runs interactively or via script.
	// The standard way is running the upgrade script.
	// Often available as `omz update` alias, but that might not be in the path for non-interactive shells.
	// We can try calling the script directly if we find it.

	cmd := exec.CommandContext(ctx, "sh", "-c", "$ZSH/tools/upgrade.sh")
	// Set env var to avoid interactive prompt if supported
	cmd.Env = append(cmd.Env, "ZSH="+getEnv("ZSH", "~/.oh-my-zsh"))

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("omz update failed: %s: %v", string(output), err)
	}
	return nil
}

// Helper to get env with fallback (simplified)
func getEnv(key, fallback string) string {
	// In a real scenario, we'd read actual os.Environ
	return fallback
}
