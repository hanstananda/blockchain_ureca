package main

import (
	"strings"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

const protocol = "udp"
const nodeVersion = 1
const commandLength = 12
const maxDatagramSize = 8192

var numNodes = 3
var selfID = ""
var nodeAddress string
var addrUDP = "224.0.0.1"
var NotaryNodes = []string{"3000"}
var mempool = make(map[string]Transaction)

func isNotary(nodeID string) bool{
	for _, n := range NotaryNodes{
		if nodeID == n{
			fmt.Println("This is the notary node!")
			return true
		}
	}
	return false
}

type tally struct {
	yes 		int
	no 			int
	nodes 		[]string
}

type tx struct {
	AddFrom     string
	Transaction []byte
}

type vote struct {
	AddFrom string
	result    string
}

var votePool = make(map[string]tally)

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

func sendData(data []byte, target_group string) {
	addr, err := net.ResolveUDPAddr(protocol, target_group)
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

func sendTx(tnx *Transaction, target_group string) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)
	sendData(request, target_group)
}

func SendID(nodeID ,target_group string){
	payload := commandToBytes(nodeID)
	request := append(commandToBytes("syn"), payload...)
	sendData(request, target_group)
}

func SendVote(nodeID, target_group string,ID []byte, result bool){
	data := vote{nodeID+","+string(ID[:])+","+strconv.FormatBool(result), "abc" }
	fmt.Println(data)
	payload := gobEncode(vote{nodeID+","+string(ID[:]), strconv.FormatBool(result) })
	fmt.Println(payload)
	request := append(commandToBytes("vot"), payload...)
	fmt.Println(request)
	var buff bytes.Buffer
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	var res vote
	err := dec.Decode(&res)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(buff)
	fmt.Println(res)
	sendData(request,target_group)
}

func SendTxs(nodeID string,target_group string){
	bc := NewBlockchain(nodeID)
	bci := bc.Iterator()
	var txs []*Transaction
	for {
		block := bci.Next()
		var curtxs[]*Transaction
		for _, tx := range block.Transactions {
			curtxs = append([]*Transaction{tx}, curtxs...)
		}
		txs = append(txs,curtxs...)
		curtxs = nil
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	// reverse transaction sending
	for i, j := 0, len(txs)-1; i < j; i, j = i+1, j-1 {
		txs[i], txs[j] = txs[j], txs[i]
	}
	for _,tx := range txs{
		// send the transactions to all parties in the group
		fmt.Println(tx)
		sendTx(tx,target_group)
	}

}

// StartServer starts a node
func StartServer(nodeID, portUDP string, h func(*net.UDPAddr, int, []byte, *Blockchain)) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	target_group := fmt.Sprintf("%s:%s",addrUDP,portUDP)
	selfID = nodeID
	addr, err := net.ResolveUDPAddr("udp", target_group) // currently will always connect to a udp port
	if err != nil {
		log.Panic(err)
	}

	if isNotary(nodeID)==false{
		SendVote(nodeID, target_group, []byte("test"), true)
	}
	//if nodeID=="3000"	{
	//	SendTxs(nodeID,target_group)
	//}

	ln, err := net.ListenMulticastUDP(protocol,nil, addr)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := NewBlockchain(nodeID)

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
	log.Println(request)  // hard debug

	if isNotary(selfID)==true{
		switch command {
		case "vt":
			handleVote(request,bc)
			return
		}
	}
	switch command {
	case "tx":
		handleTx(request, bc)
	default:
		fmt.Println("Unknown command!")
	}
}

func handleVote(request []byte, bc * Blockchain){
	var buff bytes.Buffer
	var payload vote

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(payload)
	tmp := strings.Split(payload.AddFrom,",")
	TxID := tmp[0]
	result := tmp[1]
	if val, ok := votePool[TxID]; ok {
		if(result =="true"){
			val.yes ++

		}else {
			val.no++
		}
		val.nodes = append(val.nodes, payload.AddFrom)
		fmt.Println(result,val.yes,val.no,numNodes,val.nodes)
		if val.yes> numNodes/2+1{
			fmt.Printf("Transaction %s accepted!\n",TxID)
		}
	}	else{
		// New payload
		var val tally
		if(result =="true"){
			val.yes++
		}		else{
			val.no++
		}
		val.nodes = append(val.nodes, payload.AddFrom)
		votePool[TxID]= val
		fmt.Println(result,val.yes,val.no,numNodes,val.nodes)
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

		fmt.Printf("New block %x is created\n",newBlock.Hash)

		for _, tx := range txs {
			txID := hex.EncodeToString(tx.ID)
			delete(mempool, txID)
		}

		if len(mempool) > 0 {
			goto VerifyTransactions
		}
	}
}