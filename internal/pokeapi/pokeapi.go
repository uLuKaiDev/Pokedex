package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/uLuKaiDev/Pokedex/internal/pokecache"
	"github.com/uLuKaiDev/Pokedex/shared"
)

func fetchMapData(url string, c any, ca *pokecache.Cache) error {
	if data, ok := ca.Get(url); ok {
		err := json.Unmarshal(data, c)
		fmt.Println("Data fetched from cache")
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return nil
	}

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

	ca.Add(url, body)

	return nil
}

func CommandMap(ca *pokecache.Cache, extra ...any) error {
	loc := shared.NewLocationArea()

	url := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	if loc.Next != "" {
		url = loc.Next
	}

	if err := fetchMapData(url, loc, ca); err != nil {
		return err
	}

	for _, location := range loc.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func CommandMapBack(ca *pokecache.Cache, extra ...any) error {
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

	if err := fetchMapData(previousURL, loc, ca); err != nil {
		return err
	}

	for _, location := range loc.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func CommandExplore(ca *pokecache.Cache, extra ...any) error {
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
	if err := fetchMapData(url, loc, ca); err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", loc.Name)
	fmt.Printf("Found Pokemon:\n")
	for _, encounter := range loc.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
