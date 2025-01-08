package main

import (
	"github.com/mwrks/ticket-app-broker/controllers"
	"github.com/mwrks/ticket-app-broker/initializers"
	"github.com/mwrks/ticket-app-broker/routes"
)

func init() {
	initializers.LoadEnv()         // Load environment variables
	initializers.ConnectDatabase() // Connect to database
}

func main() {
	// Initialize router
	r := routes.SetupRouter()

	// Run consumer goroutines
	go controllers.ConsumeFromQueue()

	// Listen and serve on 0.0.0.0:8080
	r.Run()
}
