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

// type Arc struct {
// 	Title   string   `json:"title"`
// 	Story   []string `json:"story"`
// 	Options []Option `json:"options"`
// }

// type Option struct {
// 	Text    string `json:"text"`
// 	NextArc string `json:"arc"`
// }

// func getArcs(file *os.File) map[string]Arc {
// 	content, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var arcs map[string]Arc
// 	json.Unmarshal([]byte(content), &arcs)
// 	return arcs
// }

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

		// OPTIONS NOT PROPERLY PARSED YET
		fmt.Println("Options: ", arc.Options)
		fmt.Println()
	}
}
