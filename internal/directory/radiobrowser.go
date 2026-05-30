package directory

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bprendie/weazltunes/internal/security"
)

func parseRadioBrowserStations(body string) ([]Station, error) {
	var payload []struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		URLResolved string `json:"url_resolved"`
		Tags        string `json:"tags"`
		Country     string `json:"country"`
		Codec       string `json:"codec"`
		Bitrate     int    `json:"bitrate"`
		ClickCount  int    `json:"clickcount"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, err
	}
	out := make([]Station, 0, len(payload))
	for _, entry := range payload {
		streamURL := strings.TrimSpace(entry.URLResolved)
		if streamURL == "" {
			streamURL = strings.TrimSpace(entry.URL)
		}
		if entry.Name == "" || streamURL == "" {
			continue
		}
		out = append(out, Station{
			Name:        security.SanitizeStationName(entry.Name),
			Description: security.SanitizeDescription(radioBrowserDescription(entry.Country, entry.Codec, entry.Bitrate, entry.Tags)),
			Genre:       security.SanitizeGenre(entry.Tags),
			URL:         security.SanitizeURL(streamURL),
			Source:      "Radio Browser",
			Listeners:   entry.ClickCount,
		})
	}
	return out, nil
}

func radioBrowserDescription(country, codec string, bitrate int, tags string) string {
	var parts []string
	if country != "" {
		parts = append(parts, country)
	}
	if codec != "" {
		parts = append(parts, codec)
	}
	if bitrate > 0 {
		parts = append(parts, fmt.Sprintf("%dkbps", bitrate))
	}
	if tags != "" {
		parts = append(parts, tags)
	}
	return strings.Join(parts, " - ")
}
