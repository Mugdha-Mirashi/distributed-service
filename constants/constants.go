package constants

import "time"

const (
	DefaultPort = 8080
)

const (
	HeartbeatInterval = 5 * time.Second
	HeartbeatTimeout  = 10 * time.Second
)
