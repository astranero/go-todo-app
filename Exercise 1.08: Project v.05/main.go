package main

import (
	"net/http"
	"log"
	"os"
)

func main(){
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)	

	log.Printf("Server started in port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}