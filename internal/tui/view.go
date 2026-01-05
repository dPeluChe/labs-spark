package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dpeluche/spark/internal/core"
)

// ViewMain renders the main dashboard view
func (m Model) ViewMain() string {
	header := m.renderHeader()

	// Show search bar if in search mode or filter active
	var searchBar string
	if m.state == stateSearch || m.searchQuery != "" {
		searchBar = m.renderSearchBar() + "\n\n"
	}

	grid := m.renderGrid()

	// Show progress bar during updates (BELOW the grid)
	var progressBar string
	if m.state == stateUpdating && m.totalUpdate > 0 {
		progressBar = "\n\n" + m.renderProgressBar() + "\n"
	}

	helpBar := m.renderHelpBar()

	content := header + "\n\n" + searchBar + grid + progressBar + helpBar
	return appStyle.Render(content)
}

// ViewSplash renders the animated splash screen
func (m Model) ViewSplash() string {
	// Animate logo color based on frame (cycles through colors)
	colors := []lipgloss.Color{
		cBlue,   // Frame 0-2
		lipgloss.Color("#4EA7FF"), // Lighter blue
		cPurple, // Purple
		lipgloss.Color("#00D9FF"), // Cyan
		cGreen,  // Green
		cBlue,   // Back to blue
	}

	frameIndex := m.splashFrame % len(colors)
	animatedStyle := lipgloss.NewStyle().
		Foreground(colors[frameIndex]).
		Bold(true).
		MarginBottom(1)

	// Add loading dots animation
	dots := strings.Repeat(".", (m.splashFrame/3)%4)

	logo := animatedStyle.Render(sparkArt)
	sub := splashSubtitleStyle.Render(
		fmt.Sprintf("\n   Surgical Precision Update Utility v0.6.0\n   Initializing System Core%s", dots),
	)
	content := lipgloss.JoinVertical(lipgloss.Center, logo, sub)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// overlayModal renders the danger zone confirmation modal
func (m Model) overlayModal(_ string) string {
	modalContent := lipgloss.NewStyle().Bold(true).Foreground(cRed).Render("⚠️  DANGER ZONE ⚠️") + "\n\n"
	modalContent += "You have selected Critical Runtimes.\nUpdating Node/Python may break your projects.\n\n"
	modalContent += lipgloss.NewStyle().Foreground(cWhite).Render("Are you sure? (y/N)")
	modalBox := modalStyle.Render(modalContent)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modalBox)
}

// --- Search Bar Rendering ---

func (m Model) renderSearchBar() string {
	cursor := ""
	if m.state == stateSearch {
		cursor = "█" // Blinking cursor
	}

	searchText := m.searchQuery + cursor
	if m.searchQuery == "" && m.state == stateSearch {
		searchText = cursor
	}

	resultCount := ""
	if m.filteredItems != nil {
		resultCount = fmt.Sprintf(" (%d results)", len(m.filteredItems))
	}

	label := lipgloss.NewStyle().
		Foreground(cYellow).
		Render("Search: ")

	input := lipgloss.NewStyle().
		Foreground(cWhite).
		Background(cDark).
		Padding(0, 1).
		Render(searchText)

	results := lipgloss.NewStyle().
		Foreground(cGray).
		Render(resultCount)

	hint := ""
	if m.state == stateSearch {
		hint = lipgloss.NewStyle().
			Foreground(cGray).
			Render("\n[ESC] Cancel • [ENTER] Confirm")
	}

	return label + input + results + hint
}

// --- Progress Bar Rendering ---

func (m Model) renderProgressBar() string {
	// Calculate progress: completed / total
	completed := m.totalUpdate - m.updating
	percent := float64(completed) / float64(m.totalUpdate)

	progressBarView := m.progress.ViewAs(percent)

	// Format: "Updating: 5/10 [=========>     ] 50%"
	label := fmt.Sprintf("Progress: %d/%d completed", completed, m.totalUpdate)

	return lipgloss.NewStyle().
		Foreground(cBlue).
		Render(label) + "\n" + progressBarView
}

// --- Header Rendering ---

func (m Model) getHeaderText() string {
	switch m.state {
	case stateUpdating:
		return fmt.Sprintf(" UPDATING (%d remaining)... ", m.updating)
	case stateSummary:
		return " UPDATE SUMMARY "
	default:
		if m.loading > 0 {
			return fmt.Sprintf(" SPARK DASHBOARD (Scanning %d...)", m.loading)
		}
		return " SPARK DASHBOARD "
	}
}

func (m Model) renderHeader() string {
	headerTxt := m.getHeaderText()
	return lipgloss.NewStyle().
		Background(cBlue).
		Foreground(cWhite).
		Bold(true).
		Padding(0, 1).
		Render(headerTxt)
}

// --- Grid Rendering ---

func (m Model) renderGrid() string {
	col1 := m.renderColumn([]core.Category{
		core.CategoryCode,
		core.CategoryTerm,
		core.CategoryIDE,
		core.CategoryProd,
	}, []string{"C", "T", "I", "P"})

	col2 := m.renderColumn([]core.Category{
		core.CategoryInfra,
		core.CategoryUtils,
		core.CategoryRuntime,
		core.CategorySys,
	}, []string{"F", "U", "R", "S"})

	return lipgloss.JoinHorizontal(lipgloss.Top, col1, "  ", col2)
}

func (m Model) renderColumn(categories []core.Category, keys []string) string {
	var cards []string
	for i, cat := range categories {
		card := m.renderCategoryCard(cat, keys[i])
		if card != "" {
			cards = append(cards, card)
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, cards...)
}

func (m Model) renderCategoryCard(targetCat core.Category, key string) string {
	var rows []string
	title := fmt.Sprintf("[%s] %s",
		lipgloss.NewStyle().Foreground(cGreen).Render(key),
		getCategoryLabel(targetCat))

	hasItems := false
	for i, item := range m.items {
		if item.Tool.Category != targetCat {
			continue
		}
		// Skip items not matching filter
		if !m.isItemVisible(i) {
			continue
		}
		hasItems = true
		lineStr := m.renderToolLine(i, item)
		rows = append(rows, lineStr)
	}

	if !hasItems {
		return ""
	}

	body := strings.Join(rows, "\n")
	return cardStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			cardTitleStyle.Render(title),
			body))
}

// --- Tool Line Rendering ---

func (m Model) renderToolLine(index int, item core.ToolState) string {
	cursor := m.getCursorIndicator(index)
	checked := m.getCheckedIndicator(index)
	status := m.renderItemStatus(index, item)
	name := m.formatToolName(item.Tool.Name)

	lineStr := fmt.Sprintf("%s %s %-22s %s", cursor, checked, name, status)

	// Apply styling based on selection
	if m.cursor == index && m.state == stateMain {
		return selectedItemStyle.Render(lineStr)
	}
	return dimmedItemStyle.Render(lineStr)
}

func (m Model) getCursorIndicator(index int) string {
	if m.cursor == index && m.state == stateMain {
		return lipgloss.NewStyle().Foreground(cGreen).Bold(true).Render("❯")
	}
	return " "
}

func (m Model) getCheckedIndicator(index int) string {
	if m.checked[index] {
		return lipgloss.NewStyle().Foreground(cGreen).Render("[✔]")
	}
	// Make empty checkbox more visible with subtle color
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#4B5563")).Render("[ ]")
}

func (m Model) formatToolName(name string) string {
	if len(name) > 22 {
		return name[:21] + "…"
	}
	return name
}

func (m Model) renderItemStatus(index int, item core.ToolState) string {
	// During update or summary phase
	if m.state == stateUpdating || m.state == stateSummary {
		switch item.Status {
		case core.StatusUpdating:
			return statusUpdating
		case core.StatusUpdated:
			return statusSuccess
		case core.StatusFailed:
			return statusFailed
		}

		// Not yet updated but selected
		if m.checked[index] {
			return lipgloss.NewStyle().Foreground(cGray).Render("⏳ Pending...")
		}

		// Not selected - show dimmed version
		return lipgloss.NewStyle().Foreground(cDark).Render(item.LocalVersion)
	}

	// Normal state - checking or displaying version
	if item.Status == core.StatusChecking {
		return statusChecking
	}

	switch item.Status {
	case core.StatusMissing:
		return statusMissing
	case core.StatusInstalled:
		// Special case: if version string is "MISSING", show it as warning
		if item.LocalVersion == "MISSING" {
			return lipgloss.NewStyle().Foreground(cYellow).Render("MISSING")
		}
		return statusUpToDate
	default:
		// Catch any "MISSING" strings that slip through
		if item.LocalVersion == "MISSING" {
			return lipgloss.NewStyle().Foreground(cYellow).Render("MISSING")
		}
		return item.LocalVersion
	}
}

// --- Help Bar Rendering ---

func (m Model) getHelpText() string {
	switch m.state {
	case stateSearch:
		return "[Type to search] • [ESC] Cancel • [ENTER] Confirm"
	case stateUpdating:
		return "[UPDATING IN PROGRESS... PLEASE WAIT]"
	case stateSummary:
		return "[UPDATE COMPLETE] Press any key to return to dashboard"
	default:
		help := "[SPACE] Select • [G/A] Group • [/] Search • [D] Dry-Run • [ENTER] Update • [Q] Quit"
		if m.searchQuery != "" {
			help = "[Filter active] " + help + " • [ESC] Clear filter"
		}
		return help
	}
}

func (m Model) renderHelpBar() string {
	help := m.getHelpText()
	return lipgloss.NewStyle().Foreground(cGray).Render("\n\n" + help)
}

// --- Utility Functions ---

func getCategoryLabel(c core.Category) string {
	switch c {
	case core.CategoryCode:
		return "AI Development"
	case core.CategoryTerm:
		return "Terminals"
	case core.CategoryIDE:
		return "IDEs & Editors"
	case core.CategoryProd:
		return "Productivity"
	case core.CategoryInfra:
		return "Infrastructure"
	case core.CategoryUtils:
		return "Utilities"
	case core.CategoryRuntime:
		return "Runtimes"
	case core.CategorySys:
		return "System"
	default:
		return string(c)
	}
}
