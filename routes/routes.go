package routes

import (
	"github.com/mwrks/ticket-app-broker/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	ticketRoutes := r.Group("/ticket")
	{
		ticketRoutes.POST("/", controllers.CreateTicket)
		ticketRoutes.GET("/", controllers.GetTickets)
		ticketRoutes.GET("/:id", controllers.GetTicketByID)
		ticketRoutes.PUT("/:id", controllers.UpdateTicket)
		ticketRoutes.DELETE("/:id", controllers.DeleteTicket)
	}

	orderRoutes := r.Group("/order")
	{
		orderRoutes.POST("/", controllers.CreateOrder)
		orderRoutes.POST("/:id/reset", controllers.DeleteOrdersByTicketID)
		orderRoutes.GET("/", controllers.GetOrders)
		orderRoutes.GET("/:id", controllers.GetOrderByID)
		orderRoutes.DELETE("/:id", controllers.DeleteOrder)
	}
	return r
}
