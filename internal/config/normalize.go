package config

import (
	"net/url"
	"strings"

	"github.com/bprendie/weazltunes/internal/security"
)

func (c *Config) normalize() {
	c.Presets = cleanPresetList(c.Presets, PresetLimit)
	c.MyStations = cleanPresetList(c.MyStations, 0)
}

func cleanPresetList(stations []Preset, limit int) []Preset {
	out := make([]Preset, 0, len(stations))
	for _, station := range stations {
		cleaned, ok := cleanPreset(station)
		if !ok || containsURL(out, cleaned.URL) {
			continue
		}
		out = append(out, cleaned)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out
}

func cleanPreset(st Preset) (Preset, bool) {
	st.Name = cleanName(st.Name)
	st.URL = security.SanitizeURL(st.URL)
	if st.URL == "" || !validStreamURL(st.URL) {
		return Preset{}, false
	}
	if st.Name == "" {
		st.Name = st.URL
	}
	return st, true
}

func cleanName(name string) string {
	return security.SanitizeStationName(name)
}

func containsURL(stations []Preset, streamURL string) bool {
	for _, existing := range stations {
		if strings.EqualFold(existing.URL, streamURL) {
			return true
		}
	}
	return false
}

func validStreamURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return parsed.Host != "" && (parsed.Scheme == "http" || parsed.Scheme == "https")
}
