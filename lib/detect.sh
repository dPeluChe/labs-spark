#!/bin/bash

get_local_version() {
    local binary=$1
    
    # Special handling for macOS Apps
    if [[ "$binary" == "iterm" ]]; then
        [ -d "/Applications/iTerm.app" ] && defaults read /Applications/iTerm.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "ghostty" ]]; then
        [ -d "/Applications/Ghostty.app" ] && defaults read /Applications/Ghostty.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "warp" ]]; then
        [ -d "/Applications/Warp.app" ] && defaults read /Applications/Warp.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "code" ]]; then
        [ -d "/Applications/Visual Studio Code.app" ] && defaults read "/Applications/Visual Studio Code.app/Contents/Info.plist" CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "cursor" ]]; then
        [ -d "/Applications/Cursor.app" ] && defaults read /Applications/Cursor.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "zed" ]]; then
        [ -d "/Applications/Zed.app" ] && defaults read /Applications/Zed.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "windsurf" ]]; then
        [ -d "/Applications/Windsurf.app" ] && defaults read /Applications/Windsurf.app/Contents/Info.plist CFBundleShortVersionString 2>/dev/null && return
    elif [[ "$binary" == "antigravity" ]]; then
        if [ -f "$HOME/.antigravity/antigravity/bin/antigravity" ]; then
             "$HOME/.antigravity/antigravity/bin/antigravity" --version 2>/dev/null | head -n 1 | awk '{print $2}' && return
        fi
        if command -v antigravity &>/dev/null; then
            antigravity --version 2>/dev/null | head -n 1 | awk '{print $2}' && return
        fi
        echo "MISSING" && return
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
         # Try direct version command first (curl installation or brew)
         ver=$($binary --version 2>/dev/null | awk '{print $1}')
         # Fallback to npm if not found and not a simple binary
         if [[ -z "$ver" ]] || [[ "$ver" == "MISSING" ]]; then
            ver=$(npm list -g @anthropic-ai/claude-code --depth=0 2>/dev/null | grep claude-code | awk -F@ '{print $NF}')
         fi
    elif [[ "$binary" == "droid" ]]; then
        ver=$(droid --version 2>/dev/null | head -n 1) || ver="Installed"
    elif [[ "$binary" == "toad" ]]; then
        # Toad doesn't support --version, use uv tool list
        ver=$(uv tool list 2>/dev/null | grep "batrachian-toad" | awk '{print $2}' | sed 's/^v//') || ver="Installed"
    elif [[ "$binary" == "opencode" ]]; then
        ver=$(opencode --version 2>/dev/null | head -n 1 | awk '{print $NF}') || ver="Installed"
    elif [[ "$binary" == "omz" ]]; then
        if [ -d "$HOME/.oh-my-zsh" ]; then
            # Safe check inside the directory
            ver=$(git --git-dir="$HOME/.oh-my-zsh/.git" --work-tree="$HOME/.oh-my-zsh" rev-parse --short HEAD 2>/dev/null)
            if [[ -z "$ver" ]]; then ver="Installed"; fi
        else
            ver="MISSING"
        fi
    elif [[ "$binary" == "aws" ]]; then
         ver=$($binary --version 2>/dev/null | awk '{print $1}' | cut -d/ -f2)
    elif [[ "$binary" == "ngrok" ]]; then
         ver=$($binary --version 2>/dev/null | awk '{print $3}')
    elif [[ "$binary" == "go" ]]; then
         ver=$($binary version 2>/dev/null | awk '{print $3}' | sed 's/go//')
    elif [[ "$binary" == "bash" ]]; then
         ver=$($binary --version | head -n 1 | awk '{print $4}' | sed 's/(.*//')
    elif [[ "$binary" == "gemini" ]]; then
         ver=$(npm list -g @google/gemini-cli --depth=0 2>/dev/null | grep gemini-cli | awk -F@ '{print $NF}') || ver="Unknown"
    elif [[ "$binary" == "codex" ]]; then
         ver=$(npm list -g @openai/codex --depth=0 2>/dev/null | grep codex | awk -F@ '{print $NF}') || ver="Unknown"
    elif [[ "$binary" == "crush" ]]; then
         ver=$(crush --version 2>/dev/null | head -n 1 | awk '{print $NF}') || ver="Installed"
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
        # For Mac Apps, first check if it is installed via Brew Cask
        if [[ "$method" == "mac_app" ]]; then
            if ! echo "$BREW_CASK_LIST" | grep -q "^$package$"; then
                # Not managed by Brew, but consult Brew for the latest version info
                local latest=$(brew info --cask "$package" 2>/dev/null | grep "^==> $package: " | awk '{print $3}' | sed 's/,.*//' | sed 's/\.stable_.*//')
                if [[ -n "$latest" ]]; then
                    echo "$latest (External)"
                else
                    echo "Unmanaged"
                fi
                return
            fi
        fi

        local update_info=$(echo "$BREW_CACHE" | grep "^$package ")
        if [[ -n "$update_info" ]]; then
            # Format is usually: name (current) < (latest)
            # Clean version string (remove commas, underscores for some packages)
            echo "$update_info" | awk '{print $NF}' | sed 's/,.*//' | sed 's/_.*//'
        else
            # If not in outdated, it's either up to date or unmanaged
            echo "$local_ver"
        fi
    elif [[ "$method" == "claude" ]]; then
        # Query npm for Claude CLI if we can't get it from brew cache
        if echo "$BREW_CASK_LIST" | grep -q "claude-code"; then
             local update_info=$(echo "$BREW_CACHE" | grep "^claude-code ")
             if [[ -n "$update_info" ]]; then
                echo "$update_info" | awk '{print $NF}'
             else
                echo "$local_ver"
             fi
        else
            npm view @anthropic-ai/claude-code version 2>/dev/null || echo "$local_ver"
        fi
    elif [[ "$method" == "toad" ]]; then
        curl -s "https://pypi.org/pypi/batrachian-toad/json" 2>/dev/null | grep -o '"version":"[^"]*"' | head -1 | cut -d'"' -f4
    elif [[ "$method" == "droid" ]]; then
        npm view factory-cli version 2>/dev/null || echo "$local_ver"
    elif [[ "$method" == "opencode" ]]; then
        echo "$local_ver"
    elif [[ "$method" == "antigravity" ]]; then
        echo "Manual"
    else
        echo "Latest"
    fi
}
