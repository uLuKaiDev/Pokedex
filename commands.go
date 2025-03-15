package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/uLuKaiDev/Pokedex/internal/pokeapi"
	"github.com/uLuKaiDev/Pokedex/internal/pokecache"
	"github.com/uLuKaiDev/Pokedex/shared"
)

type cliCommand struct {
	name         string
	description  string
	callback     func(ca *pokecache.Cache, extra ...any) error
	commandOrder int
}

var commands map[string]cliCommand

func initCommands() {
	commands = map[string]cliCommand{
		"help": {
			name:         "help",
			description:  "Displays a help message",
			callback:     commandHelp,
			commandOrder: 1,
		},
		"exit": {
			name:         "exit",
			description:  "Exit the Pokedex",
			callback:     commandExit,
			commandOrder: 2,
		},
		"map": {
			name:         "map",
			description:  "Displays next 20 location areas",
			callback:     commandMap,
			commandOrder: 3,
		},
		"mapb": {
			name:         "mapb",
			description:  "Displays previous 20 location areas",
			callback:     commandMapBack,
			commandOrder: 4,
		},
		"explore": {
			name:         "explore",
			description:  "Explore a specified location area",
			callback:     commandExplore,
			commandOrder: 5,
		},
		"catch": {
			name:         "catch",
			description:  "Catch a Pokemon",
			callback:     commandCatch,
			commandOrder: 6,
		},
	}
}

func commandExit(_ *pokecache.Cache, _extra ...any) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *pokecache.Cache, _extra ...any) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	var cmds []cliCommand
	var longestName int
	for _, cmd := range commands {
		if len(cmd.name) > longestName {
			longestName = len(cmd.name)
		}
		cmds = append(cmds, cmd)
	}
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].commandOrder < cmds[j].commandOrder
	})

	for _, command := range cmds {
		numberOfSpaces := longestName - len(command.name)
		whitespace := strings.Repeat(" ", numberOfSpaces)
		fmt.Printf("%s:%s %s\n", command.name, whitespace, command.description)
	}
	return nil
}

func commandMap(ca *pokecache.Cache, extra ...any) error {
	loc := shared.NewLocationArea()

	url := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	if loc.Next != "" {
		url = loc.Next
	}

	if err := pokeapi.FetchMapData(url, loc, ca); err != nil {
		return err
	}

	for _, location := range loc.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapBack(ca *pokecache.Cache, extra ...any) error {
	loc := shared.NewLocationArea()

	if loc.Previous == nil {
		fmt.Println("No previous page available")
		return nil
	}
	previousURL, ok := loc.Previous.(string)
	if !ok || previousURL == "" {
		fmt.Println("You're on the first page")
		return nil
	}

	if err := pokeapi.FetchMapData(previousURL, loc, ca); err != nil {
		return err
	}

	for _, location := range loc.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(ca *pokecache.Cache, extra ...any) error {
	if len(extra) == 0 {
		return fmt.Errorf("missing location area")
	}

	inputUrl, ok := extra[0].(string)
	if !ok || inputUrl == "" {
		return fmt.Errorf("invalid location area")
	}

	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	url := baseUrl + inputUrl

	loc := shared.NewLocation()
	if err := pokeapi.FetchMapData(url, loc, ca); err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", loc.Name)
	fmt.Printf("Found Pokemon:\n")
	for _, encounter := range loc.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(ca *pokecache.Cache, extra ...any) error {
	if len(extra) == 0 {
		return fmt.Errorf("missing Pokemon name")
	}

	inputUrl, ok := extra[0].(string)
	if !ok || inputUrl == "" {
		return fmt.Errorf("invalid Pokemon name")
	}

	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	url := baseUrl + inputUrl

	pok := shared.NewPokemon()
	if err := pokeapi.FetchMapData(url, pok, ca); err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pok.Name)
	if pok.CatchAttempt() {
		fmt.Printf("You caught %s!\n", pok.Name)
		shared.AddToPokedex(pok)
		fmt.Printf("%s has been added to your Pokedex!\n", pok.Name)
		return nil
	}
	fmt.Printf("%s broke free!\n", pok.Name)

	return nil
}
