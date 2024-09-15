package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// BlockChain structure contains a slice of blocks.
type BlockChain struct {
	blocks []Block
}

// Block represents a block in the blockchain.
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

// DeriveHash generates a hash for the block using the data and the previous block's hash.
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// CreateBlock creates a new block with the provided data and previous block's hash.
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

// AddBlock adds a new block with the provided data to the blockchain.
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, *newBlock)
}

// Genesis creates the first block in the blockchain (genesis block).
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// NewBlockChain creates a new blockchain with the genesis block.
func NewBlock() *BlockChain {
	return &BlockChain{[]Block{*Genesis()}}
}

func main() {
	chain := NewBlock()

	chain.AddBlock("Choi a member of this group")
	chain.AddBlock("Raffy a member of this group")
	chain.AddBlock("Blake a member of this group")

	for _, block := range chain.blocks {
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data in the block: %s\n", block.Data)
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
	}
}
