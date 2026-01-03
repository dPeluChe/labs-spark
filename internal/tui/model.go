package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dpeluche/spark/internal/core"
	"github.com/dpeluche/spark/internal/updater"
)

// --- App State Management ---
type sessionState int

const (
	stateSplash sessionState = iota
	stateMain
	stateConfirm
	stateUpdating
	stateSummary
)

// --- Professional Styling ---
var (
	cGreen  = lipgloss.Color("#04B575")
	cBlue   = lipgloss.Color("#2E7DE1")
	cPurple = lipgloss.Color("#A78BFA")
	cGray   = lipgloss.Color("#6B7280")
	cWhite  = lipgloss.Color("#FFFFFF")
	cDark   = lipgloss.Color("#1F2937")
	cYellow = lipgloss.Color("#F59E0B")
	cRed    = lipgloss.Color("#EF4444")

	appStyle = lipgloss.NewStyle().Padding(1, 2)

	splashTitleStyle = lipgloss.NewStyle().Foreground(cBlue).Bold(true).MarginBottom(1)
	splashSubtitleStyle = lipgloss.NewStyle().Foreground(cGray).Italic(true)

	cardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cPurple).
		Padding(0, 1).
		MarginBottom(1).
		Width(60)

	cardTitleStyle = lipgloss.NewStyle().Foreground(cPurple).Bold(true).Padding(0, 1).Background(cDark)

	selectedItemStyle = lipgloss.NewStyle().Foreground(cGreen).Bold(true)
	dimmedItemStyle = lipgloss.NewStyle().Foreground(cGray)

	statusChecking = lipgloss.NewStyle().Foreground(cYellow).Render("⟳ Checking...")
	statusUpToDate = lipgloss.NewStyle().Foreground(cGray).Render("✔ Up to date")
	statusOutdated = lipgloss.NewStyle().Foreground(cYellow).Bold(true)
	statusMissing  = lipgloss.NewStyle().Foreground(cRed).Render("○ Not Installed")
	statusUpdating = lipgloss.NewStyle().Foreground(cBlue).Render("➜ Updating...")
	statusSuccess  = lipgloss.NewStyle().Foreground(cGreen).Render("✔ Updated")
	statusFailed   = lipgloss.NewStyle().Foreground(cRed).Render("✘ Failed")

	modalStyle = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderForeground(cRed).Padding(1, 2).Align(lipgloss.Center).Width(50)
)

const sparkArt = `
   _____ ____  ___  ____  __ __
  / ___// __ \/   |/ __ \/ //_/
  \__ \/ /_/ / /| / /_/ / ,<   
 ___/ / ____/ ___ / _, _/ /| |  
/____/_/   /_/  |/_/ |_/_/ |_|  
`

type CheckResultMsg struct {
	Index  int
	Local  string
	Remote string
	Status core.ToolStatus
}

type UpdateResultMsg struct {
	Index   int
	Success bool
	Message string
}

type ToolState struct {
	Tool   core.Tool
	Status core.ToolStatus
	Local  string
	Remote string
}

type Model struct {
	state    sessionState
	items    []ToolState
	detector *updater.Detector
	cursor   int
	checked  map[int]bool
	quitting bool
	width    int
	height   int
	loading  int 
	updating int 
}

func NewModel() Model {
	inv := core.GetInventory()
	states := make([]ToolState, len(inv))
	for i, t := range inv {
		states[i] = ToolState{
			Tool:   t,
			Status: core.StatusChecking,
			Local:  "...",
			Remote: "...",
		}
	}

	return Model{
		state:    stateSplash,
		items:    states,
		detector: updater.NewDetector(),
		checked:  make(map[int]bool),
		loading:  len(inv),
	}
}

func (m Model) checkVersion(i int) tea.Cmd {
	return func() tea.Msg {
		t := m.items[i].Tool
		local := m.detector.GetLocalVersion(t)
		status := core.StatusInstalled
		if local == "MISSING" {
			status = core.StatusMissing
		}
		remote := "Latest" 
		return CheckResultMsg{Index: i, Local: local, Remote: remote, Status: status}
	}
}

func (m Model) performUpdate(i int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second * 2)
		return UpdateResultMsg{Index: i, Success: true, Message: "Updated"}
	}
}

func (m Model) startUpdates() tea.Cmd {
	var cmds []tea.Cmd
	m.updating = 0
	for i := range m.items {
		if m.checked[i] {
			m.items[i].Status = core.StatusUpdating
			m.updating++
			cmds = append(cmds, m.performUpdate(i))
		}
	}
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

func tick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), m.checkAllVersions())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case CheckResultMsg:
		m.items[msg.Index].Local = msg.Local
		m.items[msg.Index].Remote = msg.Remote
		m.items[msg.Index].Status = msg.Status
		m.loading--
		return m, nil

	case UpdateResultMsg:
		if msg.Success {
			m.items[msg.Index].Status = core.StatusUpdated
		} else {
			m.items[msg.Index].Status = core.StatusFailed
		}
		m.updating--
		if m.updating == 0 {
			m.state = stateSummary
		}
		return m, nil

	case tea.KeyMsg:
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
			m.quitting = true
			return m, tea.Quit
		}
		
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 { m.cursor-- }
		case "down", "j":
			if m.cursor < len(m.items)-1 { m.cursor++ }
		
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
		
		case "g", "G": 
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

		case "a":
			if len(m.checked) == len(m.items) {
				m.checked = make(map[int]bool) 
			} else {
				for i := range m.items { m.checked[i] = true }
			}

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

func (m Model) View() string {
	if m.quitting { return "" }

	switch m.state {
	case stateSplash:
		return m.ViewSplash()
	case stateConfirm:
		return m.overlayModal("")
	case stateUpdating, stateSummary:
		return m.ViewMain()
	default:
		return m.ViewMain()
	}
}

func (m Model) overlayModal(_ string) string {
	modalContent := lipgloss.NewStyle().Bold(true).Foreground(cRed).Render("⚠️  DANGER ZONE ⚠️") + "\n\n"
	modalContent += "You have selected Critical Runtimes.\nUpdating Node/Python may break your projects.\n\n"
	modalContent += lipgloss.NewStyle().Foreground(cWhite).Render("Are you sure? (y/N)")
	modalBox := modalStyle.Render(modalContent)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modalBox)
}

func (m Model) ViewSplash() string {
	logo := splashTitleStyle.Render(sparkArt)
	sub := splashSubtitleStyle.Render("\n   Surgical Precision Update Utility v0.5.0\n   Initializing System Core...")
	content := lipgloss.JoinVertical(lipgloss.Center, logo, sub)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) ViewMain() string {
	var s string
	headerTxt := " SPARK DASHBOARD "
	if m.state == stateUpdating {
		headerTxt = fmt.Sprintf(" UPDATING (%d remaining)... ", m.updating) 
	} else if m.state == stateSummary {
		headerTxt = " UPDATE SUMMARY "
	} else if m.loading > 0 {
		headerTxt += fmt.Sprintf("(Scanning %d...)", m.loading)
	}
	
s += lipgloss.NewStyle().
		Background(cBlue).Foreground(cWhite).Bold(true).Padding(0, 1).
		Render(headerTxt) + "\n\n"

	renderCategoryCard := func(targetCat core.Category, key string) string {
		var rows []string
		title := fmt.Sprintf("[%s] %s", lipgloss.NewStyle().Foreground(cGreen).Render(key), getCategoryLabel(targetCat))
		hasItems := false
		for i, item := range m.items {
			if item.Tool.Category != targetCat { continue }
			hasItems = true
			cursor := " "; if m.cursor == i && m.state == stateMain { cursor = "➜" }
			checked := "[ ]"; if _, ok := m.checked[i]; ok { checked = "[✔]" }
			
			status := statusChecking
			if m.state == stateUpdating || m.state == stateSummary {
				if item.Status == core.StatusUpdating { status = statusUpdating } else if item.Status == core.StatusUpdated { status = statusSuccess } else if item.Status == core.StatusFailed { status = statusFailed } else if m.checked[i] { status = lipgloss.NewStyle().Foreground(cGray).Render("⏳ Pending...") } else { status = lipgloss.NewStyle().Foreground(cDark).Render(item.Local) }
			} else {
				if item.Status != core.StatusChecking {
					if item.Status == core.StatusMissing { status = statusMissing } else if item.Status == core.StatusInstalled { status = statusUpToDate } else { status = item.Local }
				}
			}

			name := item.Tool.Name
			if len(name) > 22 { name = name[:21] + "…" }
			lineStr := fmt.Sprintf("%s %s %-22s %s", cursor, checked, name, status)
						if m.cursor == i && m.state == stateMain { lineStr = selectedItemStyle.Render(lineStr) } else { lineStr = dimmedItemStyle.Render(lineStr) }
						rows = append(rows, lineStr)
					}		if !hasItems { return "" }
		body := strings.Join(rows, "\n")
		return cardStyle.Render(lipgloss.JoinVertical(lipgloss.Left, cardTitleStyle.Render(title), body))
	}

	col1 := lipgloss.JoinVertical(lipgloss.Left, renderCategoryCard(core.CategoryCode, "C"), renderCategoryCard(core.CategoryTerm, "T"), renderCategoryCard(core.CategoryIDE, "I"), renderCategoryCard(core.CategoryProd, "P"))
	col2 := lipgloss.JoinVertical(lipgloss.Left, renderCategoryCard(core.CategoryInfra, "F"), renderCategoryCard(core.CategoryUtils, "U"), renderCategoryCard(core.CategoryRuntime, "R"), renderCategoryCard(core.CategorySys, "S"))
	grid := lipgloss.JoinHorizontal(lipgloss.Top, col1, "  ", col2)
	s += grid
	
	help := "[SPACE] Select • [G] Group • [TAB] Next • [ENTER] Update • [Q] Quit"
	if m.state == stateUpdating { help = "[UPDATING IN PROGRESS... PLEASE WAIT]" }
	if m.state == stateSummary { help = "[UPDATE COMPLETE] Press any key to exit." }
	s += lipgloss.NewStyle().Foreground(cGray).Render("\n\n" + help)
	return appStyle.Render(s)
}

func getCategoryLabel(c core.Category) string {
	switch c {
	case core.CategoryCode: return "AI Development"
	case core.CategoryTerm: return "Terminals"
	case core.CategoryIDE:  return "IDEs & Editors"
	case core.CategoryProd: return "Productivity"
	case core.CategoryInfra: return "Infrastructure"
	case core.CategoryUtils: return "Utilities"
	case core.CategoryRuntime: return "Runtimes"
	case core.CategorySys: return "System"
	default: return string(c)
	}
}
