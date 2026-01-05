package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
	"github.com/dpeluche/spark/internal/core"
	"github.com/dpeluche/spark/internal/updater"
)

// --- App State Management ---
type sessionState int

const (
	stateSplash sessionState = iota
	stateMain
	stateSearch   // Search/filter mode
	statePreview  // Dry-run preview mode
	stateConfirm
	stateUpdating
	stateSummary
)

// Message Types
type CheckResultMsg struct {
	Index         int
	LocalVersion  string
	RemoteVersion string
	Status        core.ToolStatus
	Message       string
}

type UpdateResultMsg struct {
	Index   int
	Success bool
	Message string
}

type Model struct {
	state         sessionState
	items         []core.ToolState // Using core.ToolState instead of local duplicate
	detector      *updater.Detector
	cursor        int
	checked       map[int]bool
	quitting      bool
	width         int
	height        int
	loading       int
	updating      int
	totalUpdate   int           // Total items to update
	progress      progress.Model // Progress bar component
	searchQuery   string        // Current search query
	filteredItems []int         // Indices of filtered items
	splashFrame   int           // Current animation frame for splash screen
}

func NewModel() Model {
	inv := core.GetInventory()
	states := make([]core.ToolState, len(inv))
	for i, t := range inv {
		states[i] = core.ToolState{
			Tool:          t,
			Status:        core.StatusChecking,
			LocalVersion:  "...",
			RemoteVersion: "...",
			Message:       "",
		}
	}

	// Initialize progress bar with theme colors
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
	)

	return Model{
		state:    stateSplash,
		items:    states,
		detector: updater.NewDetector(),
		checked:  make(map[int]bool),
		loading:  len(inv),
		progress: prog,
	}
}

func (m Model) checkVersion(i int) tea.Cmd {
	return func() tea.Msg {
		t := m.items[i].Tool
		local := m.detector.GetLocalVersion(t)
		status := core.StatusInstalled
		message := ""
		if local == "MISSING" {
			status = core.StatusMissing
			message = "Not installed"
		}
		remote := "Latest"
		return CheckResultMsg{
			Index:         i,
			LocalVersion:  local,
			RemoteVersion: remote,
			Status:        status,
			Message:       message,
		}
	}
}

func (m Model) performUpdate(i int) tea.Cmd {
	return func() tea.Msg {
		// Simulate realistic update time (3-5 seconds)
		// This makes progress bar visible during updates
		time.Sleep(time.Second * 4)
		return UpdateResultMsg{Index: i, Success: true, Message: "Updated"}
	}
}

func (m *Model) startUpdates() tea.Cmd {
	var cmds []tea.Cmd
	m.updating = 0
	m.totalUpdate = 0

	// Count and start updates
	for i := range m.items {
		if m.checked[i] {
			m.items[i].Status = core.StatusUpdating
			m.updating++
			m.totalUpdate++
			cmds = append(cmds, m.performUpdate(i))
		}
	}

	// Add refresh ticker to update progress bar smoothly
	cmds = append(cmds, refreshTick())

	return tea.Batch(cmds...)
}

func (m Model) checkAllVersions() tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.items {
		cmds = append(cmds, m.checkVersion(i))
	}
	return tea.Batch(cmds...)
}

type TickMsg time.Time
type AnimateMsg time.Time
type RefreshMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func animateSplash() tea.Cmd {
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return AnimateMsg(t)
	})
}

func refreshTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return RefreshMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), animateSplash(), m.checkAllVersions())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update progress bar width (leave some margin)
		m.progress.Width = msg.Width - 20
		if m.progress.Width < 40 {
			m.progress.Width = 40
		}

	case CheckResultMsg:
		m.items[msg.Index].LocalVersion = msg.LocalVersion
		m.items[msg.Index].RemoteVersion = msg.RemoteVersion
		m.items[msg.Index].Status = msg.Status
		m.items[msg.Index].Message = msg.Message
		m.loading--
		return m, nil

	case UpdateResultMsg:
		if msg.Success {
			m.items[msg.Index].Status = core.StatusUpdated
			m.items[msg.Index].Message = msg.Message
		} else {
			m.items[msg.Index].Status = core.StatusFailed
			m.items[msg.Index].Message = msg.Message
		}
		m.updating--
		if m.updating == 0 {
			m.state = stateSummary
		}
		return m, nil

	case tea.KeyMsg:
		if m.state == statePreview {
			switch msg.String() {
			case "enter":
				// Proceed with updates - check for dangerous runtimes first
				hasCritical := false
				for i := range m.items {
					if m.checked[i] && m.items[i].Tool.Category == core.CategoryRuntime {
						hasCritical = true
						break
					}
				}

				if hasCritical {
					m.state = stateConfirm
					return m, nil
				}

				m.state = stateUpdating
				return m, m.startUpdates()
			case "esc", "q":
				// Cancel and return to main
				m.state = stateMain
				return m, nil
			}
			return m, nil
		}

		if m.state == stateConfirm {
			switch msg.String() {
			case "y", "Y":
				m.state = stateUpdating
				return m, m.startUpdates()
			case "n", "N", "esc", "q":
				m.state = stateMain
				return m, nil
			}
			return m, nil
		}

		if m.state == stateSplash {
			m.state = stateMain
			return m, nil
		}

		if m.state == stateUpdating {
			if msg.String() == "ctrl+c" {
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}

		if m.state == stateSummary {
			// Return to main dashboard instead of quitting
			// Clear selections and reset state
			m.state = stateMain
			m.checked = make(map[int]bool)
			m.totalUpdate = 0
			m.updating = 0
			return m, nil
		}

		// Search mode handling
		if m.state == stateSearch {
			switch msg.String() {
			case "esc":
				// Exit search mode, clear filter
				m.state = stateMain
				m.searchQuery = ""
				m.filteredItems = nil
				return m, nil
			case "enter":
				// Confirm search and return to main
				m.state = stateMain
				return m, nil
			case "backspace":
				// Remove last character
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.updateFilter()
				}
				return m, nil
			default:
				// Add character to search query (only printable chars)
				if len(msg.String()) == 1 {
					char := msg.String()
					m.searchQuery += char
					m.updateFilter()
				}
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			// Clear filter if active, otherwise quit
			if m.searchQuery != "" {
				m.searchQuery = ""
				m.filteredItems = nil
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit
		case "/":
			// Enter search mode
			m.state = stateSearch
			m.searchQuery = ""
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				// Skip invisible items when filtering
				for !m.isItemVisible(m.cursor) && m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				// Skip invisible items when filtering
				for !m.isItemVisible(m.cursor) && m.cursor < len(m.items)-1 {
					m.cursor++
				}
			}
		
		case "c", "C": m.jumpToCategory(core.CategoryCode)
		case "t", "T": m.jumpToCategory(core.CategoryTerm)
		case "i", "I": m.jumpToCategory(core.CategoryIDE)
		case "p", "P": m.jumpToCategory(core.CategoryProd)
		case "f", "F": m.jumpToCategory(core.CategoryInfra)
		case "u", "U": m.jumpToCategory(core.CategoryUtils)
		case "r", "R": m.jumpToCategory(core.CategoryRuntime)
		case "s", "S": m.jumpToCategory(core.CategorySys)

		case "tab":
			currentCat := m.items[m.cursor].Tool.Category
			for i := m.cursor + 1; i < len(m.items); i++ {
				if m.items[i].Tool.Category != currentCat {
					m.cursor = i
					return m, nil
				}
			}
			m.cursor = 0
			
		case " ":
			if _, ok := m.checked[m.cursor]; ok {
				delete(m.checked, m.cursor)
			} else {
				m.checked[m.cursor] = true
			}
		
		case "g", "G", "a", "A":
			// Toggle selection for all items in current category
			currentCat := m.items[m.cursor].Tool.Category
			allSelected := true
			for i, item := range m.items {
				if item.Tool.Category == currentCat {
					if !m.checked[i] { allSelected = false; break }
					}
			}
			for i, item := range m.items {
				if item.Tool.Category == currentCat {
					if allSelected { delete(m.checked, i) } else { m.checked[i] = true }
					}
			}

		case "d", "D":
			// Dry-run preview mode - show what would be updated
			if m.loading > 0 {
				return m, nil
			}
			if len(m.checked) == 0 {
				m.checked[m.cursor] = true
			}
			m.state = statePreview
			return m, nil

		case "enter":
			if m.loading > 0 { return m, nil }
			if len(m.checked) == 0 { m.checked[m.cursor] = true }

			hasCritical := false
			for i := range m.items {
				if m.checked[i] && m.items[i].Tool.Category == core.CategoryRuntime {
					hasCritical = true
					break
				}
			}

			if hasCritical {
				m.state = stateConfirm
				return m, nil
			}

			m.state = stateUpdating
			return m, m.startUpdates()
		}

	case TickMsg:
		if m.state == stateSplash {
			m.state = stateMain
			return m, nil
		}

	case AnimateMsg:
		if m.state == stateSplash {
			m.splashFrame++
			return m, animateSplash()
		}

	case RefreshMsg:
		// Keep refreshing while updating to animate progress bar
		if m.state == stateUpdating {
			return m, refreshTick()
		}
	}
	return m, nil
}

func (m *Model) jumpToCategory(cat core.Category) {
	for i, item := range m.items {
		if item.Tool.Category == cat {
			m.cursor = i
			return
		}
	}
}

// --- Search Functionality ---

func (m *Model) updateFilter() {
	if m.searchQuery == "" {
		m.filteredItems = nil
		return
	}

	m.filteredItems = []int{}
	query := strings.ToLower(m.searchQuery)

	for i, item := range m.items {
		// Search in name, binary, package, category
		if strings.Contains(strings.ToLower(item.Tool.Name), query) ||
			strings.Contains(strings.ToLower(item.Tool.Binary), query) ||
			strings.Contains(strings.ToLower(item.Tool.Package), query) ||
			strings.Contains(strings.ToLower(string(item.Tool.Category)), query) {
			m.filteredItems = append(m.filteredItems, i)
		}
	}

	// Reset cursor to first filtered item
	if len(m.filteredItems) > 0 {
		m.cursor = m.filteredItems[0]
	}
}

func (m *Model) isItemVisible(index int) bool {
	if m.filteredItems == nil {
		return true // No filter active
	}

	for _, i := range m.filteredItems {
		if i == index {
			return true
		}
	}
	return false
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	switch m.state {
	case stateSplash:
		return m.ViewSplash()
	case statePreview:
		return m.ViewPreview()
	case stateConfirm:
		return m.overlayModal("")
	case stateUpdating:
		return m.ViewMain()
	case stateSummary:
		return m.ViewSummary()
	default:
		return m.ViewMain()
	}
}

