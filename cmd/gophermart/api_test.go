package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type MockDb struct {
	DbInterface
	storage map[string]string
}

func NewMock() *MockDb {
	return &MockDb{
		storage: make(map[string]string),
	}
}

func (mock *MockDb) Connect() {}

func (mock *MockDb) Update() {}

func (mock *MockDb) Save(data interface{}) bool {
	convertedData := *data.(*UserInfo)

	for k := range mock.storage {
		if k == convertedData.Username {
			return false
		}
	}

	mock.storage[convertedData.Username] = convertedData.Password
	return true
}

func (mock *MockDb) Get(data interface{}) (interface{}, error) {
	convertedData := *data.(*UserInfo)
	for k, v := range mock.storage {
		if k == convertedData.Username {
			return &UserInfo{
				Username: k,
				Password: v,
			}, nil
		}
	}
	return nil, ErrValidation
}

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
