package tui

import (
	"testing"

	"github.com/bprendie/weazltunes/internal/config"
	"github.com/bprendie/weazltunes/internal/directory"
	tea "github.com/charmbracelet/bubbletea"
)

func TestStationFromURL(t *testing.T) {
	st, ok := stationFromURL(" https://example.com/live/stream.mp3 ")
	if !ok {
		t.Fatal("expected URL station")
	}
	if st.Name != "stream" {
		t.Fatalf("name = %q", st.Name)
	}
	if st.URL != "https://example.com/live/stream.mp3" {
		t.Fatalf("url = %q", st.URL)
	}
	if st.Source != "My Station" {
		t.Fatalf("source = %q", st.Source)
	}
}

func TestStationFromURLRejectsNonHTTP(t *testing.T) {
	if _, ok := stationFromURL("file:///tmp/song.mp3"); ok {
		t.Fatal("expected file URL to be rejected")
	}
	if _, ok := stationFromURL("not a url"); ok {
		t.Fatal("expected invalid URL to be rejected")
	}
}

func TestStationFromYouTubeWatchURL(t *testing.T) {
	st, ok := stationFromURL("https://www.youtube.com/watch?v=abc123")
	if !ok {
		t.Fatal("expected YouTube URL station")
	}
	if st.Name != "youtube abc123" {
		t.Fatalf("name = %q", st.Name)
	}
}

func TestConfigStationHelpers(t *testing.T) {
	cfg := config.Config{Presets: []config.Preset{{Name: "Preset", URL: "https://example.com/preset"}}}
	stations := presetStations(cfg)
	if len(stations) != 1 || stations[0].Source != "Preset" {
		t.Fatalf("unexpected preset stations: %#v", stations)
	}

	preset := presetFromStation(directory.Station{Name: "Station", URL: "https://example.com/station"})
	if preset.Name != "Station" || preset.URL != "https://example.com/station" {
		t.Fatalf("unexpected preset: %#v", preset)
	}
}

func TestThreeShowsXiphWithoutTypingIntoInput(t *testing.T) {
	model := New(config.Config{Presets: []config.Preset{{Name: "Preset", URL: "https://example.com/preset"}}})
	next, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	got := next.(Model)

	if got.mode != modeXiph {
		t.Fatalf("mode = %v, want modeXiph", got.mode)
	}
	if !got.input.Focused() {
		t.Fatal("expected input to be focused")
	}
	if got.input.Value() != "" {
		t.Fatalf("input value = %q, want empty", got.input.Value())
	}
	if len(got.list.Items()) != 0 {
		t.Fatalf("list items = %d, want 0", len(got.list.Items()))
	}
}

func TestSlashFocusesInputWithoutTypingSlash(t *testing.T) {
	model := New(config.Config{})
	next, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	got := next.(Model)

	if !got.input.Focused() {
		t.Fatal("expected input to be focused")
	}
	if got.input.Value() != "" {
		t.Fatalf("input value = %q, want empty", got.input.Value())
	}
}

func TestFullPresetBankPromptsForReplacement(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	model := New(config.Config{Presets: fullPresets()})
	model.mode = modeXiph
	model.setStations([]directory.Station{{Name: "New", URL: "https://example.com/new"}})

	next, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlP, Alt: false})
	got := next.(Model)

	if got.replacing == nil {
		t.Fatal("expected replacement prompt")
	}
	if got.cfg.Presets[0].URL != "https://example.com/0" {
		t.Fatalf("preset changed before slot selection: %q", got.cfg.Presets[0].URL)
	}
}

func TestPresetReplacementUsesSelectedSlot(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	model := New(config.Config{Presets: fullPresets()})
	model.replacing = &directory.Station{Name: "New", URL: "https://example.com/new"}

	next, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	got := next.(Model)

	if got.replacing != nil {
		t.Fatal("expected replacement prompt to close")
	}
	if got.cfg.Presets[2].URL != "https://example.com/new" {
		t.Fatalf("slot 3 URL = %q", got.cfg.Presets[2].URL)
	}
}

func TestPresetHotkeyIndex(t *testing.T) {
	tests := map[string]int{
		"alt+1": 0,
		"alt+2": 1,
		"alt+3": 2,
		"alt+4": 3,
		"alt+5": 4,
		"alt+6": 5,
		"alt+7": 6,
		"alt+8": 7,
	}
	for key, want := range tests {
		got, ok := presetHotkeyIndex(key)
		if !ok || got != want {
			t.Fatalf("presetHotkeyIndex(%q) = %d, %v; want %d, true", key, got, ok, want)
		}
	}
	if _, ok := presetHotkeyIndex("ctrl+3"); ok {
		t.Fatal("ctrl+3 should not be a preset hotkey")
	}
	if _, ok := presetHotkeyIndex("alt+9"); ok {
		t.Fatal("alt+9 should not be a preset hotkey")
	}
}

func fullPresets() []config.Preset {
	presets := make([]config.Preset, config.PresetLimit)
	for i := range presets {
		presets[i] = config.Preset{Name: string(rune('A' + i)), URL: "https://example.com/" + string(rune('0'+i))}
	}
	return presets
}
