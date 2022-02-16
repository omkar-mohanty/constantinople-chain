package blockchain

import (
	"bytes"

	"github.com/omkar-mohanty/golang-blockchain/wallet"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{Value: value, PubKeyHash: nil}
	txo.Lock([]byte(address))
	return txo
}
func (in *TxInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Equal(lockingHash, publicKeyHash)
}
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}
func (out *TxOutput) IsLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, publicKeyHash)
}
