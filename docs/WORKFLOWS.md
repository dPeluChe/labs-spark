# 🔄 SPARK - User Workflows & Processes

This document describes all user interaction flows in SPARK v0.6.0.

---

## Table of Contents

1. [Basic Update Workflow](#basic-update-workflow)
2. [Search & Filter Workflow](#search--filter-workflow)
3. [Dry-Run Preview Workflow](#dry-run-preview-workflow)
4. [Dangerous Runtime Workflow](#dangerous-runtime-workflow)
5. [Navigation Patterns](#navigation-patterns)
6. [Error Handling](#error-handling)

---

## Basic Update Workflow

### Flow: Select and Update Tools

```
┌─────────────────┐
│ Launch SPARK    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Splash Screen   │
│ (2 seconds)     │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────┐
│ Main Dashboard                      │
│ - All 71 tools displayed            │
│ - Grouped by category               │
│ - Version checking in progress...   │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│ Select Tools                        │
│ - Navigate with ↑/↓ or j/k          │
│ - Press SPACE to select/deselect    │
│ - Press A/G to select category      │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│ Press ENTER to Update               │
└────────┬────────────────────────────┘
         │
         ├─── Has Runtime Tools?
         │         │
         │         ├─ YES ──┐
         │         │        │
         │         └─ NO    │
         │            │     │
         ▼            │     ▼
┌──────────────────┐  │  ┌──────────────────┐
│ Start Updating   │  │  │ Danger Zone      │
│ (Immediate)      │  │  │ Modal            │
└────────┬─────────┘  │  └────────┬─────────┘
         │            │           │
         │            │      Press Y/N?
         │            │           │
         │            │      ├─ N → Cancel
         │            │      │
         │            │      └─ Y ──┐
         │            │             │
         │            └─────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│ Updating Screen                     │
│ - Progress bar shows 0%→100%        │
│ - Each tool status updates live     │
│ - Can't quit (except Ctrl+C)        │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│ Summary Screen                      │
│ - Statistics (success rate)         │
│ - List of updated tools             │
│ - List of failures (if any)         │
└────────┬────────────────────────────┘
         │
         ▼
   Press Any Key
         │
         ▼
       EXIT
```

### Step-by-Step Instructions

1. **Launch**: Run `spark` command
2. **Wait**: Splash screen (2s), then automatic version checks (~2s)
3. **Navigate**:
   - Use `j`/`k` or arrow keys to move cursor
   - Press `C`, `T`, `I`, `P`, `F`, `U`, `R`, `S` to jump to categories
   - Press `TAB` to jump to next category
4. **Select**:
   - Press `SPACE` to toggle selection of current tool
   - Press `A` or `G` to select/deselect entire category
5. **Update**:
   - Press `ENTER` to start updates
   - If runtimes selected → Confirmation modal appears
   - Otherwise → Updates start immediately
6. **Monitor**: Watch progress bar and live status updates
7. **Review**: Check summary statistics
8. **Exit**: Press any key

**Time estimate**: 30 seconds to 5 minutes depending on number of tools.

---

## Search & Filter Workflow

### Flow: Find Specific Tools

```
┌─────────────────────────┐
│ Main Dashboard          │
└────────┬────────────────┘
         │
         │ Press /
         ▼
┌─────────────────────────────────────┐
│ Search Mode                         │
│ Search: █                           │
│ [Type to search]                    │
└────────┬────────────────────────────┘
         │
         │ Type characters
         ▼
┌─────────────────────────────────────┐
│ Filtered View                       │
│ Search: node█ (3 results)           │
│                                     │
│ [✔] Node.js        20.11.0          │
│ [ ] Nodemon        2.0.22           │
│ [ ] NodeMon CLI    1.5.0            │
│                                     │
│ (Only matching tools shown)         │
└────────┬────────────────────────────┘
         │
         ├─ Press ENTER → Keep filter, return to main
         ├─ Press ESC   → Clear filter, return to main
         └─ Continue typing → Update filter
```

### Step-by-Step Instructions

1. **Activate Search**: Press `/` from main dashboard
2. **Type Query**: Start typing (searches name, binary, package, category)
3. **See Results**: Grid updates in real-time, showing only matches
4. **Refine**: Keep typing or use Backspace
5. **Confirm**: Press `ENTER` to keep filter active
6. **Or Cancel**: Press `ESC` to clear filter and return

**Search Behavior**:
- **Case-insensitive**: "node" matches "Node.js", "NODE", etc.
- **Partial match**: "no" matches "Node", "Nodemon", "Snowflake"
- **Multi-field**: Searches in Name, Binary, Package, Category
- **Live update**: Results appear as you type
- **Cursor auto-moves**: To first matching item

**Clear Filter**: Press `ESC` from main dashboard (when filter is active)

---

## Dry-Run Preview Workflow

### Flow: Preview Before Updating

```
┌─────────────────────────┐
│ Main Dashboard          │
│ (Select tools first)    │
└────────┬────────────────┘
         │
         │ Press D
         ▼
┌─────────────────────────────────────┐
│ 🔍 UPDATE PREVIEW (DRY-RUN)         │
│                                     │
│ ╭───────────────────────╮           │
│ │ SUMMARY STATISTICS    │           │
│ │                       │           │
│ │ Total Selected: 10    │           │
│ │  • AI Dev: 2 tools    │           │
│ │  • Runtimes: 1 tool   │           │
│ │  • Productivity: 7    │           │
│ ╰───────────────────────╯           │
│                                     │
│ AI Development                      │
│  → Claude CLI (current: 1.2.3)      │
│  → Droid CLI (current: 0.5.0)       │
│                                     │
│ Runtimes                            │
│  → Node.js (current: 20.11.0)       │
│  ⚠ WARNING: Runtime detected        │
│                                     │
│ [ENTER] Proceed • [ESC] Cancel      │
└────────┬────────────────────────────┘
         │
         ├─ Press ENTER → Proceed to update
         │                (check for runtimes)
         │
         └─ Press ESC   → Return to main
                          (selections preserved)
```

### Step-by-Step Instructions

1. **Select Tools**: Use SPACE/G/A to select tools
2. **Preview**: Press `D` for dry-run preview
3. **Review Summary**:
   - See total count
   - See breakdown by category
   - See current versions
   - See warnings for dangerous tools
4. **Decide**:
   - Press `ENTER` to proceed with updates
   - Press `ESC` to cancel and modify selections
5. **If Proceed**: Flow continues to Danger Zone (if runtimes) or Update screen

**Benefits**:
- No surprises - know exactly what will update
- See current versions before updating
- Double-check selections
- Extra safety for critical changes

---

## Dangerous Runtime Workflow

### Flow: Runtime Update Confirmation

```
┌─────────────────────────────────┐
│ User tries to update runtimes   │
│ (Node, Python, Go, Ruby, Postgres)
└────────┬────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│ ╔════════════════════════════════╗   │
│ ║   ⚠️  DANGER ZONE ⚠️            ║   │
│ ╚════════════════════════════════╝   │
│                                      │
│ You have selected Critical Runtimes. │
│ Updating Node/Python may break       │
│ your projects.                       │
│                                      │
│ Are you sure? (y/N)                  │
└────────┬─────────────────────────────┘
         │
         ├─ Press Y → Confirm, start updates
         │
         └─ Press N/ESC/Q → Cancel, return to main
```

### Step-by-Step Instructions

1. **Trigger**: Happens automatically when:
   - Pressing `ENTER` with runtimes selected, OR
   - Pressing `ENTER` from Preview with runtimes

2. **Modal Appears**: Bright red warning, impossible to miss

3. **Think Carefully**: Consider impact on projects

4. **Confirm or Cancel**:
   - Press `Y` to confirm and proceed
   - Press `N`, `ESC`, or `Q` to cancel

**Which Tools Trigger This**:
- Node.js
- Python 3.13
- Go Lang
- Ruby
- PostgreSQL 16

**Why This Exists**:
- Prevents accidental breaking of development environments
- Forces conscious decision
- Gives moment to reconsider

**Best Practice**:
- Update runtimes one at a time
- Have backup/rollback plan
- Check project compatibility first
- Consider using version managers (nvm, pyenv, etc.) instead

---

## Navigation Patterns

### Pattern 1: Linear Navigation

```
Current: Claude CLI

↓ or j → Move down one item
↑ or k → Move up one item
```

**When filtering is active**: Skips invisible items automatically

---

### Pattern 2: Category Jumps

```
Press C → Jump to first item in CODE category
Press T → Jump to first item in TERM category
Press I → Jump to first item in IDE category
Press P → Jump to first item in PROD category
Press F → Jump to first item in INFRA category
Press U → Jump to first item in UTILS category
Press R → Jump to first item in RUNTIME category
Press S → Jump to first item in SYS category
```

**Mnemonic**: First letter of category name (except F=inFra)

---

### Pattern 3: Category Cycling

```
Current: In CODE category

TAB → Jump to next category (TERM)
TAB → Jump to next category (IDE)
TAB → Jump to next category (PROD)
...
TAB → Wrap around to first category (CODE)
```

---

### Pattern 4: Group Selection

```
Current: Node.js (in RUNTIME category)

Press G → Selects ALL items in RUNTIME category
  [ ] Node.js       →  [✔] Node.js
  [ ] Python 3.13   →  [✔] Python 3.13
  [ ] Go Lang       →  [✔] Go Lang
  [ ] Ruby          →  [✔] Ruby
  [ ] PostgreSQL 16 →  [✔] PostgreSQL 16

Press G again → Deselects ALL items in RUNTIME
```

---

## Error Handling

### Scenario 1: Tool Not Found

```
Tool Status: ○ Not Installed

Behavior:
- Shows in red
- Still selectable
- If selected for update → Will attempt installation (future)
```

---

### Scenario 2: Version Detection Failed

```
Tool Status: Detected

Behavior:
- Shows generic "Detected" text
- Means tool is installed but version couldn't be parsed
- Still selectable for update
```

---

### Scenario 3: Update Fails

```
During Update:
  ➜ Updating... → ✘ Failed

After Update (Summary):
  Shows in failed list with error message

User can:
- Review error in summary
- Try again manually
- Check logs (spark_debug.log)
```

---

### Scenario 4: Network/Command Timeout

```
Version Check:
- Times out after 2 seconds
- Shows as "MISSING"

Update:
- Times out after 5 minutes (future)
- Marked as failed
```

---

## Keyboard Reference

### Main Dashboard

| Key | Action |
|-----|--------|
| `↑` `↓` `j` `k` | Navigate items |
| `C` `T` `I` `P` `F` `U` `R` `S` | Jump to category |
| `TAB` | Next category |
| `SPACE` | Toggle selection |
| `A` / `G` | Toggle category |
| `/` | Search mode |
| `D` | Dry-run preview |
| `ENTER` | Start update |
| `ESC` | Clear filter / Quit |
| `Q` | Quit |
| `Ctrl+C` | Force quit |

### Search Mode

| Key | Action |
|-----|--------|
| `Type` | Add to query |
| `Backspace` | Remove character |
| `ENTER` | Confirm filter |
| `ESC` | Cancel filter |

### Preview Mode

| Key | Action |
|-----|--------|
| `ENTER` | Proceed with update |
| `ESC` `Q` | Cancel |

### Danger Zone Modal

| Key | Action |
|-----|--------|
| `Y` | Confirm |
| `N` `ESC` `Q` | Cancel |

### Updating Screen

| Key | Action |
|-----|--------|
| `Ctrl+C` | Emergency exit |
| (all others) | Ignored |

### Summary Screen

| Key | Action |
|-----|--------|
| Any key | Exit SPARK |

---

## Advanced Workflows

### Workflow: Update Only Outdated Tools

**Future Feature** - Currently "Outdated" status not implemented (shows "Latest")

```
1. Launch SPARK
2. Wait for version checks
3. Filter automatically to outdated tools
4. Press A to select all
5. Press ENTER to update
```

### Workflow: Update by Category

```
1. Launch SPARK
2. Press P to jump to PROD category
3. Press G to select all PROD tools
4. Press ENTER to update
```

### Workflow: Update Single Tool Fast

```
1. Launch SPARK
2. Press / to search
3. Type tool name
4. Press ENTER to confirm filter
5. Press ENTER to update (auto-selects if none selected)
```

---

## State Transitions Reference

See `internal/tui/states.go` for detailed state machine diagram.

**Valid Transitions**:
```
stateSplash → stateMain
stateMain → stateSearch, statePreview, stateConfirm, stateUpdating
stateSearch → stateMain
statePreview → stateMain, stateConfirm, stateUpdating
stateConfirm → stateMain, stateUpdating
stateUpdating → stateSummary
stateSummary → EXIT
```

**Invalid Transitions** (blocked by code):
- Cannot skip from Main directly to Summary
- Cannot go back from Updating to Main
- Cannot return from Summary to Main

---

## Tips & Tricks

### Tip 1: Quick Update Category

```
1. Launch SPARK
2. Press C (jump to desired category, e.g., CODE)
3. Press A or G (select all in category)
4. Press D (preview - optional)
5. Press ENTER (update)
```

### Tip 2: Update Only AI Tools

```
1. Press C (jump to CODE category)
2. Press A (select all CODE tools)
3. Press ENTER
```

### Tip 3: Find and Update

```
1. Press / (search)
2. Type "node"
3. Press ENTER (keep filter)
4. Press SPACE on each desired match
5. Press ENTER (update)
```

### Tip 4: Toggle Category Selections

```
Press A/G twice while in a category:
- First press: Select all in category
- Second press: Deselect all in category
```

---

## Next Steps

- See `docs/ARCHITECTURE.md` for technical details
- See `docs/ADDING_TOOLS.md` for extending SPARK
- See `docs/INSTALLATION.md` for setup instructions
