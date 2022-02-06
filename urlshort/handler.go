package urlshort

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	yml "gopkg.in/yaml.v2"
)

type Redirect struct {
	Path string
	URL  string
}

func buildMap(redirects []Redirect) map[string]string {
	m := make(map[string]string)
	for _, r := range redirects {
		m[r.Path] = r.URL
	}
	return m
}

var redirects = []byte("redirects")

func redirectHandler(path, newURL string) http.Handler {
	rh := func(w http.ResponseWriter, r *http.Request) {
		if path == r.URL.Path {
			http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
		}
	}
	return http.HandlerFunc(rh)
}

func getInput(filePath string) map[string]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileType := strings.SplitAfter(filePath, ".")
	var parsedYaml []Redirect

	switch fileType[1] {
	case "yaml":
		d := yml.NewDecoder(file)
		err = d.Decode(&parsedYaml)
		if err != nil {
			log.Fatalf("yaml decoding failed with '%s'\n", err)
		}
	case "json":
		d := json.NewDecoder(file)
		err = d.Decode(&parsedYaml)
		if err != nil {
			log.Fatalf("json decoding failed with '%s'\n", err)
		}
	default:
		fmt.Println("No valid file type passed")
	}
	pathMap := buildMap(parsedYaml)
	return pathMap
}

func BuildDB(file string) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	input := getInput(file)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(redirects)
		if err != nil {
			return err
		}
		for path, url := range input {
			err = bucket.Put([]byte(path), []byte(url))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func DBHandler(fallback http.Handler) (http.HandlerFunc, error) {
	mux := http.NewServeMux()
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(redirects)
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", redirects)
		}
		if e := bucket.ForEach(func(p, u []byte) error {
			rh := redirectHandler(string(p), string(u))
			mux.Handle(string(p), rh)
			return nil
		}); e != nil {
			return e
		}
		return nil
	})
	mux.Handle("/", fallback)
	if err != nil {
		return nil, err
	}
	return mux.ServeHTTP, nil
}

/*
------------------------------------------------------------------------------------
--- ORIGINAL SOLUTION COMMENTED OUT FOR LEGACY REASONS; SOLUTION IS USING BOLTDB ---
------------------------------------------------------------------------------------

func parseYAML(yaml []byte) ([]Redirect, error) {
	var redirects []Redirect
	err := yml.Unmarshal(yaml, &redirects)
	if err != nil {
		return []Redirect{}, err
	}
	return redirects, nil
} */

/*

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
} */

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

/*

func InputHandler(filePath string, fallback http.Handler) (http.HandlerFunc, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileType := strings.SplitAfter(filePath, ".")
	fmt.Println("INPUT: ", fileType[1])
	var parsedYaml []Redirect

	switch fileType[1] {
	case "yaml":
		d := yml.NewDecoder(file)
		err = d.Decode(&parsedYaml)
		if err != nil {
			log.Fatalf("yaml decoding failed with '%s'\n", err)
		}
	case "json":
		d := json.NewDecoder(file)
		err = d.Decode(&parsedYaml)
		if err != nil {
			log.Fatalf("json decoding failed with '%s'\n", err)
		}
	default:
		fmt.Println("No valid file type passed")
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
*/
