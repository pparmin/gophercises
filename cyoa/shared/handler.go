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
			w.Write([]byte("Hello, " + name + "\n"))

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
