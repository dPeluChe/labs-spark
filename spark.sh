#!/bin/bash

# SPARK v0.4.2 - Surgical Precision CLI Updater
# Codenamed: Spark (The life-force of Transformers)

# Resolve directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Load Configuration & Modules
if [ -f "$DIR/config/tools.conf" ]; then
    source "$DIR/config/tools.conf"
else
    echo "Error: config/tools.conf not found."
    exit 1
fi

for module in common detect update ui; do
    if [ -f "$DIR/lib/$module.sh" ]; then
        source "$DIR/lib/$module.sh"
    else
        echo "Error: lib/$module.sh not found."
        exit 1
    fi
done

# --- Main Logic ---
banner
analyze_system
check_active_sessions

echo -e "${BOLD}Update Modes:${RESET}"
echo -e "  ${BOLD}[1]${RESET} ${CYAN}AI Tools${RESET}           (Claude, Droid, Gemini, Toad, etc.)"
echo -e "  ${BOLD}[2]${RESET} ${MAGENTA}Terminals & IDEs${RESET}   (iTerm, Windsurf, Antigravity, VSCode)"
echo -e "  ${BOLD}[3]${RESET} ${GREEN}Utilities${RESET}          (Git, Tmux, Zellij, Oh My Zsh, etc.)"
echo -e "  ${BOLD}[4]${RESET} ${RED}Runtimes${RESET}           (Node, Python, Go, Postgres) ${RED}⚠️${RESET}"
echo -e "  ${BOLD}[5]${RESET} ${YELLOW}Full System${RESET}        (Everything included)"
echo -e "  ${BOLD}[6]${RESET} Open Manual Links    (For External/Manual apps)"
echo -e "  ${BOLD}[7]${RESET} Exit"
echo ""

read -p "Select option [1-7] or enter Tool ID/Name: " input
input=${input:-1}

TARGET_CATEGORY=""
TARGET_TOOL_ID=""
METHOD_MODE="UPDATE"

# Logic to determine if input is a Menu Option or a Tool
if [[ "$input" =~ ^[0-9]+$ ]] && [ "$input" -ge 1 ] && [ "$input" -le 7 ]; then
    # Standard Menu Selection
    if [[ "$input" == "7" ]]; then echo "Bye!"; exit 0; fi
    if [[ "$input" == "1" ]]; then TARGET_CATEGORY="CODE"; fi
    if [[ "$input" == "2" ]]; then TARGET_CATEGORY="TERM_IDE"; fi
    if [[ "$input" == "3" ]]; then TARGET_CATEGORY="UTILS"; fi
    if [[ "$input" == "4" ]]; then TARGET_CATEGORY="RUNTIME"; fi
    if [[ "$input" == "5" ]]; then TARGET_CATEGORY="ALL"; fi
    if [[ "$input" == "6" ]]; then 
        TARGET_CATEGORY="ALL"
        METHOD_MODE="OPEN_LINKS"
    fi
else
    # Targeted Tool Selection (ID or Name)
    TARGET_CATEGORY="ALL" # We search all, but filter by binary
    
    # Check if input matches an ID (supports raw number OR S-XX format)
    clean_id=$(echo "$input" | sed 's/[Ss]-//')
    if [[ "$clean_id" =~ ^[0-9]+$ ]] && [ "$clean_id" -le "$TOTAL_TOOLS" ]; then
        # Remove leading zeros if present for array index
        idx=$((10#$clean_id))
        TARGET_TOOL_ID="${TOOL_ID_MAP[$idx]}"
        TARGET_TOOL_NAME="${TOOL_ID_MAP[$idx]}" # For display
        echo -e "${DIM}Targeting Spark ID S-$(printf "%02d" $idx): $TARGET_TOOL_ID${RESET}"
    else
        # Assume input is a name (e.g., "claude")
        TARGET_TOOL_ID="$input"
        TARGET_TOOL_NAME="$input"
        echo -e "${DIM}Targeting Tool Name: $TARGET_TOOL_ID${RESET}"
    fi
fi

# Safety Check for Runtimes
if [[ "$TARGET_CATEGORY" == "RUNTIME" ]] || [[ "$TARGET_CATEGORY" == "ALL" ]]; then
    # Skip safety check if we are targeting a specific single tool that is NOT a runtime
    if [[ -z "$TARGET_TOOL_ID" ]] && [[ "$METHOD_MODE" != "OPEN_LINKS" ]]; then
        echo -e "\n${RED}${BOLD}⚠️  WARNING: You are about to update critical runtimes (Node, Python, DBs).${RESET}"
        echo -e "${RED}    This might break existing projects or virtual environments.${RESET}"
        read -p "    Are you absolutely sure? (type 'yes' to proceed): " confirm
        if [[ "$confirm" != "yes" ]]; then
            echo "    Operation aborted by user."
            exit 0
        fi
    fi
fi

echo -e "\n${BOLD}Starting Update Sequence...${RESET}\n"

for tool_entry in "${TOOLS[@]}"; do
    IFS=':' read -r category binary pkg display method <<< "$tool_entry"

    MATCH=0
    
    # If a specific tool is targeted, override category logic
    if [[ -n "$TARGET_TOOL_ID" ]]; then
        if [[ "$binary" == "$TARGET_TOOL_ID" ]]; then
            MATCH=1
        fi
    else
        # Standard Category Matching
        if [[ "$TARGET_CATEGORY" == "ALL" ]]; then MATCH=1; fi
        if [[ "$TARGET_CATEGORY" == "TERM_IDE" ]] && ([[ "$category" == "TERM" ]] || [[ "$category" == "IDE" ]]); then MATCH=1; fi
        if [[ "$category" == "$TARGET_CATEGORY" ]]; then MATCH=1; fi
    fi

    if [[ $MATCH -eq 1 ]]; then
        if command -v "$binary" &> /dev/null || [[ "$method" == "mac_app" ]] || [[ "$method" == "antigravity" ]]; then
            current=$(get_local_version "$binary")
            target=$(get_remote_version "$method" "$pkg" "$current")
            perform_update "$method" "$display" "$pkg" "$current" "$target"
        fi
    fi
done

show_summary
echo -e "${BOLD}${GREEN}✨ Spark Sequence Complete.${RESET}"