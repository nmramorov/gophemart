package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	flags := NewCliOptions()
	envs, err := NewEnvConfig()
	if err != nil {
		ErrorLog.Fatal(err)
	}
	app := NewApp(NewConfig(flags, envs))
	go app.Run()

	client := &http.Client{}

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	encoder.Encode(&UserInfo{
		Username: "test",
		Password: "test",
	})
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/login", buff)
	if err != nil {
		ErrorLog.Fatalf("Error with login request: %e", err)
	}
	response, err := client.Do(request)
	if err != nil {
		ErrorLog.Fatalf("Error with login response: %e", err)
	}
	defer response.Body.Close()

	assert.Equal(t, 401, response.StatusCode)

	buff = bytes.NewBuffer([]byte{})
	encoder = json.NewEncoder(buff)
	encoder.Encode(&UserInfo{
		Username: "test",
		Password: "test",
	})
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buff)
	if err != nil {
		ErrorLog.Fatalf("Error with registration: %e", err)
	}
	response, err = client.Do(request)
	if err != nil {
		ErrorLog.Fatalf("Error with registration response: %e", err)
	}
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)
}
