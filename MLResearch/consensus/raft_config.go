package consensus

import (
	"fmt"
	"time"
)

// RaftConfig holds the configuration parameters for the Raft consensus system.
type RaftConfig struct {
	ElectionTimeoutMin time.Duration 
	ElectionTimeoutMax time.Duration 
	HeartbeatInterval  time.Duration
	PeerCount          int           
	MaxRetries         int           
}

// Default configuration values for Raft protocol
var DefaultRaftConfig = RaftConfig{
	ElectionTimeoutMin: 150 * time.Millisecond, 
	ElectionTimeoutMax: 300 * time.Millisecond, 
	HeartbeatInterval:  100 * time.Millisecond, 
	PeerCount:          5,                      
	MaxRetries:         3,                      
}

// NewConfig creates a new configuration for the Raft system with the specified parameters.
// You can use this function if you want to modify the defaults for a specific Raft node setup.
func NewConfig(minTimeout, maxTimeout, heartbeatInterval time.Duration, peerCount, maxRetries int) RaftConfig {
	return RaftConfig{
		ElectionTimeoutMin: minTimeout,
		ElectionTimeoutMax: maxTimeout,
		HeartbeatInterval:  heartbeatInterval,
		PeerCount:          peerCount,
		MaxRetries:         maxRetries,
	}
}

// ValidateConfig checks that the provided RaftConfig is valid.
func ValidateConfig(config RaftConfig) error {
	if config.ElectionTimeoutMin <= 0 || config.ElectionTimeoutMax <= 0 {
		return fmt.Errorf("election timeouts must be positive values")
	}
	if config.ElectionTimeoutMin > config.ElectionTimeoutMax {
		return fmt.Errorf("ElectionTimeoutMin cannot be greater than ElectionTimeoutMax")
	}
	if config.HeartbeatInterval <= 0 {
		return fmt.Errorf("HeartbeatInterval must be positive")
	}
	if config.PeerCount <= 1 {
		return fmt.Errorf("PeerCount must be greater than 1")
	}
	return nil
}
