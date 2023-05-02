package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)


func TestPostOrders(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	type arguments struct {
		url         string
		number      interface{}
		contentType string
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Test Positive order post",
			want: want{
				code:     202,
				response: "new order created",
			},
			args: arguments{
				url:         "http://localhost:8080/api/user/orders",
				number:      "14242412",
				contentType: "text/plain",
			},
		},
		{
			name: "Test Positive order posted already",
			want: want{
				code:     200,
				response: "order created already",
			},
			args: arguments{
				url:         "http://localhost:8080/api/user/orders",
				number:      "14242412",
				contentType: "text/plain",
			},
		},
		{
			name: "Test Negative post order wrong number",
			want: want{
				code:     422,
				response: "wrong number format\n",
			},
			args: arguments{
				url:         "http://localhost:8080/api/user/orders",
				number:      "",
				contentType: "text/plain",
			},
		},
		{
			name: "Test Negative post order already registered by another user",
			want: want{
				code:     409,
				response: "order was uploaded already by another user\n",
			},
			args: arguments{
				url:         "http://localhost:8080/api/user/orders",
				number:      "14242412",
				contentType: "text/plain",
			},
		},
		{
			name: "Test Negative post order bad request",
			want: want{
				code:     400,
				response: "wrong content\n",
			},
			args: arguments{
				url:         "http://localhost:8080/api/user/orders",
				number:      "0",
				contentType: "application/json",
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
	handler.Post("/api/user/orders", handler.UploadOrder)
	ts := httptest.NewServer(handler)
	handler.Cursor.Save(&UserInfo{
		Username: "test",
		Password: "test",
	})

	defer ts.Close()

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buff := bytes.NewBuffer([]byte{})
			_, _ = buff.Write([]byte(tt.args.number.(string)))
			request := httptest.NewRequest(http.MethodPost, tt.args.url, buff)
			request.Header.Add("Content-Type", tt.args.contentType)

			w := httptest.NewRecorder()
			if tt.name == "Test Negative post order already registered by another user" {
				buff := bytes.NewBuffer([]byte{})
				encoder := json.NewEncoder(buff)
				encoder.Encode(&UserInfo{
					Username: "test2",
					Password: "test2",
				})
				request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buff)
				request.Header.Add("Content-Type", "application/json")

				w := httptest.NewRecorder()
				handler.ServeHTTP(w, request)
				res := w.Result()
				cookies = res.Cookies()
			}
			request.AddCookie(cookies[0])
			handler.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}




