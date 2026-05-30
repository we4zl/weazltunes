package config

import "strings"

func upsert(stations []Preset, st Preset) []Preset {
	out := []Preset{st}
	for _, existing := range stations {
		if strings.EqualFold(existing.URL, st.URL) {
			continue
		}
		out = append(out, existing)
	}
	return out
}

func renameByURL(stations []Preset, url, name string) bool {
	name = cleanName(name)
	if name == "" {
		return false
	}
	for i := range stations {
		if strings.EqualFold(stations[i].URL, url) {
			stations[i].Name = name
			return true
		}
	}
	return false
}

func deleteByURL(stations []Preset, url string) ([]Preset, bool) {
	next := removeByURL(stations, url)
	return next, len(next) != len(stations)
}

func removeByURL(stations []Preset, url string) []Preset {
	out := stations[:0]
	for _, station := range stations {
		if strings.EqualFold(station.URL, url) {
			continue
		}
		out = append(out, station)
	}
	return out
}

func moveByURL(stations []Preset, url string, delta int) ([]Preset, int) {
	for i := range stations {
		if !strings.EqualFold(stations[i].URL, url) {
			continue
		}
		next := i + delta
		if next < 0 || next >= len(stations) {
			return stations, -1
		}
		stations[i], stations[next] = stations[next], stations[i]
		return stations, next
	}
	return stations, -1
}
