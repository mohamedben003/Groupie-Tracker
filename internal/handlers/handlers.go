package handlers

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"grouping_tracker/internal"
	"grouping_tracker/internal/services"
)

type Config struct {
	Renderer *TemplateRenderer
	Service  *services.GroupieService
	Artists  []internal.Artist
}

type Handlers struct {
	renderer *TemplateRenderer
	service  *services.GroupieService
	artists  []internal.Artist
}

func New(cfg Config) *Handlers {
	return &Handlers{
		renderer: cfg.Renderer,
		service:  cfg.Service,
		artists:  cfg.Artists,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.render404(w)
		return
	}
	if err := h.renderer.ExecuteTemplate(w, "index.html", h.artists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) Artist(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/artist/"):]
	if idStr == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.render404(w)
		return
	}

	var selected internal.Artist
	found := false
	for _, a := range h.artists {
		if a.ID == id {
			selected = a
			found = true
			break
		}
	}
	if !found {
		h.render404(w)
		return
	}

	detail := internal.ArtistDetail{Artist: selected}

	// Fetch sub-resources concurrently (goroutines) so artist pages load faster.
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		loc, err := h.service.FetchLocations(selected.Locations)
		if err != nil {
			log.Printf("Error fetching locations: %v", err)
			return
		}
		detail.LocationsData = loc
	}()

	go func() {
		defer wg.Done()
		dates, err := h.service.FetchDates(selected.ConcertDates)
		if err != nil {
			log.Printf("Error fetching dates: %v", err)
			return
		}
		detail.ConcertDatesData = dates
	}()

	go func() {
		defer wg.Done()
		rel, err := h.service.FetchRelations(selected.Relations)
		if err != nil {
			log.Printf("Error fetching relations: %v", err)
			return
		}
		detail.RelationsData = rel
	}()

	wg.Wait()

	if err := h.renderer.ExecuteTemplate(w, "artist.html", detail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) render404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	if err := h.renderer.ExecuteTemplate(w, "404.html", nil); err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}


