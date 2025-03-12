package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/uLuKaiDev/Pokedex/shared"
)

func fetchMapData(url string, c *shared.Config) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get data: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code %d and body: %s", res.StatusCode, body)
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

func CommandMap(c *shared.Config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if c.Next != "" {
		url = c.Next
	}

	if err := fetchMapData(url, c); err != nil {
		return err
	}

	for _, location := range c.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func CommandMapBack(c *shared.Config) error {
	if c.Previous == "" {
		fmt.Println("No previous page available")
		return nil
	}
	previousURL, ok := c.Previous.(string)
	if !ok || previousURL == "" {
		fmt.Println("You're on the first page")
		return nil
	}

	if err := fetchMapData(previousURL, c); err != nil {
		return err
	}

	for _, location := range c.Results {
		fmt.Println(location.Name)
	}
	return nil
}
