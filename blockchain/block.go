package blockchain

type Blockchain struct {
	Blocks []*Block
}
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.nonce = nonce
	return block
}
func (chain *Blockchain) AddBlock(data string) {
	block := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, block.Hash)
	chain.Blocks = append(chain.Blocks, new)
}
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}
func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}
