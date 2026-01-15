package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	api "grouping_tracker/internal/api"
	handler "grouping_tracker/internal/handler"
	types "grouping_tracker/internal/types"
)


func main() {
	var err error
	types.Templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	if err := api.FetchArtists(); err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/artist/", handler.ArtistHandler)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
