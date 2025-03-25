package commands

import (
	. "example.com/dabates/pokedexcli/types"
	"fmt"
	"strconv"
)

func CommandInspect(cfg *Config, params []string) error {
	data, err := cfg.Pokemon[params[0]]

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
