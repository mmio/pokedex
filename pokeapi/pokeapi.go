package pokeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func encode(data any) (io.Reader, error) {
	blob, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(blob), nil
}

func decode[T any](data io.Reader) (T, error) {
	var results T
	if err := json.NewDecoder(data).Decode(&results); err != nil {
		return results, err
	}

	return results, nil
}

func callEndpoint[T, U any](method, endpoint string, data T) (U, error) {
	var zero U

	payload, err := encode(data)
	if err != nil {
		return zero, err
	}

	request, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return zero, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return zero, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 && response.StatusCode > 299 {
		return zero, fmt.Errorf("Bad status code: %v", response.StatusCode)
	}

	result, err := decode[U](response.Body)
	if err != nil {
		return zero, err
	}

	return result, nil
}

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

	locationResponse, err := callEndpoint[any, LocationResponse](
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

	locationResponse, err := callEndpoint[any, LocationResponse](
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
