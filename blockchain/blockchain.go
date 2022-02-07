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

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *Blockchain) Iterator() func() *Block {
	db := chain.Database
	lastHash := chain.LastHash
	return func() *Block {
		var block *Block
		err := db.View(func(txn *badger.Txn) error {
			item, rerr := txn.Get(lastHash)
			if rerr != nil {
				return rerr
			}
			blockEncoded, rerr := item.ValueCopy([]byte{})
			block = Dserealize(blockEncoded)
			lastHash = block.PrevHash
			return rerr
		})
		if err == badger.ErrKeyNotFound {
			return nil
		}
		return block
	}
}
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleErr(err)
		lastHash, err = item.ValueCopy([]byte{})
		return err
	})
	HandleErr(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serealize())
		HandleErr(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	HandleErr(err)
}

func InitBlockChain() *Blockchain {
	var lastHash []byte
	opt := badger.DefaultOptions(dbPath)
	opt.Truncate = true
	opt.Logger = nil
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
			HandleErr(err)
			lastHash, err = item.ValueCopy([]byte{})
			return err
		}
	})
	HandleErr(err)
	return &Blockchain{lastHash, db}
}
