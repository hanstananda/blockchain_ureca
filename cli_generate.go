package main

import "fmt"

func (cli *CLI) generate(to string, amount int) {
	bc := NewBlockchain()
	defer bc.db.Close()

	tx := NewCoinbaseTX(to, "", amount)
	bc.NewBlock([]*Transaction{tx})
	fmt.Println("Success!")
}
