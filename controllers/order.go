package controllers

import (
	"net/http"

	"github.com/mwrks/ticket-app-broker/initializers"
	"github.com/mwrks/ticket-app-broker/models"

	"github.com/gin-gonic/gin"
)

// func CreateOrder(c *gin.Context) {
// 	var order models.Order
// 	if err := c.ShouldBindJSON(&order); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var ticket models.Ticket
// 	if err := initializers.DB.First(&ticket, order.TicketID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
// 		return
// 	}

// 	if ticket.CurrentQuantity < 1 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough tickets available"})
// 		return
// 	}

// 	ticket.CurrentQuantity--
// 	initializers.DB.Save(&ticket)

// 	if result := initializers.DB.Create(&order); result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, order)
// }

func GetOrders(c *gin.Context) {
	var orders []models.Order
	if result := initializers.DB.Unscoped().Find(&orders); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func GetOrderByID(c *gin.Context) {
	ticket_id := c.Param("id")
	var order []models.Order
	if result := initializers.DB.Where("ticket_id = ?", ticket_id).Find(&order); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if len(order) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.Order{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func DeleteOrdersByTicketID(c *gin.Context) {
	ticketID := c.Param("id")

	result := initializers.DB.Unscoped().Where("ticket_id = ?", ticketID).Delete(&models.Order{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete orders"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No orders found for the given ticket ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Orders deleted successfully"})

}

func ResetOrderSequence(c *gin.Context) {
	// Execute the SQL query to reset the sequence
	query := "ALTER SEQUENCE orders_order_id_seq RESTART WITH 1;"
	if err := initializers.DB.Exec(query).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset order sequence"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order sequence reset successfully"})
}
