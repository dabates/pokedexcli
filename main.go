package main

import (
	"bufio"
	"encoding/json"
	"example.com/dabates/pokedexcli/internal/pokecache"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const locationURL = "https://pokeapi.co/api/v2/location-area/"
const pokemonURL = "https://pokeapi.co/api/v2/pokemon/"

type config struct {
	next     string
	previous string
	cache    *pokecache.Cache
	pokemon  map[string]pokemonResponse
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

type locationResponse struct {
	id        string
	name      string
	gameIndex string

	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
type pokemonResponse struct {
	Id             int    `json:"id"`
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
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
	fmt.Println("Count of parameters:", len(params))

	for _, param := range params {
		url := locationURL + param

		body, err := cfg.cache.Get(param)
		if err != nil {
			res, err := http.Get(url)
			if err != nil {
				return err
			}

			body, err = io.ReadAll(res.Body)

			if err != nil {
				return err
			}
			res.Body.Close()

			cfg.cache.Add(param, body)
		}

		data := locationResponse{}
		err = json.Unmarshal(body, &data)

		fmt.Println("Found Pokemon:")
		for _, encounter := range data.Encounters {
			fmt.Println("-", encounter.Pokemon.Name)
		}
	}
	return nil
}

func commandCatch(cfg *config, params []string) error {
	fmt.Println("Throwing a Pokeball at", params[0]+"...")

	body, err := cfg.cache.Get(params[0])
	if err != nil {
		url := pokemonURL + params[0]
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		res.Body.Close()
		cfg.cache.Add(params[0], body)
	}

	data := pokemonResponse{}
	err = json.Unmarshal(body, &data)

	// do the calc here
	chance := 1.0 / (1.0 + float64(data.BaseExperience)/50.0)
	if rand.Float64() < chance {
		fmt.Println("Pokemon:", params[0], "was caught!")
		cfg.pokemon[params[0]] = data
	} else {
		fmt.Println("Pokemon:", params[0], "escaped!")
	}

	return nil
}

func commandInspect(cfg *config, params []string) error {
	data, err := cfg.pokemon[params[0]]

	if !err {
		fmt.Println("Pokemon:", params[0], "was not found!")
		return nil
	}

	fmt.Println("Name:", data.Name)
	fmt.Println("Height:", data.Height)
	fmt.Println("Weight:", data.Weight)

	fmt.Println("Stats:")
	for _, stat := range data.Stats {
		fmt.Println("  -" + stat.Stat.Name + ": " + strconv.Itoa(stat.BaseStat))
	}

	fmt.Println("Types:")
	for _, typ := range data.Types {
		fmt.Println("  -", typ.Type.Name)
	}
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
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon",
			callback:    commandInspect,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}
	cfg.cache = pokecache.NewCache(30 * time.Second)
	cfg.pokemon = map[string]pokemonResponse{}

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		cmd, ok := commands[input[0]]
		params := input[1:]

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
