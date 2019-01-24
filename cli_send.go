package main

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int,nodeID string) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := NewCoinbaseTX(from, "",10)
	txs := []*Transaction{cbTx, tx}

	newBlock := bc.NewBlock(txs)
	UTXOSet.Update(newBlock)
	fmt.Println("Success!")
}
