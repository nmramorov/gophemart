package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	config "github.com/nmramorov/gophemart/internal/configuration"
	"github.com/nmramorov/gophemart/internal/logger"
	"github.com/nmramorov/gophemart/internal/models"
)

func TestApp(t *testing.T) {
	flags := config.NewCliOptions()
	envs, err := config.NewEnvConfig()
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}
	app, _ := NewApp(config.NewConfig(flags, envs))
	go app.Run()

	client := &http.Client{}

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	encoder.Encode(&models.UserInfo{
		Username: "test",
		Password: "test",
	})
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/login", buff)
	if err != nil {
		logger.ErrorLog.Fatalf("Error with login request: %e", err)
	}
	response, err := client.Do(request)
	if err != nil {
		logger.ErrorLog.Fatalf("Error with login response: %e", err)
	}
	defer response.Body.Close()

	assert.Equal(t, 401, response.StatusCode)

	buff = bytes.NewBuffer([]byte{})
	encoder = json.NewEncoder(buff)
	encoder.Encode(&models.UserInfo{
		Username: "test",
		Password: "test",
	})
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buff)
	if err != nil {
		logger.ErrorLog.Fatalf("Error with registration: %e", err)
	}
	response, err = client.Do(request)
	if err != nil {
		logger.ErrorLog.Fatalf("Error with registration response: %e", err)
	}
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)
}
