package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	api "grouping_tracker/internal/api"
	handler "grouping_tracker/internal/handler"
	"grouping_tracker/internal/render"
	types "grouping_tracker/internal/types"
)

func main() {
	var err error
	types.Templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	// For the artists's fetch
	if err := api.FetchData("https://groupietrackers.herokuapp.com/api/artists", &types.Artists); err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	// For the locations's fetch
	if err := api.FetchData("https://groupietrackers.herokuapp.com/api/locations", &types.AllLocations); err != nil {
		log.Fatal("Error fetching locations:", err)
	}

	fs := http.FileServer(http.Dir("static"))

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		path := "static" + r.URL.Path[len("/static"):]

		info, err := os.Stat(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			render.Render404(w)
			return
		}

		// Block directory listing
		if info.IsDir() {
			w.WriteHeader(http.StatusNotFound)
			render.Render404(w)
			return
		}

		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/artist/", handler.ArtistHandler)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
