package main

import (
	"encoding/json"
	"fmt"

	"github.com/Hitchhiker007/pokeDex/pokeapi"
)

type Config struct {
	pokeapiClient pokeapi.Client
	MapNextURL    string // URL for next page of locations
	MapPrevURL    string // URL for previous page (optional)
}

// LocationList struct
type LocationList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// Generic Unmarshal function
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func fetchLocations(cfg *Config, url string) error {

	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}

	// Use pokeapi client to fetch
	body, err := cfg.pokeapiClient.Get(url)
	if err != nil {
		return err
	}

	// unmarshel the JSON
	var data LocationList
	if err := Unmarshal(body, &data); err != nil {
		return err
	}

	// Print locations
	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	// save previous and next urls in the config for paginaiton
	cfg.MapNextURL = data.Next
	cfg.MapPrevURL = data.Previous

	return nil

}

func commandMap(cfg *Config) error {
	url := cfg.MapNextURL
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}
	return fetchLocations(cfg, url)
}

func commandMapBack(cfg *Config) error {
	if cfg.MapPrevURL == "" {
		fmt.Println("You are already at the first page!")
		return nil
	}
	return fetchLocations(cfg, cfg.MapPrevURL)
}
