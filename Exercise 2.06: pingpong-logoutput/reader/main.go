package main


import (
	"fmt"
	"log"
	"sync"
	"net/http"
	"time"
	"github.com/google/uuid"
	"io"
	"bufio"
	"os"
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
			fmt.Printf("%s: %s\n", timestamp , randomString)
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
	
	log := fmt.Sprintf("%s: %s\n", timestamp , randomString)

	requestURL := fmt.Sprintf("http://pong-svc:%d", 8081)
	resp, err := http.Get(requestURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making HTTP request: %s", err), http.StatusInternalServerError)
		return
	}

	pong, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading HTTP response: %s", err), http.StatusInternalServerError)
		return
	}

	fileContent, err := reader("/config/information.txt")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading HTTP response: %s", err), http.StatusInternalServerError)
		return
	}
	messageEnvVar := os.Getenv("MESSAGE")
	if messageEnvVar == "" {
		http.Error(w, fmt.Sprintf("Error reading HTTP response: %s", err), http.StatusInternalServerError)
	}


	fmt.Fprintln(w, "file content: "+string(fileContent)+"\n"+"env variable: MESSAGE=" + string(messageEnvVar) + "\n" + string(log)+"\n"+string(pong))
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