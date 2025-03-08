package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

const difficulty = 4

type Transaction struct {
	ID      string
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	TxID      string
	OutIdx    int
	Signature string
}

type TxOutput struct {
	Value  int
	PubKey string
}

type Block struct{
	ID        uint  `gorm:"primaryKey;autoIncrement"`
	PrevHash string
	SelfHash string
	Transaction Transaction
	Index int
	Nonce int
	CreatedAt string
}

func (b *Block) MineBlock() {
	for {
		b.SelfHash = b.CalculateHash()
		if b.SelfHash[:difficulty] == strings.Repeat("0", difficulty) {
			break
		}
		b.Nonce++
	}
}

func (b *Block) CalculateHash() string {
	blockData := string(b.Index) + b.CreatedAt + b.PrevHash
	hash := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hash[:])
}

func CreateTransaction(sender, receiver string, amount int) Transaction {
	return Transaction{
		ID: "tx1",
		Inputs: []TxInput{
			{TxID: "genesis", OutIdx: 0, Signature: "sig"},
		},
		Outputs: []TxOutput{
			{Value: amount, PubKey: receiver},
			{Value: 90, PubKey: sender}, // Сдача
		},
	}
}

func NewBlock(index int, tx Transaction, prev string) *Block {
	block := &Block{
		Index: index,
		Transaction: tx,
		PrevHash: prev,
		CreatedAt: time.Now().GoString(),
		Nonce: 0,
	}
	block.SelfHash = block.CalculateHash()
	return block
}