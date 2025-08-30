package main

import (
	"fmt"
)

func commandPokedex(cfg *Config, args []string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("You haven't caught any Pokemon yet!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.Pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}
