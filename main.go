package main

import (
	"github.com/omkar-mohanty/golang-blockchain/blockchain"
	"github.com/omkar-mohanty/golang-blockchain/cmd"
)

func main() {
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()
	cli := cmd.NewCmd(chain)
	cli.Run()
}
