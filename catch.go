package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/chzyer/readline"
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

		newPokemon := &PokemonInstance{
			Species:  pokemon,
			Nickname: pokemonName,
			Level:    1,
			Hp:       pokemon.Stats[0].BaseStat,
			Boxed:    false,
		}

		placed := false
		for i := 0; i < len(cfg.Party); i++ {
			if cfg.Party[i] == nil {
				fmt.Print("Would you like to give your pokemon a nickname? (y/n): ")

				var response string
				fmt.Scanln(&response)

				if response == "y" || response == "Y" {
					rl, _ := readline.New("")
					defer rl.Close()

					nickname, err := rl.Readline()
					if err == nil && strings.TrimSpace(nickname) != "" {
						newPokemon.Nickname = nickname
					}
				}

				cfg.Party[i] = newPokemon
				placed = true
				fmt.Printf("Gotcha! %s was caught and added to your party!\n", pokemonName)
				// wild Pokemon level is near the players level
				// wildLevel = Player Level + random offset
				// random offset = -1, 0, or 1
				// wildLevel = PlayerLevel + r, where r ∈ {-1, 0, 1}
				wildLevel := cfg.PlayerLV + rand.Intn(3) - 1
				if wildLevel < 1 {
					wildLevel = 1
				}
				// EXP = (BaseExp × Level) / 7
				EXP := (pokemon.BaseExperience * wildLevel)
				fmt.Printf("Player earnt %d xp!\n", EXP)
				cfg.PlayerXP += EXP

				// cehck for lv up occurence
				for cfg.PlayerLV < 50 && cfg.PlayerXP >= XpLevelCheck(cfg.PlayerLV+1) {
					cfg.PlayerLV++
					fmt.Printf("Player Level Up! New Player Level: %d\n", cfg.PlayerLV)
				}
				break
			}
		}

		if !placed {
			cfg.PC = append(cfg.PC, newPokemon)
			newPokemon.Boxed = true
			fmt.Printf("Gotcha! %s was caught, but party is full — sent to the PC!\n", pokemonName)
		}

	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func XpLevelCheck(level int) int {
	if level == 2 {
		return 1 // super fast first level for testing
	}
	return 6 * level * level * level
	// 6 times the xp needed for a pokemon
}
