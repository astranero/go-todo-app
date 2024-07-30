package main


import (
	"github.com/google/uuid"
	"fmt"
	"log"
	"time"
	"sync"
	"net/http"
)


var (
	randomString string
	timestamp string
	mutex sync.Mutex
)


func main(){
	log.Println("Application Started")

	ticker := time.NewTicker(5 * time.Second)
	go func(){ 
		for t := range ticker.C {
			mutex.Lock()
			timestamp = t.Format(time.RFC3339)
			randomString = uuid.New().String()
			mutex.Unlock()
			fmt.Printf("%s: %s \n ", timestamp , randomString)
		}
	}()

	http.HandleFunc("/", homeHandler)
	port := "8080"
	log.Printf("Server started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


func homeHandler(w http.ResponseWriter, r *http.Request){
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Fprintf(w, "Timestamp: %s \nRandom String: %s", timestamp, randomString)
}