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
    conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
    if err != nil {
        return err
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()

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

    body, err := json.Marshal(order)
    if err != nil {
        return err
    }

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

	msgs, err := ch.Consume(
		"orders", // Queue name should match the one used in the producer
		"",       // Consumer tag
		true,     // Auto-acknowledge messages
		false,    // Exclusive
		false,    // No-local
		false,    // No-wait
		nil,      // Additional arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		var order models.Order
		if err := json.Unmarshal(d.Body, &order); err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}

		processOrder(order)
	}
}

func processOrder(order models.Order) {
	var ticket models.Ticket
	if err := initializers.DB.First(&ticket, order.TicketID).Error; err != nil {
		log.Printf("Ticket not found: %v", err)
		return
	}

	if ticket.CurrentQuantity < 1 {
		log.Printf("No stock available for ticket ID %d", order.TicketID)
		return
	}

	ticket.CurrentQuantity--
	if err := initializers.DB.Save(&ticket).Error; err != nil {
		log.Printf("Failed to update stock: %v", err)
		return
	}

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
	if 	err := PublishToQueue(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish order"})
		return
	}

	// Respond immediately to the client
	c.JSON(http.StatusAccepted, gin.H{"message": "Order received for processing"})
}