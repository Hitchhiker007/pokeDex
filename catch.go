package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

func catch(cfg *Config, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("usage: catch <pokemon_name>")
	}

	pokemonName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	body, err := cfg.pokeapiClient.Get(url)
	if err != nil {
		return err
	}

	// parse JSON
	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return err
	}

	chance := 100 - pokemon.BaseExperience
	if chance < 10 {
		chance = 10 // flat chance of 10%
	}

	roll := rand.Intn(100) + 1

	if roll <= chance {
		cfg.Pokedex[pokemonName] = pokemon
		fmt.Printf("Gotcha! %s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}
