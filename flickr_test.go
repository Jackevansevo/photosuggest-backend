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

	flickr, err := newFlickr(client, "")
	if err != nil {
		t.Error(err)
	}

	licenses := map[string][]string{
		"1,2,3":       {"modify"},
		"4,5,6":       {"commercial"},
		"1,2,3,4,5,6": {"modify", "commercial"},
	}

	for expected, types := range licenses {
		res := flickr.buildURL("dogs", types)
		u, _ := url.Parse(res)
		urlParams := u.Query()["license"][0]
		if urlParams != expected {
			t.Errorf("expected %#v, got: %#v", expected, urlParams)
		}

	}

}
