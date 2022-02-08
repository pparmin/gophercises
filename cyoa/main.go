/* Interesting discussion of two methods to handle more complicated
JSON files: https://tutorialedge.net/golang/parsing-json-with-golang/

In this case we will use map[string]interface, since the data is
unstructured and we will not necessarily now the names of each arc

Another very interesting post about how to do this properly and then
get proper types for each field: https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data

*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options struct {
		Text    string `json:"text"`
		NextArc string `json:"arc"`
	} `json:"options"`
}

func main() {
	file, err := os.Open("gopher.json")
	if err != nil {
		fmt.Println("ERROR DURING OPENING OF FILE")
		log.Fatal(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var arcs map[string]Arc
	json.Unmarshal([]byte(content), &arcs)

	// dec := json.NewDecoder(file)
	// err = dec.Decode(&arcs)
	// if err != nil {
	// 	fmt.Println("ERROR DURING DECODING OF JSON: ")
	// 	log.Fatal(err)
	// }

	// Need to use type assertions to properly decode all fields into the proper type
	// --> WRITE FUNCTION THAT PROPERLY PARSES ALL ARCS AND RETURNS AN ARRAY OF ARCS
	for name, arc := range arcs {
		fmt.Println("ARC NAME: ", name)
		fmt.Println("title: ", arc.Title)
		fmt.Println("story: ", arc.Story)

		// OPTIONS NOT PROPERLY PARSED YET
		fmt.Println("Options text: ", arc.Options.Text)
		fmt.Println("Next arc: ", arc.Options.NextArc)
	}
}

// fmt.Println()
// }
