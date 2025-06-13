package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbit, err := connect()
	if err != nil {
		log.Println("can not connect yet", err)
		os.Exit(1)
	}
	defer rabbit.Close()
	log.Println("connected")
	// start listening to messages

	// create consumer

	// watch the queue and consume events
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
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Trying again...")
		time.Sleep(backoff)
		continue
	}

	return conn, nil
}
