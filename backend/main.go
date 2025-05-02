package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	name := os.Getenv("BACKEND_NAME")

	if port == "" {
		port = "9000"
	}
	if name == "" {
		name = "default"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %s on port %s\n", name, port)
	})

	fmt.Printf("Starting %s on port %s\n", name, port)
	http.ListenAndServe(":"+port, nil)
}
