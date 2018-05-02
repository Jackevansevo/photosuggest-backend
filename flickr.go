package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type flickr struct {
	client  http.Client
	baseURL url.URL
}

type licenseInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var flickrURL, _ = url.Parse("https://www.flickr.com/services/rest/")

var licenseTypes = map[string]licenseInfo{
	"0":  {"All Rights Reserved", ""},
	"1":  {"Attribution-NonCommercial-ShareAlike License", "https://creativecommons.org/licenses/by-nc-sa/2.0/"},
	"2":  {"Attribution-NonCommercial License", "https://creativecommons.org/licenses/by-nc/2.0/"},
	"3":  {"Attribution-NonCommercial-NoDerivs License", "https://creativecommons.org/licenses/by-nc-nd/2.0/"},
	"4":  {"Attribution License", "https://creativecommons.org/licenses/by/2.0/"},
	"5":  {"Attribution-ShareAlike License", "https://creativecommons.org/licenses/by-sa/2.0/"},
	"6":  {"Attribution-NoDerivs License", "https://creativecommons.org/licenses/by-nd/2.0/"},
	"7":  {"No known copyright restrictions", "https://www.flickr.com/commons/usage/"},
	"8":  {"United States Government Work", "http://www.usa.gov/copyright.shtml"},
	"9":  {"Public Domain Dedication (CC0)", "https://creativecommons.org/publicdomain/zero/1.0/"},
	"10": {"Public Domain Mark", "https://creativecommons.org/publicdomain/mark/1.0/"},
}

type flickrPhoto struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Server      string `json:"server"`
	Farm        int    `json:"farm"`
	Title       string `json:"title"`
	License     string `json:"license"`
	Owner       string `json:"owner"`
	OwnerName   string `json:"ownername"`
	Description struct {
		Content string `json:"_content"`
	} `json:"description"`
}

var flickrSearchParams = url.Values{
	"media":        {"photos"},
	"method":       {"flickr.photos.search"},
	"per_page":     {"30"},
	"safe_search":  {"1"},
	"sort":         {"relevance"},
	"content_type": {"1"},
	"extras":       {"description,license,owner_name"},
}

func newFlickr(client http.Client, apiKey string) (flickr, error) {

	flickrService := flickr{client, *flickrURL}

	params := url.Values{
		"format":         {"json"},
		"nojsoncallback": {"1"},
		"api_key":        {apiKey},
	}

	flickrService.baseURL.RawQuery = params.Encode()

	return flickrService, nil
}

func (f flickr) buildURL(query string, license string) string {

	newURL := url.URL(f.baseURL)

	// Get the url.Values
	params := newURL.Query()

	params["text"] = []string{query}

	var flickrLicenses []string

	switch license {
	case "public":
		flickrLicenses = []string{}
	case "share":
		flickrLicenses = []string{"1", "2", "3", "4", "5", "6"}
	case "sharecommercially":
		flickrLicenses = []string{"4", "5", "6"}
	case "modify":
		flickrLicenses = []string{"1", "2", "4", "5"}
	case "modifycommercially":
		flickrLicenses = []string{"4", "5"}
	default:
		flickrLicenses = []string{"1", "2", "3", "4", "5", "6"}
	}

	flickrLicenses = append(flickrLicenses, []string{"7", "8", "9", "10"}...)

	params["license"] = []string{strings.Join(flickrLicenses, ",")}

	for key, value := range flickrSearchParams {
		params[key] = value
	}

	newURL.RawQuery = params.Encode()
	return newURL.String()

}

func (f flickr) query(query string, license string) ([]interface{}, error) {

	searchEndpoint := f.buildURL(query, license)

	resp, err := getBytes(searchEndpoint, f.client)

	if err != nil {
		return nil, err
	}

	json, err := f.process(resp)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func (f flickr) process(body []byte) ([]interface{}, error) {

	var data struct{ Photos struct{ Photo []flickrPhoto } }

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	photos := make([]interface{}, len(data.Photos.Photo))

	urlFmt := "https://farm%d.staticflickr.com/%s/%s_%s.jpg"

	for i, p := range data.Photos.Photo {

		photos[i] = struct {
			Source      string `json:"source"`
			Owner       string `json:"owner"`
			URL         string `json:"url"`
			Thumbnail   string `json:"thumbnail"`
			Origin      string `json:"origin"`
			Title       string `json:"title"`
			Description string `json:"description"`
			licenseInfo `json:"license"`
		}{
			"flickr",
			p.OwnerName,
			fmt.Sprintf(urlFmt, p.Farm, p.Server, p.ID, p.Secret),
			fmt.Sprintf(urlFmt, p.Farm, p.Server, p.ID, p.Secret),
			fmt.Sprintf("https://flickr.com/%s/%s", p.Owner, p.ID),
			p.Title,
			p.Description.Content,
			licenseTypes[p.License],
		}

	}

	return photos, nil
}
