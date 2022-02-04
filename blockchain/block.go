package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.nonce = nonce
	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
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
		log.Panic(err)
	}
}
