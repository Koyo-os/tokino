package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

const difficulty = 4

type Block struct{
	ID        uint  `gorm:"primaryKey;autoIncrement"`
	PrevHash string
	SelfHash string
	Data string
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
	blockData := string(b.Index) + b.CreatedAt + b.Data + b.PrevHash
	hash := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hash[:])
}

func NewBlock(index int, data string, prev string) *Block {
	block := &Block{
		Index: index,
		Data: data,
		PrevHash: prev,
		CreatedAt: time.Now().GoString(),
		Nonce: 0,
	}
	block.SelfHash = block.CalculateHash()
	return block
}