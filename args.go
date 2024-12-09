package main

import (
	"errors"
	"flag"
	"fmt"
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
}

func parseArguments() (arguments, error) {
	args := arguments{}

	// Flags for starting in interactive mode.
	flag.BoolVar(&args.interactive, "i", false, "starts an interactive command terminal")
	flag.BoolVar(&args.interactive, "interactive", false, "starts an interactive command terminal")

	flag.Parse()

	if len(flag.Args()) == 0 {
		return arguments{}, errTooFewArgs
	} else if len(flag.Args()) > 2 {
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

	fmt.Println(args)
	return args, nil
}
