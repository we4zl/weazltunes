package tui

import tea "github.com/charmbracelet/bubbletea"

func handleKey(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.replacing != nil {
		if handled := m.handlePresetReplacement(msg.String()); handled {
			return m, noop
		}
	}
	if index, ok := presetHotkeyIndex(msg.String()); ok {
		m.playPresetSlot(index)
		return m, noop
	}
	switch msg.String() {
	case "ctrl+c", "q":
		m.stopMeter()
		m.player.Stop()
		return m, tea.Quit
	case "1":
		m.showPresets()
	case "2":
		m.mode = modeSoma
		m.refreshDelegate()
		spin := m.beginSearch("loading SomaFM")
		return m, tea.Batch(spin, m.loadSoma())
	case "3":
		m.showXiphSearch()
		return m, noop
	case "4":
		m.showMyStations()
	case "ctrl+p":
		promoteSelectedOrInput(&m)
	case "ctrl+r":
		m.startRenameSelected()
	case "ctrl+d":
		m.deleteSelected()
	case "[":
		m.moveSelected(-1)
	case "]":
		m.moveSelected(1)
	case " ":
		m.togglePause()
	case "enter":
		return handleEnter(m)
	case "/":
		m.input.Focus()
		return m, noop
	case "esc":
		m.renaming = nil
		m.replacing = nil
		m.input.SetValue("")
		m.input.Prompt = tunePrompt
		m.input.Blur()
		return m, noop
	case "s":
		m.stop()
	}
	return m, nil
}

func handleEnter(m Model) (Model, tea.Cmd) {
	if m.input.Focused() {
		if m.renaming != nil {
			m.finishRename()
			return m, nil
		}
		if st, ok := stationFromURL(m.input.Value()); ok {
			m.saveMyStation(st)
			m.showMyStations()
			m.play(st)
			m.input.SetValue("")
			m.input.Blur()
			return m, nil
		}
		spin := m.beginSearch("searching Xiph")
		return m, tea.Batch(spin, m.searchXiph(m.input.Value()))
	}
	if item, ok := m.list.SelectedItem().(stationItem); ok {
		m.play(directoryStation(item))
	}
	return m, nil
}

func promoteSelectedOrInput(m *Model) {
	if st, ok := stationFromURL(m.input.Value()); ok {
		m.saveMyStation(st)
		m.promotePresetOrPrompt(st)
		m.input.SetValue("")
		m.input.Blur()
		return
	}
	if item, ok := m.list.SelectedItem().(stationItem); ok {
		m.promotePresetOrPrompt(directoryStation(item))
		return
	}
	if m.playing != nil {
		m.promotePresetOrPrompt(*m.playing)
	}
}

func noop() tea.Msg {
	return nil
}

func presetHotkeyIndex(key string) (int, bool) {
	switch key {
	case "alt+1":
		return 0, true
	case "alt+2":
		return 1, true
	case "alt+3":
		return 2, true
	case "alt+4":
		return 3, true
	case "alt+5":
		return 4, true
	case "alt+6":
		return 5, true
	case "alt+7":
		return 6, true
	case "alt+8":
		return 7, true
	default:
		return 0, false
	}
}
