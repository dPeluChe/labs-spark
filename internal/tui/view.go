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
		// Render main view as background with summary overlay
		return m.overlaySummaryModal(m.ViewMain())
	default:
		return m.ViewMain()
	}
}

// overlaySummaryModal renders the results over the main content
func (m Model) overlaySummaryModal(background string) string {
	successCount := 0
	failCount := 0
	var failureDetails []string

	for _, item := range m.items {
		// Only count items that were actually in the update queue via their status
		if item.Status == core.StatusUpdated {
			successCount++
		} else if item.Status == core.StatusFailed {
			failCount++
			failureDetails = append(failureDetails, fmt.Sprintf("• %s: %s", item.Tool.Name, item.Message))
		}
	}

	title := lipgloss.NewStyle().
		Background(cPurple).
		Foreground(cWhite).
		Bold(true).
		Padding(0, 1).
		Render(" UPDATE COMPLETE ")

	stats := fmt.Sprintf("\n✔ Successful: %d\n✘ Failed:     %d\n", successCount, failCount)
	
	content := stats
	if len(failureDetails) > 0 {
		content += "\nErrors:\n" + lipgloss.NewStyle().Foreground(cRed).Render(strings.Join(failureDetails, "\n"))
	}
	
	content += "\n\n" + lipgloss.NewStyle().Foreground(cGray).Render("[Press ENTER to close]")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cPurple).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Center, title, content))

	// Center the box on the screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
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
	
	currentTool := ""
	if m.currentUpdate >= 0 && m.currentUpdate < len(m.items) {
		currentTool = fmt.Sprintf("• Processing: %s", m.items[m.currentUpdate].Tool.Name)
	}

	return lipgloss.NewStyle().Foreground(cBlue).Render(label) + 
		lipgloss.NewStyle().Foreground(cYellow).Bold(true).Render(" "+currentTool) + 
		"\n" + progressBarView
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

	lineStr := fmt.Sprintf("%s %s %-18s %s", cursor, checked, name, status)

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
			// Animated spinner: ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
			frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			frame := frames[m.splashFrame%len(frames)]
			return lipgloss.NewStyle().Foreground(cBlue).Render(frame + " Updating...")
		case core.StatusUpdated:
			return lipgloss.NewStyle().Foreground(cGreen).Render("✔ " + item.LocalVersion)
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

	// Format version display
	versionStr := item.LocalVersion
	if item.RemoteVersion != "..." && item.RemoteVersion != "Checking..." && item.RemoteVersion != "Unknown" {
		if item.RemoteVersion != item.LocalVersion && item.LocalVersion != "MISSING" {
			// Show update path: 1.0.0 -> 1.1.0
			versionStr = fmt.Sprintf("%s %s %s",
				lipgloss.NewStyle().Foreground(cGray).Render(item.LocalVersion),
				lipgloss.NewStyle().Foreground(cYellow).Render("→"),
				lipgloss.NewStyle().Foreground(cGreen).Bold(true).Render(item.RemoteVersion))
		} else if item.RemoteVersion == item.LocalVersion {
			versionStr = lipgloss.NewStyle().Foreground(cGray).Render(item.LocalVersion)
		}
	}

	switch item.Status {
	case core.StatusMissing:
		return statusMissing
	case core.StatusOutdated:
		return versionStr
	case core.StatusInstalled:
		if item.LocalVersion == "MISSING" {
			return lipgloss.NewStyle().Foreground(cYellow).Render("MISSING")
		}
		// If we have remote info and they match, it's truly up to date
		if item.RemoteVersion == item.LocalVersion && item.LocalVersion != "..." {
			return lipgloss.NewStyle().Foreground(cGray).Render(item.LocalVersion + " " + "✓")
		}
		return versionStr
	default:
		if item.LocalVersion == "MISSING" {
			return lipgloss.NewStyle().Foreground(cYellow).Render("MISSING")
		}
		return versionStr
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