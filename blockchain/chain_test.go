package blockchain

import (
	"fmt"
	"testing"

	"github.com/dgraph-io/badger"
)

func TestChain(t *testing.T) {
	chain := InitBlockChain()
	chain.AddBlock("Hello")
	db := chain.Database
	err := db.View(func(txn *badger.Txn) error {
		item, rerr := txn.Get(chain.LastHash)
		HandleTestErr(t, rerr)
		blockEncoded, rerr := item.ValueCopy([]byte{})
		block := Dserealize(blockEncoded)
		pow := NewProof(block)
		if pow.Validate() == false {
			t.Fatalf("Not validated")
		}
		HandleTestErr(t, rerr)
		return rerr
	})
	HandleTestErr(t, err)
}
func TestIter(t *testing.T) {
	chain := InitBlockChain()
	chain.AddBlock("Hello")
	chain.AddBlock("Bye")
	next := chain.Iterator()
	block := next()
	fmt.Printf("Data %s\n", string(block.Data))
}
func HandleTestErr(t *testing.T, rerr error) {
	if rerr != nil {
		t.Error(rerr)
	}
}
