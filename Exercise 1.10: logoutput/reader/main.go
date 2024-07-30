package main


import (
	"fmt"
	"log"
	"sync"
	"net/http"
	"os"
	"bufio"
)


var (
	randomString string
	timestamp string
	mutex sync.Mutex
)


func main(){
	log.Println("Application Started")

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
	data, err := reader("/usr/src/shared/files/log.txt")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(data))
}

func reader(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()
	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		lastLine = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastLine, nil
}
