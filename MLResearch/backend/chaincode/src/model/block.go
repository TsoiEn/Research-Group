package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
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

// DeriveHash generates a hash for the block using a homomorphic hashing scheme.
func (b *Block) DeriveHash() {
	// Modular arithmetic parameters for homomorphic hashing
	modulus := big.NewInt(1 << 62) // Example modulus for a large space
	info := bytes.Join([][]byte{[]byte(fmt.Sprintf("%d", b.Index)), []byte(b.Timestamp), b.Data, b.PrevHash}, []byte{})
	infoInt := new(big.Int).SetBytes(info)

	// Homomorphic operation: Modulo addition (simulated homomorphic hashing)
	hash := new(big.Int).Mod(infoInt, modulus)
	b.Hash = hash.Bytes()
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
