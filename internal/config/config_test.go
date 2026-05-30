package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPromotePresetCapsAtEightAndMovesDuplicateToTop(t *testing.T) {
	cfg := Config{}
	for i := 0; i < PresetLimit+2; i++ {
		cfg.PromotePreset(Preset{Name: string(rune('A' + i)), URL: "https://example.com/" + string(rune('a'+i))})
	}
	if len(cfg.Presets) != PresetLimit {
		t.Fatalf("expected %d presets, got %d", PresetLimit, len(cfg.Presets))
	}

	duplicate := cfg.Presets[3]
	cfg.PromotePreset(duplicate)
	if len(cfg.Presets) != PresetLimit {
		t.Fatalf("duplicate changed preset count: %d", len(cfg.Presets))
	}
	if cfg.Presets[0].URL != duplicate.URL {
		t.Fatalf("duplicate was not moved to top: got %q want %q", cfg.Presets[0].URL, duplicate.URL)
	}
}

func TestReplacePresetReplacesSelectedSlot(t *testing.T) {
	cfg := Config{Presets: []Preset{
		{Name: "One", URL: "https://example.com/1"},
		{Name: "Two", URL: "https://example.com/2"},
		{Name: "Three", URL: "https://example.com/3"},
	}}

	if !cfg.ReplacePreset(1, Preset{Name: "New", URL: "https://example.com/new"}) {
		t.Fatal("expected replace to succeed")
	}
	if cfg.Presets[1].Name != "New" {
		t.Fatalf("slot 2 = %q, want New", cfg.Presets[1].Name)
	}
	if len(cfg.Presets) != 3 {
		t.Fatalf("len = %d, want 3", len(cfg.Presets))
	}
}

func TestSaveMyStationUpsertsByURL(t *testing.T) {
	cfg := Config{}
	cfg.SaveMyStation(Preset{Name: "Old", URL: "https://example.com/stream"})
	cfg.SaveMyStation(Preset{Name: "New", URL: "https://example.com/stream"})

	if len(cfg.MyStations) != 1 {
		t.Fatalf("expected one station, got %d", len(cfg.MyStations))
	}
	if cfg.MyStations[0].Name != "New" {
		t.Fatalf("expected latest station name, got %q", cfg.MyStations[0].Name)
	}
}

func TestRenameAndDeleteByURL(t *testing.T) {
	cfg := Config{Presets: []Preset{{Name: "Old", URL: "https://example.com/stream"}}}

	if !cfg.RenamePreset("https://example.com/stream", "New") {
		t.Fatal("expected rename to succeed")
	}
	if cfg.Presets[0].Name != "New" {
		t.Fatalf("rename did not update name: %q", cfg.Presets[0].Name)
	}
	if cfg.RenamePreset("https://example.com/stream", "   ") {
		t.Fatal("expected blank rename to fail")
	}
	if !cfg.DeletePreset("https://example.com/stream") {
		t.Fatal("expected delete to succeed")
	}
	if len(cfg.Presets) != 0 {
		t.Fatalf("expected no presets after delete, got %d", len(cfg.Presets))
	}
}

func TestMoveByURLBounds(t *testing.T) {
	cfg := Config{MyStations: []Preset{
		{Name: "One", URL: "https://example.com/1"},
		{Name: "Two", URL: "https://example.com/2"},
		{Name: "Three", URL: "https://example.com/3"},
	}}

	if next, ok := cfg.MoveMyStation("https://example.com/2", -1); !ok || next != 0 {
		t.Fatalf("expected move up to index 0, got index=%d ok=%v", next, ok)
	}
	if cfg.MyStations[0].Name != "Two" {
		t.Fatalf("expected Two at top, got %q", cfg.MyStations[0].Name)
	}
	if _, ok := cfg.MoveMyStation("https://example.com/2", -1); ok {
		t.Fatal("expected moving beyond top to fail")
	}
	if _, ok := cfg.MoveMyStation("https://example.com/missing", 1); ok {
		t.Fatal("expected missing URL move to fail")
	}
}

func TestLoadNormalizesConfig(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configRoot)
	dir := filepath.Join(configRoot, "weazltunes")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	raw := `{
  "presets": [
    {"name": "One\u001b[31m", "url": "https://example.com/1"},
    {"name": "Bad", "url": "file:///tmp/song.mp3"},
    {"name": "Dup", "url": "https://example.com/1"},
    {"name": "Two", "url": "https://example.com/2"},
    {"name": "Three", "url": "https://example.com/3"},
    {"name": "Four", "url": "https://example.com/4"},
    {"name": "Five", "url": "https://example.com/5"},
    {"name": "Six", "url": "https://example.com/6"},
    {"name": "Seven", "url": "https://example.com/7"},
    {"name": "Eight", "url": "https://example.com/8"},
    {"name": "Nine", "url": "https://example.com/9"}
  ],
  "my_stations": [
    {"name": "Saved\nStation", "url": "https://example.com/saved"},
    {"name": "Nope", "url": "ftp://example.com/nope"}
  ]
}`
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Presets) != PresetLimit {
		t.Fatalf("expected %d presets, got %d", PresetLimit, len(cfg.Presets))
	}
	if cfg.Presets[0].Name != "One[31m" {
		t.Fatalf("preset name was not sanitized: %q", cfg.Presets[0].Name)
	}
	if cfg.Presets[1].URL != "https://example.com/2" {
		t.Fatalf("invalid or duplicate preset was not removed: %#v", cfg.Presets)
	}
	if len(cfg.MyStations) != 1 || cfg.MyStations[0].Name != "SavedStation" {
		t.Fatalf("my stations not normalized: %#v", cfg.MyStations)
	}
}

func TestSaveUsesPrivateConfigDir(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configRoot)

	if err := Save(Config{Presets: []Preset{{Name: "One", URL: "https://example.com/1"}}}); err != nil {
		t.Fatal(err)
	}
	path, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("config dir mode = %o, want 700", got)
	}
}
