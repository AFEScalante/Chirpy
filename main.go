package main

import (
	"fmt"
	"net/http"
)

type Server struct{}

func main() {
	port := "8080"
	filepathRoot := "."

	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/", handler)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	fmt.Println("Server running on http://localhost:8080")
	server.ListenAndServe()
}
