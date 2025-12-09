package main

import (
	"fmt"
)

func commandPokedex(cfg *Config, args []string) error {
	if cfg == nil {
		return fmt.Errorf("error! config is nil")
	}

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

func commandViewParty(cfg *Config, args []string) error {
	if cfg == nil {
		return fmt.Errorf("error! config is nil")
	}

	if len(cfg.Party) == 0 {
		fmt.Println("You currently have no pokemon in your party!")
		return nil
	}

	fmt.Println("Your Party:")
	empty := true
	for i, pokemon := range cfg.Party {
		if pokemon != nil {
			fmt.Printf(" - %d. %s\n", i+1, pokemon.Nickname)
			empty = false
		} else {
			fmt.Printf(" %d. [Empty Slot]\n", i+1)
		}
	}
	if empty {
		fmt.Println("You currently have no pokemon in your party!")
	}
	return nil
}
