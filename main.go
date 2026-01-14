package main

import (
	"fmt"
	"grouping_tracker/internal"
	"log"
	"net/http"

	"grouping_tracker/internal/handlers"
	"grouping_tracker/internal/services"
)

var (
	artists []internal.Artist
)

func main() {
	renderer, err := handlers.NewTemplateRenderer("templates/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	svc := services.NewGroupieService(nil)
	artists, err = svc.FetchArtists()
	if err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	h := handlers.New(handlers.Config{
		Renderer: renderer,
		Service:  svc,
		Artists:  artists,
	})
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/artist/", h.Artist)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
