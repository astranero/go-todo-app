package main

import (
	"net/http"
	"log"
	"os"
	"time"
	"sync"
	"io"
	"fmt"
)

var (
	mutex sync.Mutex
)

func main(){
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	imagePath := "/usr/src/shared/files/picsum.png"
	imageUrl := "https://picsum.photos/1200"
	err := downloadFile(imagePath, imageUrl)
	if err != nil {
		log.Fatalf("Failed to download image.")
	}

	ticker := time.NewTicker(3600 * time.Second)
	go func(){
		for range ticker.C {
			mutex.Lock()
			err := downloadFile(imagePath, imageUrl)
			if err != nil {
				log.Fatalf("Failed to download image.")
			}
			mutex.Unlock()
		}
	}()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)	

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w,r,imagePath)
	})

	log.Printf("Server started in port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func downloadFile(filepath string, url string)(err error){
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}