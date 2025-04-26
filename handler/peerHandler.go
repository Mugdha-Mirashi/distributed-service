package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"distributed-counter-system/models"
	"distributed-counter-system/node"
)

// adding instance of node in the handler struct

type Controller struct {
	ns *node.NodeService
}

func NewController(ns *node.NodeService) *Controller {
	return &Controller{
		ns: ns,
	}
}

// HandleJoin processes a join request from a new peer
func (handler *Controller) HandleJoin(c *gin.Context) {
	var request models.JoinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid join request"})
		return
	}

	handler.ns.RegisterPeer(request.Sender)

	handler.ns.MergePeers(request.Peers)

	peerList := handler.ns.GetPeers()
	c.JSON(http.StatusOK, models.PeerListResponse{
		Peers: peerList,
	})

        // Sync with the new peer asynchronously
    go handler.ns.SyncPeersFrom(request.Sender)
        
	// Notify other peers about the new peer asynchronously
	go handler.ns.NotifyAllPeersAboutNewPeer(request.Sender)

}

// HandleIncrement handles local increment and propagates it to peers
func (handler *Controller) HandleIncrement(c *gin.Context) {
	id := uuid.New().String()

	// Try to increment the counter
	applied := handler.ns.Counter.Increment(id)
	if !applied {
		c.JSON(http.StatusOK, gin.H{"status": "duplicate", "id": id})
		return
	}

	fmt.Printf("Increment applied from peer: %s\n", id, handler.ns.Counter.Value)

	// If applied, propagate to peers
	go handler.ns.PropagateIncrement(id)

	c.JSON(http.StatusOK, gin.H{"status": "incremented", "id": id})
}

// HandleGetCount handles GET /count
func (handler *Controller) HandleGetCount(c *gin.Context) {
	count := handler.ns.Counter.Get()
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (handler *Controller) HandlePropagateIncrement(c *gin.Context) {
	var req struct {
		ID string `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	applied := handler.ns.Counter.Increment(req.ID)
	if applied {
		fmt.Printf("Increment applied from peer: %s\n", req.ID)
	} else {
		fmt.Printf("Duplicate increment received from peer: %s\n", req.ID)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
