package consensus

import "math/rand"

// GenerateRandomTerm generates a random term within a reasonable range
// for election timeout to avoid candidates timing out simultaneously.
func GenerateRandomTerm(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// ValidateTerm checks if a given term is valid (greater than or equal to the node's term).
func ValidateTerm(nodeTerm, replyTerm int) bool {
	return replyTerm >= nodeTerm
}

// PeerCount returns the number of peers in the cluster.
func PeerCount(peers []string) int {
	return len(peers)
}

// SendVoteRequest sends a vote request to another node.
func SendVoteRequest(peerID string, currentTerm int) (*RequestVoteReply, error) {
	// Simulate sending a vote request to a peer.
	// In a real Raft implementation, this would involve sending a network RPC.
	// Here we just simulate the response for simplicity.

	// Simulated response: randomly grant the vote based on some condition.
	reply := &RequestVoteReply{
		Term:        currentTerm,
		VoteGranted: rand.Float32() > 0.5, // Simulate 50% chance of vote granted
	}

	return reply, nil
}

// SendHeartbeat simulates sending a heartbeat message to all follower peers.
func SendHeartbeat(peers []string, leaderID string) {
	for _, peer := range peers {
		// Simulate sending a heartbeat to each peer
		println("Heartbeat sent to peer:", peer, "from leader:", leaderID)
	}
}

// IsMajority checks if the given count is a majority in the cluster.
func IsMajority(count, total int) bool {
	return count > total/2
}

// GetMajorityCount returns the minimum number of votes required for a majority.
func GetMajorityCount(total int) int {
	return (total / 2) + 1
}
