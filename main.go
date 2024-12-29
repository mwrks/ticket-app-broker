package main

import (
	"github.com/mwrks/ticket-app-broker/controllers"
	"github.com/mwrks/ticket-app-broker/initializers"
	"github.com/mwrks/ticket-app-broker/routes"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
	initializers.InitQueue()
}
func main() {
	r := routes.SetupRouter()
	go controllers.ConsumeFromQueue()
	r.Run() // listen and serve on 0.0.0.0:8080
}
