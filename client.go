package main

import (
	"log"

	"github.com/streadway/amqp"
)

func Client(result chan string) {
	ch, conn, err := DeclareQueue(QUEUE_NAME, false, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
		result <- err.Error()
		return
	}
	defer ch.Close()
	defer conn.Close()

	// If we can consume the message, then we know it published successfully.
	msgs, err := ch.Consume(
		REPLY_TO_QUEUE_NAME, // queue
		"ReplyToConsumer",   // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		log.Printf("consume fail")
	}

	// clientMessage must read "client messaage {tmestamp}"
	clientMessage := "message from client " + GetTimestamp()

	err = ch.Publish(
		"",         // exchange
		QUEUE_NAME, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: GetCorrelationId(),
			Body:          []byte(clientMessage),
			ReplyTo:       REPLY_TO_QUEUE_NAME,
		})
	if err != nil {
		log.Printf("publish fail")
	}

	for m := range msgs {
		result <- string(m.Body)
	}
}
