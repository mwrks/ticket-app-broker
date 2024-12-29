package controllers

import (
	"net/http"
	"strconv"
	"github.com/mwrks/ticket-app-broker/initializers"
	"github.com/mwrks/ticket-app-broker/models"

	"github.com/gin-gonic/gin"
)

func CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := initializers.DB.Create(&ticket); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func GetTickets(c *gin.Context) {
	var tickets []models.Ticket
	if result := initializers.DB.Find(&tickets); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func GetTicketByID(c *gin.Context) {
	id := c.Param("id")
	var ticket models.Ticket
	if result := initializers.DB.First(&ticket, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func UpdateTicket(c *gin.Context) {
	ticketIDstr := c.Param("id")

	ticketID, err := strconv.ParseUint(ticketIDstr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket not found"})
	}
	var ticket models.Ticket
	if err := initializers.DB.First(&ticket, ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket.TicketID = uint(ticketID)

	if err := initializers.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func DeleteTicket(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.Ticket{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}
