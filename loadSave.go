package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func commandLoad(cfg *Config, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("usage: load <filename>")
	}

	filePath := args[0]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		pokemonName := scanner.Text()
		fmt.Println(pokemonName)
		url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

		body, err := cfg.pokeapiClient.Get(url)
		if err != nil {
			return err
		}

		var pokemon Pokemon
		if err := json.Unmarshal(body, &pokemon); err != nil {
			return err
		}

		cfg.Pokedex[pokemonName] = pokemon
		fmt.Printf("Loaded %s into your Pokédex.\n", pokemonName)

		data, err := os.ReadFile("pokemon_instance.txt")
		if err != nil {
			return err
		}

		var instance PokemonInstance
		if err := json.Unmarshal(data, &instance); err != nil {
			return err
		}

		cfg.PC = append(cfg.PC, &instance)
		fmt.Println(instance.Nickname, "Lv", instance.Level)

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading file: %v\n", err)
	}

	return nil

}
