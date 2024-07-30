package main


import (
	"fmt"
	"log"
	"sync"
	"net/http"
)


var (
	counter int
	mutex sync.Mutex
)


func main(){
	log.Println("Application Started")
	counter = 0

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/pingpong", pingHandler)

	port := "8081"
	log.Printf("Server started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


func homeHandler(w http.ResponseWriter, r *http.Request){
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Fprintf(w, "Ping / Pongs: %d", counter)
}

func pingHandler(w http.ResponseWriter, r *http.Request){
	mutex.Lock()
	defer mutex.Unlock()
	counter++
}
