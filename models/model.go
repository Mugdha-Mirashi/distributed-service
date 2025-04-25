package models

import "time"

// Represents a counter increment propagated to peers
type IncrementMessage struct {
	ID        string    `json:"id"`        // Unique UUID for deduplication
	Sender    string    `json:"sender"`    // Node ID that originally incremented
	Timestamp time.Time `json:"timestamp"` // When this was generated
}

// Represents a new node joining the cluster
type JoinRequest struct {
	Sender string   `json:"sender"` // The node ID joining (host:port)
	Peers  []string `json:"peers"`  // Initial peer list (optional)
}

// Represents a node list returned by /peers or after a join
type PeerListResponse struct {
	Peers []string `json:"peers"`
}

// Response format for /count
type CountResponse struct {
	Count int `json:"count"`
}


