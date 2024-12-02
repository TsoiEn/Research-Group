package consensus

import (
	"log"
	"sync"
	"time"

	model "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/chaincode/src/model"
	"golang.org/x/exp/rand"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type RaftNode struct {
	NodeID        string
	State         State
	CurrentTerm   int
	VotedFor      string
	Log           []Block
	CommitIndex   int
	LastApplied   int
	Peers         []string
	BlockChain    []*model.Block
	LeaderID      string
	ElectionTimer *time.Timer
	Mutex         sync.Mutex
	electionChan  chan bool
	heartbeat     time.Duration
}

type Block struct {
	Index     int
	Data      string
	Timestamp time.Time
	Hash      []byte
	PrevHash  []byte
}

func NewRaftNode(nodeID string, peers []string) *RaftNode {
	// Initialize the Raft node
	node := &RaftNode{
		NodeID:      nodeID,
		State:       Follower,
		CurrentTerm: 0,
		VotedFor:    "",
		Log:         []Block{},
		CommitIndex: 0,
		LastApplied: 0,
		Peers:       peers,
		LeaderID:    "",
	}

	// Create and append the genesis block
	genesisBlock := &model.Block{
		Index:     0,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      []byte("Genesis Block"),
		PrevHash:  nil,
	}
	genesisBlock.DeriveHash()
	node.BlockChain = append(node.BlockChain, genesisBlock)

	log.Printf("Node %s: Genesis block created with hash: %x", nodeID, genesisBlock.Hash)

	node.ResetElectionTimer()
	return node
}

func (rn *RaftNode) ResetElectionTimer() {
	rn.Mutex.Lock()
	defer rn.Mutex.Unlock()

	if rn.ElectionTimer != nil {
		rn.ElectionTimer.Stop()
	}
	rn.ElectionTimer = time.AfterFunc(rn.getRandomElectionTimeout(), func() {
		rn.startElection()
	})
}

func (rn *RaftNode) getRandomElectionTimeout() time.Duration {
	return 150*time.Millisecond + time.Duration(rand.Intn(150))*time.Millisecond
}

func (rn *RaftNode) Start() error {
	rn.Mutex.Lock()
	defer rn.Mutex.Unlock()

	// Initialize node state
	rn.State = Follower
	rn.CurrentTerm = 0
	rn.VotedFor = ""

	log.Printf("Node %s initialized as Follower", rn.NodeID)

	go rn.startElectionTimer()

	return nil
}

func (rn *RaftNode) startElection() {
	rn.Mutex.Lock()
	rn.State = Candidate
	rn.CurrentTerm++
	rn.VotedFor = rn.NodeID
	rn.Mutex.Unlock()

	log.Printf("Node %s: Transitioned to Candidate for term %d.", rn.NodeID, rn.CurrentTerm)

	// Simulate requesting votes from peers
	voteCount := 1 // Vote for self

	for _, peer := range rn.Peers {
		go func(peerID string) {
			if rn.requestVote(peerID) {
				rn.Mutex.Lock()
				voteCount++
				if voteCount > len(rn.Peers)/2 {
					rn.becomeLeader()
				}
				rn.Mutex.Unlock()
			}
		}(peer)
	}
}

func (rn *RaftNode) startElectionTimer() {
	for {
		// Randomize election timeout to reduce split votes.
		timeout := rn.getRandomElectionTimeout()

		select {
		case <-time.After(timeout):
			// No heartbeat received; start an election.
			log.Printf("Node %s: Election timeout reached. Starting election.", rn.NodeID)
			rn.startElection()
		case <-rn.electionChan:
			// Received heartbeat or reset signal; continue as Follower.
			log.Printf("Node %s: Election timer reset.", rn.NodeID)
		}
	}
}

func (rn *RaftNode) requestVote(peerID string) bool {
	log.Printf("Node %s: Requesting vote from %s", rn.NodeID, peerID)
	time.Sleep(50 * time.Millisecond)
	return true
}

func (rn *RaftNode) becomeLeader() {
	rn.Mutex.Lock()
	defer rn.Mutex.Unlock()

	rn.State = Leader
	rn.LeaderID = rn.NodeID

	log.Printf("Node %s became the leader for term %d.", rn.NodeID, rn.CurrentTerm)

	// Start sending heartbeats
	go rn.sendHeartbeats()
}

func (rn *RaftNode) sendHeartbeats() {
	for rn.State == Leader {
		log.Printf("Leader %s: Sending heartbeats.", rn.NodeID)
		time.Sleep(rn.heartbeat)
	}
}

func (rn *RaftNode) ProposeBlock(block *model.Block) bool {
	rn.Mutex.Lock()
	defer rn.Mutex.Unlock()

	if rn.State != Leader {
		log.Printf("Node %s: Cannot propose block as it is not the leader", rn.NodeID)
		return false
	}

	// Append the block to the local blockchain
	rn.BlockChain = append(rn.BlockChain, block)

	log.Printf("Node %s: Block proposed successfully", rn.NodeID)
	return true
}
