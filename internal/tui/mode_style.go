package tui

import "github.com/charmbracelet/lipgloss"

func modeTitle(m mode) string {
	label := "top 8 presets"
	color := crushGold
	switch m {
	case modeSoma:
		label = "SomaFM directory"
		color = crushPink
	case modeXiph:
		label = "Xiph search"
		color = crushPurple
	case modeMyStations:
		label = "my stations"
		color = crushMint
	}
	rail := lipgloss.NewStyle().Foreground(crushPurple).Bold(true).Render("╱╱")
	text := lipgloss.NewStyle().Foreground(color).Bold(true).Render(" " + label)
	return rail + text
}
