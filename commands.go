package main

import (
	"fmt"
	"os"

	"github.com/uLuKaiDev/Pokedex/internal/pokeapi"
	"github.com/uLuKaiDev/Pokedex/shared"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *shared.Config) error
}

var commands map[string]cliCommand

func initCommands() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    CommandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    CommandExit,
		},
		"map": {
			name:        "map",
			description: " Displays next 20 location areas", //extra space for formatting
			callback:    pokeapi.CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location areas",
			callback:    pokeapi.CommandMapBack,
		},
	}
}

func CommandExit(_ *shared.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(_ *shared.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
