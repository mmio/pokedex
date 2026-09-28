package pokeapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/mmio/pokedex/utilities"
)

type PokeAPIState struct {
	Next          *string
	Previous      *string
	EmptyGetCache *utilities.Cache
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationResponse struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

func NewPokeAPI() (PokeAPIState, error) {
	duration_30s, err := time.ParseDuration("30s")
	if err != nil {
		return PokeAPIState{}, err
	}

	return PokeAPIState{
		EmptyGetCache: utilities.NewCache(duration_30s),
	}, nil
}

func (pokeAPIState *PokeAPIState) CallMap() (LocationResponse, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/"

	// TODO: This is most likely incorrect, e.g. when we get to the end we wrap, which we probably don't want (e.g. we go to the beginning)
	if pokeAPIState.Next != nil {
		endpoint = *pokeAPIState.Next
	}

	// check cache
	response, ok := pokeAPIState.EmptyGetCache.Get(endpoint)
	if ok {
		locationResponse, ok := response.(LocationResponse)
		if !ok {
			return LocationResponse{}, errors.New("Bad type in cache")
		}

		pokeAPIState.Previous = locationResponse.Previous
		pokeAPIState.Next = locationResponse.Next
		return locationResponse, nil
	}

	// cache miss
	locationResponse, err := utilities.CallEndpoint[any, LocationResponse](
		http.MethodGet,
		endpoint,
		struct{}{},
	)

	if err != nil {
		return LocationResponse{}, err
	}

	// update
	pokeAPIState.EmptyGetCache.Add(endpoint, locationResponse)

	pokeAPIState.Previous = locationResponse.Previous
	pokeAPIState.Next = locationResponse.Next

	return locationResponse, nil
}

func (pokeAPIState *PokeAPIState) CallMapBack() (LocationResponse, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/"

	if pokeAPIState.Previous != nil {
		endpoint = *pokeAPIState.Previous
	}

	// check cache
	response, ok := pokeAPIState.EmptyGetCache.Get(endpoint)
	if ok {
		locationResponse, ok := response.(LocationResponse)
		if !ok {
			return LocationResponse{}, errors.New("Bad type in cache")
		}

		pokeAPIState.Previous = locationResponse.Previous
		pokeAPIState.Next = locationResponse.Next
		return locationResponse, nil
	}

	// cache miss
	locationResponse, err := utilities.CallEndpoint[any, LocationResponse](
		http.MethodGet,
		endpoint,
		struct{}{},
	)

	if err != nil {
		return LocationResponse{}, err
	}

	// update
	pokeAPIState.EmptyGetCache.Add(endpoint, locationResponse)

	pokeAPIState.Previous = locationResponse.Previous
	pokeAPIState.Next = locationResponse.Next

	return locationResponse, nil
}
