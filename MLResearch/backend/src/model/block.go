package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// BlockChain structure contains a slice of blocks.
type BlockChain struct {
	Blocks []Block
}

// Block represents a block in the blockchain.
type Block struct {
	Index     int
	Timestamp string
	Data      []byte
	Hash      []byte
	PrevHash  []byte
}

// Serialize serializes the block into a JSON byte slice.
func (b *Block) Serialize() ([]byte, error) {
	return json.Marshal(b)
}

// DeriveHash generates a hash for the block using its index, timestamp, data, and previous hash.
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{[]byte(fmt.Sprintf("%d", b.Index)), []byte(b.Timestamp), b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// CreateBlock creates a new block with the given data and previous hash.
func CreateBlock(index int, blockData []byte, prevHash []byte) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      blockData,
		PrevHash:  prevHash,
	}
	block.DeriveHash()
	return block
}

// AddBlock adds a new block to the blockchain.
func (chain *BlockChain) AddBlock(blockData []byte) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newIndex := prevBlock.Index + 1
	newBlock := CreateBlock(newIndex, blockData, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, *newBlock)
}

// Genesis creates the first block in the blockchain.
func Genesis() *Block {
	return CreateBlock(0, []byte("Genesis Block"), []byte{})
}

// NewBlockChain creates a blockchain with the genesis block.
func NewBlockChain() *BlockChain {
	return &BlockChain{Blocks: []Block{*Genesis()}}
}

// FindCredentialByID searches the blockchain for a credential with the given ID.
func (chain *BlockChain) FindCredentialByID(id string) (*Credential, error) {
	for _, block := range chain.Blocks {
		var cred Credential
		if err := json.Unmarshal(block.Data, &cred); err == nil {
			if cred.ID == id {
				return &cred, nil
			}
		}
	}
	return nil, fmt.Errorf("credential with ID %s not found", id)
}
