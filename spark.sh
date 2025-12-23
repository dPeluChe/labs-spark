#!/bin/bash

# SPARK v2.0 - Intelligent CLI Updater
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
    "CODE:droid:factory-cli:Droid CLI:droid"
    "CODE:gemini:@google/gemini-cli:Gemini CLI:npm_pkg"
    "CODE:opencode:opencode-ai:OpenCode:opencode"
    "CODE:codex:@openai/codex:Codex CLI:npm_pkg"
    "SYS:brew:homebrew:Homebrew:brew"
    "SYS:npm:npm:NPM Globals:npm_sys"
    "SYS:crush:crush:Crush CLI:brew_pkg"
)

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
    echo -e "${BLUE}  System Intelligence & Update Utility v2.0${RESET}"
    echo -e "${DIM}  =========================================${RESET}\n"
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
    elif [[ "$binary" == "droid" ]]; then
        # Droid doesn't always expose version easily, strictly checking
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
    echo "$ver"
}

get_remote_version() {
    local method=$1
    local package=$2
    
    # Only fetch remote versions for NPM packages to keep script fast.
    # Fetching Brew or Curl versions remotely is too slow/complex for a quick CLI check.
    if [[ "$method" == "npm_pkg" ]] || [[ "$method" == "npm_sys" ]]; then
        echo $(npm view "$package" version 2>/dev/null)
    else
        echo "Latest"
    fi
}

perform_update() {
    local method=$1
    local name=$2
    local pkg=$3

    echo -e "${DIM}   Executing update strategy for $name...${RESET}"

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
}

# --- Main Logic ---

banner

analyze_system() {
    echo -e "${BOLD}Analyzing System Components...${RESET}"
    printf "${BOLD}%-4s %-18s %-15s %-15s${RESET}\n" "Sts" "Tool" "Current" "Target"
    echo "--------------------------------------------------------"

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

                if [[ "$current" == "MISSING" ]]; then
                    icon="${DIM}○${RESET}"
                    color="${DIM}"
                    current="Not Installed"
                else
                    icon="${GREEN}●${RESET}"
                    color="${RESET}"
                    # Only check target if installed
                    target=$(get_remote_version "$method" "$pkg")
                fi

                # Visual Table Row
                printf "%-13b ${color}%-18s %-15s ${MAGENTA}%-15s${RESET}\n" "$icon" "$display" "$current" "$target"
            fi
        done
    }

    print_group "CODE" "AI Development Tools"
    echo ""
    print_group "SYS" "System Tools"
    echo ""
}

# --- Main Logic ---

banner

# 1. Analyze System (Display Table)
analyze_system


# 2. Selection Menu
echo -e "${BOLD}Update Modes:${RESET}"
echo -e "  ${BOLD}[1]${RESET} ${CYAN}Code AI Tools Only${RESET} (Droid, Gemini, OpenCode...) ${DIM}(Recommended)${RESET}"
echo -e "  ${BOLD}[2]${RESET} ${YELLOW}Full System Update${RESET} (Include Homebrew & NPM Globals)"
echo -e "  ${BOLD}[3]${RESET} Exit"
echo ""
read -p "Select option [1]: " mode
mode=${mode:-1} # Default to 1

# 3. Execution
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

    # Check if installed before updating (Skip missing tools unless we want to force install logic, which we don't for now)
    if command -v "$binary" &> /dev/null; then
        echo -e "${BOLD}${CYAN}⚡ Updating $display...${RESET}"
        perform_update "$method" "$display" "$pkg"
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}   ✔ Success${RESET}\n"
        else
            echo -e "${RED}   ✘ Error updating $display${RESET}\n"
        fi
    fi
done

echo -e "${BOLD}${GREEN}✨ Spark Sequence Complete.${RESET}"