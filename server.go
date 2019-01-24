package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

const protocol = "udp"
const nodeVersion = 1
const commandLength = 12
const maxDatagramSize = 8192

var nodeAddress string
var knownNodes = []string{"224.0.0.1:9999"}
var mempool = make(map[string]Transaction)


type tx struct {
	AddFrom     string
	Transaction []byte
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}


func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func sendData(data []byte) {
	addr, err := net.ResolveUDPAddr(protocol, knownNodes[0])
	conn, err := net.DialUDP(protocol,nil, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendTx(tnx *Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(request)
}

// StartServer starts a node
func StartServer(nodeID string, h func(*net.UDPAddr, int, []byte, *Blockchain)) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	addr, err := net.ResolveUDPAddr("udp", knownNodes[0]) // currently will always connect to a udp port
	if err != nil {
		log.Panic(err)
	}

	if(nodeID=="3000")	{
		bc := NewBlockchain(nodeID)
		bci := bc.Iterator()
		for {
			block := bci.Next()

			for _, tx := range block.Transactions {
				// send the transactions to all parties in the group
				sendTx(tx)
			}

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
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

func handleConnection(conn *net.UDPAddr, n int, b []byte, bc *Blockchain) {
	//request, err := ioutil.ReadAll(conn)
	//if err != nil {
	//	log.Panic(err)
	//}
	request := b
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)
	// log.Println(request)  // hard debug

	switch command {
	case "tx":
		handleTx(request, bc)
	default:
		fmt.Println("Unknown command!")
	}
}

func handleTx(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if len(mempool) >= 1 {
	VerifyTransactions:
		var txs []*Transaction

		for id := range mempool {
			tx := mempool[id]
			if bc.VerifyTransaction(&tx) {
				txs = append(txs, &tx)
			}
		}

		if len(txs) == 0 {
			fmt.Println("All transactions are invalid! Waiting for new ones...")
			return
		}

		//cbTx := NewCoinbaseTX(miningAddress, "")
		// txs = append(txs, cbTx)

		newBlock := bc.NewBlock(txs)
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()

		fmt.Println("New block %x is created",newBlock.Hash)

		for _, tx := range txs {
			txID := hex.EncodeToString(tx.ID)
			delete(mempool, txID)
		}

		if len(mempool) > 0 {
			goto VerifyTransactions
		}
	}
}