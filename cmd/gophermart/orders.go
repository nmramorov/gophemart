package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) UploadOrder(rw http.ResponseWriter, r *http.Request) {
	val := r.Header.Get("Content-Type")
	if val != "text/plain" {
		http.Error(rw, "wrong content", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, _ := r.Cookie("session_token")
	sessionToken := cookie.Value
	username, err := h.Cursor.GetUsernameByToken(sessionToken)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	requestNumber := string(body)

	if _, err := strconv.Atoi(requestNumber); err != nil {
		http.Error(rw, "wrong number format", http.StatusUnprocessableEntity)
		return
	}

	order, err := GetOrderFromDB(h.Cursor, requestNumber)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if order == nil {
		h.Cursor.SaveOrder(&Order{
			Number:     requestNumber,
			Username:   username,
			UploadedAt: time.Now(),
			Status:     "NEW",
		})
		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte(`new order created`))
		return
	}

	if order.Username != username {
		http.Error(rw, "order was uploaded already by another user", http.StatusConflict)
		return
	}
	InfoLog.Printf("request number: %s", requestNumber)
	//ToDo: call Accrual Worker to get updates for specific order in goroutine
	h.Manager.AddJob(requestNumber)

	if order.Number == requestNumber {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`order created already`))
	}
}

func GetOrderFromDB(cursor *Cursor, requestOrder string) (*Order, error) {
	order, err := cursor.GetOrder(requestOrder)
	if order == nil {
		return nil, err
	}
	return order, nil
}

func (h *Handler) GetOrders(rw http.ResponseWriter, r *http.Request) {
	// h.Manager.mu.Lock()
	// defer h.Manager.mu.Unlock()
	orders, err := h.Cursor.GetOrders()
	// InfoLog.Println(&orders)
	// defer h.Manager.mu.Unlock()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	if orders == nil {
		rw.WriteHeader(http.StatusNoContent)
		rw.Write([]byte(`no orders found`))
	} else {
		InfoLog.Println("Writing orders to JSON")
		body := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(body)
		encoder.Encode(&orders)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(body.Bytes())
	}
}
