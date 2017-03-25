package main

import (
	"log"
	"simple-blockchain/blockchain"
)

func main() {
	blockchain := blockchain.GetBlockchain()
	log.Println(blockchain.Head)
	blockchain.AddNewBlock("Hello World!")
	log.Println(blockchain.Head)
	log.Println("Valid? ", blockchain.IsValid())

}
