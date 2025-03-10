package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func CleanInput(text string) []string {
	words := strings.Split(strings.ToLower(text), " ")
	var result []string

	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}

	return result
}
