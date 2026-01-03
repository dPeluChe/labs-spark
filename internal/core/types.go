package core

// UpdateMethod defines how a tool is updated
type UpdateMethod string

const (
	MethodBrew      UpdateMethod = "brew"
	MethodNpmSys    UpdateMethod = "npm_sys"
	MethodNpmPkg    UpdateMethod = "npm_pkg"
	MethodBrewPkg   UpdateMethod = "brew_pkg"
	MethodMacApp    UpdateMethod = "mac_app"
	MethodClaude    UpdateMethod = "claude"
	MethodDroid     UpdateMethod = "droid"
	MethodToad      UpdateMethod = "toad"
	MethodOpencode  UpdateMethod = "opencode"
	MethodOmz       UpdateMethod = "omz"
	MethodManual    UpdateMethod = "manual" // For tools like Antigravity
)

// Category groups tools logically
type Category string

const (
	CategoryCode    Category = "CODE"    // AI Tools
	CategoryTerm    Category = "TERM"    // Terminals
	CategoryIDE     Category = "IDE"     // IDEs
	CategoryProd    Category = "PROD"    // Productivity
	CategoryInfra   Category = "INFRA"   // Infrastructure
	CategoryUtils   Category = "UTILS"   // Utilities
	CategoryRuntime Category = "RUNTIME" // High Risk Runtimes
	CategorySys     Category = "SYS"     // System Managers
)

// Tool represents a software component managed by Spark
type Tool struct {
	ID          string       // Unique internal ID (S-01, etc.)
	Name        string       // Display Name (e.g., "Claude CLI")
	Binary      string       // Binary command (e.g., "claude") or App Name
	Package     string       // Package name (e.g., "@anthropic-ai/claude-code")
	Category    Category     // Grouping category
	Method      UpdateMethod // How to update it
	Description string       // Optional description
}

// ToolStatus represents the current state of a tool
type ToolStatus int

const (
	StatusChecking    ToolStatus = iota // Currently checking version
	StatusInstalled                     // Installed and up to date
	StatusOutdated                      // Update available
	StatusMissing                       // Not installed
	StatusUnmanaged                     // Installed but not managed by Spark/Brew
	StatusManualCheck                   // Requires manual intervention
	StatusUpdating                      // Update in progress
	StatusUpdated                       // Successfully updated
	StatusFailed                        // Update failed
)

// ToolState holds the runtime data for a tool
type ToolState struct {
	Tool         Tool
	Status       ToolStatus
	LocalVersion string
	RemoteVersion string
	Message       string // Error message or status detail
}
