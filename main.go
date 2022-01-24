package main

import (
	"github.com/snkim/sncoin/blockchain"
)

func main(){
	blockchain.Blockchain().AddBlock("First")
	blockchain.Blockchain().AddBlock("Second")
	blockchain.Blockchain().AddBlock("Third")
}