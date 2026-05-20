package tui

import (
	"github.com/bprendie/weazltunes/internal/directory"
	"github.com/charmbracelet/bubbles/list"
)

func (m *Model) showPresets() {
	m.mode = modePresets
	m.refreshDelegate()
	m.setStations(presetStations(m.cfg))
	m.status = "presets"
	m.err = ""
}

func (m *Model) showMyStations() {
	m.mode = modeMyStations
	m.refreshDelegate()
	m.setStations(myStations(m.cfg))
	m.status = "my stations"
	m.err = ""
}

func (m *Model) setStations(stations []directory.Station) {
	m.refreshDelegate()
	items := make([]list.Item, 0, len(stations))
	for _, st := range stations {
		items = append(items, stationItem{station: st})
	}
	m.list.SetItems(items)
}

func (m *Model) refreshDelegate() {
	m.list.SetDelegate(stationDelegate{styles: m.styles, mode: m.mode})
	m.list.Title = modeTitle(m.mode)
}

func directoryStation(item stationItem) directory.Station {
	return item.station
}
