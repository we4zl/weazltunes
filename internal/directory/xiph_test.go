package directory

import "testing"

func TestParseRadioBrowserStations(t *testing.T) {
	body := `[
		{
			"name": "Rock Antenne",
			"url": "http://example.com/listen.pls",
			"url_resolved": "http://example.com/stream.mp3",
			"tags": "rock,classic rock",
			"country": "Germany",
			"codec": "MP3",
			"bitrate": 128,
			"clickcount": 42
		}
	]`

	stations, err := parseRadioBrowserStations(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(stations) != 1 {
		t.Fatalf("len = %d, want 1", len(stations))
	}
	st := stations[0]
	if st.Name != "Rock Antenne" {
		t.Fatalf("name = %q", st.Name)
	}
	if st.URL != "http://example.com/stream.mp3" {
		t.Fatalf("url = %q", st.URL)
	}
	if st.Source != "Radio Browser" {
		t.Fatalf("source = %q", st.Source)
	}
	if st.Listeners != 42 {
		t.Fatalf("listeners = %d", st.Listeners)
	}
	if st.Description != "Germany - MP3 - 128kbps - rock,classic rock" {
		t.Fatalf("description = %q", st.Description)
	}
}

func TestParseRadioBrowserStationsFallsBackToURL(t *testing.T) {
	body := `[{"name":"Station","url":"http://example.com/live","url_resolved":""}]`

	stations, err := parseRadioBrowserStations(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(stations) != 1 {
		t.Fatalf("len = %d, want 1", len(stations))
	}
	if stations[0].URL != "http://example.com/live" {
		t.Fatalf("url = %q", stations[0].URL)
	}
}
