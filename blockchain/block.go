package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"

	"github.com/omkar-mohanty/golang-blockchain/blockchain/transaction"
)

type Block struct {
	Hash         []byte
	Transactions []*transaction.Transaction
	PrevHash     []byte
	Nonce        int
}

func CreateBlock(transactions []*transaction.Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, transactions, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func Genesis(coinbase *transaction.Transaction) *Block {
	return CreateBlock([]*transaction.Transaction{coinbase}, []byte{})
}
func (b *Block) Serealize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	HandleErr(err)
	return res.Bytes()
}

func Dserealize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	HandleErr(err)
	return &block
}

func HandleErr(err error) {
	if err != nil {
		log.Fatalf("Error: %s ", err.Error())
	}
}

func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	data := bytes.Join(txHashes, []byte{})
	txHash = sha256.Sum256(data)
	return txHash[:]
}
