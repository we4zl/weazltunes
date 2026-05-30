package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Presets    []Preset `json:"presets"`
	MyStations []Preset `json:"my_stations"`
}

type Preset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

const PresetLimit = 8

func Load() (Config, error) {
	cfg := defaultConfig()
	path, err := Path()
	if err != nil {
		return cfg, err
	}
	if err := ensureConfigDir(filepath.Dir(path)); err != nil {
		return cfg, err
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, Save(cfg)
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	cfg.normalize()
	if len(cfg.Presets) == 0 {
		cfg.Presets = defaultConfig().Presets
	}
	return cfg, nil
}

func (c *Config) SaveMyStation(st Preset) {
	st, ok := cleanPreset(st)
	if !ok {
		return
	}
	c.MyStations = upsert(c.MyStations, st)
}

func (c *Config) PromotePreset(st Preset) {
	st, ok := cleanPreset(st)
	if !ok {
		return
	}
	presets := upsert(c.Presets, st)
	c.Presets = presets
	if len(c.Presets) > PresetLimit {
		c.Presets = c.Presets[:PresetLimit]
	}
}

func (c *Config) ReplacePreset(index int, st Preset) bool {
	st, ok := cleanPreset(st)
	if !ok || index < 0 || index >= PresetLimit || index >= len(c.Presets) {
		return false
	}
	if strings.EqualFold(c.Presets[index].URL, st.URL) {
		c.Presets[index] = st
		return true
	}
	c.Presets = removeByURL(c.Presets, st.URL)
	if index >= len(c.Presets) {
		return false
	}
	c.Presets[index] = st
	return true
}

func (c *Config) RenamePreset(url, name string) bool {
	return renameByURL(c.Presets, url, name)
}

func (c *Config) RenameMyStation(url, name string) bool {
	return renameByURL(c.MyStations, url, name)
}

func (c *Config) DeletePreset(url string) bool {
	next, deleted := deleteByURL(c.Presets, url)
	c.Presets = next
	return deleted
}

func (c *Config) DeleteMyStation(url string) bool {
	next, deleted := deleteByURL(c.MyStations, url)
	c.MyStations = next
	return deleted
}

func (c *Config) MovePreset(url string, delta int) (int, bool) {
	next, moved := moveByURL(c.Presets, url, delta)
	c.Presets = next
	return moved, moved >= 0
}

func (c *Config) MoveMyStation(url string, delta int) (int, bool) {
	next, moved := moveByURL(c.MyStations, url, delta)
	c.MyStations = next
	return moved, moved >= 0
}

func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := ensureConfigDir(filepath.Dir(path)); err != nil {
		return err
	}
	cfg.normalize()
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o600)
}

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "weazltunes", "config.json"), nil
}

func ensureConfigDir(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return os.Chmod(dir, 0o700)
}

func defaultConfig() Config {
	return Config{
		Presets: []Preset{
			{Name: "Groove Salad", URL: "https://somafm.com/groovesalad.pls"},
			{Name: "Drone Zone", URL: "https://somafm.com/dronezone.pls"},
			{Name: "DEF CON Radio", URL: "https://somafm.com/defcon.pls"},
			{Name: "Indie Pop Rocks", URL: "https://somafm.com/indiepop.pls"},
			{Name: "Deep Space One", URL: "https://somafm.com/deepspaceone.pls"},
		},
	}
}
