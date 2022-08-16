package rabbitMQ

import (
	"context"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Publish(queueName string, message models.QueueMessage) error {

	msg, err := json.Marshal(message)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to marshal queue message to json: %v", err))
	}

	err = r.conn.Channel.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msg,
		},
	)

	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot publish message to the queue: %v", err))
	}

	log.Printf("message sent to: %v", queueName)
	return nil
}
