package main

import (
	"github.com/snkim/sncoin/blockchain"
	"github.com/snkim/sncoin/cli"
)

func main(){
	blockchain.Blockchain()
	cli.Start()
}