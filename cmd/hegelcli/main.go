package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jacalz/hegelmote/remote"
)

func main() {
	args, err := parseArguments()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.PrintDefaults()
		os.Exit(2)
	}

	command := &remote.Control{}
	defer command.Disconnect()

	err = command.Connect(args.ip+":"+args.port, args.amplifier)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if args.interactive {
		runInteractiveMode(command, args.amplifier)
		return
	}

	flag.PrintDefaults()
}
