package mycli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/omkar-mohanty/golang-blockchain/blockchain"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage")
	fmt.Println("getbalance -address ADDRESS - get the balance for the address")
	fmt.Println("createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount of coins")
}
func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}
func (cli *CommandLine) printChain() {
	fmt.Println()
	chain := *blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	next := chain.Iterator()
	block := next()
	for {

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
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
func (cli *CommandLine) createBlockchain(address string) {
	chain := *blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished")
}
func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer chain.Database.Close()
	unspentUTXOs := chain.FindUTXO(address)
	balance := 0
	for _, out := range unspentUTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of %s : %d \n", address, balance)
}
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()
	txn := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{txn})
	fmt.Println("Success")
}
func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
func NewCmd() *CommandLine {
	return &CommandLine{}
}
