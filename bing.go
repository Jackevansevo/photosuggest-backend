package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type bingPhoto struct {
	Name             string `json:"name"`
	URL              string `json:"contentUrl"`
	Thumbnail        string `json:"thumbnailUrl"`
	Origin           string `json:"hostPageUrl"`
	InsightsMetadata struct {
		BestRepresentativeQuery struct {
			Text string `json:"text"`
		} `json:"bestRepresentativeQuery"`
	} `json:"insightsMetadata"`
}

type bing struct {
	client  http.Client
	baseURL url.URL
}

var bingSearchParams = url.Values{
	"count":   {"150"},
	"Type":    {"Photo"},
	"modules": {"Tags"},
}

var domainBlackList = []*regexp.Regexp{
	regexp.MustCompile(".*pixabay\\.com"),
}

func newBing(clinet http.Client, apiKey string) (source, error) {

	var bingService bing

	bingURL, err := url.Parse("https://api.cognitive.microsoft.com/bing/v7.0/images/search")

	if err != nil {
		return bingService, err
	}

	bingService.baseURL = *bingURL

	return bingService, nil

}

func (b bing) buildURL(query string, license string) string {

	newURL := url.URL(b.baseURL)

	params := newURL.Query()

	params["q"] = []string{query}

	switch license {

	case "public":
		params["license"] = []string{"Public"}
	case "share":
		params["license"] = []string{"Share"}
	case "sharecommercially":
		params["license"] = []string{"ShareCommercially"}
	case "modify":
		params["license"] = []string{"Modify"}
	case "modifycommercially":
		params["license"] = []string{"ModifyCommercially"}
	default:
		params["license"] = []string{"Any"}
	}

	for key, value := range bingSearchParams {
		params[key] = value
	}

	newURL.RawQuery = params.Encode()

	return newURL.String()
}

func (b bing) query(query string, license string) ([]interface{}, error) {

	url := b.buildURL(query, license)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", Env.BingAPIKey)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// Report any server errors
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	json, err := b.process(body)

	if err != nil {
		return nil, err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return json, nil
}

func isBlackListedDomain(target string) bool {
	for _, regex := range domainBlackList {
		if regex.MatchString(target) {
			return true
		}
	}
	return false
}

func (b bing) process(body []byte) ([]interface{}, error) {

	var data struct {
		Photos []bingPhoto `json:"value"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	photos := make([]interface{}, 0, len(data.Photos))

	for _, photo := range data.Photos {

		domain := strings.Split(photo.URL, "/")[2]

		if !isBlackListedDomain(domain) {
			photos = append(photos, struct {
				Source      string `json:"source"`
				Origin      string `json:"origin"`
				Title       string `json:"title"`
				Thumbnail   string `json:"thumbnail"`
				URL         string `json:"url"`
				Description string `json:"description"`
			}{
				"bing",
				photo.Origin,
				photo.Name,
				photo.Thumbnail,
				photo.URL,
				photo.InsightsMetadata.BestRepresentativeQuery.Text,
			})
		}
	}

	return photos, nil

}
