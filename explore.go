package main

import (
	"encoding/json"
	"fmt"
)

type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandExplore(cfg *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: explore <area_name>")
	}

	areaName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", areaName)

	fmt.Printf("Exploring %s...\n", areaName)

	body, err := cfg.pokeapiClient.Get(url)
	if err != nil {
		return err
	}

	// parse JSON
	var data ExploreResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	// print pokemon
	fmt.Println("Found Pokemon:")
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
