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
var target_group = ""

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

type TallyResult struct{
	ID []byte
	Yes int
	No int
	Result bool
}

type Tx struct {
	AddFrom     string
	Transaction []byte
}

type Vote struct {
	AddFrom string
	ID []byte
	Result  bool
}

type RequestVote struct {
	TxID []byte
}

type InitVote struct {
	TxID []byte
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

func sendData(data []byte) {
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

func sendTx(tnx *Transaction) {
	data := Tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)
	sendData(request)
}

func SendID(nodeID string){
	payload := commandToBytes(nodeID)
	request := append(commandToBytes("syn"), payload...)
	sendData(request)
}

func SendVote(nodeID string,ID []byte, result bool){
	data := Vote{
		AddFrom: nodeID,
		ID: ID,
		Result: result,}
	//fmt.Println(data)
	payload := gobEncode(data)
	//fmt.Println(payload)
	request := append(commandToBytes("vote"), payload...)
	//fmt.Println(request)
	sendData(request)
}

func SendTxs(nodeID string){
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
		sendTx(tx)
	}

}

// StartServer starts a node
func StartServer(nodeID, portUDP string, h func(*net.UDPAddr, int, []byte, *Blockchain)) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	target_group = fmt.Sprintf("%s:%s",addrUDP,portUDP)
	selfID = nodeID
	addr, err := net.ResolveUDPAddr("udp", target_group) // currently will always connect to a udp port
	if err != nil {
		log.Panic(err)
	}

	ln, err := net.ListenMulticastUDP(protocol,nil, addr)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := NewBlockchain(nodeID)

	//if isNotary(nodeID)==false{
	//	SendVote(nodeID, []byte("test"), true)
	//}
	if nodeID=="3001"	{
		txid, err1 := hex.DecodeString("430344e951a895eafaa1889acac1fd1c6ea11f421aedbc7f3c3aae7a11e6a00f")
		if err1 != nil {
			log.Panic(err)
		}
		tx, err2 := bc.FindTransaction(txid)
		if err2 != nil {
			log.Panic(err)
		}
		sendTx(&tx)
	}

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
		case "tx":
			handleTx(request, bc)
		case "rv":
			handleRequestVote(request,bc)
		case "iv":
		case "tally":
			handleTallyResult(request,bc)
		default:
			fmt.Println("Unknown command!")
		}
		return
	}
	switch command {
	case "vote": // No OP
	case "tx":
		handleTx(request, bc)
	case "rv":
	case "iv":
		handleInitVote(request,bc)
	case "tally":
		handleTallyResult(request,bc)
	default:
		fmt.Println("Unknown command!")
	}
}


func sendRequestVote(tx *Transaction){
	data := RequestVote{
		TxID: tx.ID,}
	payload := gobEncode(data)
	request := append(commandToBytes("rv"), payload...)
	sendData(request)
}

func handleRequestVote(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload RequestVote
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txid_str := hex.EncodeToString(payload.TxID)
	if _, ok := votePool[txid_str]; ok {
		// Transaction is currently being voted, no need to initiate another one
		return
	}
	if tx, ok := mempool[txid_str]; ok {
		// Find the transaction in the memory then initiate the vote
		sendInitVote(&tx)
		var val Tally
		val.Nodes = make(map[string]bool)
		votePool[txid_str]= val
	}
}

func sendInitVote(tx *Transaction){
	data := InitVote{
		TxID: tx.ID,}
	payload := gobEncode(data)
	request := append(commandToBytes("iv"), payload...)
	sendData(request)
}

func handleInitVote(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload InitVote
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txID := payload.TxID
	_, er := bc.FindTransaction(txID)
	tx := mempool[hex.EncodeToString(txID)]
	if(er!=nil){
		if(bc.VerifyTransaction(&tx)){
			SendVote(selfID, txID, true)
		} else{
			SendVote(selfID, txID, false)
		}
	}
	// Send the vote
	SendVote(selfID, txID, true)
}

func sendTallyResult(res *TallyResult){
	payload := gobEncode(res)
	request := append(commandToBytes("tally"), payload...)
	sendData(request)
}

func handleTallyResult(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload TallyResult
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.Result ==true{
		// Transaction is accepted by all majority, put in blockchain
		_,err1 := bc.FindTransaction(payload.ID)
		fmt.Printf("Transaction %s accepted!\n",hex.EncodeToString(payload.ID))
		if err1 != nil {
			var txs []*Transaction
			tx := mempool[hex.EncodeToString(payload.ID)]
			txs = append(txs, &tx)
			newBlock := bc.NewBlock(txs)
			UTXOSet := UTXOSet{bc}
			UTXOSet.Reindex()
			fmt.Printf("New block %x is created\n",newBlock.Hash)
		}
	}	else{
		fmt.Printf("Transaction %s rejected!\n",hex.EncodeToString(payload.ID))
	}
	// Delete the transaction from memory after voting is done
	delete(mempool, hex.EncodeToString(payload.ID))
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
	txid_str := hex.EncodeToString(txid)
	result := payload.Result
	if val, ok := votePool[txid_str]; ok {
		if val.Nodes[voter]==true {
			fmt.Printf("Transaction %s has been voted by %s!\n", txid_str, voter)
			return
		}

		if(result == true){
			val.Yes ++

		}else {
			val.No++
		}
		val.Nodes[voter] = true
		fmt.Println(result,val.Yes,val.No,numNodes,val.Nodes)
		if err != nil {
			log.Fatal(err)
		}
		if val.Yes> numNodes/2+1 {
			fmt.Printf("Transaction %s accepted!\n", txid_str)
			res := TallyResult{
				ID: txid,
				Yes: val.Yes,
				No:val.No,
				Result:true,
			}
			sendTallyResult(&res)
			// Remove from pool, free up memory
			delete(votePool, txid_str)
			// Do the acceptance here
		}	else if(val.Yes + val.No == numNodes)		{
			fmt.Printf("Transaction %s rejected!\n",txid_str)
			res := TallyResult{
				ID: txid,
				Yes: val.Yes,
				No:val.No,
				Result:false,
			}
			sendTallyResult(&res)
			// Remove from pool, free up memory
			delete(votePool, txid_str)
			// Do the transaction deletion here
		}
		// Update Pool
		votePool[txid_str] = val
	}	else{
		// VotePool not initialized, initvote is not called before voting commences
		log.Fatal(ok)
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

	var txs []*Transaction

	for id := range mempool {
		tx := mempool[id]
		_,err := bc.FindTransaction(tx.ID)
		if err != nil {
			// Transaction does not exist yet
			if isNotary(payload.AddFrom)==true{
				// if the transaction if from notary, it is considered always correct
				txs = append(txs, &tx)
				delete(mempool, hex.EncodeToString(tx.ID))
			} else{
				fmt.Printf(tx.String())
				//  check with other nodes whether transaction is valid
				sendRequestVote(&tx)
			}
		}
	}

	//cbTx := NewCoinbaseTX(miningAddress, "")
	// txs = append(txs, cbTx)
	if len(txs) == 0 {
		return;
	}
	newBlock := bc.NewBlock(txs)
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Printf("New block %x is created\n",newBlock.Hash)

}