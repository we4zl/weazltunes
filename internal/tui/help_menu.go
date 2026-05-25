package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type helpAction struct {
	key   string
	label string
	color lipgloss.Color
}

func renderHelpMenu(width int) string {
	groups := [][]helpAction{
		{
			{"1", "presets", crushGold},
			{"2", "SomaFM", crushPink},
			{"3", "Xiph", crushPurple},
			{"4", "my stations", crushMint},
		},
		{
			{"/", "tune", crushPurple},
			{"brackets", "move", crushGold},
			{"ctrl+p", "preset", crushGold},
			{"ctrl+r", "rename", crushPurple},
			{"ctrl+d", "delete", crushPink},
		},
		{
			{"space", "pause", crushMint},
			{"enter", "play/add", crushMint},
			{"ctrl+1-8", "play preset", crushMint},
			{"s", "stop", crushGold},
		},
		{
			{"q", "quit", muted},
		},
	}
	var parts []string
	for _, group := range groups {
		var actions []string
		for _, action := range group {
			actions = append(actions, renderHelpAction(action))
		}
		parts = append(parts, strings.Join(actions, "  "))
	}
	separator := lipgloss.NewStyle().Foreground(crushPurple).Faint(true).Render("  ╱  ")
	return wrapStyled(strings.Join(parts, separator), width)
}

func renderHelpAction(action helpAction) string {
	key := lipgloss.NewStyle().Foreground(crushGold).Bold(true).Render("[" + action.key + "]")
	label := lipgloss.NewStyle().Foreground(action.color).Render(" " + action.label)
	return key + label
}

func wrapStyled(s string, width int) string {
	if lipgloss.Width(s) <= width {
		return s
	}
	var out strings.Builder
	lineWidth := 0
	for _, part := range strings.Split(s, "  ") {
		partWidth := lipgloss.Width(part)
		sep := 0
		if lineWidth > 0 {
			sep = 2
		}
		if lineWidth > 0 && lineWidth+sep+partWidth > width {
			out.WriteByte('\n')
			lineWidth = 0
			sep = 0
		}
		if sep > 0 {
			out.WriteString("  ")
			lineWidth += sep
		}
		out.WriteString(part)
		lineWidth += partWidth
	}
	return out.String()
}
