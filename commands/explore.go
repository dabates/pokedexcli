package commands

import (
	"encoding/json"
	. "example.com/dabates/pokedexcli/types"
	"fmt"
	"io"
	"net/http"
)

func CommandExplore(cfg *Config, params []string) error {
	fmt.Println("Count of parameters:", len(params))

	for _, param := range params {
		url := LocationURL + param

		body, err := cfg.Cache.Get(param)
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

			cfg.Cache.Add(param, body)
		}

		data := LocationResponse{}
		err = json.Unmarshal(body, &data)

		fmt.Println("Found Pokemon:")
		for _, encounter := range data.Encounters {
			fmt.Println("-", encounter.Pokemon.Name)
		}
	}
	return nil
}
