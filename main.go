package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// --- Structs ---

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
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

// --- Global Variables ---

var (
	artists   []Artist
	templates *template.Template
)

// --- Main Logic ---

func main() {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	if err := fetchArtists(); err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist/", artistHandler)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchArtists() error {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&artists)
}

// --- Handlers ---

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := templates.ExecuteTemplate(w, "index.html", artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/artist/"):]
	if idStr == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	id,err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > len(artists) {
		http.NotFound(w, r)
		return
	}
	// Safe ID Lookup
	var selectedArtist Artist
	found := false
	for _, a := range artists {
		if a.ID == id {
			selectedArtist = a
			found = true
			break
		}
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	detail := ArtistDetail{Artist: selectedArtist}

	fetchData(selectedArtist.Locations, &detail.LocationsData)
	fetchData(selectedArtist.ConcertDates, &detail.ConcertDatesData)
	fetchData(selectedArtist.Relations, &detail.RelationsData)

	
	if err := templates.ExecuteTemplate(w, "artist.html", detail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper function to reduce repetitive HTTP code
func fetchData(url string, target interface{}) error {
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