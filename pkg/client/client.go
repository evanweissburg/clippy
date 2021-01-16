package client

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func Upload(baseURL string, input io.Reader) (string, error) {
	resp, err := http.Post(baseURL, "text/plain", input)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	clipcode := string(body)
	return clipcode, nil
}

func Download(baseURL string, clipcode string) (io.ReadCloser, error) {
	resp, err := http.Get(baseURL + clipcode)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp.Body, nil

	case http.StatusInternalServerError:
		return nil, errors.New("Server Error")

	case http.StatusBadRequest:
		return nil, errors.New("Invalid Clipcode")

	default:
		return nil, errors.New("Unknown Status Code")
	}
}
