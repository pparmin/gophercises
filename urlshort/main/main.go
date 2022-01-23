package main

import (
	"flag"
	"fmt"
	"net/http"

	urlshort "github.com/pparmin/gophercises/urlshort"
)

func getHandler(fileType string) http.HandlerFunc {
	// var handler http.HandlerFunc
	switch fileType {
	case "yaml":

	}
	return nil
}

func main() {
	mux := defaultMux()
	yamlFile := flag.String("r", "yaml/default.yaml", "Specify an input file which holds a number of paths and redirects (currently accepted file types: .yaml, .json")
	//jsonFile := flag.String("j", "json/default.json", "Specify a json file which holds a number of paths and redirects")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/google":         "https://google.com",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	yamlHandler, err := urlshort.InputHandler(*yamlFile, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
