package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/uLuKaiDev/Pokedex/internal/pokecache"
)

func FetchMapData(url string, c any, ca *pokecache.Cache) error {
	if data, ok := ca.Get(url); ok {
		err := json.Unmarshal(data, c)
		fmt.Println("*** Data fetched from cache")
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
