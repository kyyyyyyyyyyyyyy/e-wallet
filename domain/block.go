package domain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Block struct {
	Header       *Header        `json:"header"`
	Transactions []*Transaction `json:"transactions"`
}

type BlockSerialized struct {
	Key   string `json:"key"`
	Value *Block `json:"value"`
}

func NewBlock(nonce int, prevHash string, transactions []*Transaction) *Block {
	b := new(Block)
	b.Header = &Header{
		Nonce:    nonce,
		PrevHash: prevHash,
		Time:     time.Now().Unix(),
	}
	b.Transactions = transactions
	b.Persist()
	return b
}

func (b *Block) Hash() string {
	m, _ := json.Marshal(b)
	return fmt.Sprintf("%x", sha256.Sum256(m))
}

func (b *Block) Persist() {
	blockSerialized, _ := json.Marshal(BlockSerialized{b.Hash(), b})
	f, _ := os.OpenFile("database/blockchain.db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	_, _ = f.Write(blockSerialized)
	_, _ = f.Write([]byte("\n"))
}
