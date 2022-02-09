package blockchain

import (
	"fmt"
	"testing"
)

func TestChain(t *testing.T) {
	chain := InitBlockChain("me")
	fmt.Printf("%x", chain.LastHash)
}

func HandleTestErr(t *testing.T, rerr error) {
	if rerr != nil {
		t.Error(rerr)
	}
}
