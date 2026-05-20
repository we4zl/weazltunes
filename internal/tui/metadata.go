package tui

import (
	"strings"

	"github.com/bprendie/weazltunes/internal/player"
	tea "github.com/charmbracelet/bubbletea"
)

func fetchTitle(p *player.Player) tea.Cmd {
	return func() tea.Msg {
		title, err := p.Title()
		if err != nil {
			return nil
		}
		title = strings.TrimSpace(title)
		if title == "" {
			return nil
		}
		return titleMsg{title: title}
	}
}
