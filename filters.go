package main

import (
	"net/http"
	"strconv"
	"strings"
)

// FilterData holds the values to send back to the HTML inputs
// so they don't reset when the page reloads.
type FilterData struct {
	CreationDateMin int
	CreationDateMax int
	FirstAlbumMin   int
	FirstAlbumMax   int
	Members         []string
	Location        string
}

// FilterArtists filters the list based on the form data
func FilterArtists(artists []Artist, r *http.Request) ([]Artist, FilterData) {
	var filtered []Artist

	// 1. Parse Form Data
	if err := r.ParseForm(); err != nil {
		return artists, FilterData{} // Return original list on error
	}

	// 2. Get Form Values
	minC, _ := strconv.Atoi(r.FormValue("creationDateMin"))
	maxC, _ := strconv.Atoi(r.FormValue("creationDateMax"))
	minA, _ := strconv.Atoi(r.FormValue("firstAlbumMin"))
	maxA, _ := strconv.Atoi(r.FormValue("firstAlbumMax"))
	location := strings.ToLower(r.FormValue("location"))
	members := r.Form["members"] // Checkboxes return a slice

	// Set Defaults if values are empty (0)
	if minC == 0 {
		minC = 1950
	}
	if maxC == 0 {
		maxC = 2025
	}
	if minA == 0 {
		minA = 1950
	}
	if maxA == 0 {
		maxA = 2025
	}

	// Prepare data to send back to UI
	filterData := FilterData{
		CreationDateMin: minC,
		CreationDateMax: maxC,
		FirstAlbumMin:   minA,
		FirstAlbumMax:   maxA,
		Members:         members,
		Location:        location,
	}

	// 3. Loop and Check
	for _, a := range artists {

		// --- Check 1: Creation Date ---
		if a.CreationDate < minC || a.CreationDate > maxC {
			continue
		}

		// --- Check 2: First Album Date ---
		// API format is "DD-MM-YYYY". We want the Year (last 4 chars)
		parts := strings.Split(a.FirstAlbum, "-")
		if len(parts) == 3 {
			year, _ := strconv.Atoi(parts[2])
			if year < minA || year > maxA {
				continue
			}
		}

		// --- Check 3: Number of Members ---
		if len(members) > 0 {
			matches := false
			currentMembers := strconv.Itoa(len(a.Members))
			for _, m := range members {
				if m == currentMembers {
					matches = true
					break
				}
			}
			if !matches {
				continue
			}
		}

		// --- Check 4: Location ---
		// NOTE: In the API, "Locations" is just a URL string.
		// You cannot filter by City name using just the basic Artist struct
		// unless you fetch all location data on startup.
		// For now, this is ignored or requires advanced data fetching.

		filtered = append(filtered, a)
	}

	return filtered, filterData
}
