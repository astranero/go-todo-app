package main

import (
	"net/http"
	"fmt"
	"log"
	"os"
)

func main(){
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server started in port %s", port)
	})

	log.Printf("Server started in port %s", port)
	if err := http.ListenAndServe(":"+port,nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}