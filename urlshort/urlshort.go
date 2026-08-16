package urlshort

import (
	"fmt"
	"net/http"

	"github.com/kevinjimenez96/url-shortener/internal/database"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var yamlMappings []UrlMapping

	err := yaml.Unmarshal([]byte(yml), &yamlMappings)

	if err != nil {
		return nil, err
	}

	urlMappings := make(map[string]string)

	for _, mapping := range yamlMappings {
		urlMappings[mapping.Path] = mapping.Url
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := urlMappings[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil
}

func DbHandler(db *database.Queries, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mapping, err := db.GetMapping(r.Context(), r.URL.Path)

		if err != nil {
			fmt.Print(err)
		}

		if err == nil && len(mapping.Url) > 0 {
			http.Redirect(w, r, mapping.Url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type UrlMapping struct {
	Path string
	Url  string
}
