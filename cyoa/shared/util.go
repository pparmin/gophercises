package cyoa

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	NextArc string `json:"arc"`
}

func GetArcs(file *os.File) map[string]Arc {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	var arcs map[string]Arc
	json.Unmarshal([]byte(content), &arcs)
	return arcs
}
