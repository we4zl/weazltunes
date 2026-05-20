package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (m Model) View() string {
	var b strings.Builder
	contentWidth := max(40, m.width-4)
	if contentWidth >= maxLineWidth(logo) {
		b.WriteString(renderLogo(logo, contentWidth))
	} else {
		b.WriteString(m.styles.header.Render("WeazlTunes"))
	}
	b.WriteString("\n\n")
	b.WriteString(renderHelpMenu(contentWidth))
	b.WriteString("\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(m.styles.panel.Width(contentWidth).Render(m.visualizer.View(m.styles)))
	b.WriteString("\n")
	b.WriteString(m.trackPanel(contentWidth))
	b.WriteString("\n")
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(m.statusLine())
	if m.err != "" {
		b.WriteString("\n" + m.styles.error.Render(m.err))
	}
	return m.styles.frame.Render(b.String())
}

func (m Model) trackPanel(width int) string {
	title := m.trackTitle
	if title == "" {
		title = "waiting for stream title"
	}
	if m.playing == nil {
		title = "no track metadata"
	}
	label := m.styles.help.Render("track: ")
	titleWidth := max(8, width-lipgloss.Width("track: ")-6)
	title = ansi.Truncate(title, titleWidth, "...")
	line := label + gradientText(title, trackStops)
	return m.styles.track.Width(width).Height(1).Render(line)
}

func (m Model) statusLine() string {
	if m.playing != nil {
		state := "now: "
		if m.paused {
			state = "paused: "
		}
		meter := " meter:fallback"
		if m.energy.Live {
			meter = " " + meterBadge("LIVE", crushMint)
		} else {
			meter = " " + meterBadge("SYNTH", crushGold)
		}
		return m.styles.status.Render(ansi.Wordwrap(state+m.playing.Name+meter+" <"+m.playing.URL+">", max(20, m.width-4), " /_-<>"))
	}
	return m.styles.status.Render(ansi.Wordwrap(m.status, max(20, m.width-4), " /_-"))
}

func meterBadge(label string, color lipgloss.Color) string {
	return lipgloss.NewStyle().
		Foreground(color).
		Border(lipgloss.NormalBorder(), false, true, false, true).
		BorderForeground(crushPurple).
		Padding(0, 1).
		Render(label)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
