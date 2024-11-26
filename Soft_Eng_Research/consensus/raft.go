package consensus

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"time"

	"github.com/TsoiEn/Research-Group/Soft_Eng_Research/blockchain_core/chaincode"
)

type RaftNode struct {
	id                 string
	state              string
	term               int
	log                []LogEntry
	commitIndex        int
	lastApplied        int
	nextIndex          map[string]int
	matchIndex         map[string]int
	timeout            time.Duration
	heartbeatTimeout   time.Duration
	heartbeatTimer     *time.Timer
	electionTimer      *time.Timer
	leaderID           string
	votedFor           string
	votesReceived      int
	peers              []string
	mu                 sync.Mutex
	ElectionTimeoutMin time.Duration // Minimum timeout for elections
	ElectionTimeoutMax time.Duration // Maximum timeout for elections
	HeartbeatInterval  time.Duration // Interval for heartbeats
	PeerCount          int           // Total number of peers in the cluster
	MaxRetries         int
	currentState       map[string]interface{} // Simulated blockchain state
	otherNodes         []string               // Addresses of other nodes
	transactionCh      chan map[string]interface{}
}

type LogEntry struct {
	term    int
	command string
	Args    []interface{}
}

type AppendEntriesArgs struct {
	Term         int
	LeaderID     string
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

type RequestVoteArgs struct {
	Term         int
	CandidateID  string
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

func (node *RaftNode) waitForElectionResult() {
	votesReceived := 1

	for {
		select {
		case <-node.heartbeatTimer.C:
			fmt.Println("Heartbeat timeout")
			node.StartElection()
			return
		default:
			if votesReceived > len(node.nextIndex)/2 {
				node.mu.Lock()
				node.state = "leader"
				node.leaderID = node.id
				node.mu.Unlock()
				node.heartbeat()
				return
			}
		}
	}
}

func (node *RaftNode) heartbeat() {
	fmt.Println("Sending heartbeat...")
	node.heartbeatTimer.Reset(node.heartbeatTimeout)
	// Implement RPC to send heartbeats to peers
}

func (node *RaftNode) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	node.mu.Lock()
	defer node.mu.Unlock()
	if args.Term < node.term {
		reply.Term = node.term
		reply.Success = false
		return nil
	}
	node.term = args.Term
	node.leaderID = args.LeaderID
	node.state = "follower"
	node.electionTimer.Reset(node.timeout)

	if args.LeaderCommit > node.commitIndex {
		node.commitIndex = min(args.LeaderCommit, len(node.log)-1)
		node.Commit()
	}

	reply.Success = true
	reply.Term = node.term
	return nil
}

func (node *RaftNode) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
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
	return nil
}

func isCandidateLogUpToDate(node *RaftNode, lastLogIndex int, lastLogTerm int) bool {
	if len(node.log) == 0 {
		return true
	}
	lastLogEntry := node.log[len(node.log)-1]
	if lastLogTerm > lastLogEntry.term {
		return true
	}
	if lastLogTerm == lastLogEntry.term && lastLogIndex >= len(node.log)-1 {
		return true
	}
	return false
}

func (node *RaftNode) electionTimeout() {
	for range node.electionTimer.C {
		node.mu.Lock()
		node.state = "candidate"
		node.term++
		node.votedFor = node.id
		node.votesReceived = 1
		node.mu.Unlock()
		node.StartElection()
	}
}

func (node *RaftNode) StartElection() {
	fmt.Printf("%s is starting an election for term %d\n", node.id, node.term)
	node.votesReceived = 1 // Node votes for itself

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, peer := range node.peers {
		wg.Add(1)
		go func(peer string) {
			defer wg.Done()
			args := &RequestVoteArgs{
				Term:         node.term,
				CandidateID:  node.id,
				LastLogIndex: len(node.log) - 1,
				LastLogTerm:  node.getLastLogTerm(),
			}
			reply := &RequestVoteReply{}

			client, err := rpc.Dial("tcp", peer)
			if err != nil {
				log.Printf("Failed to connect to peer %s: %v\n", peer, err)
				return
			}
			defer client.Close()

			err = client.Call("RaftNode.RequestVote", args, reply)
			if err != nil {
				log.Printf("Error during RequestVote RPC to %s: %v\n", peer, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			if reply.VoteGranted {
				node.votesReceived++
				if node.votesReceived > len(node.peers)/2 {
					fmt.Printf("%s became the leader for term %d\n", node.id, node.term)
					node.mu.Lock()
					node.state = "leader"
					node.leaderID = node.id
					node.mu.Unlock()
					node.heartbeat()
				}
			} else if reply.Term > node.term {
				node.mu.Lock()
				node.term = reply.Term
				node.state = "follower"
				node.votedFor = ""
				node.mu.Unlock()
				node.electionTimer.Reset(node.timeout)
			}
		}(peer)
	}

	// Wait for responses from all peers before ending the election.
	wg.Wait()
	go node.waitForElectionResult()
}

func (node *RaftNode) getLastLogTerm() int {
	if len(node.log) == 0 {
		return 0
	}
	return node.log[len(node.log)-1].term
}

func (node *RaftNode) RequestVoteResponse(voteGranted bool) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if voteGranted {
		node.votesReceived++
		if node.votesReceived > len(node.peers)/2 {
			node.state = "leader"
			node.leaderID = node.id
			node.heartbeat()
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (node *RaftNode) AppendEntriesResponse(success bool, followerID string, matchIndex int) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if success {
		node.matchIndex[followerID] = matchIndex
		node.nextIndex[followerID] = matchIndex + 1

		// Update commitIndex if a majority of matchIndex are greater than commitIndex
		for i := node.commitIndex + 1; i < len(node.log); i++ {
			count := 1 // count itself
			for _, matchIdx := range node.matchIndex {
				if matchIdx >= i {
					count++
				}
			}
			if count > len(node.nextIndex)/2 {
				node.commitIndex = i
				node.Commit()
			}
		}
	} else {
		node.nextIndex[followerID] = max(1, node.nextIndex[followerID]-1)
	}

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (node *RaftNode) Commit() {
	for node.lastApplied < node.commitIndex {
		node.lastApplied++
		entry := node.log[node.lastApplied]
		node.ApplyLog(entry)
	}
}

func (node *RaftNode) ApplyLog(entry LogEntry) {
	fmt.Printf("Applying log entry: %v\n", entry)

	switch entry.command {
	case "AddNewStudent":
		// Extract arguments from the log entry
		if len(entry.Args) != 7 {
			fmt.Println("Invalid arguments for AddNewStudent")
			return
		}

		id := entry.Args[0].(int)
		firstName := entry.Args[1].(string)
		lastName := entry.Args[2].(string)
		age := entry.Args[3].(int)
		birthDate := entry.Args[4].(time.Time)
		studentNum := entry.Args[5].(int)
		chain := entry.Args[6].(*chaincode.StudentChain)

		// Execute the chaincode function
		student := chaincode.AddNewStudent(id, firstName, lastName, age, birthDate, studentNum, chain)
		fmt.Printf("Added new student: %v\n", student)

	case "UpdateStudentCredentials":
		if len(entry.Args) != 3 {
			fmt.Println("Invalid arguments for UpdateStudentCredentials")
			return
		}

		id := entry.Args[0].(int)
		newCredential := entry.Args[1].(chaincode.Credential)
		chain := entry.Args[2].(*chaincode.StudentChain)

		success := chain.UpdateStudentCredentials(id, newCredential)
		if success {
			fmt.Printf("Updated credentials for student ID %d\n", id)
		} else {
			fmt.Printf("Failed to update credentials for student ID %d\n", id)
		}

	default:
		fmt.Printf("Unknown command: %s\n", entry.command)
	}
}

func (rn *RaftNode) ProposeTransaction(transaction map[string]interface{}) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// Leadership check
	if !rn.isLeader() {
		return errors.New("this node is not the leader")
	}

	// Append transaction to local log
	entry := LogEntry{
		term:    rn.term,
		command: fmt.Sprintf("%v", transaction),
	}
	rn.log = append(rn.log, entry)

	// Replicate to other nodes
	err := rn.replicateTransaction(transaction)
	if err != nil {
		return err
	}

	// Commit and apply transaction
	rn.commitIndex++
	rn.applyTransaction(transaction)
	return nil
}

func (rn *RaftNode) replicateTransaction(transaction map[string]interface{}) error {
	for _, nodeAddress := range rn.otherNodes {
		// Send the transaction to the other nodes (simulate with HTTP or gRPC)
		// Example: HTTP POST to replicate the transaction
		// Simulated here:
		go func(addr string) {
			// Simulate replication logic
		}(nodeAddress)
	}
	return nil
}

func (rn *RaftNode) applyTransaction(transaction map[string]interface{}) {
	// Simulate applying the transaction to the current state
	rn.currentState[transaction["action"].(string)] = transaction["data"]
}

func (rn *RaftNode) isLeader() bool {
	// Simulate checking if the current node is the leader
	return rn.leaderID == rn.id // Check if this node is the leader
}

func (node *RaftNode) SubmitTransaction(command string, args []interface{}) {
	entry := LogEntry{
		term:    node.term,
		command: command,
		Args:    args,
	}

	node.mu.Lock()
	node.log = append(node.log, entry)
	node.mu.Unlock()

	fmt.Printf("Transaction submitted: %v\n", entry)

	chain := &StudentChain{} // Initialize or reference your blockchain state
	node.SubmitTransaction("AddNewStudent", []interface{}{1, "John", "Doe", 20, time.Now(), 12345, chain})

	newCredential := Credential{ /* Fill in the credential details */ }
	chain = &StudentChain{}
	node.SubmitTransaction("UpdateStudentCredentials", []interface{}{1, newCredential, chain})

}

func (node *RaftNode) Metrics() map[string]interface{} {
	node.mu.Lock()
	defer node.mu.Unlock()

	metrics := map[string]interface{}{
		"id":            node.id,
		"state":         node.state,
		"term":          node.term,
		"commitIndex":   node.commitIndex,
		"lastApplied":   node.lastApplied,
		"leaderID":      node.leaderID,
		"votedFor":      node.votedFor,
		"votesReceived": node.votesReceived,
		"peers":         node.peers,
	}

	return metrics
}
