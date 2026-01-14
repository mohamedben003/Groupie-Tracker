package main

import (
	"fmt"
	"log"
	"net/http"

	"grouping_tracker/internal/api"
	"grouping_tracker/internal/render"
	"grouping_tracker/internal/web"
)

func main() {
	if err := render.Init("templates/*.html"); err != nil {
		log.Fatal("Error loading templates:", err)
	}

	artists, err := api.FetchArtists()
	if err != nil {
		log.Fatal("Error fetching artists:", err)
	}
	web.SetArtists(artists)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", web.HomeHandler)
	http.HandleFunc("/artist/", web.ArtistHandler)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
