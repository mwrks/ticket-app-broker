package initializers

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func InitQueue() {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()
}
