package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bprendie/weazltunes/internal/config"
	"github.com/bprendie/weazltunes/internal/directory"
)

func (m *Model) saveMyStation(st directory.Station) {
	m.cfg.SaveMyStation(presetFromStation(st))
	if err := config.Save(m.cfg); err != nil {
		m.err = err.Error()
		return
	}
	m.status = "saved to my stations"
}

func (m *Model) promotePreset(st directory.Station) {
	m.cfg.PromotePreset(presetFromStation(st))
	if err := config.Save(m.cfg); err != nil {
		m.err = err.Error()
		return
	}
	m.status = "added to presets"
}

func (m *Model) promotePresetOrPrompt(st directory.Station) {
	if canPromoteDirectly(m.cfg.Presets, st.URL) {
		m.promotePreset(st)
		m.showPresets()
		return
	}
	m.replacing = &st
	m.showPresets()
	m.status = "preset bank full - press 1-8 to replace a slot"
	m.err = ""
}

func (m *Model) handlePresetReplacement(key string) bool {
	if key == "esc" || key == "q" || key == "ctrl+c" {
		return false
	}
	slot, err := strconv.Atoi(key)
	if err != nil || slot < 1 || slot > config.PresetLimit {
		m.err = "press 1-8 to choose a preset slot, or esc to cancel"
		return true
	}
	if m.replacing == nil {
		return false
	}
	st := *m.replacing
	if !m.cfg.ReplacePreset(slot-1, presetFromStation(st)) || !m.saveConfig() {
		return true
	}
	m.replacing = nil
	m.showPresets()
	m.list.Select(slot - 1)
	m.status = fmt.Sprintf("replaced preset %d with %s", slot, st.Name)
	m.err = ""
	return true
}

func canPromoteDirectly(presets []config.Preset, url string) bool {
	if len(presets) < config.PresetLimit {
		return true
	}
	for _, preset := range presets {
		if strings.EqualFold(preset.URL, url) {
			return true
		}
	}
	return false
}

func (m *Model) startRenameSelected() {
	item, ok := m.list.SelectedItem().(stationItem)
	if !ok || (m.mode != modePresets && m.mode != modeMyStations) {
		m.err = "select a preset or my station to rename"
		return
	}
	st := directoryStation(item)
	m.renaming = &st
	m.input.SetValue(st.Name)
	m.input.Prompt = "rename > "
	m.input.Focus()
	m.status = "enter a new station name"
	m.err = ""
}

func (m *Model) finishRename() {
	if m.renaming == nil {
		return
	}
	renamed := false
	if m.mode == modePresets {
		renamed = m.cfg.RenamePreset(m.renaming.URL, m.input.Value())
	} else if m.mode == modeMyStations {
		renamed = m.cfg.RenameMyStation(m.renaming.URL, m.input.Value())
	}
	if !renamed || !m.saveConfig() {
		return
	}
	m.input.SetValue("")
	m.input.Prompt = tunePrompt
	m.input.Blur()
	m.renaming = nil
	m.refreshEditableList()
	m.status = "renamed"
}

func (m *Model) deleteSelected() {
	item, ok := m.list.SelectedItem().(stationItem)
	if !ok || (m.mode != modePresets && m.mode != modeMyStations) {
		m.err = "select a preset or my station to delete"
		return
	}
	st := directoryStation(item)
	deleted := false
	if m.mode == modePresets {
		deleted = m.cfg.DeletePreset(st.URL)
	} else {
		deleted = m.cfg.DeleteMyStation(st.URL)
	}
	if !deleted || !m.saveConfig() {
		return
	}
	if m.playing != nil && m.playing.URL == st.URL {
		m.stop()
	}
	m.refreshEditableList()
	m.status = "deleted " + st.Name
}

func (m *Model) moveSelected(delta int) {
	item, ok := m.list.SelectedItem().(stationItem)
	if !ok || (m.mode != modePresets && m.mode != modeMyStations) {
		m.err = "select a preset or my station to move"
		return
	}
	st := directoryStation(item)
	next := -1
	moved := false
	if m.mode == modePresets {
		next, moved = m.cfg.MovePreset(st.URL, delta)
	} else {
		next, moved = m.cfg.MoveMyStation(st.URL, delta)
	}
	if !moved || !m.saveConfig() {
		return
	}
	m.refreshEditableList()
	if next >= 0 {
		m.list.Select(next)
	}
	m.status = "moved " + st.Name
}

func (m *Model) saveConfig() bool {
	if err := config.Save(m.cfg); err != nil {
		m.err = err.Error()
		return false
	}
	return true
}

func (m *Model) refreshEditableList() {
	if m.mode == modePresets {
		m.showPresets()
	} else {
		m.showMyStations()
	}
}
