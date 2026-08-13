package main

import (
	"fmt"
	"net/http"
)

type Server struct{}

func main() {
	port := "8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "index.html")
	})

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	fmt.Println("Server running on http://localhost:8080")
	server.ListenAndServe()
}
