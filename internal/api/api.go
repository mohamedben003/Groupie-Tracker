package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"grouping_tracker/internal"
)

func FetchArtists() ([]internal.Artist, error) {
	var artists []internal.Artist
	if err := FetchJSON("https://groupietrackers.herokuapp.com/api/artists", &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func FetchJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}


