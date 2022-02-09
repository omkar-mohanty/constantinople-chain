package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
)

//Proof of work steps
//Take a block of data
//Start a nonce with initial value of 0
//Compute the hash of the data + nonce
//See if the hash meets a number of requirements
//Requirements:
//The hash must start with a certain number of zeros
const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	return &ProofOfWork{Block: b, Target: target}
}
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.HashTransactions(),
		pow.Block.PrevHash,
		ToHex(int64(nonce)),
		ToHex(int64(Difficulty)),
	}, []byte{})
	return data
}
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}
