#!/bin/bash

# SPARK v0.2.4 - Surgical Precision CLI Updater
# Codenamed: Spark (The life-force of Transformers)

# --- Configuration & Styling ---
BOLD="\033[1m"
DIM="\033[2m"
RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[34m"
MAGENTA="\033[35m"
CYAN="\033[36m"
RESET="\033[0m"

# Tools Configuration
# Format: "CATEGORY:BinaryName:PackageName:DisplayName:UpdateMethod"
TOOLS=(
    # AI Development
    "CODE:claude:@anthropic-ai/claude-code:Claude CLI:npm_pkg"
    "CODE:droid:factory-cli:Droid CLI:droid"
    "CODE:gemini:@google/gemini-cli:Gemini CLI:npm_pkg"
    "CODE:opencode:opencode-ai:OpenCode:opencode"
    "CODE:codex:@openai/codex:Codex CLI:npm_pkg"
    "CODE:crush:crush:Crush CLI:brew_pkg"

    # Terminal Emulators
    "TERM:iterm:iterm2:iTerm2:mac_app"
    "TERM:ghostty:ghostty:Ghostty:mac_app"
    "TERM:warp:warp:Warp Terminal:mac_app"

    # Safe Utilities (Low Risk)
    "UTILS:zellij:zellij:Zellij:brew_pkg"
    "UTILS:tmux:tmux:Tmux:brew_pkg"
    "UTILS:git:git:Git:brew_pkg"
    "UTILS:bash:bash:Bash:brew_pkg"
    "UTILS:sqlite3:sqlite:SQLite:brew_pkg"
    "UTILS:watchman:watchman:Watchman:brew_pkg"
    "UTILS:direnv:direnv:Direnv:brew_pkg"
    "UTILS:heroku:heroku:Heroku CLI:brew_pkg"
    "UTILS:pre-commit:pre-commit:Pre-commit:brew_pkg"

    # Critical Runtimes (High Risk)
    "RUNTIME:node:node:Node.js:brew_pkg"
    "RUNTIME:python3:python@3.13:Python 3.13:brew_pkg"
    "RUNTIME:go:go:Go Lang:brew_pkg"
    "RUNTIME:ruby:ruby:Ruby:brew_pkg"
    "RUNTIME:psql:postgresql@16:PostgreSQL 16:brew_pkg"

    # System Managers
    "SYS:brew:homebrew:Homebrew Core:brew"
    "SYS:npm:npm:NPM Globals:npm_sys"
)

# Global Storage
CODE_UPDATES_COUNT=0
TERM_UPDATES_COUNT=0
UTILS_UPDATES_COUNT=0
RUNTIME_UPDATES_COUNT=0
SYS_UPDATES_COUNT=0
UPDATED_TOOLS=()
BREW_CACHE=""

# --- Helper Functions ---

banner() {
    clear
    echo -e "${CYAN}${BOLD}"
    echo "   _____ ____  ___  ____  __ __"
    echo "  / ___// __ \/   |/ __ \/ //_/"
    echo "  \__ \/ /_/ / /| / /_/ / ,<   "
    echo " ___/ / ____/ ___ / _, _/ /| |  "
    echo "/____/_/   /_/  |/_/ |_/_/ |_|  "
    echo -e "${RESET}"
    echo -e "${BLUE}  Surgical Precision Update Utility v0.2.4${RESET}"
    echo -e "${DIM}  ========================================${RESET}\n"
}

get_local_version() {
    local binary=$1
    
    # Special handling for macOS Apps
    if [[ "$binary" == "iterm" ]]; then
        [ -d "/Applications/iTerm.app" ] && defaults read /Applications/iTerm.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "ghostty" ]]; then
        [ -d "/Applications/Ghostty.app" ] && defaults read /Applications/Ghostty.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "warp" ]]; then
        [ -d "/Applications/Warp.app" ] && defaults read /Applications/Warp.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    fi

    if ! command -v "$binary" &> /dev/null; then
        echo "MISSING"
        return
    fi

    local ver=""
    if [[ "$binary" == "brew" ]]; then
        ver=$(brew --version | head -n 1 | awk '{print $2}')
    elif [[ "$binary" == "npm" ]]; then
        ver=$(npm --version)
    elif [[ "$binary" == "claude" ]]; then
         ver=$($binary --version 2>/dev/null | awk '{print $NF}')
         if [[ -z "$ver" ]]; then
            ver=$(npm list -g @anthropic-ai/claude-code --depth=0 2>/dev/null | grep claude-code | awk -F@ '{print $NF}')
         fi
    elif [[ "$binary" == "droid" ]]; then
        ver="Installed"
    elif [[ "$binary" == "opencode" ]]; then
        ver=$(opencode --version 2>/dev/null | head -n 1 | awk '{print $3}') || ver="Installed"
    elif [[ "$binary" == "gemini" ]]; then
         ver=$(npm list -g @google/gemini-cli --depth=0 2>/dev/null | grep gemini-cli | awk -F@ '{print $NF}') || ver="Unknown"
    elif [[ "$binary" == "codex" ]]; then
         ver=$(npm list -g @openai/codex --depth=0 2>/dev/null | grep codex | awk -F@ '{print $NF}') || ver="Unknown"
    elif [[ "$binary" == "sqlite3" ]]; then
         ver=$(sqlite3 --version | awk '{print $1}')
    else
        ver=$($binary --version 2>/dev/null | head -n 1 | awk '{print $NF}') || ver="Detected"
    fi
    
    if [[ -z "$ver" ]]; then ver="Detected"; fi
    echo "$ver"
}

get_remote_version() {
    local method=$1
    local package=$2
    local local_ver=$3
    
    if [[ "$method" == "npm_pkg" ]] || [[ "$method" == "npm_sys" ]]; then
        npm view "$package" version 2>/dev/null
    elif [[ "$method" == "brew_pkg" ]] || [[ "$method" == "mac_app" ]]; then
        # Intelligent Brew Check using pre-fetched cache
        # grep format from 'brew outdated --verbose': "package (old) < new"
        local update_info=$(echo "$BREW_CACHE" | grep "^$package ")
        
        if [[ -n "$update_info" ]]; then
            # Extract the last field (target version)
            echo "$update_info" | awk '{print $NF}'
        else
            # Not in outdated list = Up to date
            echo "$local_ver"
        fi
    else
        echo "Latest"
    fi
}

check_active_sessions() {
    local active_found=0
    echo -e "${DIM}Checking for active sessions...${RESET}"
    
    for tool_entry in "${TOOLS[@]}"; do
        IFS=':' read -r category binary pkg display method <<< "$tool_entry"
        if command -v "$binary" &> /dev/null || [[ "$method" == "mac_app" ]]; then
            local proc=$binary
            [[ "$binary" == "iterm" ]] && proc="iTerm2"
            [[ "$binary" == "warp" ]] && proc="Warp"
            [[ "$binary" == "ghostty" ]] && proc="Ghostty"
            [[ "$binary" == "python3" ]] && proc="python"
            
            if pgrep -fi "$proc" > /dev/null;
 then
                if [[ $active_found -eq 0 ]]; then
                    echo -e "${YELLOW}${BOLD}⚠️  Active Sessions Detected:${RESET}"
                    active_found=1
                fi
                echo -e "   - $display is currently running"
            fi
        fi
    done
    
    if [[ $active_found -eq 1 ]]; then
        echo -e "${YELLOW}   Updating running tools may cause interruptions.${RESET}\n"
    fi
}

analyze_system() {
    echo -e "${BOLD}Analyzing System Components...${RESET}"
    
    # Pre-fetch Homebrew data
    echo -e "${DIM}   Fetching Homebrew intelligence (this may take a moment)...${RESET}"
    BREW_CACHE=$(brew outdated --verbose)
    
    printf "${BOLD}%-4s %-18s %-15s %-15s${RESET}\n" "Sts" "Tool" "Current" "Target"
    echo "--------------------------------------------------------"

    CODE_UPDATES_COUNT=0
    TERM_UPDATES_COUNT=0
    UTILS_UPDATES_COUNT=0
    RUNTIME_UPDATES_COUNT=0
    SYS_UPDATES_COUNT=0

    print_group() {
        local target_cat=$1
        local title=$2
        echo -e "${DIM}--- $title ---${RESET}"
        
        for tool_entry in "${TOOLS[@]}"; do
            IFS=':' read -r category binary pkg display method <<< "$tool_entry"
            if [[ "$category" == "$target_cat" ]]; then
                local current=$(get_local_version "$binary")
                local icon=""
                local target="-"
                local color=""
                local needs_update=0

                if [[ "$current" == "MISSING" ]]; then
                    icon="${DIM}○${RESET}"
                    color="${DIM}"
                    current="Not Installed"
                else
                    icon="${GREEN}●${RESET}"
                    color="${RESET}"
                    target=$(get_remote_version "$method" "$pkg" "$current")
                    
                    if [[ "$target" == "Latest" ]]; then
                        needs_update=0 
                    elif [[ "$target" != "-" ]] && [[ "$current" != "$target" ]]; then
                        needs_update=1
                        color="${YELLOW}"
                        icon="${YELLOW}↑${RESET}"
                    fi
                fi
                
                if [[ $needs_update -eq 1 ]]; then
                    [[ "$category" == "CODE" ]] && ((CODE_UPDATES_COUNT++))
                    [[ "$category" == "TERM" ]] && ((TERM_UPDATES_COUNT++))
                    [[ "$category" == "UTILS" ]] && ((UTILS_UPDATES_COUNT++))
                    [[ "$category" == "RUNTIME" ]] && ((RUNTIME_UPDATES_COUNT++))
                    [[ "$category" == "SYS" ]] && ((SYS_UPDATES_COUNT++))
                fi

                printf "% -13b ${color}%-18s %-15s ${MAGENTA}%-15s${RESET}\n" "$icon" "$display" "$current" "$target"
            fi
        done
    }

    print_group "CODE" "AI Development Tools"
    echo ""
    print_group "TERM" "Terminal Emulators"
    echo ""
    print_group "UTILS" "Safe Utilities"
    echo ""
    print_group "RUNTIME" "Critical Runtimes (High Risk)"
    echo ""
    print_group "SYS" "System Managers"
    echo ""
}

perform_update() {
    local method=$1
    local name=$2
    local pkg=$3
    local current=$4
    local target=$5

    if [[ "$target" != "Latest" ]] && [[ "$target" != "-" ]] && [[ "$current" == "$target" ]]; then
         echo -e "${DIM}   ○ $name is up to date ($current). Skipped.${RESET}"
         return
    fi

    echo -e "${BOLD}${CYAN}⚡ Updating $name...${RESET}"

    local success=0
    case $method in
        brew) brew update && brew upgrade && brew cleanup && success=1 ;;
        npm_sys) npm update -g && success=1 ;;
        npm_pkg) npm install -g "$pkg@latest" && success=1 ;;
        droid) curl -fsSL https://app.factory.ai/cli | sh && success=1 ;;
        opencode) (opencode upgrade || curl -fsSL https://opencode.ai/install | bash) && success=1 ;;
        brew_pkg) (brew upgrade "$pkg" 2>/dev/null || echo -e "     ${YELLOW}No update needed or package not pinned.${RESET}") && success=1 ;;
        mac_app)
            if brew list --cask "$pkg" &>/dev/null;
 then
                brew upgrade --cask "$pkg" && success=1
            else
                echo -e "${YELLOW}   ! $name is not managed by Homebrew.${RESET}"
            fi
            ;; 
        *) echo "   No update method found." ;; 
    esac

    if [ $success -eq 1 ]; then
        echo -e "${GREEN}   ✔ Success${RESET}\n"
        local ver_msg="$current -> $target"
        UPDATED_TOOLS+=("$name ($ver_msg)")
    else
        echo -e "${RED}   ✘ Error updating $name${RESET}\n"
    fi
}

show_summary() {
    echo -e "${BOLD}--- SPARK UPDATE SUMMARY ---${RESET}"
    if [ ${#UPDATED_TOOLS[@]} -eq 0 ]; then
        echo -e "${DIM}No tools were updated.${RESET}"
    else
        for tool in "${UPDATED_TOOLS[@]}"; do
            echo -e "${GREEN}[✔] $tool${RESET}"
        done
    fi
    echo ""
}

# --- Main Logic ---
banner
analyze_system
check_active_sessions

echo -e "${BOLD}Update Modes:${RESET}"
echo -e "  ${BOLD}[1]${RESET} ${CYAN}AI & Terminals${RESET}     (Code Tools + iTerm/Warp)"
echo -e "  ${BOLD}[2]${RESET} ${GREEN}Utilities${RESET}          (Git, Tmux, Zellij, etc.)"
echo -e "  ${BOLD}[3]${RESET} ${RED}Runtimes${RESET}           (Node, Python, Go, Postgres) ${RED}⚠️${RESET}"
echo -e "  ${BOLD}[4]${RESET} ${YELLOW}Full System${RESET}        (Everything included)"
echo -e "  ${BOLD}[5]${RESET} Exit"
echo ""

read -p "Select option [1]: " mode
mode=${mode:-1}

if [[ "$mode" == "5" ]]; then echo "Bye!"; exit 0; fi

TARGET_CATEGORY=""
if [[ "$mode" == "1" ]]; then TARGET_CATEGORY="CODE_TERM"; fi
if [[ "$mode" == "2" ]]; then TARGET_CATEGORY="UTILS"; fi
if [[ "$mode" == "3" ]]; then TARGET_CATEGORY="RUNTIME"; fi
if [[ "$mode" == "4" ]]; then TARGET_CATEGORY="ALL"; fi

# Safety Check for Runtimes
if [[ "$TARGET_CATEGORY" == "RUNTIME" ]] || [[ "$TARGET_CATEGORY" == "ALL" ]]; then
    echo -e "\n${RED}${BOLD}⚠️  WARNING: You are about to update critical runtimes (Node, Python, DBs).${RESET}"
    echo -e "${RED}    This might break existing projects or virtual environments.${RESET}"
    read -p "    Are you absolutely sure? (type 'yes' to proceed): " confirm
    if [[ "$confirm" != "yes" ]]; then
        echo "    Operation aborted by user."
        exit 0
    fi
fi

echo -e "\n${BOLD}Starting Update Sequence...${RESET}\n"

for tool_entry in "${TOOLS[@]}"; do
    IFS=':' read -r category binary pkg display method <<< "$tool_entry"
    
    # Matching Logic
    MATCH=0
    if [[ "$TARGET_CATEGORY" == "ALL" ]]; then MATCH=1; fi
    if [[ "$TARGET_CATEGORY" == "CODE_TERM" ]] && ([[ "$category" == "CODE" ]] || [[ "$category" == "TERM" ]]); then MATCH=1; fi
    if [[ "$category" == "$TARGET_CATEGORY" ]]; then MATCH=1; fi

    if [[ $MATCH -eq 1 ]]; then
        if command -v "$binary" &> /dev/null || [[ "$method" == "mac_app" ]]; then
            current=$(get_local_version "$binary")
            target=$(get_remote_version "$method" "$pkg" "$current")
            perform_update "$method" "$display" "$pkg" "$current" "$target"
        fi
    fi
done

show_summary
echo -e "${BOLD}${GREEN}✨ Spark Sequence Complete.${RESET}"