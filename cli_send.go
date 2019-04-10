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

	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet)
	// Currently disable generating new coins for each transaction performed
	//cbTx := NewCoinbaseTX(from, "",10)
	//txs := []*Transaction{cbTx, tx}
	txs := []*Transaction{tx}
	newBlock := bc.NewBlock(txs)
	UTXOSet.Update(newBlock)
	fmt.Println("Adding to local blockchain success!")
}
