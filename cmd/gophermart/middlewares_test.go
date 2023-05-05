package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookiesMiddleware(t *testing.T) {
	handler := NewHandler("http://localhost:8081", &Cursor{NewMock()})
	ts := httptest.NewServer(handler)

	defer ts.Close()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/balance", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, request)
	res := w.Result()
	assert.Equal(t, 401, res.StatusCode)
}
