package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/uLuKaiDev/Pokedex/shared"
)

func main() {
	initCommands()

	scanner := bufio.NewScanner(os.Stdin)
	cfg := &shared.Config{}
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
			err := command.callback(cfg)
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
