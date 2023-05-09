package main

import (
	"log"

	"github.com/streadway/amqp"
)

func Server(result chan string) {
	ch, conn, err := DeclareQueue(QUEUE_NAME, false, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
		result <- err.Error()
		return
	}

	defer ch.Close()
	defer conn.Close()

	// If we can consume the message, then we know it published successfully.
	log.Println("Server is running...")
	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	msgs, err := ch.Consume(
		QUEUE_NAME,        // queue
		"ReplyToConsumer", // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		log.Printf("consume fail")
		log.Println(err.Error())
	}

	// reply message is a timestamp
	replyMessage := "server response " + GetTimestamp()

	for m := range msgs {

		err = ch.Publish(
			"",        // exchange
			m.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: GetCorrelationId(),
				Body:          []byte(replyMessage),
			})
		if err != nil {
			log.Printf("publish fail")
		}

		err = ch.Ack(m.DeliveryTag, false)
		if err != nil {
			log.Printf("ack fail")
		}

		result <- string(m.Body)
	}
}
