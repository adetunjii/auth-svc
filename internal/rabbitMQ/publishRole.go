package rabbitMQ

func PublishToRoleQueue(userID string) error {
	//RabbitmqHost := os.Getenv("RABBITMQ_HOST")
	//RabbitmqPort := os.Getenv("RABBITMQ_PORT")
	//RabbitmqUser := os.Getenv("RABBITMQ_USER")
	//RabbitmqPass := os.Getenv("RABBITMQ_PASS")
	//rabbitMQURL := os.Getenv("CLOUDAMQP_URL")
	//
	//if rabbitMQURL == "" {
	//	rabbitMQURL = fmt.Sprintf("amqp://%s:%s@%s:%s/", RabbitmqUser, RabbitmqPass, RabbitmqHost, RabbitmqPort)
	//}
	//
	//conn, err := amqp.Dial(rabbitMQURL)
	//helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	//defer func(conn *amqp.Connection) {
	//	err := conn.Close()
	//	if err != nil {
	//		log.Printf("Error closing connection: %s", err)
	//	}
	//}(conn)
	//
	//ch, err := conn.Channel()
	//helpers.FailOnError(err, "Failed to open a channel")
	//defer func(ch *amqp.Channel) {
	//	err := ch.Close()
	//	if err != nil {
	//		log.Printf("Error closing channel: %s", err)
	//	}
	//}(ch)
	//
	//err = ch.ExchangeDeclare(
	//	"user-role", // name
	//	"fanout",    // type
	//	true,        // durable
	//	false,       // auto-deleted
	//	false,       // internal
	//	false,       // no-wait
	//	nil,         // arguments
	//)
	//helpers.FailOnError(err, "could not declare queue")
	//
	//body, err := json.Marshal(userID)
	//if err != nil {
	//	log.Printf("err :%v", err)
	//	return err
	//}
	//
	//err = ch.Publish(
	//	"user-role", // exchange
	//	"",          // routing key
	//	false,       // mandatory
	//	false,       // immediate
	//	amqp.Publishing{
	//		ContentType: "json",
	//		Body:        body,
	//	})
	//
	//log.Printf("Sent %s\n", body)
	//if err != nil {
	//	log.Fatalf("err: %v", err)
	//	return err
	//}
	//
	//log.Println("published role")
	return nil
}
