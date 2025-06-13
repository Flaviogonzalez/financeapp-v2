package event

import (
	"encoding/json"
	"log"

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
	// y aca dependiendo q tenga en name ejecuto X función
	}
}

type Payload struct {
	Name string
	Data string
}
