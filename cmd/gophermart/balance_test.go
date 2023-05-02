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

func TestBalanceGet(t *testing.T) {
	expectedBalance := &Balance{
		Current:   500.5,
		Withdrawn: 42,
	}
	type want struct {
		code int
	}
	type arguments struct {
		url string
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Test Positive balance get",
			want: want{
				code: 200,
			},
			args: arguments{
				url: "http://localhost:8080/api/user/balance",
			},
		},
	}
	handler := &Handler{
		Mux: chi.NewMux(),
		Cursor: &Cursor{
			NewMock(),
		},
	}
	handler.Post("/api/user/register", handler.RegisterUser)
	handler.Post("/api/user/login", handler.Login)
	handler.Get("/api/user/balance", handler.GetBalance)
	ts := httptest.NewServer(handler)
	handler.Cursor.Save(&UserInfo{
		Username: "test",
		Password: "test",
	})

	result := handler.Cursor.UpdateUserBalance(
		"test", expectedBalance,
	)
	assert.Equal(t, expectedBalance, result)

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	encoder.Encode(&UserInfo{
		Username: "test",
		Password: "test",
	})
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/login", buff)
	request.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, request)

	res := w.Result()

	cookies := res.Cookies()

	assert.Equal(t, res.StatusCode, 200)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.args.url, nil)
			request.AddCookie(cookies[0])
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()

			actualBalance := &Balance{}
			if err := json.NewDecoder(res.Body).Decode(&actualBalance); err != nil {
				panic(err)
			}
			assert.Equal(t, expectedBalance, actualBalance)

		})
	}
}
