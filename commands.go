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
			description:  "Displays all available commands, including this one!",
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
			description:  "Displays the next 20 location areas",
			callback:     commandMap,
			commandOrder: 3,
		},
		"mapb": {
			name:         "mapb",
			description:  "Displays the previous 20 location areas",
			callback:     commandMapBack,
			commandOrder: 4,
		},
		"explore": {
			name:         "explore",
			description:  "Explore a specified location. Usage explore <location name> (e.g. Inspect mt-coronet-6f)",
			callback:     commandExplore,
			commandOrder: 5,
		},
		"catch": {
			name:         "catch",
			description:  "Try to catch a Pokemon. Usage catch <Pokemon name> (e.g. catch pikachu)",
			callback:     commandCatch,
			commandOrder: 6,
		},
		"inspect": {
			name:         "inspect",
			description:  "Inspect a caught Pokemon. Usage inspect <Pokemon name> (e.g. inspect pikachu)",
			callback:     commandInspect,
			commandOrder: 7,
		},
		"pokedex": {
			name:         "pokedex",
			description:  "Displays all caught Pokemon",
			callback:     commandPokedex,
			commandOrder: 8,
		},
		"save": {
			name:         "save",
			description:  "Save the current state of the Pokedex:",
			callback:     commandSave,
			commandOrder: 9,
		},
		"new": {
			name:         "new",
			description:  "Clear your current Pokedex and start a new adventure!",
			callback:     commandNew,
			commandOrder: 10,
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
		fmt.Printf("You may now inspect it with the inspect command.\n")
		return nil
	}
	fmt.Printf("%s broke free!\n", pok.Name)

	return nil
}

func commandInspect(_ *pokecache.Cache, extra ...any) error {
	if len(extra) == 0 {
		return fmt.Errorf("specify which Pokemon you want to inspect")
	}

	inputName, ok := extra[0].(string)
	if !ok || inputName == "" {
		return fmt.Errorf("invalid Pokemon name")
	}

	pok, ok := shared.Pokedex[inputName]
	if !ok {
		return fmt.Errorf("you have not caught a %s yet", inputName)
	}

	fmt.Printf("Name: %s\n", pok.Name)
	fmt.Printf("Height: %d\n", pok.Height)
	fmt.Printf("Weight: %d\n", pok.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pok.Stats {
		fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, poke_type := range pok.Types {
		fmt.Printf("- %s\n", poke_type.Type.Name)
	}
	return nil
}

func commandPokedex(_ *pokecache.Cache, extra ...any) error {
	if len(shared.Pokedex) == 0 {
		fmt.Printf("You have not caught any Pokemon yet!\n")
		return nil
	}
	fmt.Print("Your Pokedex:\n")
	for name := range shared.Pokedex {
		fmt.Printf("- %s\n", name)
	}
	return nil
}

func commandSave(ca *pokecache.Cache, extra ...any) error {
	fmt.Printf("Saving the current state of the Pokedex...\n")
	shared.SavePokedex()
	fmt.Println("Pokedex saved!")
	return nil
}

func commandNew(ca *pokecache.Cache, extra ...any) error {
	fmt.Printf("Starting a new adventure...\n")
	shared.ClearPokedex()
	return nil
}
