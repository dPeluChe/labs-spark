#!/bin/bash

# SPARK v0.1.1 - Intelligent CLI Updater
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
# Categories: CODE (AI/Dev tools), SYS (System managers like brew/npm)
TOOLS=(
    "CODE:claude:@anthropic-ai/claude-code:Claude CLI:npm_pkg"
    "CODE:droid:factory-cli:Droid CLI:droid"
    "CODE:gemini:@google/gemini-cli:Gemini CLI:npm_pkg"
    "CODE:opencode:opencode-ai:OpenCode:opencode"
    "CODE:codex:@openai/codex:Codex CLI:npm_pkg"
    "CODE:crush:crush:Crush CLI:brew_pkg"
    "SYS:brew:homebrew:Homebrew:brew"
    "SYS:npm:npm:NPM Globals:npm_sys"
)

# Global Counters
CODE_UPDATES_COUNT=0
SYS_UPDATES_COUNT=0

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
    echo -e "${BLUE}  System Intelligence & Update Utility v0.1.1${RESET}"
    echo -e "${DIM}  ===========================================${RESET}\n"
}

get_local_version() {
    local binary=$1
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
         # Try direct version flag first, fallback to npm list
         ver=$($binary --version 2>/dev/null | awk '{print $NF}')
         if [[ -z "$ver" ]]; then
            ver=$(npm list -g @anthropic-ai/claude-code --depth=0 2>/dev/null | grep claude-code | awk -F@ '{print $NF}')
         fi
    elif [[ "$binary" == "droid" ]]; then
        # Droid doesn't always expose version easily
        ver="Installed"
    elif [[ "$binary" == "opencode" ]]; then
        ver=$(opencode --version 2>/dev/null | head -n 1 | awk '{print $3}') || ver="Installed"
    elif [[ "$binary" == "gemini" ]]; then
         ver=$(npm list -g @google/gemini-cli --depth=0 2>/dev/null | grep gemini-cli | awk -F@ '{print $NF}') || ver="Unknown"
    elif [[ "$binary" == "codex" ]]; then
         ver=$(npm list -g @openai/codex --depth=0 2>/dev/null | grep codex | awk -F@ '{print $NF}') || ver="Unknown"
    else
        ver=$($binary --version 2>/dev/null | head -n 1 | awk '{print $NF}') || ver="Detected"
    fi
    
    # Final cleanup
    if [[ -z "$ver" ]]; then ver="Detected"; fi
    echo "$ver"
}

get_remote_version() {
    local method=$1
    local package=$2
    
    # Only fetch remote versions for NPM packages to keep script fast.
    if [[ "$method" == "npm_pkg" ]] || [[ "$method" == "npm_sys" ]]; then
        local remote=$(npm view "$package" version 2>/dev/null)
        echo "$remote"
    else
        echo "Latest"
    fi
}

check_active_sessions() {
    local active_found=0
    
    echo -e "${DIM}Checking for active sessions...${RESET}"
    
    for tool_entry in "${TOOLS[@]}"; do
        IFS=':' read -r category binary pkg display method <<< "$tool_entry"
        
        if command -v "$binary" &> /dev/null; then
            # pgrep -f matches the full command line
            if pgrep -f "$binary" > /dev/null; then
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
    printf "${BOLD}%-4s %-18s %-15s %-15s${RESET}\n" "Sts" "Tool" "Current" "Target"
    echo "--------------------------------------------------------"

    # Reset counters
    CODE_UPDATES_COUNT=0
    SYS_UPDATES_COUNT=0

    # Function to print a category group
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
                    target=$(get_remote_version "$method" "$pkg")
                    
                    # Logic to determine if update is needed
                    if [[ "$target" == "Latest" ]]; then
                        # Cannot determine version parity, assume update might be needed or manual check
                        needs_update=1 
                    elif [[ "$target" != "-" ]] && [[ "$current" != "$target" ]]; then
                        needs_update=1
                        color="${YELLOW}"
                        icon="${YELLOW}↑${RESET}"
                    elif [[ "$current" == "$target" ]]; then
                        color="${GREEN}"
                    fi
                fi

                # Increment Counters
                if [[ $needs_update -eq 1 ]]; then
                    if [[ "$category" == "CODE" ]]; then ((CODE_UPDATES_COUNT++)); fi
                    if [[ "$category" == "SYS" ]]; then ((SYS_UPDATES_COUNT++)); fi
                fi

                # Visual Table Row
                printf "% -13b ${color}%-18s %-15s ${MAGENTA}%-15s${RESET}\n" "$icon" "$display" "$current" "$target"
            fi
        done
    }

    print_group "CODE" "AI Development Tools"
    echo ""
    print_group "SYS" "System Tools"
    echo ""
}

perform_update() {
    local method=$1
    local name=$2
    local pkg=$3
    local current=$4
    local target=$5

    # Smart Skip Logic
    if [[ "$target" != "Latest" ]] && [[ "$target" != "-" ]] && [[ "$current" == "$target" ]]; then
         echo -e "${DIM}   ○ $name is up to date ($current). Skipped.${RESET}"
         return
    fi

    echo -e "${BOLD}${CYAN}⚡ Updating $name...${RESET}"
    echo -e "${DIM}   Strategy: $method${RESET}"

    case $method in
        brew)
            brew update && brew upgrade && brew cleanup
            ;;
        npm_sys)
            npm update -g
            ;;
        npm_pkg)
            npm install -g "$pkg@latest"
            ;;
        droid)
            curl -fsSL https://app.factory.ai/cli | sh
            ;;
        opencode)
            opencode upgrade || curl -fsSL https://opencode.ai/install | bash
            ;;
        brew_pkg) 
             brew upgrade "$pkg" 2>/dev/null || echo "   Already latest or not managed by brew"
            ;;
        *)
            echo "   No update method found."
            ;;
    esac
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}   ✔ Success${RESET}\n"
    else
        echo -e "${RED}   ✘ Error updating $name${RESET}\n"
    fi
}

# --- Main Logic ---

banner

# 1. Analyze System (Display Table & Count Updates)
analyze_system

# 2. Check for active sessions
check_active_sessions

# 3. Selection Menu
TOTAL_UPDATES=$((CODE_UPDATES_COUNT + SYS_UPDATES_COUNT))
echo -e "${BOLD}Update Modes:${RESET}"
echo -e "  ${BOLD}[1]${RESET} ${CYAN}Code AI Tools Only${RESET} (${BOLD}${CODE_UPDATES_COUNT}${RESET} updates available)"
echo -e "  ${BOLD}[2]${RESET} ${YELLOW}Full System Update${RESET} (${BOLD}${TOTAL_UPDATES}${RESET} updates available)"
echo -e "  ${BOLD}[3]${RESET} Exit"
echo ""

read -p "Select option [1]: " mode
mode=${mode:-1} # Default to 1

# 4. Execution
if [[ "$mode" == "3" ]]; then
    echo "Bye!"
    exit 0
fi

TARGET_CATEGORY="CODE"
if [[ "$mode" == "2" ]]; then
    TARGET_CATEGORY="ALL"
fi

echo -e "\n${BOLD}Starting Update Sequence...${RESET}\n"

for tool_entry in "${TOOLS[@]}"; do
    IFS=':' read -r category binary pkg display method <<< "$tool_entry"

    # Filter logic
    if [[ "$TARGET_CATEGORY" != "ALL" ]] && [[ "$category" != "$TARGET_CATEGORY" ]]; then
        continue
    fi

    # Check if installed before updating
    if command -v "$binary" &> /dev/null; then
        # Fetch versions again or pass them (re-fetching locally is fast, remote is cached in var logic ideally but we just fetch again for simplicity in execution loop)
        # To match the logic of 'analyze', we need the target again.
        # Optimization: Just re-run get_local/remote.
        
        current=$(get_local_version "$binary")
        target=$(get_remote_version "$method" "$pkg")
        
        perform_update "$method" "$display" "$pkg" "$current" "$target"
    fi
done

echo -e "${BOLD}${GREEN}✨ Spark Sequence Complete.${RESET}"
