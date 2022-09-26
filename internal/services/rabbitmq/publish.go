package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/adetunjii/auth-svc/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Connection) Publish(queueName string, message model.Notification) error {

	if ok := c.IsValidQueue(queueName); !ok {
		return errors.New("invalid queue name")
	}

	// check for connection error if none, it goes to the default which does nothing
	select {
	case err := <-c.err:
		if err != nil {
			c.Reconnect()
		}
	default:
	}

	msg, err := json.Marshal(message)
	if err != nil {
		c.logger.Error("failed to marshal queue message to json", err)
		return err
	}

	if err := c.channel.PublishWithContext(
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
	); err != nil {
		c.logger.Error("failed to publish message to the queue", err)
		return err
	}

	c.logger.Info(fmt.Sprintf("message sent to: %v", queueName))
	return nil
}
