package main

import (
	"io/ioutil"
	"net/http"
)

// [TODO] Need to be able to mock this out

// getBytes Makes a Http Get request and returns body as bytes
func getBytes(url string, client http.Client) ([]byte, error) {

	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil

}
