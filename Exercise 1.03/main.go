package main

import (
	"github.com/google/uuid"
	"fmt"
	"log"
	"time"
)

func main(){
	randomString := uuid.New().String()
	log.Println("Application Started")

	ticker := time.NewTicker(5 * time.Second)
	for t := range ticker.C {
		fmt.Printf("%s: %s\n", t.Format(time.RFC3339), randomString)
	}
}