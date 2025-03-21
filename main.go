package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/uLuKaiDev/Pokedex/internal/pokecache"
	"github.com/uLuKaiDev/Pokedex/shared"
)

func main() {
	initCommands()
	shared.CheckSavedPokedex()

	cache := pokecache.NewCache(5 * time.Second)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		text := scanner.Text()
		cleanedInput := CleanInput(text)
		if len(cleanedInput) == 0 {
			fmt.Print("Pokedex > ")
			continue
		}

		command, exists := commands[cleanedInput[0]]
		if exists {
			var err error
			if len(cleanedInput) > 1 {
				err = command.callback(cache, cleanedInput[1])
			} else {
				err = command.callback(cache)
			}

			if err != nil {
				fmt.Println("Error executing command:", err)
			}
		} else {
			fmt.Println("Unknown command:", cleanedInput[0])
		}

		fmt.Print("Pokedex > ")
	}
}

func CleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	var result []string

	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}

	return result
}
