package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type bingPhoto struct {
	Name             string `json:"name"`
	URL              string `json:"contentUrl"`
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

func (b bing) process(body []byte) ([]interface{}, error) {

	var data struct {
		Photos []bingPhoto `json:"value"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	photos := make([]interface{}, len(data.Photos))

	for i, p := range data.Photos {

		if strings.Split(p.URL, "/")[2] == "cdn.pixabay.com" {
			continue
		}

		photos[i] = struct {
			Source      string `json:"source"`
			Name        string `json:"name"`
			URL         string `json:"url"`
			Description string `json:"description"`
		}{
			"bing",
			p.Name,
			p.URL,
			p.InsightsMetadata.BestRepresentativeQuery.Text,
		}
	}

	return photos, nil

}
