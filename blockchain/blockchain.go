package blockchain

import (
	"crypto/sha256"
	"log"
	"time"
)

type Blockchain struct {
	Head   Block
	blocks map[string]Block
}

func GetBlockchain() Blockchain {
	seed := SeedBlock()
	return Blockchain{
		Head:   seed,
		blocks: map[string]Block{seed.Hash: seed},
	}
}

func (bc *Blockchain) InsertBlock(block Block) {
	bc.Head = block
	bc.blocks[block.Hash] = block
}

func (bc *Blockchain) ReplaceChain(blockchain Blockchain) bool {
	seed := SeedBlock()
	if seed != bc.blocks[seed.Hash] {
		log.Println("blockchain: seed block is invalid")
		return false
	}

	if len(blockchain.blocks) > len(bc.blocks) {
		log.Println("blockchain: replaced blockchain with new one")
		bc = &blockchain
		return true
	}
	return false
}

func (bc *Blockchain) IsValid() bool {
	curBlock := bc.Head
	seed := SeedBlock()
	for seed != curBlock {
		prevBlock, found := bc.blocks[curBlock.PrevHash]
		if !found || !bc.IsValidBlock(curBlock, prevBlock) {
			log.Println("blockchain: invalid blocks")
			return false
		}
		curBlock = prevBlock
	}

	_, found := bc.blocks[seed.PrevHash]
	return !found
}

type Block struct {
	Index     int       `json:"Index"`
	Timestamp time.Time `json:"Timestamp"`
	Data      string    `json:"Data"`
	Hash      string    `json:"Hash"`
	PrevHash  string    `json:"PrevHash"`
}

func (bc *Blockchain) AddNewBlock(data string) Block {
	block := Block{
		Index:     len(bc.blocks) + 1,
		Timestamp: time.Now(),
		Data:      data,
		PrevHash:  bc.Head.Hash,
	}
	block.Hash = generateHash(block)
	bc.InsertBlock(block)
	return block
}

func (bc *Blockchain) IsValidBlock(block Block, prevBlock Block) bool {
	if block.Index != len(bc.blocks) {
		log.Printf("block: invalid index %d, expected %d", block.Index, len(bc.blocks))
		return false
	} else if prevBlock.Hash != block.PrevHash {
		log.Printf("block: invalid previous hash on block with index %d", block.Index)
		return false
	} else if block.Timestamp.Before(prevBlock.Timestamp) {
		log.Printf("block: invalid timestamp on block with index %d", block.Index)
		return false
	} else if generateHash(block) != block.Hash {
		log.Printf("block: invalid hash on block with index %d", block.Index)
		return false
	}
	return true
}

func SeedBlock() Block {
	return Block{
		Index:     1,
		Timestamp: time.Unix(0, 0),
		Data:      "",
		Hash:      "BASE_HASH",
		PrevHash:  "BASE_PREV_HASH",
	}
}

func generateHash(block Block) string {
	encodedTime, _ := block.Timestamp.MarshalBinary()
	hashSeed := []byte{}
	hashSeed = append(hashSeed, byte(block.Index))
	hashSeed = append(hashSeed, []byte(block.PrevHash)...)
	hashSeed = append(hashSeed, encodedTime...)
	hashSeed = append(hashSeed, []byte(block.Data)...)
	hash := sha256.Sum256(hashSeed)
	return string(hash[:len(hash)])
}
