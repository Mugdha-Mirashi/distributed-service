package main

import (
	"distributed-counter-system/constants"
	handlers "distributed-counter-system/handler"
	"distributed-counter-system/node"

	"flag"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define flags
	port := flag.String("port", "8080", "Port to run this node on")
	peersFlag := flag.String("peers", "", "Comma-separated list of peers (e.g., localhost:8081,localhost:8082)")
	flag.Parse()

	selfID := "localhost:" + *port

	// Split peer list
	var peers []string
	if *peersFlag != "" {
		peers = strings.Split(*peersFlag, ",")
	}

	nodeService := node.NewNodeService(selfID, peers)

	// Create a new node service
	// nodeService := node.NewNodeService("localhost:8080", []string{})
	controller := handlers.NewController(nodeService)

	nodeService.StartHeartbeats()

	router.POST(constants.JoinPath, controller.HandleJoin)
	router.POST(constants.IncrementPath, controller.HandleIncrement)
	router.GET(constants.CountPath, controller.HandleGetCount)
	router.GET(constants.PingPath, handlers.HandlePing)
	// Run the server
	router.Run(":" + *port)
	// router.Run(":8080")
}
