package main

import (
	handlers "distributed-counter-system/handler"
	"distributed-counter-system/node"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Create a new node service
	nodeService := node.NewNodeService("localhost:8081", []string{"localhost:8080"})
	controller := handlers.NewController(nodeService)

	nodeService.StartHeartbeats() 

	router.POST("/join", controller.HandleJoin)
	router.POST("/increment", controller.HandleIncrement)
	router.GET("/count", controller.HandleGetCount)
	router.GET("/ping", handlers.HandlePing)
	// Run the server
	router.Run(":8081")
}
