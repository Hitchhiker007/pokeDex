package main

import (
	"fmt"
)

func commandInspect(cfg *Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: inspect <pokemon_name>")
	}

	if len(cfg.Pokedex) == 0 {
		fmt.Println("You have not caught any pokemon yet!")
		return nil
	}

	name := args[0]

	pokemon, found := cfg.Pokedex[name]
	if !found {
		fmt.Printf("You have not caught %s yet!\n", name)
	}

	if err := displaySprite(pokemon.Sprites.FrontDefault); err != nil {
		fmt.Println("sprite error:", err)
	}
	fmt.Println()

	// print detailed stats
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Base Experience: %d\n", pokemon.BaseExperience)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Abilities:")
	for _, ability := range pokemon.Abilities {
		fmt.Println(" -", ability.Ability.Name)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Println(" -", t.Type.Name)
	}
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	return nil
}
