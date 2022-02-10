package cyoa

import (
	"fmt"
	"net/http"
)

func ArcHandler(arcs map[string]Arc, fallback http.Handler) (http.HandlerFunc, error) {
	// mux := http.NewServeMux()
	// mux.Handle("/", fallback)
	fmt.Println("I AM A HANDLER!")
	return fallback.ServeHTTP, nil
}
