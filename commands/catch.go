package commands

import (
	"encoding/json"
	. "example.com/dabates/pokedexcli/types"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

func CommandCatch(cfg *Config, params []string) error {
	fmt.Println("Throwing a Pokeball at", params[0]+"...")

	body, err := cfg.Cache.Get(params[0])
	if err != nil {
		url := PokemonURL + params[0]
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		res.Body.Close()
		cfg.Cache.Add(params[0], body)
	}

	data := PokemonResponse{}
	err = json.Unmarshal(body, &data)

	// do the calc here
	chance := 1.0 / (1.0 + float64(data.BaseExperience)/50.0)
	if rand.Float64() < chance {
		fmt.Println("Pokemon:", params[0], "was caught!")
		cfg.Pokemon[params[0]] = data
	} else {
		fmt.Println("Pokemon:", params[0], "escaped!")
	}

	return nil
}
