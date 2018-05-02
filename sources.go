package main

import (
	"net/http"
	"time"
)

type source interface {
	query(string, string) ([]interface{}, error)
	process([]byte) ([]interface{}, error)
}

var client = http.Client{Timeout: time.Duration(10 * time.Second)}
var sources = make(map[string]source)
var defaultSources []string

func setupSources() error {

	flickr, err := newFlickr(client, Env.FlickrAPIKey)
	if err != nil {
		return err
	}
	sources["flickr"] = flickr

	bing, err := newBing(client, Env.BingAPIKey)
	if err != nil {
		return err
	}
	sources["bing"] = bing

	defaultSources = []string{"bing", "flickr"}

	return nil
}
