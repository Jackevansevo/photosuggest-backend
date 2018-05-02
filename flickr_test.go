package main

import (
	"net/url"
	"testing"
)

func TestProcessJSON(t *testing.T) {
	flickr, err := newFlickr(client, "")
	if err != nil {
		t.Errorf(err.Error())
	}
	photos, err := stubSource(flickr, "fixtures/flickr/", "dogs")
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(photos) == 0 {
		t.Errorf("Returned no photos")
	}
}

func TestProcessMalformedJSON(t *testing.T) {
	flickr, err := newFlickr(client, "")
	if err != nil {
		t.Error(err)
	}
	_, err = stubSource(flickr, "fixtures/flickr/", "malformed")
	if err == nil {
		t.Errorf("Malformed JSON should return an error")
	}
}

func TestBuildURL(t *testing.T) {

	flickrClient, err := newFlickr(client, "")
	if err != nil {
		t.Error(err)
	}

	licenses := map[string]string{
		"public":             "7,8,9,10",
		"share":              "1,2,3,4,5,6,7,8,9,10",
		"sharecommercially":  "4,5,6,7,8,9,10",
		"modify":             "1,2,4,5,7,8,9,10",
		"modifycommercially": "4,5,7,8,9,10",
		"":                   "1,2,3,4,5,6,7,8,9,10",
	}

	for category, expected := range licenses {
		parsed, _ := url.Parse(flickrClient.buildURL("dogs", category))
		urlParams := parsed.Query()["license"][0]
		if urlParams != expected {
			t.Errorf("expected %#v, got: %#v", expected, urlParams)
		}

	}

}
