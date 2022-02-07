package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/omkar-mohanty/golang-blockchain/blockchain"
)

type CommandLine struct {
	Blockchain *blockchain.Blockchain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage")
	fmt.Println("add -block BLOCK DATA adds block to the chain")
	fmt.Println("print -Prints blocks in the chain")
}
func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}
func (cli *CommandLine) AddBlock(data string) {
	cli.Blockchain.AddBlock(data)
	fmt.Println("Added Block")
}
func (cli *CommandLine) printChain() {
	fmt.Println()
	next := cli.Blockchain.Iterator()
	block := next()
	for {

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("\n\n")
		block = next()
		if block == nil {
			break
		}
	}
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
func NewCmd(chain *blockchain.Blockchain) *CommandLine {
	return &CommandLine{chain}
}
