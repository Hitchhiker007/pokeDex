package main

import (
	"time"

	"github.com/Hitchhiker007/pokeDex/pokeapi"
)

func main() {
	client := pokeapi.NewClient(5 * time.Second)
	cfg := &Config{
		pokeapiClient: client,
	}

	startRepl(cfg)
}
