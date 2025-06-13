package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	// Connect to RabbitMQ
	conn, err := connect()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		"asd",   // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	// Register handler
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Publish to RabbitMQ
		body, _ := json.Marshal(req)
		err = ch.Publish(
			"asd",      // exchange
			"register", // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Registration request received",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var conn *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			log.Println("Error trying to connect: ", err)
			counts++
		} else {
			conn = c
			break
		}

		if counts > 5 {
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Trying again...")
		time.Sleep(backoff)
		continue
	}

	return conn, nil
}
