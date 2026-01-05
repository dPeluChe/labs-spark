package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dpeluche/spark/internal/core"
)

// ViewSummary renders the update summary screen
func (m Model) ViewSummary() string {
	stats := m.calculateSummaryStats()

	// Title
	title := lipgloss.NewStyle().
		Background(cGreen).
		Foreground(cWhite).
		Bold(true).
		Padding(0, 1).
		Render(" ✓ UPDATE COMPLETE ")

	// Statistics box
	statsBox := m.renderStatisticsBox(stats)

	// Updated tools list
	updatedList := m.renderUpdatedToolsList(stats)

	// Failed tools list (if any)
	failedList := ""
	if stats.Failed > 0 {
		failedList = "\n" + m.renderFailedToolsList(stats)
	}

	// Help text
	help := lipgloss.NewStyle().
		Foreground(cGray).
		Render("\n\nPress any key to exit...")

	content := title + "\n\n" + statsBox + updatedList + failedList + help
	return appStyle.Render(content)
}

// SummaryStats holds statistics about the update session
type SummaryStats struct {
	Total     int
	Successful int
	Failed    int
	Skipped   int
	UpdatedTools []core.ToolState
	FailedTools  []core.ToolState
}

// calculateSummaryStats computes statistics from the update session
func (m Model) calculateSummaryStats() SummaryStats {
	stats := SummaryStats{
		Total:        m.totalUpdate,
		UpdatedTools: []core.ToolState{},
		FailedTools:  []core.ToolState{},
	}

	for i, item := range m.items {
		if !m.checked[i] {
			continue
		}

		switch item.Status {
		case core.StatusUpdated:
			stats.Successful++
			stats.UpdatedTools = append(stats.UpdatedTools, item)
		case core.StatusFailed:
			stats.Failed++
			stats.FailedTools = append(stats.FailedTools, item)
		default:
			stats.Skipped++
		}
	}

	return stats
}

// renderStatisticsBox creates a visual box with summary statistics
func (m Model) renderStatisticsBox(stats SummaryStats) string {
	var lines []string

	// Header
	lines = append(lines, lipgloss.NewStyle().
		Foreground(cPurple).
		Bold(true).
		Render("SUMMARY STATISTICS"))

	lines = append(lines, "")

	// Success rate
	successRate := 0.0
	if stats.Total > 0 {
		successRate = float64(stats.Successful) / float64(stats.Total) * 100
	}

	lines = append(lines, fmt.Sprintf("Total Updates:   %d", stats.Total))
	lines = append(lines, lipgloss.NewStyle().
		Foreground(cGreen).
		Render(fmt.Sprintf("✓ Successful:    %d", stats.Successful)))

	if stats.Failed > 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(cRed).
			Render(fmt.Sprintf("✘ Failed:        %d", stats.Failed)))
	}

	if stats.Skipped > 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(cGray).
			Render(fmt.Sprintf("○ Skipped:       %d", stats.Skipped)))
	}

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Success Rate:    %.1f%%", successRate))

	box := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cPurple).
		Padding(1, 2).
		Width(40).
		Render(box)
}

// renderUpdatedToolsList creates a list of successfully updated tools
func (m Model) renderUpdatedToolsList(stats SummaryStats) string {
	if stats.Successful == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(cGreen).
		Bold(true).
		Render("UPDATED TOOLS"))
	lines = append(lines, "")

	for _, tool := range stats.UpdatedTools {
		line := fmt.Sprintf("  ✓ %s", tool.Tool.Name)
		if tool.LocalVersion != "" && tool.LocalVersion != "..." {
			line += lipgloss.NewStyle().
				Foreground(cGray).
				Render(fmt.Sprintf(" (%s)", tool.LocalVersion))
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderFailedToolsList creates a list of failed updates
func (m Model) renderFailedToolsList(stats SummaryStats) string {
	if stats.Failed == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(cRed).
		Bold(true).
		Render("FAILED UPDATES"))
	lines = append(lines, "")

	for _, tool := range stats.FailedTools {
		line := fmt.Sprintf("  ✘ %s", tool.Tool.Name)
		if tool.Message != "" {
			line += lipgloss.NewStyle().
				Foreground(cGray).
				Render(fmt.Sprintf(" - %s", tool.Message))
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
