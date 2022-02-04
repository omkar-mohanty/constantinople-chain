package main

import (
	"github.com/omkar-mohanty/golang-blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()
	chain.AddBlock("First Block")
	chain.AddBlock("Second Block")
}
