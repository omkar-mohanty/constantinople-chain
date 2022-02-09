package transaction

import (
	"errors"
	"testing"
)

func TestCoinbaseTxn(t *testing.T) {
	testTxn := CoinbaseTx("testAdd", "***")
	txnInput := testTxn.Inputs[0]
	txnOutput := testTxn.Outputs[0]
	if len(txnInput.ID) != 0 {
		err := errors.New("ID must be empty")
		t.Fatalf("Failed %s", err.Error())
	}
	if !txnOutput.CanBeUnlocked("testAdd") {
		err := errors.New("account must unlock Output")
		t.Fatalf("Failed %s", err.Error())
	}
}
