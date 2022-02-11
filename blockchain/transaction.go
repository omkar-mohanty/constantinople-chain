package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (txn *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(txn)
	Handle(err)
	hash = sha256.Sum256(encoded.Bytes())
	txn.ID = hash[:]
}
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = to
	}
	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{1000000, to}
	txn := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	txn.SetId()
	return &txn
}
func (txn *Transaction) IsCoinbase() bool {
	return len(txn.Inputs) == 0 && len(txn.Inputs[0].ID) == 0 && txn.Inputs[0].Out == -1
}
func NewTransaction(
	from, to string,
	amount int,
	chain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	acc, validOutputs := chain.FindSpendabaleOutputs(from, amount)
	if acc < amount {
		log.Panic("Insufficient funds in account")
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		HandleErr(err)
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}
	txn := &Transaction{nil, inputs, outputs}
	txn.SetId()
	return txn
}
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
