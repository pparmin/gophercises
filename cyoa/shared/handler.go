package cyoa

import (
	"fmt"
	"net/http"
	"strings"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, test")
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
			fmt.Fprintf(w, "<h2>How do you want to proceed?</h2>")
			for _, option := range arcs[name].Options {
				pathToNext := "/" + option.NextArc
				fmt.Fprintf(w, "<li><a href=%s>%s</li>", pathToNext, option.Text)
			}
		} else {
			w.Write([]byte("NO MATCHING PATH FOUND"))
		}
	}
	return http.HandlerFunc(fn)
}

func ArcHandler(arcs map[string]Arc, fallback http.Handler) (http.HandlerFunc, error) {
	mux := http.NewServeMux()
	for name := range arcs {
		path := "/" + name
		fmt.Println("NEW PATH: ", path)
		nh := nameHandler(arcs, path)
		mux.Handle(path, nh)
	}
	mux.HandleFunc("/test", test)
	mux.HandleFunc("/", fallback.ServeHTTP)
	fmt.Println()
	return mux.ServeHTTP, nil
}

func TestHandler(fallback http.Handler) (http.HandlerFunc, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/intro", test)
	mux.HandleFunc("/", fallback.ServeHTTP)
	return mux.ServeHTTP, nil
}
