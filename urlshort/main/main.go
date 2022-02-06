package main

import (
	"flag"
	"fmt"
	"net/http"

	urlshort "github.com/pparmin/gophercises/urlshort"
)

func main() {
	mux := defaultMux()
	inputFile := flag.String("r", "yaml/default.yaml", "Specify an input file which holds a number of paths and redirects (currently accepted file types: .yaml, .json")
	flag.Parse()

	urlshort.BuildDB(*inputFile)
	dbHandler, err := urlshort.DBHandler(mux)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

/*
------------------------------------------------------------------------------------
--- ORIGINAL SOLUTION COMMENTED OUT FOR LEGACY REASONS; SOLUTION IS USING BOLTDB ---
------------------------------------------------------------------------------------

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/google":         "https://google.com",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	yamlHandler, err := urlshort.InputHandler(*inputFile, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("YAML: ", yamlHandler, "\n")
*/
