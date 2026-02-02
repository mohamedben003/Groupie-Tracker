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

type IndexPageData struct {
	Artists []Artist
	Filters FilterData
}

// --- Err pages ---

type ErrorPageData struct {
	Code    int
	Title   string
	Message string
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
		render404(w)
		return
	}

	// Call the filter logic from filters.go
	filteredArtists, currentFilters := FilterArtists(artists, r)

	// Prepare the data for the template
	data := IndexPageData{
		Artists: filteredArtists,
		Filters: currentFilters,
	}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		render500(w)
		return
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/artist/"):]
	if idStr == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > len(artists) {
		render404(w)
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
		render404(w)
		return
	}

	detail := ArtistDetail{Artist: selectedArtist}

	if err := fetchData(selectedArtist.Locations, &detail.LocationsData); err != nil {
		log.Printf("Error fetching locations: %v", err)
		// Continue anyway - partial data is better than nothing
	}

	if err := fetchData(selectedArtist.ConcertDates, &detail.ConcertDatesData); err != nil {
		log.Printf("Error fetching dates: %v", err)
	}

	if err := fetchData(selectedArtist.Relations, &detail.RelationsData); err != nil {
		log.Printf("Error fetching relations: %v", err)
	}

	if err := templates.ExecuteTemplate(w, "artist.html", detail); err != nil {
		log.Printf("Template execution error: %v", err)
		render500(w)
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

func renderError(w http.ResponseWriter, code int, title, message string) {
	w.WriteHeader(code)
	data := ErrorPageData{
		Code:    code,
		Title:   title,
		Message: message,
	}
	err := templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%d %s", code, title), code)
	}
}

func render404(w http.ResponseWriter) {
	renderError(w, http.StatusNotFound, "Page Not Found", "Got lost? It seems that the page you're looking for doesn't exist.")
}

func render500(w http.ResponseWriter) {
	renderError(w, http.StatusInternalServerError, "Internal Server Error", "Something went wrong on our end. Please try again later.")
}