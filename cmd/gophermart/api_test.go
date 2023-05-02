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

type TestHandler struct {
	*chi.Mux
	DbApi *MockDb
}

func TestRegistration(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	type userinfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type arguments struct {
		url     string
		payload userinfo
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Test Positive registration",
			want: want{
				code:     200,
				response: `user created successfully`,
			},
			args: arguments{
				url: "http://localhost:8080/api/user/register",
				payload: userinfo{
					Username: "test",
					Password: "test",
				},
			},
		},
		{
			name: "Test Negative registration user exists",
			want: want{
				code:     409,
				response: "user already exists\n",
			},
			args: arguments{
				url: "http://localhost:8080/api/user/register",
				payload: userinfo{
					Username: "test",
					Password: "test",
				},
			},
		},
		{
			name: "Test Negative registration invalid request",
			want: want{
				code:     400,
				response: "validation error\n",
			},
			args: arguments{
				url:     "http://localhost:8080/api/user/register",
				payload: userinfo{},
			},
		},
		{
			name: "Test Negative registration user exists and different password",
			want: want{
				code:     409,
				response: "user already exists\n",
			},
			args: arguments{
				url: "http://localhost:8080/api/user/register",
				payload: userinfo{
					Username: "test",
					Password: "new test",
				},
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
	ts := httptest.NewServer(handler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buff := bytes.NewBuffer([]byte{})
			encoder := json.NewEncoder(buff)
			encoder.Encode(&tt.args.payload)
			request := httptest.NewRequest(http.MethodPost, tt.args.url, buff)
			request.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

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
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
		})
	}
}

func TestAuthentication(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	type userinfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type arguments struct {
		url     string
		payload userinfo
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Test Positive authentication",
			want: want{
				code:     200,
				response: "success",
			},
			args: arguments{
				url: "http://localhost:8080/api/user/login",
				payload: userinfo{
					Username: "test",
					Password: "test",
				},
			},
		},
		{
			name: "Test Negative authentication wrong username",
			want: want{
				code:     401,
				response: "wrong password/username\n",
			},
			args: arguments{
				url: "http://localhost:8080/api/user/login",
				payload: userinfo{
					Username: "t",
					Password: "test",
				},
			},
		},
		{
			name: "Test Negative authentication wrong password",
			want: want{
				code:     401,
				response: "wrong password/username\n",
			},
			args: arguments{
				url: "http://localhost:8080/api/user/login",
				payload: userinfo{
					Username: "test",
					Password: "testtest",
				},
			},
		},
		{
			name: "Test Negative authentication wrong password and login",
			want: want{
				code:     401,
				response: "wrong password/username\n",
			},
			args: arguments{
				url: "http://localhost:8080/api/user/login",
				payload: userinfo{
					Username: "testtest",
					Password: "testtest",
				},
			},
		},
		{
			name: "Test Negative authentication bad request",
			want: want{
				code:     400,
				response: "validation error\n",
			},
			args: arguments{
				url:     "http://localhost:8080/api/user/login",
				payload: userinfo{},
			},
		},
	}
	handler := &Handler{
		Mux: chi.NewMux(),
		Cursor: &Cursor{
			NewMock(),
		},
	}
	handler.Post("/api/user/login", handler.Login)
	ts := httptest.NewServer(handler)
	handler.Cursor.Save(&UserInfo{
		Username: "test",
		Password: "test",
	})

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buff := bytes.NewBuffer([]byte{})
			encoder := json.NewEncoder(buff)
			encoder.Encode(&tt.args.payload)
			request := httptest.NewRequest(http.MethodPost, tt.args.url, buff)
			request.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

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
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
		})
	}
}

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
