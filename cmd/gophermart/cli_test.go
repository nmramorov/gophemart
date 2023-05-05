package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentCLI(t *testing.T) {
	os.Args = []string{"main.go", "-a", "localhost:4444", "-d=11s", "-r=5m"}
	var address = flag.String("a", "localhost:8080", "server address")
	var accrual = flag.String("r", "10s", "accrual address")
	var database = flag.String("d", "2s", "database address")
	flag.Parse()

	args := &CLIOptions{
		Address:     *address,
		DatabaseURI: *database,
		Accrual:     *accrual,
	}

	assert.Equal(t, "localhost:4444", args.Address)
	assert.Equal(t, "11s", args.DatabaseURI)
	assert.Equal(t, "5m", args.Accrual)
}
