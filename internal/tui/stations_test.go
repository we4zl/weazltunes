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
