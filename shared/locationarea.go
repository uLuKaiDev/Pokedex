package shared

import (
	"sync"
)

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var (
	locationAreaInstance *LocationArea
	onceLocationArea     sync.Once
)

func NewLocationArea() *LocationArea {
	onceLocationArea.Do(func() {
		locationAreaInstance = &LocationArea{}
	})
	return locationAreaInstance
}
