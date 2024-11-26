package consensus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"time"
)

// CreateRaftNode is for creating a new Raft node.
func CreateRaftNode(id string, timeout time.Duration, heartbeatTimeout time.Duration) *RaftNode {
	node := &RaftNode{
		id:               id,
		state:            "follower",
		term:             0,
		log:              make([]LogEntry, 0),
		commitIndex:      0,
		lastApplied:      0,
		nextIndex:        make(map[string]int),
		matchIndex:       make(map[string]int),
		timeout:          timeout,
		heartbeatTimeout: heartbeatTimeout,
		heartbeatTimer:   time.NewTimer(heartbeatTimeout),
		electionTimer:    time.NewTimer(timeout),
		votedFor:         "",
		votesReceived:    0,
	}
	return node
}

// RaftState represents the persistent state of a Raft node.
type RaftState struct {
	Term        int
	VotedFor    string
	Log         []LogEntry
	CommitIndex int
	LastApplied int
}

/*
RaftNode is a struct that represents a Raft node. State is one of "follower", "candidate", or "leader".
*/

func (node *RaftNode) Start() {
	go node.electionTimeout()
	rpc.Register(node)
	listener, err := net.Listen("tcp", ":0") // Use appropriate address.
	if err != nil {
		log.Fatal("Failed to start RPC server: ", err)
	}
	go rpc.Accept(listener)
}

// AddLogEntry is for adding a new log entry to the Raft node.
func (node *RaftNode) AddLogEntry(command string) {
	node.mu.Lock()
	defer node.mu.Unlock()
	entry := LogEntry{
		term:    node.term,
		command: command,
	}
	node.log = append(node.log, entry)
}

// InitializeNode initializes the Raft node.
func (node *RaftNode) InitializeNode() {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.state = "follower"
	node.term = 0
	node.votedFor = ""
	node.votesReceived = 0
	node.commitIndex = 0
	node.lastApplied = 0
	node.nextIndex = make(map[string]int)
	node.matchIndex = make(map[string]int)
	node.log = make([]LogEntry, 0)
	node.heartbeatTimer.Reset(node.heartbeatTimeout)
	node.electionTimer.Reset(node.timeout)
}

// StopNode stops the Raft node.
func (node *RaftNode) StopNode() {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.heartbeatTimer.Stop()
	node.electionTimer.Stop()
}

// SendRequestVote is for sending a RequestVote RPC to another Raft node.
func (node *RaftNode) SendRequestVote(target string, args *RequestVoteArgs, reply *RequestVoteReply) error {
	client, err := rpc.Dial("tcp", target)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Call("RaftNode.RequestVote", args, reply)
}

// SendAppendEntries is for sending an AppendEntries RPC to another Raft node.
func (node *RaftNode) SendAppendEntries(target string, args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	client, err := rpc.Dial("tcp", target)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Call("RaftNode.AppendEntries", args, reply)
}

// BecomeFollower transitions the node to a follower state.
func (node *RaftNode) BecomeFollower(term int) {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.state = "follower"
	node.term = term
	node.votedFor = ""
	node.votesReceived = 0
	node.heartbeatTimer.Reset(node.heartbeatTimeout)
	node.electionTimer.Reset(node.timeout)
}

// BecomeCandidate transitions the node to a candidate state.
func (node *RaftNode) BecomeCandidate() {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.state = "candidate"
	node.term++
	node.votedFor = node.id
	node.votesReceived = 1
	node.electionTimer.Reset(node.timeout)
}

// BecomeLeader transitions the node to the leader state.
func (node *RaftNode) BecomeLeader() {
	node.mu.Lock()
	defer node.mu.Unlock()

	// Ensure the node is in candidate state when transitioning to leader.
	if node.state != "candidate" {
		return
	}

	// Transition to leader state
	node.state = "leader"
	log.Printf("Node %s became the leader for term %d", node.id, node.term)

	// Initialize leader-specific structures, such as nextIndex and matchIndex
	for peer := range node.nextIndex {
		node.nextIndex[peer] = len(node.log)
		node.matchIndex[peer] = 0
	}

	// Send initial empty AppendEntries (heartbeat) to peers to assert leadership
	go node.sendHeartBeat()
}

// HandleTimeout handles the election timeout event.
func (node *RaftNode) HandleTimeout() {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.state == "leader" {
		return
	}
	node.BecomeCandidate()
	node.StartElection()
}

// HandleHeartbeatResponse handles the response to a heartbeat.
func (node *RaftNode) HandleHeartbeatResponse(success bool, term int) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if success {
		return
	}
	if term > node.term {
		node.BecomeFollower(term)
	}
}

// ProcessElectionResult processes the result of an election.
func (node *RaftNode) ProcessElectionResult(votes int) {
	node.mu.Lock()
	defer node.mu.Unlock()

	if node.state == "candidate" && votes > len(node.nextIndex)/2 {
		node.BecomeLeader() // Call BecomeLeader after winning an election.
	}
}

// HandleAppendEntriesResponse handles the response to an AppendEntries RPC.
func (node *RaftNode) HandleAppendEntriesResponse(success bool, term int) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if success {
		return
	}
	if term > node.term {
		node.BecomeFollower(term)
	}
}

// HandleAppendEntries handles an AppendEntries RPC.
func (node *RaftNode) HandleAppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if args.Term < node.term {
		reply.Term = node.term
		reply.Success = false
		return
	}
	node.term = args.Term
	node.votedFor = args.LeaderID
	node.state = "follower"
	node.electionTimer.Reset(node.timeout)
	if args.LeaderCommit > node.commitIndex {
		node.commitIndex = min(args.LeaderCommit, len(node.log)-1)
		node.Commit()
	}
	reply.Success = true
	reply.Term = node.term
}

// HandleRequestVote handles a RequestVote RPC.
func (node *RaftNode) HandleRequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if args.Term > node.term {
		node.term = args.Term
		node.votedFor = ""
		node.state = "follower"
	}
	if (node.votedFor == "" || node.votedFor == args.CandidateID) && isCandidateLogUpToDate(node, args.LastLogIndex, args.LastLogTerm) {
		node.votedFor = args.CandidateID
		node.electionTimer.Reset(node.timeout)
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}
	reply.Term = node.term
}

// LogStatus logs the current status of the Raft node.
func (node *RaftNode) LogStatus() {
	node.mu.Lock()
	defer node.mu.Unlock()
	log.Printf("Node ID: %s, State: %s, Term: %d, CommitIndex: %d, LastApplied: %d, VotesReceived: %d",
		node.id, node.state, node.term, node.commitIndex, node.lastApplied, node.votesReceived)
}

// UpdateMetrics updates the metrics of the Raft node.
func (node *RaftNode) UpdateMetrics() {
	node.mu.Lock()
	defer node.mu.Unlock()
	// Example metric updates, replace with actual metric collection logic
	log.Printf("Updating metrics for Node ID: %s", node.id)
	// Add your metric update logic here
}

// UpdatePeerList updates the list of peers for the Raft node.
func (node *RaftNode) UpdatePeerList(peers []string) {
	node.mu.Lock()
	defer node.mu.Unlock()
	for _, peer := range peers {
		if _, exists := node.nextIndex[peer]; !exists {
			node.nextIndex[peer] = len(node.log)
			node.matchIndex[peer] = 0
		}
	}
	for peer := range node.nextIndex {
		if !contains(peers, peer) {
			delete(node.nextIndex, peer)
			delete(node.matchIndex, peer)
		}
	}
}

// UpdateTimeouts updates the election and heartbeat timeouts for the Raft node.
func (node *RaftNode) UpdateTimeouts(electionTimeout, heartbeatTimeout time.Duration) {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.timeout = electionTimeout
	node.heartbeatTimeout = heartbeatTimeout
	node.electionTimer.Reset(electionTimeout)
	node.heartbeatTimer.Reset(heartbeatTimeout)
}

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SaveState saves the current state of the Raft node to disk.
func (node *RaftNode) SaveState() error {
	node.mu.Lock()
	defer node.mu.Unlock()
	state := &RaftState{
		Term:        node.term,
		VotedFor:    node.votedFor,
		Log:         node.log,
		CommitIndex: node.commitIndex,
		LastApplied: node.lastApplied,
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("raft_state.json", data, 0644)

}

// LoadState loads the state of the Raft node from disk.
func (node *RaftNode) LoadState() error {
	node.mu.Lock()
	defer node.mu.Unlock()
	data, err := ioutil.ReadFile("raft_state.json")
	if err != nil {
		return err
	}
	state := &RaftState{}
	if err := json.Unmarshal(data, state); err != nil {
		return err
	}
	node.term = state.Term
	node.votedFor = state.VotedFor
	node.log = state.Log
	node.commitIndex = state.CommitIndex
	node.lastApplied = state.LastApplied
	return nil
}

// SendRequestVoteWithRetry sends a RequestVote RPC with retry logic.
func (node *RaftNode) SendRequestVoteWithRetry(target string, args *RequestVoteArgs, reply *RequestVoteReply, retries int) error {
	for i := 0; i < retries; i++ {
		err := node.SendRequestVote(target, args, reply)
		if err == nil {
			return nil
		}
		log.Printf("Failed to send RequestVote to %s: %v. Retrying (%d/%d)...", target, err, i+1, retries)
		time.Sleep(time.Second)
	}
	return fmt.Errorf("failed to send RequestVote to %s after %d retries", target, retries)
}

// SendAppendEntriesWithRetry sends an AppendEntries RPC with retry logic.
func (node *RaftNode) SendAppendEntriesWithRetry(target string, args *AppendEntriesArgs, reply *AppendEntriesReply, retries int) error {
	for i := 0; i < retries; i++ {
		err := node.SendAppendEntries(target, args, reply)
		if err == nil {
			return nil
		}
		log.Printf("Failed to send AppendEntries to %s: %v. Retrying (%d/%d)...", target, err, i+1, retries)
		time.Sleep(time.Second)
	}
	return fmt.Errorf("failed to send AppendEntries to %s after %d retries", target, retries)
}

// sendHeartBeat sends a heartbeat to all other nodes in the cluster.
func (node *RaftNode) sendHeartBeat() {
	node.mu.Lock()
	if node.state != "leader" {
		node.mu.Unlock()
		return
	}
	node.mu.Unlock()

	for peer := range node.nextIndex {
		go func(target string) {
			args := AppendEntriesArgs{
				Term:         node.term,
				LeaderID:     node.id,
				Entries:      nil, // Heartbeat contains no log entries
				LeaderCommit: node.commitIndex,
			}
			var reply AppendEntriesReply
			err := node.SendAppendEntries(target, &args, &reply)
			if err != nil {
				log.Printf("Error sending heartbeat to %s: %v", target, err)
				return
			}

			node.HandleHeartbeatResponse(reply.Success, reply.Term)
		}(peer)
	}

	// Schedule the next heartbeat
	time.AfterFunc(node.heartbeatTimeout, node.sendHeartBeat)
}

// HandleElectionWin handles the event when a node wins the election.
func (node *RaftNode) HandleElectionWin() {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.state == "candidate" {
		node.BecomeLeader()
	}
}

// HandleVoteResponse handles the response to a RequestVote RPC.
func (node *RaftNode) HandleVoteResponse(reply *RequestVoteReply) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if reply.Term > node.term {
		node.BecomeFollower(reply.Term)
		return
	}
	if node.state == "candidate" && reply.VoteGranted {
		node.votesReceived++
		if node.votesReceived > len(node.nextIndex)/2 {
			node.HandleElectionWin()
		}
	}
}
