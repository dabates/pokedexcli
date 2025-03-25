package commands

import (
	"example.com/dabates/pokedexcli/types"
	"fmt"
)

func CommandHelp(cfg *types.Config, params []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\n")

	for _, c := range GetCommands() {
		fmt.Println(c.Name+":", c.Description)
	}

	return nil
}
