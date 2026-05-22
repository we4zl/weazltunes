package directory

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

func (c Client) XiphGenres(ctx context.Context) ([]string, error) {
	body, err := c.getString(ctx, "https://dir.xiph.org/genres")
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(?i)<a[^>]+href="/genres/([^"#?]+)"[^>]*>([^<]+)</a>`)
	seen := map[string]bool{}
	var genres []string
	for _, m := range re.FindAllStringSubmatch(body, -1) {
		g := html.UnescapeString(strings.TrimSpace(m[2]))
		if g == "" || seen[strings.ToLower(g)] {
			continue
		}
		seen[strings.ToLower(g)] = true
		genres = append(genres, g)
	}
	sort.Strings(genres)
	return genres, nil
}

func (c Client) XiphSearch(ctx context.Context, query string) ([]Station, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		if stations, err := c.xiphXML(ctx, ""); err != nil || len(stations) > 0 {
			return stations, err
		}
		return c.radioBrowserSearch(ctx, "")
	}
	if stations, err := c.xiphGenrePage(ctx, query); err == nil && len(stations) > 0 {
		return stations, nil
	}
	if stations, err := c.xiphXML(ctx, query); err != nil || len(stations) > 0 {
		return stations, err
	}
	return c.radioBrowserSearch(ctx, query)
}

func (c Client) xiphXML(ctx context.Context, query string) ([]Station, error) {
	body, err := c.getString(ctx, "https://dir.xiph.org/yp.xml")
	if err != nil {
		return nil, err
	}
	var payload struct {
		Entries []struct {
			ServerName        string `xml:"server_name"`
			ListenURL         string `xml:"listen_url"`
			Genre             string `xml:"genre"`
			CurrentSong       string `xml:"current_song"`
			ServerDescription string `xml:"server_description"`
			Listeners         int    `xml:"listeners"`
		} `xml:"entry"`
	}
	if err := xml.Unmarshal([]byte(body), &payload); err != nil {
		return nil, err
	}
	return filterXiphEntries(payload.Entries, query), nil
}

func (c Client) xiphGenrePage(ctx context.Context, genre string) ([]Station, error) {
	body, err := c.getString(ctx, "https://dir.xiph.org/genres/"+url.PathEscape(genre))
	if err != nil {
		return nil, err
	}
	blocks := regexp.MustCompile(`(?is)<h5[^>]*>(.*?)</h5>(.*?)(?:<h5|$)`).FindAllStringSubmatch(body, -1)
	var out []Station
	for _, block := range blocks {
		name := strings.TrimSpace(stripTags(block[1]))
		link := firstHref(block[2])
		if name == "" || link == "" {
			continue
		}
		out = append(out, Station{
			Name:        html.UnescapeString(name),
			Description: strings.TrimSpace(html.UnescapeString(stripTags(block[2]))),
			Genre:       genre,
			URL:         link,
			Source:      "Xiph",
		})
	}
	return out, nil
}

func filterXiphEntries(entries []struct {
	ServerName        string `xml:"server_name"`
	ListenURL         string `xml:"listen_url"`
	Genre             string `xml:"genre"`
	CurrentSong       string `xml:"current_song"`
	ServerDescription string `xml:"server_description"`
	Listeners         int    `xml:"listeners"`
}, query string) []Station {
	q := strings.ToLower(query)
	var out []Station
	for _, e := range entries {
		haystack := strings.ToLower(e.ServerName + " " + e.Genre + " " + e.ServerDescription + " " + e.CurrentSong)
		if e.ListenURL == "" || (q != "" && !strings.Contains(haystack, q)) {
			continue
		}
		out = append(out, Station{Name: e.ServerName, Description: e.ServerDescription, Genre: e.Genre, URL: e.ListenURL, Source: "Xiph", Listeners: e.Listeners, NowPlaying: e.CurrentSong})
		if len(out) >= 200 {
			break
		}
	}
	return out
}

func (c Client) radioBrowserSearch(ctx context.Context, query string) ([]Station, error) {
	params := url.Values{}
	params.Set("hidebroken", "true")
	params.Set("limit", "200")
	params.Set("order", "clickcount")
	params.Set("reverse", "true")
	if query != "" {
		params.Set("tag", query)
	}
	body, err := c.getString(ctx, "https://de1.api.radio-browser.info/json/stations/search?"+params.Encode())
	if err != nil {
		return nil, err
	}
	stations, err := parseRadioBrowserStations(body)
	if err != nil || len(stations) > 0 || query == "" {
		return stations, err
	}
	params.Del("tag")
	params.Set("name", query)
	body, err = c.getString(ctx, "https://de1.api.radio-browser.info/json/stations/search?"+params.Encode())
	if err != nil {
		return nil, err
	}
	return parseRadioBrowserStations(body)
}

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
			Name:        entry.Name,
			Description: radioBrowserDescription(entry.Country, entry.Codec, entry.Bitrate, entry.Tags),
			Genre:       entry.Tags,
			URL:         streamURL,
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
