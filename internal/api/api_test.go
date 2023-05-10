package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/mocks"
	"github.com/nmramorov/gophemart/internal/models"
)

func TestCookies(t *testing.T) {
	handler := &Handler{
		Mux: chi.NewMux(),
		Cursor: &db.Cursor{
			DBInterface: mocks.NewMock(),
		},
	}
	handler.Post("/api/user/register", handler.RegisterUser)
	handler.Post("/api/user/login", handler.Login)
	ts := httptest.NewServer(handler)

	defer ts.Close()

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	encoder.Encode(&models.UserInfo{
		Username: "test",
		Password: "test",
	})
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buff)
	request.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, request)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, 200)
}
