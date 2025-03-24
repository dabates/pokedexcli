package main

import (
	"bufio"
	"encoding/json"
	"example.com/dabates/pokedexcli/internal/pokecache"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const locationURL = "https://pokeapi.co/api/v2/location-area/"

type config struct {
	next     string
	previous string
	cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, params []string) error
}

type locationAreaResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func cleanInput(text string) []string {
	stringVals := strings.Fields(strings.ToLower(text))
	return stringVals
}

func commandExit(cfg *config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, params []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMapb(cfg *config, params []string) error {
	var next string

	if cfg.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		next = cfg.previous
	}

	body, err := cfg.cache.Get(next)

	fmt.Println(err)

	if err != nil {
		res, err := http.Get(next)
		if err != nil {
			return err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		cfg.cache.Add(next, body)
	}

	data := locationAreaResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	cfg.next = data.Next
	cfg.previous = data.Previous

	return nil
}

func commandMap(cfg *config, params []string) error {
	var next string

	if cfg.next == "" {
		next = locationURL
	} else {
		next = cfg.next
	}

	body, err := cfg.cache.Get(next)

	if err != nil {
		res, err := http.Get(next)
		if err != nil {
			return err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		cfg.cache.Add(next, body)
	}

	data := locationAreaResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	cfg.next = data.Next
	cfg.previous = data.Previous

	return nil
}

func commandExplore(cfg *config, params []string) error {
	return nil
}

func main() {
	commands := map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Map back one page",
			callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Shows this help message",
			callback:    commandHelp,
		},
		"explore": {
			name:        "explore",
			description: "Explore the Pokedex",
			callback:    commandExplore,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}
	cfg.cache = pokecache.NewCache(30 * time.Second)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		cmd, ok := commands[input[0]]
		params := input[:1]

		fmt.Println(cmd.name)

		if !ok {
			fmt.Printf("Command \"%s\" not found.\n", input[0])
		} else {
			err := cmd.callback(cfg, params)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
