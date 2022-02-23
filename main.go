package main

import (
	"github.com/snkim/sncoin/cli"
	"github.com/snkim/sncoin/db"
)



func main() {
  defer db.Close()
  db.InitDB()
  cli.Start()    
}
