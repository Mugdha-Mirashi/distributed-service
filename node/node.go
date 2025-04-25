package node

import (
	"bytes"
	"distributed-counter-system/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type NodeService struct {
	SelfID  string               // This node's ID (e.g., "localhost:8081")
	Peers   map[string]time.Time // Peers with last seen heartbeat time
	Mutex   sync.RWMutex         // Protects Peers map
	Counter *Counter
}

// NewNodeService creates a NodeService with initial known peers
func NewNodeService(selfID string, initialPeers []string) *NodeService {
	ns := &NodeService{
		SelfID:  selfID,
		Peers:   make(map[string]time.Time),
		Counter: NewCounter(),
	}

	now := time.Now()
	for _, peer := range initialPeers {
		if peer != selfID {
			ns.Peers[peer] = now
		}
	}

	return ns
}

// RegisterPeer adds or updates a peer's heartbeat timestamp
func (ns *NodeService) RegisterPeer(peerID string) {
	ns.Mutex.Lock()
	defer ns.Mutex.Unlock()
	if peerID != ns.SelfID {
		ns.Peers[peerID] = time.Now()
	}
}

// GetPeers returns a copy of current peer list
func (ns *NodeService) GetPeers() []string {
	ns.Mutex.RLock()
	defer ns.Mutex.RUnlock()

	peers := make([]string, 0, len(ns.Peers))
	for peer := range ns.Peers {
		peers = append(peers, peer)
	}
	return peers
}

// RemovePeer deletes a peer from the list
func (ns *NodeService) RemovePeer(peerID string) {
	ns.Mutex.Lock()
	defer ns.Mutex.Unlock()
	delete(ns.Peers, peerID)
}

func (ns *NodeService) UpdatePeerTimestamp(peer string) {
	ns.Mutex.Lock()
	defer ns.Mutex.Unlock()
	ns.Peers[peer] = time.Now()
}



// PropagateIncrement sends an increment message to all known peers
func (ns *NodeService) PropagateIncrement(id string) {
	message := models.IncrementMessage{
		ID:        id,
		Sender:    ns.SelfID,
		Timestamp: time.Now(),
	}

	body, _ := json.Marshal(message)

	ns.Mutex.RLock()
	peers := ns.GetPeers()
	ns.Mutex.RUnlock()

	for _, peer := range peers {
		go func(peer string) {
			url := fmt.Sprintf("http://%s/propagate", peer)
			resp, err := http.Post(url, "application/json", bytes.NewReader(body))
			if err != nil {
				fmt.Printf("❌ Failed to propagate increment to %s: %v\n", peer, err)
				return
			}
			defer resp.Body.Close()
		}(peer)
	}
}
