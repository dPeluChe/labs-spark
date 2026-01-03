#!/bin/bash

banner() {
    clear
    echo -e "${CYAN}${BOLD}"
    echo "   _____ ____  ___  ____  __ __"
    echo "  / ___// __ \/   |/ __ \/ //_/"
    echo "  \__ \/ /_/ / /| / /_/ / ,<   "
    echo " ___/ / ____/ ___ / _, _/ /| |  "
    echo "/____/_/   /_/  |/_/ |_/_/ |_|  "
    echo -e "${RESET}"
    echo -e "${BLUE}  Surgical Precision Update Utility v0.4.2${RESET}"
    echo -e "${DIM}  ========================================${RESET}\n"
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
            [[ "$binary" == "code" ]] && proc="Visual Studio Code"
            [[ "$binary" == "cursor" ]] && proc="Cursor"
            [[ "$binary" == "zed" ]] && proc="Zed"
            [[ "$binary" == "windsurf" ]] && proc="Windsurf"
            [[ "$binary" == "python3" ]] && proc="python"
            [[ "$binary" == "omz" ]] && proc="zsh"
            
            if pgrep -fi "$proc" > /dev/null; then
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
    echo -e "${DIM}   Fetching Homebrew intelligence...${RESET}"
    
    # Cache Brew data
    BREW_CACHE=$(HOMEBREW_NO_AUTO_UPDATE=1 HOMEBREW_NO_ENV_HINTS=1 brew outdated --verbose 2>&1 | grep -v "^==>")
    BREW_CASK_LIST=$(brew list --cask -1)

    printf "${BOLD}%-4s %-6s %-24s %-15s %-15s${RESET}\n" "Sts" "ID" "Tool" "Current" "Target"
    echo "------------------------------------------------------------------"

    CODE_UPDATES_COUNT=0
    TERM_UPDATES_COUNT=0
    IDE_UPDATES_COUNT=0
    PROD_UPDATES_COUNT=0
    INFRA_UPDATES_COUNT=0
    UTILS_UPDATES_COUNT=0
    RUNTIME_UPDATES_COUNT=0
    SYS_UPDATES_COUNT=0
    
    # Reset Tool Indexing
    TOOL_ID_MAP=()
    TOTAL_TOOLS=0

    print_group() {
        local target_cat=$1
        local title=$2
        echo -e "${DIM}--- $title ---${RESET}"
        
        for tool_entry in "${TOOLS[@]}"; do
            IFS=':' read -r category binary pkg display method <<< "$tool_entry"
            
            # Map every tool to an ID regardless of category display
            ((TOTAL_TOOLS++))
            TOOL_ID_MAP[$TOTAL_TOOLS]="$binary"
            local tool_id_display=$(printf "S-%02d" $TOTAL_TOOLS)

            if [[ "$category" == "$target_cat" ]]; then
                local current=$(get_local_version "$binary")
                local icon=""
                local target="-"
                local target_display="-"
                local color=""
                local needs_update=0

                if [[ "$current" == "MISSING" ]]; then
                    icon="${DIM}○${RESET}"
                    color="${DIM}"
                    current="Not Installed"
                else
                    target=$(get_remote_version "$method" "$pkg" "$current")
                    
                    # Clean the target version string for comparison (remove "(External)")
                    local clean_target="${target// (External)/}"
                    
                    if [[ "$target" == "Unmanaged" ]]; then
                        icon="${CYAN}?${RESET}"
                        color="${DIM}"
                        target_display="${DIM}(Not in Brew)${RESET}"
                        needs_update=0
                    elif [[ "$target" == "Manual" ]]; then
                        icon="${CYAN}M${RESET}"
                        color="${DIM}"
                        target_display="${DIM}(Manual Check)${RESET}"
                        needs_update=0
                    elif [[ "$target" == "Latest" ]]; then
                        icon="${GREEN}●${RESET}"
                        target_display="Manual Check"
                        needs_update=0 
                    elif [[ "$target" == *"External"* ]]; then
                         # Unmanaged but we found a version
                         if [[ "$clean_target" != "$current" ]]; then
                            # Versions differ
                            needs_update=0 # Don't count as auto-update
                            color="${CYAN}"
                            icon="${CYAN}↑${RESET}"
                            target_display="${CYAN}$clean_target (Manual)${RESET}"
                         else
                            icon="${GREEN}●${RESET}"
                            target_display="${DIM}✔ Up to date (Ext)${RESET}"
                         fi
                    elif [[ "$target" != "-" ]] && [[ "$current" != "$target" ]]; then
                        needs_update=1
                        color="${YELLOW}"
                        icon="${YELLOW}↑${RESET}"
                        target_display="${MAGENTA}$target${RESET}"
                    else
                        icon="${GREEN}●${RESET}"
                        # Current == Target (Up to date)
                        target_display="${DIM}✔ Up to date${RESET}"
                    fi
                fi
                
                if [[ $needs_update -eq 1 ]]; then
                    [[ "$category" == "CODE" ]] && ((CODE_UPDATES_COUNT++))
                    [[ "$category" == "TERM" ]] && ((TERM_UPDATES_COUNT++))
                    [[ "$category" == "IDE" ]] && ((IDE_UPDATES_COUNT++))
                    [[ "$category" == "PROD" ]] && ((PROD_UPDATES_COUNT++))
                    [[ "$category" == "INFRA" ]] && ((INFRA_UPDATES_COUNT++))
                    [[ "$category" == "UTILS" ]] && ((UTILS_UPDATES_COUNT++))
                    [[ "$category" == "RUNTIME" ]] && ((RUNTIME_UPDATES_COUNT++))
                    [[ "$category" == "SYS" ]] && ((SYS_UPDATES_COUNT++))
                fi

                # Display with Spark ID (S-XX)
                printf " %-4b ${DIM}%-6s${RESET} ${color}%-24s${RESET} %-15s %b${RESET}\n" "$icon" "$tool_id_display" "$display" "$current" "$target_display"
            fi
        done
    }

    print_group "CODE" "AI Development Tools"
    echo ""
    print_group "TERM" "Terminal Emulators"
    echo ""
    print_group "IDE" "IDEs and Code Editors"
    echo ""
    print_group "PROD" "Developer Productivity"
    echo ""
    print_group "INFRA" "Infrastructure & Cloud"
    echo ""
    print_group "UTILS" "Safe Utilities"
    echo ""
    print_group "RUNTIME" "Critical Runtimes (High Risk)"
    echo ""
    print_group "SYS" "System Managers"
    echo ""
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
