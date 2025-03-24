package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type config struct {
	next     string
	previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMapb(cfg *config) error {
	var next string

	if cfg.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		next = cfg.previous
	}

	res, err := http.Get(next)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

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

func commandMap(cfg *config) error {
	var next string

	if cfg.next == "" {
		next = "https://pokeapi.co/api/v2/location-area/"
	} else {
		next = cfg.next
	}

	res, err := http.Get(next)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

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
	}

	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		cmd, ok := commands[input[0]]
		fmt.Println(cmd.name)

		if !ok {
			fmt.Printf("Command \"%s\" not found.\n", input[0])
		} else {
			err := cmd.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
