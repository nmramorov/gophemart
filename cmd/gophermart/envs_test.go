package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvConfig(t *testing.T) {
	testConfig, _ := NewEnvConfig()
	assert.Equal(t, testConfig.Address, "127.0.0.1:8080")
	assert.Equal(t, testConfig.DatabaseURI, "127.0.0.1:5432")
	assert.Equal(t, testConfig.Accrual, "127.0.0.1:8081")
}
