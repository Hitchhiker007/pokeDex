package main

import (
	"errors"
	"strings"
)

func movePokemon(cfg *Config, args []string) error {
	if len(args) < 2 {
		return errors.New("usage: move <pokemon> <pc|party>")
	}

	name := strings.ToLower(args[0])
	direction := strings.ToLower(args[1])

	if direction == "party" {

		var pokemon *PokemonInstance
		var pcIndex int
		found := false

		for i, p := range cfg.PC {
			if strings.ToLower(p.Nickname) == name {
				pokemon = p
				pcIndex = i
				found = true
				break
			}
		}

		if !found {
			return errors.New("Pokemon not found in PC")
		}

		slot := -1
		for i, p := range cfg.Party {
			if p == nil {
				slot = i
				break
			}
		}

		if slot == -1 {
			return errors.New("Party is full")
		}

		// remove from pc and make the slot nil like the game
		cfg.PC[pcIndex] = nil

		// add to the party
		cfg.Party[slot] = pokemon

	} else if direction == "pc" {

	}

	return nil
}
