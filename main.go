package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type config struct {
	NextURL     *string
	PreviousURL *string
}

type LocationArea struct {
	Name string `json:"name"`
}

type PokeAPIResponse struct {
	Results  []LocationArea `json:"results"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
}

const baseURL = "https://pokeapi.co/api/v2/location-area/"

func fetchLocationAreas(url string) (*PokeAPIResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PokeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func commandMap(cfg *config) error {
	url := baseURL
	if cfg.NextURL != nil {
		url = *cfg.NextURL
	}

	response, err := fetchLocationAreas(url)
	if err != nil {
		return err
	}

	for _, location := range response.Results {
		fmt.Println(location.Name)
	}

	cfg.NextURL = response.Next
	cfg.PreviousURL = response.Previous
	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.PreviousURL == nil {
		fmt.Println("No previous page available.")
		return nil
	}

	response, err := fetchLocationAreas(*cfg.PreviousURL)
	if err != nil {
		return err
	}

	for _, location := range response.Results {
		fmt.Println(location.Name)
	}

	cfg.NextURL = response.Next
	cfg.PreviousURL = response.Previous
	return nil
}

func commandHelp() error {
	fmt.Print("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("\nhelp: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Display the next 20 location areas")
	fmt.Println("mapb: Display the previous 20 location areas\n")
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func initializeCommands(cfg *config) map[string]cliCommand {
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
			description: "Display the next 20 location areas",
			callback: func() error {
				return commandMap(cfg)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback: func() error {
				return commandMapBack(cfg)
			},
		},
	}
}

func main() {
	cfg := &config{}
	commands := initializeCommands(cfg)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		// Process input
		if command, exists := commands[input]; exists {
			if err := command.callback(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
