package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func stubSource(s source, path string, query string) ([]interface{}, error) {
	fname := path + strings.ToLower(query) + ".json"
	if _, err := os.Stat(fname); err == nil {
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			return nil, err
		}
		results, err := s.process(data)
		if err != nil {
			return nil, err
		}
		return results, err
	}
	return nil, nil
}
