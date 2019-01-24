package main

import (
	"fmt"
	"log"
	"net"
)

const protocol = "udp"
const nodeVersion = 1
const commandLength = 12
const maxDatagramSize = 8192

var nodeAddress string
var knownNodes = []string{"224.0.0.1:9999"}


// StartServer starts a node
func StartServer(nodeID, minerAddress string, h func(*net.UDPAddr, int, []byte, *Blockchain)) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	addr, err := net.ResolveUDPAddr("udp", knownNodes[0]) // currently will always connect to a udp port
	if err != nil {
		log.Panic(err)
	}

	ln, err := net.ListenMulticastUDP(protocol,nil, addr)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := NewBlockchain(nodeID)

	// Disabled synchronization for now
	/*if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}*/

	ln.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := ln.ReadFromUDP(b)
		if err != nil {
			log.Panic("ReadFromUDP failed:", err)
		}
		h(src, n, b, bc)
	}

}