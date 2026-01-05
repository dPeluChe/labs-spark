package tui

import "github.com/charmbracelet/lipgloss"

// Color Palette
var (
	cGreen  = lipgloss.Color("#04B575")
	cBlue   = lipgloss.Color("#2E7DE1")
	cPurple = lipgloss.Color("#A78BFA")
	cGray   = lipgloss.Color("#6B7280")
	cWhite  = lipgloss.Color("#FFFFFF")
	cDark   = lipgloss.Color("#1F2937")
	cYellow = lipgloss.Color("#F59E0B")
	cRed    = lipgloss.Color("#EF4444")
)

// Layout Styles
var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cPurple).
			Padding(0, 1).
			MarginBottom(1).
			Width(60)

	cardTitleStyle = lipgloss.NewStyle().
			Foreground(cPurple).
			Bold(true).
			Padding(0, 1).
			Background(cDark)

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(cRed).
			Padding(1, 2).
			Align(lipgloss.Center).
			Width(50)
)

// Splash Screen Styles
var (
	splashTitleStyle = lipgloss.NewStyle().
				Foreground(cBlue).
				Bold(true).
				MarginBottom(1)

	splashSubtitleStyle = lipgloss.NewStyle().
				Foreground(cGray).
				Italic(true)
)

// Item Rendering Styles
var (
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(cWhite).
				Background(lipgloss.Color("#2D3748")). // Dark blue-gray background
				Bold(true).
				Padding(0, 1)

	dimmedItemStyle = lipgloss.NewStyle().
			Foreground(cGray)
)

// Status Indicators (Pre-rendered)
var (
	statusChecking = lipgloss.NewStyle().
			Foreground(cYellow).
			Render("⟳ Checking...")

	statusUpToDate = lipgloss.NewStyle().
			Foreground(cGray).
			Render("✔ Up to date")

	statusOutdated = lipgloss.NewStyle().
			Foreground(cYellow).
			Bold(true)

	statusMissing = lipgloss.NewStyle().
			Foreground(cRed).
			Render("○ Not Installed")

	statusUpdating = lipgloss.NewStyle().
			Foreground(cBlue).
			Render("➜ Updating...")

	statusSuccess = lipgloss.NewStyle().
			Foreground(cGreen).
			Render("✔ Updated")

	statusFailed = lipgloss.NewStyle().
			Foreground(cRed).
			Render("✘ Failed")
)

// ASCII Art
const sparkArt = `
   _____ ____  ___  ____  __ __
  / ___// __ \/   |/ __ \/ //_/
  \__ \/ /_/ / /| / /_/ / ,<
 ___/ / ____/ ___ / _, _/ /| |
/____/_/   /_/  |/_/ |_/_/ |_|
`
