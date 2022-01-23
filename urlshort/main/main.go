package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	urlshort "github.com/pparmin/gophercises/urlshort"
)

func getHandler(fileType string) http.HandlerFunc {
	var handler http.HandlerFunc
	switch fileType {
	case "yaml":

	}
	return nil
}

func main() {
	mux := defaultMux()
	yamlFile := flag.String("y", "yaml/default.yaml", "Specify a yaml file which holds a number of paths and redirects")
	//jsonFile := flag.String("j", "json/default.json", "Specify a json file which holds a number of paths and redirects")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/google":         "https://google.com",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// - path: /amazon
	//   url: https://amazon.de/
	// `
	// yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	// if err != nil {
	// 	panic(err)
	// }
	fileType := strings.SplitAfter(*yamlFile, ".")
	fmt.Println("INPUT: ", fileType[1])
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
