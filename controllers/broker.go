package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mwrks/ticket-app-broker/initializers"
	"github.com/mwrks/ticket-app-broker/models"
	"github.com/streadway/amqp"
)

// Publish
func PublishToQueue(order models.Order) error {
	// Connect to AMQP
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Declare channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare queue
	q, err := ch.QueueDeclare(
		"orders", // queue name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	// Encode order to JSON
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Publish the order message to AMQP server
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

// Consume
func ConsumeFromQueue() {
	// Connect to AMQP
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Declare channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue to ensure it exists before we consume from it
	q, err := ch.QueueDeclare(
		"orders", // Queue name should match the one used in the producer
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // additional arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer tag
		false,  // Auto-acknowledge
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Additional arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Processing messages
	for d := range msgs {
		var order models.Order
		log.Printf("Received a message: %s", d.Body)

		if err := json.Unmarshal(d.Body, &order); err != nil {
			log.Printf("Error decoding message: %v", err)
			d.Nack(false, true) // Negative acknowledge, requeue the message
			continue
		}

		// Process the order
		processOrder(order)

		// Acknowledge the message after successfully processing
		d.Ack(false)
	}
}

// Process the order
func processOrder(order models.Order) {
	var ticket models.Ticket

	// Retrieve the associated ticket from the database
	if err := initializers.DB.First(&ticket, order.TicketID).Error; err != nil {
		log.Printf("Ticket not found: %v", err)
		return
	}

	// Check if enough tickets are available
	if ticket.CurrentQuantity < 1 {
		log.Printf("No stock available for ticket ID %d", order.TicketID)
		return
	}

	// Decrement ticket quantity
	ticket.CurrentQuantity--
	if err := initializers.DB.Save(&ticket).Error; err != nil {
		log.Printf("Failed to update stock: %v", err)
		return
	}

	// Create the order in the database
	if result := initializers.DB.Create(&order); result.Error != nil {
		log.Printf("Failed to create order: %v", result.Error)
		return
	}

	log.Printf("Order created successfully for ticket ID %d", order.TicketID)
}

// Broker Handler
func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the order to JSON
	// jsonData, err := json.Marshal(order)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process order"})
	// 	return
	// }

	// Publish to RabbitMQ
	if err := PublishToQueue(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish order"})
		return
	}

	// Respond immediately to the client
	c.JSON(http.StatusAccepted, gin.H{"message": "Order received for processing"})
}
