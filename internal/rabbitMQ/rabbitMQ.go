package rabbitMQ

import (
	"dh-backend-auth-sv/internal/helpers"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type client struct{}

type connection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

type RabbitMQ struct {
	conn connection
}

func New(dialUrl string) *RabbitMQ {
	client := client{}
	connection := client.Connect(dialUrl)
	return &RabbitMQ{*connection}
}

func (c *client) Connect(dialUrl string) *connection {
	conn, err := amqp.Dial(dialUrl)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to connect to rabbitMQ: %v", err))
		log.Fatalf("failed to connect to rabbitMQ %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to open a channel: %v", err))
		log.Fatalf("failed to open a channel %v", err)
	}

	log.Println("successfully connected to rabbitMQ")
	return &connection{
		Connection: conn,
		Channel:    channel,
	}
}
