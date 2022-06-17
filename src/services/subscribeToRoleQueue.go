package services

import (
	"dh-backend-auth-sv/src/helpers"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func (s *Server) SubscribeToRoleQueue() {

	RabbitmqHost := os.Getenv("RABBITMQ_HOST")
	RabbitmqPort := os.Getenv("RABBITMQ_PORT")
	RabbitmqUser := os.Getenv("RABBITMQ_USER")
	RabbitmqPass := os.Getenv("RABBITMQ_PASS")
	rabbitMQURL := os.Getenv("CLOUDAMQP_URL")

	if rabbitMQURL == "" {
		rabbitMQURL = fmt.Sprintf("amqp://%s:%s@%s:%s/", RabbitmqUser, RabbitmqPass, RabbitmqHost, RabbitmqPort)
	}

	conn, err := amqp.Dial(rabbitMQURL)
	helpers.LogEvent("ERROR", fmt.Sprintf("%s", err))
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}(conn)

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	helpers.LogEvent("ERROR", fmt.Sprintf("%s", err))
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Printf("Error closing channel: %s", err)
		}
	}(ch)

	err = ch.ExchangeDeclare(
		"roles",  // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")
	helpers.LogEvent("ERROR", fmt.Sprintf("%s", err))

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,  // queue name
		"",      // routing key
		"roles", // exchange
		false,
		nil,
	)
	helpers.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	helpers.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("our consumer received a message from: %v", os.Getpid())
			log.Printf("Received a message: %s", d.Body)
			//var user []*models.UserRole
			//err := json.Unmarshal(d.Body, &user)

			if err != nil {
				log.Fatalf("error%v", err)
			}
			err = s.RedisCache.SaveRoleChannel("roles", d.Body)
			if err != nil {
				return
			}
			//log.Printf("%v", user)
			//log.Println(err)

		}
	}()

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
