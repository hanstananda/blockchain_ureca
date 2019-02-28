# Blockchain URECA

## Getting started

These instructions will get you a copy of the project up and running on your local machine.

## Prerequisites

* The program has been tested to run on *Ubuntu Linux 18.04 x86_64*. It may not work on Windows operating system due to package dependencies

* Golang (the project is developed using version `go1.10.4`)
    >   Download from [Official Golang website](https://golang.org/dl/)

## Setup configuration

1. Copy `blockchain_genesis.db` into `blockchain_NODEID.db`. 
(NODEID) is the ID set after running `export NODE_ID=...`
    
    Example: 
    
    ```bash
    export NODE_ID=3000
    ```

2. Run `go build` inside the directory

3. Execute the executable file (by default is `./blockchain_ureca`) with parameters given below

## Available parameters 

*   No parameter
    
    Shows the list of the available parameters

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
7. The UDP Address is currently set to `224.0.0.1`, which is the typical address used for UDP multicasting in routers


## Author(s)

* [**Hans Tananda**](https://github.com/hanstananda) - *Initial work*

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Contribution

This project opens to any kind of contribution that can be done to enhance this project.


## Acknowledgments

*   This blockchain platform is implemented based on the blockchain created by Ivan Kuznetsov - [Jeiwan](https://github.com/Jeiwan)  
    Original repository: [Jeiwan/blockchain_go](https://github.com/Jeiwan/blockchain_go/tree/master)
*   `Readme.md` and `Contributing.md` are written based on templates created by Billie Thompson - [PurpleBooth](https://github.com/PurpleBooth)
*   [Elmo Huang Xuyun](elmohuang@ntu.edu.sg) as my supervisor of this project
*   [Prof. Lam Kwok Yan](kwokyan.lam@ntu.edu.sg) and [Assoc Prof. Wang Huaxiong](hxwang@ntu.edu.sg)as the Professor in charge of my project
*   [Bedi Jannat](jannat001@e.ntu.edu.sg), my colleague who works on similar project under the same supervisor
