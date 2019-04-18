package main

import "fmt"

func (cli *CLI) generate(to string, amount int, nodeID string, offline bool) {
	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()
	tx := NewCoinbaseTX(to, "", amount)
	if offline {
		newBlock := bc.NewBlock([]*Transaction{tx})
		UTXOSet.Update(newBlock)
		fmt.Println("Success!")
	} else {
		sendTx(tx)
	}

}
