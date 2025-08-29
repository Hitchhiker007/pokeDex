package main

import (
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
	}

	startRepl(cfg)
}
