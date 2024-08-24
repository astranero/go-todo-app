package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL is not set")
	}

	discordURL := os.Getenv("DISCORD_URL")
	if discordURL == "" {
		log.Fatal("DISCORD_URL is not set")
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	_, err = nc.QueueSubscribe("todos", "workers", func(m *nats.Msg) {
		log.Printf("Message received: %s\n", string(m.Data))

		payload := map[string]string{"content": string(m.Data)}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			return
		}

		req, err := http.NewRequest("POST", discordURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request to Discord: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Discord server error: %d", resp.StatusCode)
			return
		}
		log.Println("Message successfully sent to Discord")
	})

	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	select {}
}
