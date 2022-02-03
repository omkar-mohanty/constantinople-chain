package main

import (
	"fmt"

	"github.com/omkar-mohanty/golang-blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()
	chain.AddBlock("First Block")
	chain.AddBlock("Second Block")
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash %x\n", block.PrevHash)
		fmt.Printf("Block Data %s\n", block.Data)
		fmt.Printf("Block Hash %x\n", block.Hash)

	}
}
