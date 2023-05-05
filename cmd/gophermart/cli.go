package main

import "flag"

type CLIOptions struct {
	Address     string
	DatabaseURI string
	Accrual     string
}

func NewCliOptions() *CLIOptions {
	var address = flag.String("a", "localhost:8080", "server address")
	var accrual = flag.String("r", "10s", "accrual address")
	var database = flag.String("d", "2s", "database address")
	flag.Parse()

	return &CLIOptions{
		Address:     *address,
		DatabaseURI: *accrual,
		Accrual:     *database,
	}
}
