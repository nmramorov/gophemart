package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type MockAccrualHandler struct {
	*chi.Mux
	mu             sync.RWMutex
	OrdersStorage  map[string]*AccrualResponse
	ShutdownButton chan struct{}
}

func (mock *MockAccrualHandler) GetAccrual(rw http.ResponseWriter, r *http.Request) {
	number := chi.URLParam(r, "number")
	// mock.mu.Lock()
	// defer mock.mu.Unlock()
	response, ok := mock.OrdersStorage[number]
	if !ok {
		rw.WriteHeader(http.StatusNoContent)
	}
	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	encoder.Encode(response)
	rw.Write(buff.Bytes())
}

func (mock *MockAccrualHandler) CreateOrder(rw http.ResponseWriter, r *http.Request) {
	number := chi.URLParam(r, "number")

	go func() {
		mock.mu.Lock()
		mock.OrdersStorage[number] = &AccrualResponse{
			Order:  number,
			Status: "REGISTERED",
		}
		mock.mu.Unlock()
		time.Sleep(2 * time.Second)
		n, _ := strconv.Atoi(number)
		mock.mu.Lock()
		if n%2 != 0 {
			mock.OrdersStorage[number].Status = "PROCESSING"
		} else {
			mock.OrdersStorage[number].Status = "INVALID"
		}
		mock.mu.Unlock()

		time.Sleep(2 * time.Second)

		mock.mu.Lock()
		mock.OrdersStorage[number].Accrual = 100
		mock.OrdersStorage[number].Status = "PROCESSED"
		mock.mu.Unlock()
	}()
	rw.WriteHeader(http.StatusOK)
}

func NewMockAccrualHandler() *MockAccrualHandler {
	handler := &MockAccrualHandler{
		Mux:            chi.NewMux(),
		OrdersStorage:  make(map[string]*AccrualResponse),
		ShutdownButton: make(chan struct{}),
	}
	handler.Get("/api/orders/{number}", handler.GetAccrual)
	handler.Post("/api/orders/{number}", handler.CreateOrder)
	return handler
}

func initMockAccrual() *http.Server {
	handler := NewMockAccrualHandler()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}
	handler.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		server.Shutdown(context.Background())
	})
	err := server.ListenAndServe()
	if err != nil {
		ErrorLog.Printf("Failed to launch mock accrual server:%+v\n", err)
	}
	return server
}

func TestAccrualValidOrder(t *testing.T) {
	defer httptest.NewRequest(http.MethodGet, "http://localhost:8080/shutdown", nil)

	go func() {
		initMockAccrual()
	}()
	client := &http.Client{}

	createOrder, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/orders/1", nil)
	resp, err := client.Do(createOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	getOrder, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/orders/1", nil)
	resp, err = client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	result := &AccrualResponse{}
	json.NewDecoder(resp.Body).Decode(result)
	assert.Equal(t, &AccrualResponse{
		Order:  "1",
		Status: "REGISTERED",
	}, result)

	time.Sleep(2 * time.Second)

	getOrder, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/api/orders/1", nil)
	resp, err = client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	result = &AccrualResponse{}
	json.NewDecoder(resp.Body).Decode(result)
	assert.Equal(t, &AccrualResponse{
		Order:  "1",
		Status: "PROCESSING",
	}, result)

	time.Sleep(2 * time.Second)

	getOrder, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/api/orders/1", nil)
	resp, err = client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	result = &AccrualResponse{}
	json.NewDecoder(resp.Body).Decode(result)
	assert.Equal(t, &AccrualResponse{
		Order:   "1",
		Status:  "PROCESSED",
		Accrual: 100,
	}, result)
}

func TestAccrualInvalidOrder(t *testing.T) {
	defer httptest.NewRequest(http.MethodGet, "http://localhost:8080/shutdown", nil)

	go func() {
		initMockAccrual()
	}()
	client := &http.Client{}

	createOrder, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/orders/2", nil)
	resp, err := client.Do(createOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	getOrder, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/orders/2", nil)
	resp, err = client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	result := &AccrualResponse{}
	json.NewDecoder(resp.Body).Decode(result)
	assert.Equal(t, &AccrualResponse{
		Order:  "2",
		Status: "REGISTERED",
	}, result)

	time.Sleep(2 * time.Second)

	getOrder, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/api/orders/2", nil)
	resp, err = client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	result = &AccrualResponse{}
	json.NewDecoder(resp.Body).Decode(result)
	assert.Equal(t, &AccrualResponse{
		Order:  "2",
		Status: "INVALID",
	}, result)
}

func TestAccrualNoSuchOrder(t *testing.T) {
	defer httptest.NewRequest(http.MethodGet, "http://localhost:8080/shutdown", nil)

	go func() {
		initMockAccrual()
	}()
	client := &http.Client{}

	getOrder, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/orders/3", nil)
	resp, err := client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	assert.Equal(t, 204, resp.StatusCode)
}
