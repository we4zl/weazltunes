package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type stationDelegate struct {
	styles styles
	mode   mode
}

func newStationDelegate() stationDelegate {
	return stationDelegate{styles: newStyles()}
}

func (d stationDelegate) Height() int  { return 2 }
func (d stationDelegate) Spacing() int { return 0 }
func (d stationDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d stationDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	st, ok := item.(stationItem)
	if !ok || m.Width() <= 0 {
		return
	}
	width := max(10, m.Width()-4)
	titleStyle := d.styles.item
	descStyle := d.styles.help
	prefix := "  "
	if d.mode == modePresets {
		prefix = d.presetButton(index)
	}
	title := ansi.Truncate(st.Title(), max(8, width-lipgloss.Width(prefix)-2), "...")
	desc := ansi.Wordwrap(st.Description(), width, " /_-")
	desc = firstLines(desc, 1)
	if index == m.Index() {
		descStyle = lipgloss.NewStyle().Foreground(crushPurple)
		title = d.styles.selected.Render("> ") + prefix + gradientText(title, trackStops)
		fmt.Fprintf(w, "%s\n%s", title, descStyle.Render("  "+desc))
		return
	} else {
		title = prefix + title
	}
	fmt.Fprintf(w, "%s\n%s", titleStyle.Render(title), descStyle.Render("  "+desc))
}

func (d stationDelegate) presetButton(index int) string {
	if index < 0 || index > 7 {
		return "  "
	}
	return lipgloss.NewStyle().
		Foreground(crushGold).
		Border(lipgloss.NormalBorder(), false, true, false, true).
		BorderForeground(crushPurple).
		Padding(0, 1).
		Render(fmt.Sprintf("%d", index+1)) + " "
}

func firstLines(s string, limit int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= limit {
		return s
	}
	return strings.Join(lines[:limit], "\n")
}
