package client

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Upload(server string, input io.Reader) (string, error) {
	url := url.URL{Scheme: "http", Host: server}
	if url.Port() == "" {
		url.Host = url.Host + ":8090"
	}

	resp, err := http.Post(url.String(), "text/plain", input)
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

func Download(server string, clipcode string) (io.ReadCloser, error) {
	url := url.URL{Scheme: "http", Host: server, Path: clipcode}
	if url.Port() == "" {
		url.Host = url.Host + ":8090"
	}

	resp, err := http.Get(url.String())
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
