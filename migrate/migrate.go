package main

import (
	"fmt"
	"github.com/mwrks/ticket-app-broker/initializers"
	"github.com/mwrks/ticket-app-broker/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
}

func main() {
	err := initializers.DB.AutoMigrate(&models.Ticket{}, &models.Order{})
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}
}
