package commands

import (
	"example.com/dabates/pokedexcli/types"
	"fmt"
)

func CommandHelp(cfg *types.Config, params []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}
