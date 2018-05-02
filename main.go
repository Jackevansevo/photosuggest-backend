package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type envSpec struct {
	BingAPIKey   string `envconfig:"BING_API_KEY" required:"true"`
	FlickrAPIKey string `envconfig:"FLICKR_API_KEY" required:"true"`
}

// Env contains environment variables
var Env envSpec

type result struct {
	photos []interface{}
	source string
	err    error
}

func main() {
	envconfig.MustProcess("", &Env)
	err := setupSources()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/", query)
	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func query(c *gin.Context) {

	// Make a results channel
	resultsChan := make(chan result)

	type Params struct {
		Query   string   `form:"q"`
		License string   `form:"license"`
		Sources []string `form:"sources"`
	}

	var params Params
	if err := c.BindQuery(&params); err != nil {
		return
	}

	if params.Query == "" {
		c.AbortWithStatusJSON(405, gin.H{"error": "specify query"})
		return
	}

	switch params.License {
	case "":
	case "any":
	case "public":
	case "share":
	case "sharecommercially":
	case "modify":
	case "modifycommercially":
	default:
		c.AbortWithStatusJSON(405, gin.H{"error": "unknown license"})
		return
	}

	if len(params.Sources) == 0 {
		params.Sources = defaultSources
	}

	// Start a Goroutine to Query each source
	for _, source := range params.Sources {
		go func(query string, source string) {
			results, err := sources[source].query(query, params.License)
			resultsChan <- result{results, source, err}
		}(params.Query, source)
	}

	results := make([]interface{}, 0)
	status := make(map[string]string)

	// Collect the results
	for range params.Sources {
		result := <-resultsChan
		if result.photos != nil {
			results = append(results, result.photos...)
		}
		if result.err != nil {
			status[result.source] = result.err.Error()
		} else {
			status[result.source] = "ok"
		}
	}

	c.JSON(200, gin.H{"status": status, "results": results})

}
