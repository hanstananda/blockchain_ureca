package main

import "fmt"

func (cli *CLI) generate(to string, amount int, nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	tx := NewCoinbaseTX(to, "", amount)
	bc.NewBlock([]*Transaction{tx})
	fmt.Println("Success!")
}
