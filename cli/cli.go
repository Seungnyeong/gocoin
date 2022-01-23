package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/snkim/sncoin/explorer"
	"github.com/snkim/sncoin/rest"
)

func usage() {
	fmt.Printf("Welcome to sncoin\n")
	fmt.Printf("Please use the following flags:\n")
	fmt.Printf("-port: Set the PORT of the server\n")
	fmt.Printf("-mode: Choose between 'html' and 'rest'\n\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}
	port := flag.Int("port", 4000, "Set port of the server")

	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		//start rest api
		rest.Start(*port)
	case "html":
		// start html explorer
		explorer.Start(*port)
	default:
		usage()
	}

	fmt.Println(*port, *mode)
}