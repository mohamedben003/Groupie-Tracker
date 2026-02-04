package api

import (
	"encoding/json"
	"fmt"
	helper "grouping_tracker/internal/helper"
	types "grouping_tracker/internal/types"

	"net/http"
)

func FetchData(url string, target interface{}) error {


	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func FetchStartData() error {
	FetchData("https://groupietrackers.herokuapp.com/api/artists", &types.Artists);
fmt.Println("artists length:", len(types.Artists))

if err:=FetchData("https://groupietrackers.herokuapp.com/api/locations", &types.DataLocations) ;err!=nil{
fmt.Println("error here")
	return err}


 err := helper.ParseArtiste(); 

	if err != nil {
		return fmt.Errorf("bad status: ")
	}

	fmt.Println("artiste in the end", types.Artists[0])

return nil
}
