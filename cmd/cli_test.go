package cmd

import (
	"testing"

	"github.com/omkar-mohanty/golang-blockchain/blockchain"
)

func TestCli(t *testing.T) {
	chain := blockchain.InitBlockChain()
	cli := CommandLine{chain}
	cli.AddBlock("NMS")
}
