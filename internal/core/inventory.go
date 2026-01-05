package core

import "fmt"

// GetInventory returns the master list of supported tools
func GetInventory() []Tool {
	// Helper to generate IDs sequentially could be added, but manual for now is fine
	// or we can auto-assign IDs at runtime.
	
	tools := []Tool{
		// AI Development
		{Name: "Claude CLI", Binary: "claude", Package: "@anthropic-ai/claude-code", Category: CategoryCode, Method: MethodClaude},
		{Name: "Droid CLI", Binary: "droid", Package: "factory-cli", Category: CategoryCode, Method: MethodDroid},
		{Name: "Gemini CLI", Binary: "gemini", Package: "@google/gemini-cli", Category: CategoryCode, Method: MethodNpmPkg},
		{Name: "OpenCode", Binary: "opencode", Package: "opencode-ai", Category: CategoryCode, Method: MethodOpencode},
		{Name: "Codex CLI", Binary: "codex", Package: "@openai/codex", Category: CategoryCode, Method: MethodNpmPkg},
		{Name: "Crush CLI", Binary: "crush", Package: "crush", Category: CategoryCode, Method: MethodBrewPkg},
		{Name: "Toad CLI", Binary: "toad", Package: "batrachian-toad", Category: CategoryCode, Method: MethodToad},
		{Name: "Ollama", Binary: "ollama", Package: "ollama", Category: CategoryCode, Method: MethodManual},

		// Terminal Emulators
		{Name: "iTerm2", Binary: "iterm", Package: "iterm2", Category: CategoryTerm, Method: MethodMacApp},
		{Name: "Ghostty", Binary: "ghostty", Package: "ghostty", Category: CategoryTerm, Method: MethodMacApp},
		{Name: "Warp Terminal", Binary: "warp", Package: "warp", Category: CategoryTerm, Method: MethodMacApp},

		// IDEs
		{Name: "VS Code", Binary: "code", Package: "visual-studio-code", Category: CategoryIDE, Method: MethodMacApp},
		{Name: "Cursor IDE", Binary: "cursor", Package: "cursor", Category: CategoryIDE, Method: MethodMacApp},
		{Name: "Zed Editor", Binary: "zed", Package: "zed", Category: CategoryIDE, Method: MethodMacApp},
		{Name: "Windsurf", Binary: "windsurf", Package: "windsurf", Category: CategoryIDE, Method: MethodMacApp},
		{Name: "Antigravity", Binary: "antigravity", Package: "antigravity", Category: CategoryIDE, Method: MethodManual},

		// Productivity
		{Name: "JQ", Binary: "jq", Package: "jq", Category: CategoryProd, Method: MethodBrewPkg},
		{Name: "FZF", Binary: "fzf", Package: "fzf", Category: CategoryProd, Method: MethodBrewPkg},
		{Name: "Ripgrep", Binary: "rg", Package: "ripgrep", Category: CategoryProd, Method: MethodBrewPkg},
		{Name: "Bat", Binary: "bat", Package: "bat", Category: CategoryProd, Method: MethodBrewPkg},
		{Name: "HTTPie", Binary: "http", Package: "httpie", Category: CategoryProd, Method: MethodBrewPkg},
		{Name: "LazyGit", Binary: "lazygit", Package: "lazygit", Category: CategoryProd, Method: MethodBrewPkg},
		{Name: "TLDR", Binary: "tldr", Package: "tldr", Category: CategoryProd, Method: MethodBrewPkg},

		// Infrastructure
		{Name: "Docker Desktop", Binary: "docker", Package: "docker", Category: CategoryInfra, Method: MethodMacApp},
		{Name: "Kubernetes CLI", Binary: "kubectl", Package: "kubernetes-cli", Category: CategoryInfra, Method: MethodBrewPkg},
		{Name: "Helm", Binary: "helm", Package: "helm", Category: CategoryInfra, Method: MethodBrewPkg},
		{Name: "Terraform", Binary: "terraform", Package: "terraform", Category: CategoryInfra, Method: MethodBrewPkg},
		{Name: "AWS CLI", Binary: "aws", Package: "awscli", Category: CategoryInfra, Method: MethodBrewPkg},
		{Name: "Ngrok", Binary: "ngrok", Package: "ngrok", Category: CategoryInfra, Method: MethodBrewPkg},

		// Utilities
		{Name: "Oh My Zsh", Binary: "omz", Package: "oh-my-zsh", Category: CategoryUtils, Method: MethodOmz},
		{Name: "Zellij", Binary: "zellij", Package: "zellij", Category: CategoryUtils, Method: MethodBrewPkg},
		{Name: "Tmux", Binary: "tmux", Package: "tmux", Category: CategoryUtils, Method: MethodBrewPkg},
		{Name: "Git", Binary: "git", Package: "git", Category: CategoryUtils, Method: MethodBrewPkg},
		{Name: "Bash", Binary: "bash", Package: "bash", Category: CategoryUtils, Method: MethodBrewPkg},
		{Name: "SQLite", Binary: "sqlite3", Package: "sqlite", Category: CategoryUtils, Method: MethodBrewPkg},
		{Name: "Watchman", Binary: "watchman", Package: "watchman", Category: CategoryUtils, Method: MethodBrewPkg},
		{Name: "Direnv", Binary: "direnv", Package: "direnv", Category: CategoryUtils, Method: MethodBrewPkg},
		{Name: "Heroku CLI", Binary: "heroku", Package: "heroku", Category: CategoryUtils, Method: MethodBrewPkg},
		{Name: "Pre-commit", Binary: "pre-commit", Package: "pre-commit", Category: CategoryUtils, Method: MethodBrewPkg},

		// Runtimes
		{Name: "Node.js", Binary: "node", Package: "node", Category: CategoryRuntime, Method: MethodBrewPkg},
		{Name: "Python 3.13", Binary: "python3", Package: "python@3.13", Category: CategoryRuntime, Method: MethodBrewPkg},
		{Name: "Go Lang", Binary: "go", Package: "go", Category: CategoryRuntime, Method: MethodBrewPkg},
		{Name: "Ruby", Binary: "ruby", Package: "ruby", Category: CategoryRuntime, Method: MethodBrewPkg},
		{Name: "PostgreSQL 16", Binary: "psql", Package: "postgresql@16", Category: CategoryRuntime, Method: MethodBrewPkg},

		// System
		{Name: "Homebrew Core", Binary: "brew", Package: "homebrew", Category: CategorySys, Method: MethodBrewPkg},
		{Name: "NPM Globals", Binary: "npm", Package: "npm", Category: CategorySys, Method: MethodNpmSys},
	}

	// Assign IDs automatically S-01, S-02, etc.
	for i := range tools {
		tools[i].ID = fmt.Sprintf("S-%02d", i+1)
	}

	return tools
}
