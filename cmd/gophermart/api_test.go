package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCookies(t *testing.T) {
	handler := &Handler{
		Mux: chi.NewMux(),
		Cursor: &Cursor{
			NewMock(),
		},
	}
	handler.Post("/api/user/register", handler.RegisterUser)
	handler.Post("/api/user/login", handler.Login)
	ts := httptest.NewServer(handler)

	defer ts.Close()

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	encoder.Encode(&UserInfo{
		Username: "test",
		Password: "test",
	})
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buff)
	request.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, request)
	res := w.Result()

	assert.Equal(t, res.StatusCode, 200)
}
