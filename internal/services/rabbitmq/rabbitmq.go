package rabbitmq

import (
	"errors"
	"fmt"

	"github.com/adetunjii/auth-svc/internal/port"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	name     string
	dialUrl  string
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queues   []string
	err      chan error
	logger   port.AppLogger
}

func NewConnection(name string, exchange string, queues []string, logger port.AppLogger, dialUrl string) (*Connection, error) {

	if dialUrl == "" {
		return nil, errors.New("invalid dial url")
	}

	connection := &Connection{
		name:     name,
		dialUrl:  dialUrl,
		exchange: exchange,
		queues:   queues,
		err:      make(chan error),
		logger:   logger,
	}

	if err := connection.Connect(); err != nil {
		return nil, err
	}

	if err := connection.BindQueue(); err != nil {
		return nil, err
	}

	return connection, nil
}

func (c *Connection) Connect() error {
	var err error
	c.conn, err = amqp.Dial(c.dialUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %v", err)
	}

	// listening for a close notification on the connection
	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error))
		c.err <- errors.New("connection closed")
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create a rabbitmq channel: %v", err)
	}

	// if err := c.channel.ExchangeDeclare(
	// 	c.exchange,
	// 	"direct",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// ); err != nil {
	// 	return fmt.Errorf("failed to declare an exchange: %v", err)
	// }

	c.logger.Info(fmt.Sprintf("RabbitMq connected successfully on %v", c.dialUrl))

	return nil
}

func (c *Connection) BindQueue() error {
	for _, queue := range c.queues {

		_, err := c.channel.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			return errors.New("failed to declare queue")
		}

		// if err := c.channel.QueueBind(queue, "my_routing_key", c.exchange, false, nil); err != nil {
		// 	return fmt.Errorf("queue  Bind error: %s", err)
		// }
	}

	return nil
}

func (c *Connection) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}

	if err := c.BindQueue(); err != nil {
		return err
	}

	return nil
}

// check if a queue exists in the rabbitmq connection
func (c *Connection) IsValidQueue(queueName string) bool {
	for _, queue := range c.queues {
		if queue == queueName {
			return true
		}
	}

	return false
}
