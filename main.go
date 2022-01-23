package main

import (
	"github.com/snkim/sncoin/explorer"
	"github.com/snkim/sncoin/rest"
)


func main(){
	go explorer.Start(3000)
	rest.Start(4000)
}