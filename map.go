package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	pokeapi "github.com/Hitchhiker007/pokeDex/pokeapi"
	"github.com/google/uuid"
)

type Config struct {
	pokeapiClient   *pokeapi.Client `json:"-"` // skip this field when saving
	MapNextURL      string          // URL for next page of locations
	MapPrevURL      string          // URL for previous page (optional)
	Pokedex         map[string]Pokemon
	Party           [6]*PokemonInstance // fixed 6-slot party
	PC              []*PokemonInstance  // PC storage
	PlayerXP        int
	PlayerLV        int
	LastCloudSaveID string
	SaveDir         string // defaults to ~/.pokedex, overridable for tests
	Token           *TokenResponse
}

// LocationList struct
type LocationList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// SaveFile = everything needed to RESTORE the game
type SaveFile struct {
	SaveID     string
	BaseSaveID string
	Timestamp  time.Time
	Pokedex    map[string]Pokemon
	Party      [6]*PokemonInstance
	PC         []*PokemonInstance
	PlayerXP   int
	PlayerLV   int
}

// Generic Unmarshal function
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func fetchLocations(cfg *Config, url string) error {

	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}

	// Use pokeapi client to fetch
	body, err := cfg.pokeapiClient.Get(url)
	if err != nil {
		return err
	}

	// unmarshel the JSON
	var data LocationList
	if err := Unmarshal(body, &data); err != nil {
		return err
	}

	// Print locations
	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	// save previous and next urls in the config for paginaiton
	cfg.MapNextURL = data.Next
	cfg.MapPrevURL = data.Previous

	return nil

}

func commandMap(cfg *Config, args []string) error {
	url := cfg.MapNextURL
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}
	return fetchLocations(cfg, url)
}

func commandMapBack(cfg *Config, args []string) error {
	if cfg.MapPrevURL == "" {
		fmt.Println("You are already at the first page!")
		return nil
	}
	return fetchLocations(cfg, cfg.MapPrevURL)
}

func saveGameState(cfg *Config, args []string) error {
	saveFile := SaveFile{
		SaveID:     uuid.New().String(),
		BaseSaveID: cfg.LastCloudSaveID,
		Timestamp:  time.Now().UTC(),
		Pokedex:    cfg.Pokedex,
		Party:      cfg.Party,
		PC:         cfg.PC,
		PlayerXP:   cfg.PlayerXP,
		PlayerLV:   cfg.PlayerLV,
	}
	saveData, err := json.Marshal(saveFile)
	if err != nil {
		return fmt.Errorf("failed to marshal save file: %w", err)
	}

	var dirPath string

	if cfg.SaveDir != "" {
		dirPath = cfg.SaveDir
	} else {
		dirPath, err = getPokedexDir()
		if err != nil {
			return err
		}
	}

	// create the folder if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	// build the full file path ~/.pokedex/save.json
	savePath := filepath.Join(dirPath, "save.json")

	if err := os.WriteFile(savePath, saveData, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}
	fmt.Println("Game progress saved!")
	return nil
}

func loadGameState(cfg *Config, args []string) error {
	var dirPath string

	if cfg.SaveDir != "" {
		dirPath = cfg.SaveDir // use the override (e.g. from tests)
	} else {
		var err error
		dirPath, err = getPokedexDir()
		if err != nil {
			return err
		}
	}
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	savePath := filepath.Join(dirPath, "save.json")

	data, err := os.ReadFile(savePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no save file found. Start playing to create one!")
			return nil
		}
		return fmt.Errorf("failed to read save file: %w", err)
	}

	var saveFile SaveFile

	if err := Unmarshal(data, &saveFile); err != nil {
		return fmt.Errorf("failed to unmarshal save file: %w", err)
	}

	cfg.Pokedex = saveFile.Pokedex
	cfg.Party = saveFile.Party
	cfg.PC = saveFile.PC
	cfg.PlayerXP = saveFile.PlayerXP
	cfg.PlayerLV = saveFile.PlayerLV
	cfg.LastCloudSaveID = saveFile.SaveID

	fmt.Println("Successfully loaded save file!")
	fmt.Printf("Player Level: %d\n", cfg.PlayerLV)
	fmt.Printf("Player XP: %d\n", cfg.PlayerXP)

	return nil
}
