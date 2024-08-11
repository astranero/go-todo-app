package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"errors"

	"github.com/joho/godotenv"
)

var (
	mutex sync.Mutex
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "8080"
	}

	backendPort := os.Getenv("BACKEND_PORT")
	if backendPort == "" {
		backendPort = "8081"
	}

	imagePath := os.Getenv("IMAGE_PATH")
	imageURL := os.Getenv("IMAGE_URL")

	if _ , err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist){
		err = downloadFile(imagePath, imageURL)
		if err != nil {
			log.Fatalf("Failed to download image.")
		}
	}

	ticker := time.NewTicker(3600 * time.Second)
	go func() {
		for range ticker.C {
			mutex.Lock()
			err := downloadFile(imagePath, imageURL)
			mutex.Unlock()
			if err != nil {
				log.Printf("Failed to download image: %v", err)
			}
		}
	}()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w,r,imagePath)
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Invalid request method: %s", r.Method)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		todo := r.FormValue("todo")
		if todo == "" {
			log.Printf("Todo cannot be empty.")
			http.Error(w, "Todo cannot be empty", http.StatusBadRequest)
			return
		}

		requestURL := fmt.Sprintf("http://todo-backend:%s", backendPort)
		resp, err := http.Post(requestURL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("todo=%s", todo)))
		if err != nil {
			log.Printf("Error sending request to %s: %v", requestURL, err)
			http.Error(w, fmt.Sprintf("Error sending request: %v", err), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		todoBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response from %s: %v", requestURL, err)
			http.Error(w, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Received submission: Todo=%s", todo)
		w.Header().Set("Content-Type", "application/json")
		w.Write(todoBody)
	})

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		requestURL := fmt.Sprintf("http://todo-backend:%s", backendPort)
		resp, err := http.Get(requestURL)
		if err != nil {
			log.Printf("Error sending request to %s: %v", requestURL, err)
			http.Error(w, fmt.Sprintf("Error sending request: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		todoBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response from %s: %v", requestURL, err)
			http.Error(w, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(todoBody)
	})

	log.Printf("Server started on port %s", port)
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