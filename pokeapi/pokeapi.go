package pokeapi

import (
	"bytes"
	"encoding/json"
	"http"
	"io"
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
	defer response.body.Close()

	result, err := decode[U](response.body)
	if err != nil {
		return zero, err
	}

	return result, nil
}

type LocationRequest struct {
}

type Location struct {
}

type LocationResponse struct {
}

func CallMap(requestData LocationRequest) (LocationResponse, error) {
	endpoint := "https://www.pokeapi.com/location"

	locationResponse, err := callEndpoint[LocationRequest, LocationResponse](
		http.POSTMethod,
		endpoint,
		requestData,
	)

	if err != nil {
		return LocationResponse{}, err
	}

	return locationResponse, nil
}
