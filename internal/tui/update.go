package tui

import (
	"time"

	"github.com/bprendie/weazltunes/internal/audio"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		next, cmd := handleKey(m, msg)
		m = next
		if cmd != nil {
			return m, cmd
		}
	case spinner.TickMsg:
		if m.searching {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	case loadedMsg:
		m.setStations(msg.stations)
		m.status = msg.status
		m.err = ""
		m.searching = false
		m.input.Blur()
	case errMsg:
		m.err = msg.err.Error()
		m.searching = false
	case titleMsg:
		if msg.title != "" {
			m.trackTitle = msg.title
		}
	case tickMsg:
		m.drainMeter()
		m.visualizer.Step(m.playing != nil && !m.paused, m.energy)
		if m.playing != nil && !m.paused && time.Time(msg).Sub(m.titlePoll) > 2*time.Second {
			m.titlePoll = time.Time(msg)
			cmds = append(cmds, fetchTitle(m.player))
		}
		cmds = append(cmds, tick())
	}

	next, cmd := m.updateFocused(msg)
	m = next
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) beginSearch(status string) tea.Cmd {
	m.status = status
	m.err = ""
	m.searching = true
	return m.spinner.Tick
}

func (m *Model) drainMeter() {
	if m.meter == nil {
		m.energy = audio.Sample{}
		return
	}
	for {
		select {
		case sample, ok := <-m.meter.Samples():
			if !ok {
				m.meter = nil
				return
			}
			m.energy = sample
		default:
			return
		}
	}
}

func (m *Model) resize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(max(20, width-8), max(8, height-21))
	m.input.Width = max(20, width-20)
}

func (m Model) updateFocused(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == " " {
		return m, nil
	}
	if m.input.Focused() {
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
