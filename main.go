package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cleanInput(text string) []string {
	stringVals := strings.Fields(strings.ToLower(text))
	return stringVals
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func main() {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Shows this help message",
			callback:    commandHelp,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		cmd, ok := commands[input[0]]
		fmt.Println(input)
		if !ok {
			fmt.Printf("Command \"%s\" not found.\n", input[0])
		} else {
			err := cmd.callback()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
