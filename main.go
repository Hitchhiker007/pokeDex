package main

import (
	"time"

	pokeapi "github.com/Hitchhiker007/pokeDex/internal"
)

func main() {
	client := pokeapi.NewClient(5 * time.Second)
	cfg := &Config{
		pokeapiClient: client,
	}

	startRepl(cfg)
}
