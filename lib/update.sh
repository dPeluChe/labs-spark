#!/bin/bash

perform_update() {
    local method=$1
    local name=$2
    local pkg=$3
    local current=$4
    local target=$5

    # If in OPEN_LINKS mode, silence non-relevant tools
    if [[ "$METHOD_MODE" == "OPEN_LINKS" ]]; then
        if [[ "$method" == "mac_app" ]] && [[ "$target" == *"External"* ]]; then
             echo -e "${CYAN}   ↗ Opening download page for $name...${RESET}"
             brew info --cask "$pkg" 2>/dev/null | grep "https://" | head -n 1 | xargs open
             return
        elif [[ "$method" == "antigravity" ]]; then
             echo -e "${CYAN}   ↗ Opening Antigravity portal...${RESET}"
             open "https://www.antigravity.ai" 2>/dev/null || echo "   Could not open link."
             return
        else
            return 
        fi
    fi

    if [[ "$target" == "Unmanaged" ]]; then
        # Silent in targeted mode unless it's the target
        if [[ -z "$TARGET_TOOL_ID" ]]; then
            echo -e "${DIM}   ○ $name is not managed by Brew. Skipping.${RESET}"
        fi
        return
    fi
    
    if [[ "$target" == "Manual" ]]; then
        if [[ -z "$TARGET_TOOL_ID" ]]; then
             echo -e "${DIM}   ○ $name requires manual update. Skipping.${RESET}"
        fi
        return
    fi

    if [[ "$target" == "$current" ]]; then
         if [[ -z "$TARGET_TOOL_ID" ]]; then
             echo -e "${DIM}   ○ $name is already up to date. Skipped.${RESET}"
         elif [[ -n "$TARGET_TOOL_ID" ]] && [[ "$name" == "$TARGET_TOOL_NAME" ]]; then
             # If specifically targeted but up to date, tell the user
             echo -e "${GREEN}   ✔ $name is already up to date.${RESET}"
         fi
         return
    fi

    echo -e "${BOLD}${CYAN}⚡ Updating $name...${RESET}"

    local success=0
    case $method in
        brew) brew update && brew upgrade && brew cleanup && success=1 ;; 
        npm_sys) npm update -g && success=1 ;; 
        npm_pkg) npm install -g "$pkg@latest" && success=1 ;; 
        claude) 
            # Support Homebrew Cask, curl, and npm installations
            if brew list --cask claude-code &>/dev/null; then
                brew upgrade --cask claude-code && success=1
            elif [ -f "$HOME/.claude/local/claude" ]; then
                # Curl installation - reinstall via curl
                curl -fsSL https://claude.ai/install.sh | bash && success=1
            else
                # NPM installation
                npm install -g "$pkg@latest" && success=1
            fi
            ;; 
        droid) curl -fsSL https://app.factory.ai/cli | sh && success=1 ;; 
        toad) curl -fsSL https://batrachian.ai/install | sh && success=1 ;; 
        opencode) (opencode upgrade || curl -fsSL https://opencode.ai/install | bash) && success=1 ;; 
        omz) (cd ~/.oh-my-zsh && git pull) && success=1 ;; 
        brew_pkg) (brew upgrade "$pkg" 2>/dev/null || echo -e "     ${YELLOW}No update needed or package not pinned.${RESET}") && success=1 ;; 
        mac_app) 
            if brew list --cask "$pkg" &>/dev/null; then
                brew upgrade --cask "$pkg" && success=1
            else
                echo -e "${YELLOW}   ! $name is not managed by Homebrew.${RESET}"
            fi
            ;; 
        antigravity)
            echo -e "   Please use the internal updater within Antigravity."
            ;; 
        *) echo "   No update method found." ;; 
    esac

    if [ $success -eq 1 ]; then
        echo -e "${GREEN}   ✔ Success${RESET}\n"
        UPDATED_TOOLS+=("$name ($current -> $target)")
    else
        echo -e "${RED}   ✘ Error updating $name${RESET}\n"
    fi
}
