package event

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn        *amqp.Connection
	QueueString string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}
	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"asd",
			false,
			nil,
		)
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload

			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	log.Printf("waiting for messages on [Exchange, queue] [asd, %s]", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "register":
		handleRegistration(payload.Data)
	}
}

func handleRegistration(data string) {
	// Make HTTP request to auth service
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://auth-service:8080/register", strings.NewReader(data))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to auth service: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read and log the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}
	log.Printf("Auth service response: %s", string(body))
}

type Payload struct {
	Name string
	Data string
}
