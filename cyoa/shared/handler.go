package cyoa

import (
	"fmt"
	"html/template"
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
			t, err := template.ParseFiles("templates/arc.gohtml")
			if err != nil {
				panic(err)
			}

			err = t.Execute(w, arcs[name])
			if err != nil {
				panic(err)
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
