package shared

import (
	"encoding/json"
	"fmt"
	"os"
)

var Pokedex map[string]Pokemon

func init() {
	Pokedex = make(map[string]Pokemon)
}

func ClearPokedex() {
	Pokedex = make(map[string]Pokemon)
	SavePokedex()
}

func AddToPokedex(p *Pokemon) {
	Pokedex[p.Name] = *p
}

func SavePokedex() {
	data, err := json.Marshal(Pokedex)
	if err != nil {
		fmt.Println("Error marshalling Pokedex:", err)
		return
	}

	err = os.WriteFile("saves/pokedex_save.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing Pokedex:", err)
		return
	}

}

func CheckSavedPokedex() {
	if err := os.MkdirAll("saves", 0755); err != nil {
		fmt.Println("Error creating saves directory:", err)
		return
	}

	_, err := os.Stat("saves/pokedex_save.json")
	if err != nil {
		fmt.Println("No Pokedex save file found, starting a new adventure!")
		return
	}

	data, err := os.ReadFile("saves/pokedex_save.json")
	if err != nil {
		fmt.Println("Error reading Pokedex save file:", err)
		return
	}
	err = json.Unmarshal(data, &Pokedex)
	if err != nil {
		fmt.Println("failed to unmarshal JSON:", err)
	}

	if len(Pokedex) == 0 {
		fmt.Println("No Pokedex found, starting a new adventure!")
		return
	}
	fmt.Println("Pokedex loaded, you have caught", len(Pokedex), "Pokemon!")
}
