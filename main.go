package main

import (
	"bufio"
	"example.com/dabates/pokedexcli/commands"
	"example.com/dabates/pokedexcli/internal/pokecache"
	. "example.com/dabates/pokedexcli/types"
	"fmt"
	"os"
	"strings"
	"time"
)

func cleanInput(text string) []string {
	stringVals := strings.Fields(strings.ToLower(text))
	return stringVals
}

func main() {
	cmds := commands.GetCommands()

	scanner := bufio.NewScanner(os.Stdin)
	cfg := &Config{}
	cfg.Cache = pokecache.NewCache(30 * time.Second)
	cfg.Pokemon = map[string]PokemonResponse{}

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		cmd, ok := cmds[input[0]]
		params := input[1:]

		fmt.Println(cmd.Name)

		if !ok {
			fmt.Printf("Command \"%s\" not found.\n", input[0])
		} else {
			err := cmd.Callback(cfg, params)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
