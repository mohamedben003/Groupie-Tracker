package web

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"grouping_tracker/internal"
	"grouping_tracker/internal/api"
	"grouping_tracker/internal/render"
)

var artists []internal.Artist

func SetArtists(a []internal.Artist) {
	artists = a
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Render404(w)
		return
	}
	if err := render.HTML(w, "index.html", artists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/artist/"):]
	if idStr == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		Render404(w)
		return
	}

	var selected internal.Artist
	found := false
	for _, a := range artists {
		if a.ID == id {
			selected = a
			found = true
			break
		}
	}
	if !found {
		Render404(w)
		return
	}

	detail := internal.ArtistDetail{Artist: selected}

	// Simple goroutines example: 3 requests in parallel.
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := api.FetchJSON(selected.Locations, &detail.LocationsData); err != nil {
			log.Printf("Error fetching locations: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := api.FetchJSON(selected.ConcertDates, &detail.ConcertDatesData); err != nil {
			log.Printf("Error fetching dates: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := api.FetchJSON(selected.Relations, &detail.RelationsData); err != nil {
			log.Printf("Error fetching relations: %v", err)
		}
	}()

	wg.Wait()

	if err := render.HTML(w, "artist.html", detail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Render404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	if err := render.HTML(w, "404.html", nil); err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}


