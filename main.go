package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Artist represents the main artist data
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

// ArtistDetail contains all artist information for detail view
type ArtistDetail struct {
	Artist
	LocationsData   *Locations
	ConcertDatesData *Dates
	RelationsData   *Relations
}

// Locations represents locations data
type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// Dates represents concert dates data
type Dates struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// Relations represents relations data
type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var artists []Artist
var templates *template.Template

func main() {
	// Load templates
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	// Fetch artists data
	err = fetchArtists()
	if err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &artists)
	if err != nil {
		return err
	}

	return nil
}

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

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > len(artists) {
		http.NotFound(w, r)
		return
	}

	// Get the artist
	artist := artists[id-1]

	// Create ArtistDetail and fetch additional data
	detail := ArtistDetail{Artist: artist}

	// Fetch locations
	locResp, err := http.Get(artist.Locations)
	if err == nil {
		defer locResp.Body.Close()
		body, _ := ioutil.ReadAll(locResp.Body)
		var locations Locations
		json.Unmarshal(body, &locations)
		detail.LocationsData = &locations
	}

	// Fetch concert dates
	dateResp, err := http.Get(artist.ConcertDates)
	if err == nil {
		defer dateResp.Body.Close()
		body, _ := ioutil.ReadAll(dateResp.Body)
		var dates Dates
		json.Unmarshal(body, &dates)
		detail.ConcertDatesData = &dates
	}

	// Fetch relations
	relResp, err := http.Get(artist.Relations)
	if err == nil {
		defer relResp.Body.Close()
		body, _ := ioutil.ReadAll(relResp.Body)
		var relations Relations
		json.Unmarshal(body, &relations)
		detail.RelationsData = &relations
	}

	err = templates.ExecuteTemplate(w, "artist.html", detail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}