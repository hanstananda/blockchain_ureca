# Blockchain URECA

## Setup configuration

1. Copy `blockchain_genesis.db` into `blockchain_NODEID.db`. 
(NODEID) is the ID set after running `export NODE_ID=...`
    
    Example: 
    
    ```bash
    export NODE_ID=3000
    ```

2. Run `go build` inside the directory

3. Execute the executable file (by default is `./blockchain_ureca`)

## Available parameters 

*   `createblockchain -address ADDRESS`

    Create a blockchain and send genesis block reward to `ADDRESS`

*   `createwallet`
 
    Generates a new key-pair and saves it into the wallet file

*   `getbalance address ADDRESS`
    
    Get balance of `ADDRESS`

*   `listaddresses` 
    
    Lists all addresses from the wallet file

*   `printchain` 

    Print all the blocks of the blockchain

*   `reindexutxo` 

    Rebuilds the UTXO set

*   `send -from FROM -to TO -amount AMOUNT` 

    Send `AMOUNT` of coins from `FROM` address to `TO`

*   `generate -to TO -amount AMOUNT` 

    Generate `AMOUNT` of coins to `TO`

*   `startnode -port PORT_NO` 

    Start a node with ID specified in `NODE_ID` env at target port `PORT_NO`

## Additional Notes

1. Currently, nodes within the same group must manually connect to same port number to perform multicasting
2. There is no security measurements if outside party sniff into the current multicast port
3. Everyone within the multicast group will receive everything on the multicast port, including the packets that the sender has sent
4. No consensus protocol is currently implemented. Therefore, it will accept every transactions it received from the multicast group
5. No synchronization protocol is implemented yet
6. Currently, the node that sends all the transactions to the parties within the group is hardcoded to node ID **3000**
