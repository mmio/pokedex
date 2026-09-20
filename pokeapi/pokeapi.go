package pokeapi

import (
	"net/http"

	"github.com/mmio/pokedex/utilities"
)

type PokeAPIState struct {
	Next     *string
	Previous *string
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

func NewPokeAPI() PokeAPIState {
	return PokeAPIState{}
}

func (pokeAPIState *PokeAPIState) CallMap() (LocationResponse, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/"

	if pokeAPIState.Next != nil {
		endpoint = *pokeAPIState.Next
	}

	locationResponse, err := utilities.CallEndpoint[any, LocationResponse](
		http.MethodGet,
		endpoint,
		struct{}{},
	)

	if err != nil {
		return LocationResponse{}, err
	}

	pokeAPIState.Next = &locationResponse.Next

	return locationResponse, nil
}

func (pokeAPIState *PokeAPIState) CallMapBack() (LocationResponse, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/"

	if pokeAPIState.Previous != nil {
		endpoint = *pokeAPIState.Previous
	}

	locationResponse, err := utilities.CallEndpoint[any, LocationResponse](
		http.MethodGet,
		endpoint,
		struct{}{},
	)

	if err != nil {
		return LocationResponse{}, err
	}

	pokeAPIState.Previous = locationResponse.Previous

	return locationResponse, nil
}
