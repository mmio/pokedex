package utilities

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

func CallEndpoint[T, U any](method, endpoint string, data T) (U, error) {
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

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return zero, fmt.Errorf("Bad status code: %v", response.StatusCode)
	}

	result, err := decode[U](response.Body)
	if err != nil {
		return zero, err
	}

	return result, nil
}
