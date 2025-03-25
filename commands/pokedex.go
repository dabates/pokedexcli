package commands

import (
	. "example.com/dabates/pokedexcli/types"
	"fmt"
)

func CommandPokedex(cfg *Config, params []string) error {
	fmt.Println("Your Pokedex:")

	p := cfg.Pokemon

	if len(p) == 0 {
		return nil
	}

	for _, pokemon := range p {
		fmt.Println("-", pokemon.Name)
	}

	return nil
}
