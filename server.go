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

var numNodes = 3
var selfID = ""
var nodeAddress string
var addrUDP = "224.0.0.1"
var NotaryNodes = []string{"3000"}
var mempool = make(map[string]Transaction)
var notary_checked = false

func isNotary(nodeID string) bool{
	for _, n := range NotaryNodes{
		if nodeID == n{
			if !notary_checked{
				fmt.Println("This is the notary node!")
				notary_checked = true
			}

			return true
		}
	}
	return false
}

type Tally struct {
	Yes 		int
	No 			int
	Nodes 		map[string]bool
}

type Tx struct {
	AddFrom     string
	Transaction []byte
}

type Vote struct {
	AddFrom string
	ID string
	Result  bool
}

var votePool = make(map[string]Tally)

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
	data := Tx{nodeAddress, tnx.Serialize()}
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
	data := Vote{
		AddFrom: nodeID,
		ID: string(ID),
		Result: result,}
	//fmt.Println(data)
	payload := gobEncode(data)
	//fmt.Println(payload)
	request := append(commandToBytes("vote"), payload...)
	//fmt.Println(request)
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
	//log.Println(request)  // hard debug

	if isNotary(selfID)==true{
		switch command {
		case "vote":
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
	var payload Vote
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(payload)
	voter := payload.AddFrom
	txid := payload.ID
	result := payload.Result
	if val, ok := votePool[txid]; ok {
		if val.Nodes[voter]==true {
			fmt.Printf("Transaction %s has been voted by %s!\n", txid, voter)
			return
		}

		if(result == true){
			val.Yes ++

		}else {
			val.No++
		}
		val.Nodes[voter] = true
		fmt.Println(result,val.Yes,val.No,numNodes,val.Nodes)
		if val.Yes> numNodes/2+1{
			fmt.Printf("Transaction %s accepted!\n",txid)
			// Remove from pool, free up memory
			delete(votePool, txid)
			// Do the acceptance here
		}
		// Update Pool
		votePool[txid] = val
	}	else{
		// New payload
		var val Tally
		if(result == true){
			val.Yes++
		}		else{
			val.No++
		}
		val.Nodes = make(map[string]bool)
		val.Nodes[voter] = true
		votePool[txid]= val
		fmt.Println(result,val.Yes,val.No,numNodes,val.Nodes)
	}

}

func handleTx(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload Tx

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