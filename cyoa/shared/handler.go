package cyoa

import (
	"fmt"
	"net/http"
	"strings"
)

func redirectHandler() http.Handler {
	rh := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/intro", http.StatusPermanentRedirect)
	}
	return http.HandlerFunc(rh)
}

func nameHandler(arcs map[string]Arc, path string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimLeft(path, "/")
		if path == r.URL.Path {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "<h1>%s</h1>", arcs[name].Title)
			for _, text := range arcs[name].Story {
				fmt.Fprintf(w, "<p>%s</p>", text)
			}
			fmt.Fprintf(w, "<hr>")
			if path == "/home" {
				fmt.Fprintf(w, "<h2>Congratulations. You made it home! ")
				fmt.Fprintf(w, "<a href=/intro>Back to the start</a>")
			} else {
				fmt.Fprintf(w, "<h2>How do you want to proceed?</h2>")
			}
			for _, option := range arcs[name].Options {
				pathToNext := "/" + option.NextArc
				fmt.Fprintf(w, "<li><a href=%s>%s</a></li>", pathToNext, option.Text)
			}
		} else {
			w.Write([]byte("NO MATCHING PATH FOUND"))
		}
	}
	return http.HandlerFunc(fn)
}

func ArcHandler(arcs map[string]Arc) (http.HandlerFunc, error) {
	mux := http.NewServeMux()
	for name := range arcs {
		path := "/" + name
		fmt.Println("NEW PATH: ", path)
		nh := nameHandler(arcs, path)
		mux.Handle(path, nh)
	}
	rh := redirectHandler()
	mux.HandleFunc("/", rh.ServeHTTP)
	fmt.Println()
	return mux.ServeHTTP, nil
}
