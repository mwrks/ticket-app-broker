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

	_, err = ch.QueueDeclare(
		"orders", // Queue name should be the same as in consumer
		true,     // Durable (survives broker restarts)
		false,    // Auto-delete when unused
		false,    // Exclusive to this connection
		false,    // No-wait
		nil,      // Additional arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	log.Println("Queue created successfully: order-queue")
}
