package filters

import (
	"net/http"
	"strconv"
	"strings"

	"grouping_tracker/internal/types"
)

// FilterArtists filters the list based on the form data
func FilterArtists(artists []types.Artist, r *http.Request) ([]types.Artist, types.FilterData) {
	var filtered []types.Artist

	if r.URL.RawQuery == "" {
		return artists, types.FilterData{}
	}

	// 1. Parse Form Data
	if err := r.ParseForm(); err != nil {
		// On error, return empty filter data and original list
		return artists, types.FilterData{}
	}

	// 2. Get Form Values
	minC, _ := strconv.Atoi(r.FormValue("creationDateMin"))
	maxC, _ := strconv.Atoi(r.FormValue("creationDateMax"))
	minA, _ := strconv.Atoi(r.FormValue("firstAlbumMin"))
	maxA, _ := strconv.Atoi(r.FormValue("firstAlbumMax"))
	members := r.Form["members"]
	location := strings.TrimSpace(r.FormValue("location"))

	//to send back the original location typed in the search
	locationLower := strings.ToLower(location)

	// Set Defaults
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
	filterData := types.FilterData{
		CreationDateMin: minC,
		CreationDateMax: maxC,
		FirstAlbumMin:   minA,
		FirstAlbumMax:   maxA,
		Members:         members,
		Location:        location,
	}

	// 3. Loop and Check
	for _, a := range artists {

		// Check 1: Creation Date
		if a.CreationDate < minC || a.CreationDate > maxC {
			continue
		}

		// Check 2: First Album Date
		parts := strings.Split(a.FirstAlbum, "-")
		if len(parts) == 3 {
			year, _ := strconv.Atoi(parts[2])
			if year < minA || year > maxA {
				continue
			}
		}

		// Check 3: Number of Members
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

		// Check 4: Locations
		if location != "" {
			matchesLoc := false

			for _, locData := range types.AllLocations.Index {
				if locData.ID == a.ID {
					for _, city := range locData.Locations {

						cleanCity := strings.ReplaceAll(city, "_", " ")
						cleanCity = strings.ReplaceAll(cleanCity, "-", " ")
						cleanCity = strings.ToLower(cleanCity)

						if strings.Contains(cleanCity, locationLower) {
							matchesLoc = true
							break
						}
					}
					break
				}
			}
			if !matchesLoc {
				continue
			}
		}
		filtered = append(filtered, a)
	}

	return filtered, filterData
}
