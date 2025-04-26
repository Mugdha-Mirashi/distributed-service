package node

import (
	"distributed-counter-system/constants"
	"fmt"
	"net/http"
	"time"
)

func (ns *NodeService) StartHeartbeats() {
	go func() {
		for {
			time.Sleep(constants.HeartbeatInterval)

			ns.Mutex.RLock()
			peers := ns.GetPeers()
			ns.Mutex.RUnlock()

			for _, peer := range peers {
				go func(peer string) {
					url := fmt.Sprintf("http://%s/ping", peer)

					client := http.Client{
						Timeout: 2 * time.Second,
					}
					resp, err := client.Get(url)
					if err != nil || resp.StatusCode != http.StatusOK {
						fmt.Printf("Peer %s failed heartbeat. Removing.\n", peer)
						ns.RemovePeer(peer)
						return
					}
					defer resp.Body.Close()

					ns.UpdatePeerTimestamp(peer)
				}(peer)
			}
		}
	}()
}
