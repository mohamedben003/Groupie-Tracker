package internal

import "html/template"

var (
	Artists   []Artist
	Templates *template.Template
	DataLocations LocationsResponse
)

type LocationsResponse struct {
	Index []Locations `json:"index"`
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	LocationsData []string 
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}




type ArtistDetail struct {
	Artist
	LocationsData    *Locations
	ConcertDatesData *Dates
	RelationsData    *Relations
}

type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// --- Err pages ---

type ErrorPageData struct {
	Code    int
	Title   string
	Message string
}
