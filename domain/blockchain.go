package domain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type BlockChain struct {
	Pool  []*Transaction `json:"pool"`
	Chain []*Block       `json:"chain"`
}

func NewBlockChain() *BlockChain {
	bc := LoadDatabase()
	if len(bc.Chain) == 0 {
		bc.CreateGenesisBlock()
		bc.CreateBlock(0, fmt.Sprintf("%x", [32]byte{}))
	}
	return bc
}

func LoadDatabase() *BlockChain {
	f, err := os.OpenFile("database/blockchain.db", os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	blockChain := BlockChain{}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}

		var blockSerialized BlockSerialized
		err = json.Unmarshal(scanner.Bytes(), &blockSerialized)
		if err != nil {
			os.Exit(1)
		}

		if len(blockChain.Chain) > 0 && (blockChain.LatestBlock().Hash() != blockSerialized.Value.Header.PrevHash) {
			log.Fatal("Invalid Blockchain Databse")
		}

		blockChain.Chain = append(blockChain.Chain, blockSerialized.Value)
	}

	return &blockChain
}

func (bc *BlockChain) LatestBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *BlockChain) CreateBlock(nonce int, prevHash string) *Block {
	b := NewBlock(nonce, prevHash, bc.Pool)
	bc.Chain = append(bc.Chain, b)
	bc.Pool = []*Transaction{}
	return b
}

func (bc *BlockChain) CreateGenesisBlock() {
	t := NewTransaction(0, "GOD", "Kayu")
	bc.Pool = append(bc.Pool, t)
}

func (bc *BlockChain) GiveBalance(from, to string, amount int64) bool {
	if bc.CalculateBalance(from) < amount {
		return false
	}

	tx := NewTransaction(amount, from, to)
	bc.Pool = append(bc.Pool, tx)

	return true
}

func (bc *BlockChain) ToUpBalance(toPublicKey string, amount int64) bool {
	if amount <= 0 {
		return false
	}

	t := NewTransaction(amount, "Admin", toPublicKey)
	bc.Pool = append(bc.Pool, t)

	return true
}

func (bc *BlockChain) CalculateBalance(publicKey string) int64 {
	var balance int64 = 0

	for _, block := range bc.Chain {
		for _, transaction := range block.Transactions {
			if transaction.To == publicKey {
				balance += transaction.Amount
			}
			if transaction.From == publicKey {
				balance -= transaction.Amount
			}
		}
	}

	return balance
}

func (bc *BlockChain) MineBlock() bool {
	if len(bc.Pool) == 0 {
		return false
	}

	transactionsCopy := make([]*Transaction, len(bc.Pool))
	copy(transactionsCopy, bc.Pool)

	block := &Block{
		Header: &Header{
			Time:     time.Now().Unix(),
			PrevHash: bc.LatestBlock().Hash(),
			Nonce:    0,
		},
		Transactions: transactionsCopy,
	}
	bc.Chain = append(bc.Chain, block)
	bc.Pool = []*Transaction{} // Kosongkan Pool
	return true
}
