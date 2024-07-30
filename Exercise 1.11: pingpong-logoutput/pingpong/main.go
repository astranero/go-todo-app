package main


import (
	"fmt"
	"log"
	"sync"
	"net/http"
	"os"
)


var (
	counter int
	mutex sync.Mutex
)


func main(){
	log.Println("Application Started")
	counter = 0

	_, err := overwriteToFile("/usr/src/shared/files/pong.txt", counter) 
	if err != nil {
		log.Printf("Failed to write to file: %v", err)
		return
	}

	http.HandleFunc("/", homeHandler)
	port := "8081"
	log.Printf("Server started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


func homeHandler(w http.ResponseWriter, r *http.Request){
	mutex.Lock()
	defer mutex.Unlock()
	_,err := overwriteToFile("/usr/src/shared/files/pong.txt", counter)
	if err != nil {
		log.Fatalf("Failure: %v", err)
		return
	}
	counter++
}


func overwriteToFile(filename string, counter int) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	msg := []byte(fmt.Sprintf("Ping / Pongs: %d", counter))
	_, err = file.Write(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}