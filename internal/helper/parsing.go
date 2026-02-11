package helper

import (
	"strings"
)

func CleanCityName(city string) string {
	cleanCity := strings.ToLower(city)
	cleanCity = strings.ReplaceAll(city, "_", " ")
	cleanCity = strings.ReplaceAll(cleanCity, "-", " ")
	cleanCity = strings.ToLower(cleanCity)
	return cleanCity
}

func CleanEntry(entry string) string {
	entry = strings.ToLower(entry)

	entry = strings.ReplaceAll(entry, ", ", " ")

	entry = strings.TrimSpace(entry)
	entry = strings.Join(strings.Fields(entry), " ")

	return entry
}

func CheckLocation(cityName,LocationToFind string) bool {

	return cityName == LocationToFind || strings.Contains(cityName, LocationToFind)
}