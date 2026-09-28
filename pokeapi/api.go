package pokeapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/mmio/pokedex/utilities"
)

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

func (pokeAPIState *PokeAPIState) CallExplore(name string) (AreaResponse, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/" + name

	// check cache
	response, ok := pokeAPIState.EmptyGetCache.Get(endpoint)
	if ok {
		areaResponse, ok := response.(AreaResponse)
		if !ok {
			return AreaResponse{}, errors.New("Bad type in cache")
		}

		return areaResponse, nil
	}

	// cache miss
	areaResponse, err := utilities.CallEndpoint[any, AreaResponse](
		http.MethodGet,
		endpoint,
		struct{}{},
	)

	if err != nil {
		return AreaResponse{}, err
	}

	// update
	pokeAPIState.EmptyGetCache.Add(endpoint, areaResponse)

	return areaResponse, nil
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
