package blockchain

import (
	"testing"
	"time"
)

func NewBlock(params map[string]interface{}) Block {
	index := 2
	if i, ok := params["Index"]; ok {
		index = i
	}
	timestamp := time.Now()
	if t, ok := params["Timestamp"]; ok {
		timestamp = t
	}
	data := "test data"
	if d, ok := params["Data"]; ok {
		data = d
	}
	prevHash := SeedBlock().Hash
	if p, ok := params["PrevHash"]; ok {
		prevHash = p
	}
	hash := "test hash"
	if h, ok := params["Hash"]; ok {
		hash = h
	}
	return Block{
		Index:     index,
		Timestamp: timestamp,
		Data:      data,
		Hash:      hash,
		PrevHash:  prevHash,
	}
}

func TestIsValid(t *testing.T) {
	bc := GetBlockchain()
	if !bc.IsValid() {
		t.Error("Expected base blockchain to be valid, but it isn't")
	}
}

func TestAddNewBlock(t *testing.T) {
	bc := GetBlockchain()
	base := bc.Head
	b := bc.AddNewBlock("testing")
	if !bc.IsValid() {
		t.Error("Expected blockchain to be valid, but it isn't")
	}
	if !bc.IsValidBlock(b, base) {
		t.Error("Expected newly added block to be valid, but it isn't")
	}
}

func TestShouldReplaceWithChain(t *testing.T) {
	b1 := GetBlockchain()
	b2 := GetBlockchain()
	b2.AddNewBlock("extra block")
	if b2.ShouldReplaceWithChain(b1) {
		t.Error("Expected b1 to not replace b2, but it should")
	}
	if !b1.ShouldReplaceWithChain(b2) {
		t.Error("Expected b1 to be replaced by b2, but it wasn't recommended")
	}
}

func TestIsValidBlock(t *testing.T) {
	bc := GetBlockchain()
	block := NewBlock(map[string]interface{}{"Index": 4})
	if bc.IsValidBlock(block) {
		t.Error("Expected block with invalid index to be invalid, but it isn't")
	}
	block := NewBlock(map[string]interface{}{"PrevHash": "SOME_BAD_HASH"})
	if bc.IsValidBlock(block) {
		t.Error("Expected block with invalid PrevHash to be invalid, but it isn't")
	}
	block := NewBlock(map[string]interface{}{"Hash": "BAD_HASH"})
	if bc.IsValidBlock(block) {
		t.Error("Expected block with invalid Hash to be invalid, but it isn't")
	}
	bc.AddNewBlock("testing")
	if !bc.IsValidBlock(bc.Head) {
		t.Error("Expected Head of blockchain to be valid, but it isn't")
	}
	block := NewBlock(map[string]interface{}{"Timestamp": time.Now().AddDate(0, 0, -1)})
	if bc.IsValidBlock(block) {
		t.Error("Expected block with invalid Timestamp to be invalid, but it isn't")
	}
}
