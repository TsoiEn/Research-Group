package blockchain

import (
	"bytes"
	"crypto/sha256"
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

// DeriveHash generates a hash for the block using the data and the previous block's hash.
func (b *Block) DeriveHash() {
	// Concatenate data, previous hash, index, and timestamp
	info := bytes.Join([][]byte{[]byte(fmt.Sprintf("%d", b.Index)), []byte(b.Timestamp), b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:] // Store the hash as a byte slice
}

// CreateBlock creates a new block with the provided data and previous block's hash.
func CreateBlock(index int, blockData []byte, prevHash []byte) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      blockData, // Store the serialized credential data
		PrevHash:  prevHash,
	}
	block.DeriveHash()
	return block
}

// AddBlock adds a new block with the provided data to the blockchain.
func (chain *BlockChain) AddBlock(data string, blockData []byte) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newIndex := prevBlock.Index + 1
	// Use blockData instead of data for the block
	newBlock := CreateBlock(newIndex, blockData, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, *newBlock)
}

// Genesis creates the first block in the blockchain (genesis block).
func Genesis() *Block {
	return CreateBlock(0, []byte("Genesis"), []byte{}) // Make sure to pass a []byte for Data
}

// NewBlockChain creates a new blockchain with the genesis block.
func NewBlockChain() *BlockChain {
	return &BlockChain{[]Block{*Genesis()}}
}
