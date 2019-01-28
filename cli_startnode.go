package main

import (
	"fmt"
)

func (cli *CLI) startNode(nodeID, portUDP string) {
	fmt.Printf("Starting node %s\n", nodeID)
	StartServer(nodeID, portUDP, handleConnection)
}
