package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"distributed-counter-system/models"
)

// JoinCluster notifies all peers about this node and merges their known peers
func (ns *NodeService) JoinCluster() {
	ns.Mutex.RLock()
	peers := ns.GetPeers()
	ns.Mutex.RUnlock()

	for _, peer := range peers {
		go func(peer string) {
			joinBody := models.JoinRequest{
				Sender: ns.SelfID,
				Peers:  peers, // send your current known list
			}

			body, _ := json.Marshal(joinBody)
			url := fmt.Sprintf("http://%s/join", peer)

			resp, err := http.Post(url, "application/json", bytes.NewReader(body))
			if err != nil {
				fmt.Printf(" Failed to join peer %s: %v\n", peer, err)
				return
			}
			defer resp.Body.Close()

			// Expecting PeerListResponse in return
			var respData models.PeerListResponse
			if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
				fmt.Printf("Failed to decode peer list from %s\n", peer)
				return
			}

			// Merge returned peer list
			ns.MergePeers(respData.Peers)

		}(peer)
	}
}

// mergePeers updates peer map with new peers
func (ns *NodeService) MergePeers(received []string) {
	ns.Mutex.Lock()
	defer ns.Mutex.Unlock()

	now := time.Now()
	for _, peer := range received {
		if peer != ns.SelfID {
			ns.Peers[peer] = now
		}
	}
}

func (ns *NodeService) SyncPeersFrom(peer string) {
	url := fmt.Sprintf("http://%s/peers", peer)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Failed to sync from %s: %v\n", peer, err)
		return
	}
	defer resp.Body.Close()

	var peerList []string
	if err := json.NewDecoder(resp.Body).Decode(&peerList); err != nil {
		fmt.Printf("Failed to decode peer list from %s: %v\n", peer, err)
		return
	}

	// Merge the peer list into our discovery service
	for _, peer := range peerList {
		if peer != ns.SelfID {
			ns.RegisterPeer(peer)
		}
	}

	fmt.Printf("Synced peers from %s: %v\n", peer, peerList)
}

func (ns *NodeService) NotifyAllPeersAboutNewPeer(newPeer string) {
	// Get the current peer list
	peers := ns.GetPeers()

	// For each peer in our list
	for _, peer := range peers {
		// Don't notify the new peer about itself
		if peer == newPeer {
			continue
		}

		// Prepare the notification
		url := fmt.Sprintf("http://%s/notify-new-peer", peer)
		request := struct {
			NewPeer string `json:"new_peer"`
		}{
			NewPeer: newPeer,
		}

		jsonData, err := json.Marshal(request)
		if err != nil {
			fmt.Printf("Error marshaling peer notification: %v\n", err)
			continue
		}

		// Send the notification
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Failed to notify %s about new peer %s: %v\n", peer, newPeer, err)
			continue
		}
		resp.Body.Close()

		fmt.Printf("Notified %s about new peer %s\n", peer, newPeer)
	}
}
