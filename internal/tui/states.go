package tui

import (
	"fmt"

	"github.com/dpeluche/spark/internal/core"
)

/*
STATE MACHINE DIAGRAM FOR SPARK TUI

┌──────────────┐
│ stateSplash  │ (Initial animated logo)
│   (2s auto)  │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  stateMain   │ ◄─────────────────┐ (Main dashboard)
│              │                   │
└──┬───┬───┬───┘                   │
   │   │   │                       │
   │   │   │ [D] Dry-run           │
   │   │   └─────────►┌────────────┴─────┐
   │   │              │  statePreview     │
   │   │              │  (Dry-run view)   │
   │   │              └────────┬──────────┘
   │   │                       │
   │   │              [ENTER]  │ [ESC]
   │   │                       │    │
   │   │                       │    └──────────┐
   │   │ [/] Search            │               │
   │   └─────────►┌────────────▼─────┐         │
   │              │  stateSearch      │         │
   │              │  (Filter mode)    │         │
   │              └────────┬──────────┘         │
   │                       │                    │
   │              [ENTER]  │ [ESC]              │
   │                       │    └───────────────┤
   │                       │                    │
   │ [ENTER]               ▼                    │
   │ (has Runtimes)   ┌─────────────┐           │
   └──────────────────►│stateConfirm │           │
                      │ (Danger Zone)│           │
                      └──┬────────┬──┘           │
                         │        │              │
                    [Y]  │        │ [N]/[ESC]    │
                         │        └──────────────┤
                         ▼                       │
                    ┌──────────────┐             │
                    │stateUpdating │             │
                    │ (Executing)  │             │
                    └──────┬───────┘             │
                           │                     │
                      (auto complete)            │
                           │                     │
                           ▼                     │
                    ┌──────────────┐             │
                    │ stateSummary │             │
                    │  (Results)   │             │
                    └──────┬───────┘             │
                           │                     │
                      (any key)                  │
                           │                     │
                           ▼                     │
                         [EXIT] ◄────────────────┘


STATE DESCRIPTIONS:

1. stateSplash
   - Entry: Application start
   - Duration: 2 seconds (auto-advance via TickMsg)
   - Exit: -> stateMain (automatic)
   - Actions: Display logo and initialization message

2. stateMain
   - Entry: From splash, search, preview, or confirm (cancel)
   - User Actions:
     * Navigation: ↑/↓, j/k, C/T/I/P/F/U/R/S (category jumps), TAB
     * Selection: SPACE (toggle), G (group), A (all)
     * Search: / (enter search mode)
     * Preview: D (dry-run preview)
     * Update: ENTER (check for dangerous runtimes)
     * Quit: Q, Ctrl+C, ESC (if no filter active)
   - Exit Paths:
     * -> stateSearch (/)
     * -> statePreview (D)
     * -> stateConfirm (ENTER + has runtimes)
     * -> stateUpdating (ENTER + no runtimes)
     * -> EXIT (Q, Ctrl+C, ESC)

3. stateSearch
   - Entry: From stateMain (/)
   - User Actions:
     * Type: Add characters to search query
     * Backspace: Remove last character
     * ENTER: Confirm and return to main with filter active
     * ESC: Cancel and return to main without filter
   - Filter Logic:
     * Searches in: Tool.Name, Tool.Binary, Tool.Package, Tool.Category
     * Case-insensitive
     * Live updates as user types
   - Exit Paths:
     * -> stateMain (ENTER or ESC)

4. statePreview
   - Entry: From stateMain (D)
   - Display:
     * Total selected tools count
     * Breakdown by category
     * List of tools to be updated
     * Current versions
     * Warning if runtimes included
   - User Actions:
     * ENTER: Proceed with update (check for runtimes)
     * ESC/Q: Cancel and return to main
   - Exit Paths:
     * -> stateConfirm (ENTER + has runtimes)
     * -> stateUpdating (ENTER + no runtimes)
     * -> stateMain (ESC/Q)

5. stateConfirm
   - Entry: From stateMain or statePreview (when runtimes selected)
   - Display: "DANGER ZONE" modal
   - Purpose: Prevent accidental runtime updates
   - Runtimes Checked:
     * Node.js
     * Python 3.13
     * Go Lang
     * Ruby
     * PostgreSQL 16
   - User Actions:
     * Y: Confirm and proceed
     * N/ESC/Q: Cancel
   - Exit Paths:
     * -> stateUpdating (Y)
     * -> stateMain (N/ESC/Q)

6. stateUpdating
   - Entry: From stateMain, statePreview, or stateConfirm (confirmed)
   - Behavior:
     * Execute updates in parallel via Goroutines
     * Display progress bar
     * Show live status updates
     * Dim non-selected items
     * Highlight currently updating items
   - User Actions:
     * Ctrl+C: Emergency exit (kills program)
     * All other keys: Ignored
   - Exit Paths:
     * -> stateSummary (when all updates complete)

7. stateSummary
   - Entry: From stateUpdating (automatic when complete)
   - Display:
     * Success rate percentage
     * Count of successful/failed/skipped
     * List of updated tools with versions
     * List of failed tools with error messages
   - User Actions:
     * Any key: Exit application
   - Exit Paths:
     * -> EXIT (any key)

INVARIANTS:
- Only ONE item can have cursor at a time
- Cursor must always point to a valid item index
- During stateUpdating, state cannot change except to stateSummary or EXIT
- Search filter persists across state transitions until cleared
- Selected items (checked map) persist across all states
*/

// validateStateTransition checks if a state transition is valid
func (m *Model) validateStateTransition(from, to sessionState) bool {
	validTransitions := map[sessionState][]sessionState{
		stateSplash: {stateMain},
		stateMain: {
			stateSearch,
			statePreview,
			stateConfirm,
			stateUpdating,
		},
		stateSearch: {stateMain},
		statePreview: {
			stateMain,
			stateConfirm,
			stateUpdating,
		},
		stateConfirm: {
			stateMain,
			stateUpdating,
		},
		stateUpdating: {stateSummary},
		stateSummary:  {}, // Terminal state (only exits)
	}

	allowed := validTransitions[from]
	for _, validState := range allowed {
		if validState == to {
			return true
		}
	}
	return false
}

// safeTransition performs a validated state transition
func (m *Model) safeTransition(to sessionState) bool {
	if m.validateStateTransition(m.state, to) {
		m.state = to
		return true
	}
	return false
}

// getStateName returns human-readable state name for debugging
func getStateName(s sessionState) string {
	names := map[sessionState]string{
		stateSplash:   "SPLASH",
		stateMain:     "MAIN",
		stateSearch:   "SEARCH",
		statePreview:  "PREVIEW",
		stateConfirm:  "CONFIRM",
		stateUpdating: "UPDATING",
		stateSummary:  "SUMMARY",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return "UNKNOWN"
}

// validateModel performs runtime validation of model state
func (m *Model) validateModel() error {
	// Cursor bounds check
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return fmt.Errorf("invalid cursor position: %d (items: %d)", m.cursor, len(m.items))
	}

	// Checked items validation
	for idx := range m.checked {
		if idx < 0 || idx >= len(m.items) {
			return fmt.Errorf("invalid checked index: %d", idx)
		}
	}

	// Filtered items validation
	if m.filteredItems != nil {
		for _, idx := range m.filteredItems {
			if idx < 0 || idx >= len(m.items) {
				return fmt.Errorf("invalid filtered index: %d", idx)
			}
		}
	}

	// State-specific validation
	switch m.state {
	case stateUpdating:
		if m.totalUpdate == 0 {
			return fmt.Errorf("stateUpdating with no updates queued")
		}
	case stateSummary:
		if m.totalUpdate == 0 {
			return fmt.Errorf("stateSummary with no updates executed")
		}
	}

	return nil
}

// Helper function to count selected runtimes
func (m *Model) countSelectedRuntimes() int {
	count := 0
	for i := range m.items {
		if m.checked[i] && m.items[i].Tool.Category == core.CategoryRuntime {
			count++
		}
	}
	return count
}

// Helper function to get selected tools
func (m *Model) getSelectedTools() []core.ToolState {
	selected := []core.ToolState{}
	for i, item := range m.items {
		if m.checked[i] {
			selected = append(selected, item)
		}
	}
	return selected
}
