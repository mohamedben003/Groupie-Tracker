package handlers

import (
	"fmt"
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
		log.Printf("Template execution error: %v", err)
		h.render500(w)
		return
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

	// Safe ID lookup - find artist by ID
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
			// Continue anyway - partial data is better than nothing
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
		log.Printf("Template execution error: %v", err)
		h.render500(w)
		return
	}
}

// renderError renders a generic error page with custom code, title, and message
func (h *Handlers) renderError(w http.ResponseWriter, code int, title, message string) {
	w.WriteHeader(code)
	data := internal.ErrorPageData{
		Code:    code,
		Title:   title,
		Message: message,
	}
	err := h.renderer.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		// Fallback if template rendering fails
		http.Error(w, fmt.Sprintf("%d %s", code, title), code)
	}
}

func (h *Handlers) render404(w http.ResponseWriter) {
	h.renderError(w, http.StatusNotFound, "Page Not Found", "Got lost? It seems that the page you're looking for doesn't exist.")
}

func (h *Handlers) render500(w http.ResponseWriter) {
	h.renderError(w, http.StatusInternalServerError, "Internal Server Error", "Something went wrong on our end. Please try again later.")
}
