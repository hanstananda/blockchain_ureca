# Current Server Implementation

## Notary nodes

Currently, the notary nodes are hardcoded by denoting its NODE_ID. 

There is no checking about the authenticity of the NODE_ID itself as well. 

## Sending Commands 

### Commands available: 

*   **vote**
    
    Vote for a transaction to be accepted or rejected. Follows the datastructure described below. 

*   **tx**

    Send a new transaction to be included into the block. 
    

### Serialization and Deserialization

For command serialization and deserialization, there are two additional functions named `commandToBytes` and `bytesToCommand`. 
As its name suggests, it will convert the command to be sent, which is in string, into series of bytes and vice-versa. 

Meanwhile, the data is encoded and decoded using `gob` library. It requires an encoder function named `gobEncode`, 
which takes arbitrary datatype and return series of bytes. 

### Network Packet Sending 

Each packet will consists of 2 sections, namely command and data. The command length is **12 bytes**. 
Meanwhile, the size of the data is currently arbitrary. 
However, the size of the packet is currently limited to **8192 bytes**. 

## Other Internal Implementations 

### Datastructures used

*   **Tally**

    This is used by the Notary Node to tally the voting of a transaction.
    ```go 
    type Tally struct {
        Yes 	int
        No 	int
        Nodes 	map[string]bool
    }
    ```

*   **Transaction**
    
    This is used to send and receive transactions from other nodes. 
    ```go 
    type Tx struct {
        AddFrom     string
        Transaction []byte
    }
    ```

*   **Vote**
    
    This is used as a format to send the voting to the Notary nodes.
    ```go 
    type Vote struct {
        AddFrom string
        ID      string
        Result  bool
    }
    ```
