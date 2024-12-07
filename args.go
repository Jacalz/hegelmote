package main

import (
	"errors"
	"flag"
)

var (
	errTooFewArgs  = errors.New("too few arguments")
	errTooManyArgs = errors.New("too many arguments")

	errNoTargetIP = errors.New("no target ip found")
)

type arguments struct {
	ip   string
	port string
}

func parseArguments() (arguments, error) {
	flag.Parse()

	if len(flag.Args()) == 0 {
		return arguments{}, errTooFewArgs
	} else if len(flag.Args()) > 2 {
		return arguments{}, errTooManyArgs
	}

	ip := flag.Arg(0)
	if ip == "" {
		return arguments{}, errNoTargetIP
	}

	port := flag.Arg(1)
	if port == "" {
		port = "50001"
	}

	return arguments{ip: ip, port: port}, nil
}
