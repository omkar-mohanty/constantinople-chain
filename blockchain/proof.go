package blockchain

import "math/big"

//Proof of work steps
//Take a block of data
//Start a nonce with initial value of 0
//Compute the hash of the data + nonce
//See if the hash meets a number of requirements
//Requirements:
//The hash must start with a certain number of zeros
const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	BigInt *big.Int
}
