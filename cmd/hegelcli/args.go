package main

import (
	"errors"
	"flag"

	"github.com/Jacalz/hegelmote/device"
)

var (
	errTooFewArgs  = errors.New("too few arguments")
	errTooManyArgs = errors.New("too many arguments")

	errNoTargetIP = errors.New("no target ip found")
)

type arguments struct {
	ip          string
	port        string
	interactive bool
	amplifier   device.Device
}

func parseArguments() (arguments, error) {
	args := arguments{}

	// Flags for starting in interactive mode.
	flag.BoolVar(&args.interactive, "i", false, "starts an interactive command terminal")
	flag.BoolVar(&args.interactive, "interactive", false, "starts an interactive command terminal")
	flag.UintVar(&args.amplifier, "device", 1, "sets the device to use for input mappings")
	flag.Parse()

	if flag.NArg() == 0 {
		return arguments{}, errTooFewArgs
	} else if flag.NArg() > 2 {
		return arguments{}, errTooManyArgs
	}

	args.ip = flag.Arg(0)
	if args.ip == "" {
		return arguments{}, errNoTargetIP
	}

	args.port = flag.Arg(1)
	if args.port == "" {
		args.port = "50001"
	}

	return args, nil
}
