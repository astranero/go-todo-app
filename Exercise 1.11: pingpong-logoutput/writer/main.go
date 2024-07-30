package main


import (
	"github.com/google/uuid"
	"fmt"
	"log"
	"time"
	"sync"
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
			_, err := appendToFile("/usr/src/shared/files/log.txt", timestamp, randomString)
			if err != nil {
				log.Fatal(err)
				fmt.Print("Error: %s", err)
				return
			}
			mutex.Unlock()
			fmt.Printf("%s: %s\n", timestamp , randomString)
		}
	}()

	select {}
}


func appendToFile(filename, timestamp, randomString string) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	msg := []byte(fmt.Sprintf("%s: %s\n", timestamp, randomString))
	_, err = file.Write(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}