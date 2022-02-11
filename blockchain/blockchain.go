package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp"
	dbFile      = "./tmp/MANIFEST"
	genesisData = "First transaction in genesis block"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
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
func (chain *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleErr(err)
		lastHash, err = item.ValueCopy([]byte{})
		return err
	})
	HandleErr(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serealize())
		HandleErr(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	HandleErr(err)
}
func ContinueBlockchain(address string) *Blockchain {
	var lastHash []byte
	if !DBexists() {
		fmt.Println("Blockchain does not exist")
		runtime.Goexit()
	}
	opt := badger.DefaultOptions(dbPath)
	opt.Truncate = true
	opt.Logger = nil
	db, err := badger.Open(opt)
	HandleErr(err)
	err = db.Update(func(txn *badger.Txn) error {
		item, rerr := txn.Get([]byte("lh"))
		HandleErr(rerr)
		lastHash, rerr = item.ValueCopy([]byte{})
		return rerr
	})
	HandleErr(err)
	return &Blockchain{lastHash, db}
}
func InitBlockChain(address string) *Blockchain {
	var lastHash []byte
	if DBexists() {
		fmt.Println("Blockchain Already exists")
		runtime.Goexit()
	}
	opt := badger.DefaultOptions(dbPath)
	opt.Truncate = true
	opt.Logger = nil
	db, err := badger.Open(opt)
	HandleErr(err)
	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serealize())
		lastHash = genesis.Hash
		txn.Set([]byte("lh"), lastHash)
		return err
	})
	HandleErr(err)
	return &Blockchain{lastHash, db}
}

func (chain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int)
	next := chain.Iterator()
	for {
		block := next()
		if block == nil {
			break
		}
		for _, txs := range block.Transactions {
			txID := hex.EncodeToString(txs.ID)
		Outputs:
			for outIdx, out := range txs.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *txs)
				}
			}
			if !txs.IsCoinbase() {
				for _, in := range txs.Inputs {
					if in.CanUnlock(address) {
						inTxId := hex.EncodeToString(in.ID)
						spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Out)
					}
				}
			}

		}
	}
	return unspentTxs
}

func (chain *Blockchain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTxns := chain.FindUnspentTransactions(address)
	for _, unspentTxn := range unspentTxns {
		for _, output := range unspentTxn.Outputs {
			if output.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, output)
			}
		}
	}
	return UTXOs
}

func (chain *Blockchain) FindSpendabaleOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0
	unspentTxn := chain.FindUnspentTransactions(address)
Work:
	for _, txn := range unspentTxn {
		txId := hex.EncodeToString(txn.ID)
		for outIdx, output := range txn.Outputs {
			if output.CanBeUnlocked(address) && accumulated < amount {
				accumulated += output.Value
				unspentOuts[txId] = append(unspentOuts[txId], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}
