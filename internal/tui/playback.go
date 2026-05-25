package tui

import (
	"time"

	"github.com/bprendie/weazltunes/internal/audio"
	"github.com/bprendie/weazltunes/internal/directory"
)

func (m *Model) play(st directory.Station) {
	if err := m.player.Play(st.URL); err != nil {
		m.err = err.Error()
		return
	}
	m.playing = &st
	m.paused = false
	m.trackTitle = ""
	m.titlePoll = time.Time{}
	m.startMeter(st.URL)
	m.status = "playing " + st.Name
	m.err = ""
}

func (m *Model) playPresetSlot(index int) {
	if index < 0 || index >= len(m.cfg.Presets) {
		m.err = "preset slot is empty"
		return
	}
	m.showPresets()
	m.list.Select(index)
	m.play(directory.Station{
		Name:        m.cfg.Presets[index].Name,
		URL:         m.cfg.Presets[index].URL,
		Source:      "Preset",
		Description: m.cfg.Presets[index].URL,
	})
}

func (m *Model) togglePause() {
	paused, err := m.player.TogglePause()
	if err != nil {
		m.err = err.Error()
		return
	}
	m.paused = paused
	if paused {
		m.stopMeter()
		m.status = "paused"
		return
	}
	if m.playing != nil {
		m.startMeter(m.playing.URL)
		m.status = "playing " + m.playing.Name
	}
}

func (m *Model) stop() {
	m.player.Stop()
	m.stopMeter()
	m.playing = nil
	m.paused = false
	m.trackTitle = ""
	m.titlePoll = time.Time{}
	m.status = "stopped"
}

func (m *Model) startMeter(url string) {
	m.stopMeter()
	meter, err := audio.StartMeter(url)
	if err != nil {
		m.energy = audio.Sample{}
		return
	}
	m.meter = meter
}

func (m *Model) stopMeter() {
	if m.meter != nil {
		m.meter.Stop()
	}
	m.meter = nil
	m.energy = audio.Sample{}
}
