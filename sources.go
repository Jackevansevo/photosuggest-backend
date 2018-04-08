package main

type source interface {
	query(string, string) ([]interface{}, error)
	buildURL(string, string) string
	process([]byte) ([]interface{}, error)
}

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

	for k := range sources {
		defaultSources = append(defaultSources, k)
	}

	return nil
}
