package handler

import (
	"log"
	"net/http"
	"strconv"

	api "grouping_tracker/internal/api"
	render "grouping_tracker/internal/render"
	types "grouping_tracker/internal/types"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		render.Render404(w)
		return
	}
	err := types.Templates.ExecuteTemplate(w, "index.html", types.Artists)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		render.Render500(w)
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/artist/"):]
	if idStr == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	id, err := strconv.Atoi(idStr)

	if err != nil || id < 1 || id > len(types.Artists) {
		render.Render404(w)
		return
	}
	
	var selectedArtist types.Artist
	found := false
	for _, a := range types.Artists {
		if a.ID == id {
			selectedArtist = a
			found = true
			break
		}
	}

	if !found {
		render.Render404(w)
		return
	}

	detail := types.ArtistDetail{Artist: selectedArtist}

	if err := api.FetchData(selectedArtist.Locations, &detail.LocationsData); err != nil {
		log.Printf("Error fetching locations: %v", err)
		// Continue anyway - partial data is better than nothing
	}

	if err := api.FetchData(selectedArtist.ConcertDates, &detail.ConcertDatesData); err != nil {
		log.Printf("Error fetching dates: %v", err)
	}

	if err := api.FetchData(selectedArtist.Relations, &detail.RelationsData); err != nil {
		log.Printf("Error fetching relations: %v", err)
	}

	if err := types.Templates.ExecuteTemplate(w, "artist.html", detail); err != nil {
		log.Printf("Template execution error: %v", err)
		render.Render500(w)
	}
}
