package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dpeluche/spark/internal/core"
)

// ViewPreview renders the dry-run preview screen
func (m Model) ViewPreview() string {
	// Title
	title := lipgloss.NewStyle().
		Background(cYellow).
		Foreground(cDark).
		Bold(true).
		Padding(0, 1).
		Render(" 🔍 UPDATE PREVIEW (DRY-RUN) ")

	// Introduction
	intro := lipgloss.NewStyle().
		Foreground(cGray).
		Render("Review the tools that will be updated. No changes will be made yet.\n")

	// Count selected tools by category
	selectedByCategory := make(map[core.Category][]core.ToolState)
	totalSelected := 0
	hasDangerous := false

	for i, item := range m.items {
		if m.checked[i] {
			totalSelected++
			selectedByCategory[item.Tool.Category] = append(selectedByCategory[item.Tool.Category], item)
			if item.Tool.Category == core.CategoryRuntime {
				hasDangerous = true
			}
		}
	}

	// Summary box
	summaryLines := []string{
		lipgloss.NewStyle().Foreground(cPurple).Bold(true).Render("SUMMARY"),
		"",
		fmt.Sprintf("Total Tools Selected: %d", totalSelected),
	}

	// Add category breakdown
	for cat, tools := range selectedByCategory {
		categoryName := getCategoryLabel(cat)
		summaryLines = append(summaryLines, fmt.Sprintf("  • %s: %d tools", categoryName, len(tools)))
	}

	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cPurple).
		Padding(1, 2).
		Width(50).
		Render(strings.Join(summaryLines, "\n"))

	// Tools list by category
	var toolsList string
	for cat, tools := range selectedByCategory {
		categoryName := getCategoryLabel(cat)
		categoryTitle := lipgloss.NewStyle().
			Foreground(cGreen).
			Bold(true).
			Render(fmt.Sprintf("\n%s", categoryName))

		toolsList += categoryTitle + "\n"

		for _, tool := range tools {
			statusIcon := "→"
			versionInfo := ""

			if tool.LocalVersion != "" && tool.LocalVersion != "..." && tool.LocalVersion != "MISSING" {
				versionInfo = lipgloss.NewStyle().
					Foreground(cGray).
					Render(fmt.Sprintf(" (current: %s)", tool.LocalVersion))
			} else if tool.Status == core.StatusMissing {
				versionInfo = lipgloss.NewStyle().
					Foreground(cYellow).
					Render(" (will install)")
			}

			line := fmt.Sprintf("  %s %s%s\n", statusIcon, tool.Tool.Name, versionInfo)
			toolsList += line
		}
	}

	// Danger warning if runtimes selected
	dangerWarning := ""
	if hasDangerous {
		dangerWarning = "\n" + lipgloss.NewStyle().
			Background(cRed).
			Foreground(cWhite).
			Bold(true).
			Padding(0, 1).
			Render(" ⚠ WARNING: Runtime updates detected - confirmation will be required ") + "\n"
	}

	// Actions
	actions := lipgloss.NewStyle().
		Foreground(cGray).
		Render("\n[ENTER] Proceed with Updates • [ESC] Cancel")

	content := title + "\n\n" + intro + "\n" + summaryBox + "\n" + toolsList + dangerWarning + actions
	return appStyle.Render(content)
}
