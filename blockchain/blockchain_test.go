package blockchain

import "testing"

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

func TestReplaceChain(t *testing.T) {
	b1 := GetBlockchain()
	b2 := GetBlockchain()
	b2.AddNewBlock("extra block")
	b2.ReplaceChain(b1)
	if b1.Head != b2.Head {
		t.Error("Expected b1 to not replace b2, but it did")
	}
	b1.ReplaceChain(b2)
	if b1.Head != b2.Head {
		t.Error("Expected b1 to be replaced by b2, but it wasn't")
	}
}
