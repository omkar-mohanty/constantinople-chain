package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/omkar-mohanty/golang-blockchain/blockchain"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOut
}
type TxOut struct {
	Value  int
	PubKey string
}
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (txn *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(txn)
	blockchain.HandleErr(err)
	hash = sha256.Sum256(encoded.Bytes())
	txn.ID = hash[:]
}
func CoinbaseTx(to, from string) *Transaction {
	if from == "" {
		from = fmt.Sprintf("%s", to)
	}
	txIn := TxInput{[]byte{}, 100, from}
	txOut := TxOut{100, to}
	txn := Transaction{nil, []TxInput{txIn}, []TxOut{txOut}}
	txn.SetId()
	return &txn
}
func (txn *Transaction) IsCoinbase() bool {
	return len(txn.Inputs) == 0 && len(txn.Inputs[0].ID) == 0 && txn.Inputs[0].Out == -1
}
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}
func (out *TxOut) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
