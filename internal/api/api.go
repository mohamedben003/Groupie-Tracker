package api

import (
	"encoding/json"
	"fmt"
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
