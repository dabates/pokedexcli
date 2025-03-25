package commands

import (
	"example.com/dabates/pokedexcli/types"
	"fmt"
	"os"
)

func CommandExit(cfg *types.Config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
