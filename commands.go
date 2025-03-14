package main

import (
	"fmt"
	"os"

	"github.com/uLuKaiDev/Pokedex/internal/pokeapi"
	"github.com/uLuKaiDev/Pokedex/internal/pokecache"
	//"github.com/uLuKaiDev/Pokedex/shared"
)

type cliCommand struct {
	name        string
	description string
	callback    func(ca *pokecache.Cache, extra ...any) error
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
		"explore": {
			name:        "explore",
			description: "Explore a specified location area",
			callback:    pokeapi.CommandExplore,
		},
	}
}

func CommandExit(_ *pokecache.Cache, _extra ...any) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(_ *pokecache.Cache, _extra ...any) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
