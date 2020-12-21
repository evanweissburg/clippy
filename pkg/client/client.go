package client

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	BaseURL string
}

func (c *Client) Upload(input io.Reader) (string, error) {
	resp, err := http.Post(c.BaseURL, "text/plain", input)
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

func (c *Client) Download(clipcode string) (io.ReadCloser, error) {
	resp, err := http.Get(c.BaseURL + clipcode)
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
