/* Interesting discussion of two methods to handle more complicated
JSON files: https://tutorialedge.net/golang/parsing-json-with-golang/

In this case we will use map[string]interface, since the data is
unstructured and we will not necessarily now the names of each arc

Another very interesting post about how to do this properly and then
get proper types for each field: https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data

*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	cyoa "github.com/pparmin/gophercises/cyoa/shared"
)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func main() {
	file, err := os.Open("gopher.json")
	if err != nil {
		fmt.Println("ERROR DURING OPENING OF FILE")
		log.Fatal(err)
	}
	defer file.Close()

	var arcs = cyoa.GetArcs(file)
	for name, arc := range arcs {
		fmt.Println("ARC NAME: ", name)
		fmt.Println("title: ", arc.Title)
		fmt.Println("story: ", arc.Story)
		fmt.Println("Options: ", arc.Options)
		fmt.Println()
	}

	mux := defaultMux()
	arcHandler, err := cyoa.ArcHandler(arcs, mux)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", arcHandler)
}
