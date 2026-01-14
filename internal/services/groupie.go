package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"grouping_tracker/internal"
)

type GroupieService struct {
	client  *http.Client
	baseURL string
}

func NewGroupieService(client *http.Client) *GroupieService {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &GroupieService{
		client:  client,
		baseURL: "https://groupietrackers.herokuapp.com/api",
	}
}

func (s *GroupieService) FetchArtists() ([]internal.Artist, error) {
	var artists []internal.Artist
	if err := s.fetchJSON(s.baseURL+"/artists", &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func (s *GroupieService) FetchLocations(url string) (*internal.Locations, error) {
	var loc *internal.Locations
	if err := s.fetchJSON(url, &loc); err != nil {
		return nil, err
	}
	return loc, nil
}

func (s *GroupieService) FetchDates(url string) (*internal.Dates, error) {
	var dates *internal.Dates
	if err := s.fetchJSON(url, &dates); err != nil {
		return nil, err
	}
	return dates, nil
}

func (s *GroupieService) FetchRelations(url string) (*internal.Relations, error) {
	var rel *internal.Relations
	if err := s.fetchJSON(url, &rel); err != nil {
		return nil, err
	}
	return rel, nil
}

func (s *GroupieService) fetchJSON(url string, target any) error {
	resp, err := s.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}


