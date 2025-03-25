package main

import (
	. "example.com/dabates/pokedexcli/commands"
	. "example.com/dabates/pokedexcli/types"
)

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"map": {
			Name:        "map",
			Description: "Map",
			Callback:    CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Map back one page",
			Callback:    CommandMapb,
		},
		"exit": {
			Name:        "exit",
			Description: "Exits the Pokedex",
			Callback:    CommandExit,
		},
		"help": {
			Name:        "help",
			Description: "Shows this help message",
			Callback:    CommandHelp,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore the Pokedex",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a Pokemon",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a Pokemon",
			Callback:    CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Display the caught Pokemon",
			Callback:    CommandPokedex,
		},
	}
}
