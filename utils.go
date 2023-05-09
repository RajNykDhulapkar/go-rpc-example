package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	// AMQP URI
	AMQP_URI            = "amqp://guest:guest@localhost:5672/"
	QUEUE_NAME          = "rpc_queue"
	REPLY_TO_QUEUE_NAME = "amq.rabbitmq.reply-to"
)

func DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (*amqp.Channel, *amqp.Connection, error) {
	conn, err := amqp.Dial(AMQP_URI)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, errors.New("error getting channel")
	}

	_, err = ch.QueueDeclare(
		name,       // name
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		noWait,     // noWait
		args,       // arguments
	)

	return ch, conn, err
}

func GetCorrelationId() string {
	return uuid.New().String()
}

func GetTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
