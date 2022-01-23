package urlshort

import (
	"net/http"

	yml "gopkg.in/yaml.v2"
)

type Redirect struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func redirectHandler(path, newURL string) http.Handler {
	rh := func(w http.ResponseWriter, r *http.Request) {
		if path == r.URL.Path {
			http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
		}
	}
	return http.HandlerFunc(rh)
}

func parseYAML(yaml []byte) ([]Redirect, error) {
	var redirects []Redirect
	err := yml.Unmarshal(yaml, &redirects)
	if err != nil {
		return []Redirect{}, err
	}
	return redirects, nil
}

func buildMap(redirects []Redirect) map[string]string {
	m := make(map[string]string)
	for _, r := range redirects {
		m[r.Path] = r.URL
	}
	return m
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	mux := http.NewServeMux()
	for p, u := range pathsToUrls {
		rh := redirectHandler(p, u)
		mux.Handle(p, rh)
	}
	mux.Handle("/", fallback)
	return mux.ServeHTTP
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
