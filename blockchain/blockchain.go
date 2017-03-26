package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
func BlockchainFromJSON(jsonBlob []byte) Blockchain {
	blocks := []Block{}
	bc := GetBlockchain()
	if err := json.Unmarshal(jsonBlob, &blocks); err != nil {
		log.Println("unable to unmarshal JSON", string(jsonBlob), err)
		return bc
	}
	for _, b := range blocks {
		bc.InsertBlock(b)
	}
	if len(blocks) > 0 {
		bc.Head = blocks[0]
	}
	return bc
}

func (bc *Blockchain) InsertBlock(block Block) {
	bc.Head = block
	bc.blocks[block.Hash] = block
}

func (bc *Blockchain) ShouldReplaceWithChain(blockchain Blockchain) bool {
	seed := SeedBlock()
	if seed != bc.blocks[seed.Hash] {
		log.Println("blockchain: seed block is invalid")
		return false
	}

	if len(blockchain.blocks) > len(bc.blocks) {
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

func (bc *Blockchain) Blocks() []Block {
	blocks := []Block{}
	block := bc.Head
	for block != SeedBlock() {
		blocks = append(blocks, block)
		block = bc.blocks[block.PrevHash]
	}
	return blocks
}

func (bc *Blockchain) IsValidBlock(block Block, prevBlock Block) bool {
	if block.Index != len(bc.blocks)+1 {
		log.Printf("block: invalid index %d, expected %d", block.Index, len(bc.blocks)+1)
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
	h := sha256.New()
	h.Write(hashSeed)
	return hex.EncodeToString(h.Sum(nil))
}
