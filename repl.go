package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

// let all commands have a second param args so function signatures are uniform
// allows the REPL to call any command dynamically, even if some commands like help or exit dont use them
type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

// replaced bufio.Scanner with readline for a better user experience
// up and down arrow keys for navigate previous commands
// command history stored in /tmp/pokedex_history.tmp
// credit to https://github.com/chzyer/readline
func startRepl(cfg *Config) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "Pokedex > ",
		HistoryFile:     "/tmp/pokedex_history.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start REPL: %v\n", err)
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			} else {
				break
			}
		}

		words := cleanInput(line)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		command, exists := getCommands()[commandName]

		if exists {
			// now parse 2nd arg
			err := command.callback(cfg, words[1:])
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}

}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display 20 locations at a time",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Diplay 20 previous locations",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "explore the current area for pokemon",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch a specified pokemon!",
			callback:    catch,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Display all Pokémon you have caught",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect",
			description: "It takes the name of a Pokemon and prints the name, height, weight, stats and type(s) of the Pokemon.",
			callback:    commandInspect,
		},
		"load": {
			name: "load",
			description: "displays a list of save files to load from",
			callback: commandLoad,
		},
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit(cfg *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
