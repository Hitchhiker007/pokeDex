package main

import (
	"fmt"
	"time"

	pokeapi "github.com/Hitchhiker007/pokeDex/pokeapi"
)

func main() {
	// create the client once
	client := pokeapi.NewClient(5 * time.Second)

	cfg := &Config{
		pokeapiClient: client, // assign to the named field
		MapNextURL:    "",
		MapPrevURL:    "",
		Pokedex:       make(map[string]Pokemon),
		Party:         [6]*PokemonInstance{},
		PC:            []*PokemonInstance{},
		PlayerXP:      0,
		PlayerLV:      0,
		// Token is nil by default we don't need to set it here
	}

	if err := loadToken(cfg); err != nil {
		fmt.Println("Warning: could not load token:", err)
	}

	startRepl(cfg)
}
