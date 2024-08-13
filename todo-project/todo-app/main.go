package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	mutex sync.Mutex
)

func main() {

	router := gin.Default()

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

	if _, err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist) {
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

	router.GET("/image", func(c *gin.Context) {
		c.File(imagePath)
	})

	router.POST("/submit", func(c *gin.Context) {
		todo := c.PostForm("todo")

		if todo == "" {
			log.Printf("Todo cannot be empty.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo cannot be empty"})
			return
		}

		requestURL := fmt.Sprintf("http://todo-backend:%s", backendPort)
		resp, err := http.Post(requestURL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("todo=%s", todo)))
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error sending request to %s: %v", requestURL, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process your request"})
			return
		}

		defer resp.Body.Close()

		todoBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response from %s: %v", requestURL, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error reading response: %v", err)})
			return
		}

		log.Printf("Received submission: Todo=%s", todo)
		c.Data(http.StatusOK, "application/json", todoBody)
	})

	router.PUT("/todos/:id", func(c *gin.Context) {
		todo := c.PostForm("todo")

		if todo == "" {
			log.Printf("Todo cannot be empty.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo cannot be empty"})
			return
		}

		requestURL := fmt.Sprintf("http://todo-backend:%s", backendPort)
		req, err := http.NewRequest(http.MethodPut, requestURL, strings.NewReader(fmt.Sprintf("todo=%s", todo)))
		if err != nil {
			log.Printf("Error creating PUT request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating request: %v", err)})
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error sending request to %s: %v", requestURL, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error sending request: %v", err)})
			return
		}
		defer resp.Body.Close()

		todoBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response from %s: %v", requestURL, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error reading response: %v", err)})
			return
		}

		log.Printf("Updated the submission: Todo=%s", todo)
		c.Data(http.StatusOK, "application/json", todoBody)
	})

	router.GET("/todos", func(c *gin.Context) {
		requestURL := fmt.Sprintf("http://todo-backend:%s", backendPort)
		resp, err := http.Get(requestURL)
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error sending request to %s: %v", requestURL, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error sending request: %v", err)})
			return
		}
		defer resp.Body.Close()

		todoBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response from %s: %v", requestURL, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error reading response: %v", err)})
			return
		}

		c.Data(http.StatusOK, "application/json", todoBody)
	})

	router.Static("/home", "./static")

	log.Printf("Server started on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func downloadFile(filepath string, url string) (err error) {
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
