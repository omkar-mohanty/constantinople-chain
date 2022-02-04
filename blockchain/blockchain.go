package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

func (chain *Blockchain) AddBlock(data string) {
	err := chain.Database.Update(func(txn *badger.Txn) error {
		new := CreateBlock(data, chain.LastHash)
		err := txn.Set(new.Hash, new.Serealize())
		chain.LastHash = new.Hash
		return err
	})
	HandleErr(err)
}

func InitBlockChain() *Blockchain {
	var lastHash []byte
	opt := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opt)
	HandleErr(err)
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serealize())
			lastHash = genesis.Hash
			txn.Set([]byte("lh"), lastHash)
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			lastHash, err = item.ValueCopy([]byte{})
			return err
		}
	})
	HandleErr(err)
	return &Blockchain{lastHash, db}
}
