package commands

import (
	"encoding/json"
	. "example.com/dabates/pokedexcli/types"
	"fmt"
	"io"
	"net/http"
)

func CommandMapb(cfg *Config, params []string) error {
	var next string

	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		next = cfg.Previous
	}

	body, err := cfg.Cache.Get(next)

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

		cfg.Cache.Add(next, body)
	}

	data := LocationAreaResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous

	return nil
}

func CommandMap(cfg *Config, params []string) error {
	var next string

	if cfg.Next == "" {
		next = LocationURL
	} else {
		next = cfg.Next
	}

	body, err := cfg.Cache.Get(next)

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

		cfg.Cache.Add(next, body)
	}

	data := LocationAreaResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous

	return nil
}
