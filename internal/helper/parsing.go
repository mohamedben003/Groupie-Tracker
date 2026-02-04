package helper

import (
	"fmt"
	types "grouping_tracker/internal/types"
)
func ParseArtiste() error {
		fmt.Println("artist befaure parsing ", types.Artists[0])

	
		dataLocation:=types.DataLocations.Index
	for i := range types.Artists {
		println("artist id ", types.Artists[i].ID)


		types.Artists[i].LocationsData =
			dataLocation[types.Artists[i].ID-1].Locations
	}

return nil
}