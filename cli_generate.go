package main

import "fmt"

func (cli *CLI) generate(to string, amount int, nodeID string) {
	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	tx := NewCoinbaseTX(to, "", amount)
	newBlock := bc.NewBlock([]*Transaction{tx})
	UTXOSet.Update(newBlock)
	fmt.Println("Success!")
}
