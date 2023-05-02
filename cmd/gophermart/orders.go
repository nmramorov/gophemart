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

	order, err := GetOrderFromDb(h.Cursor, requestNumber)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if order == nil {
		h.Cursor.SaveOrder(&Order{
			Number:     requestNumber,
			Username:   username,
			UploadedAt: time.Now(),
		})
		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte(`new order created`))
		return
	}

	if order.Username != username {
		http.Error(rw, "order was uploaded already by another user", http.StatusConflict)
		return
	}

	if order.Number == requestNumber {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`order created already`))
	}
}

func GetOrderFromDb(cursor *Cursor, requestOrder string) (*Order, error) {
	dbData, err := cursor.GetOrder(requestOrder)
	if dbData == nil {
		return nil, nil
	}
	order := dbData.(*Order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (h *Handler) GetOrders(rw http.ResponseWriter, r *http.Request) {
	orders, err := h.Cursor.GetOrders()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	if orders == nil {
		rw.WriteHeader(http.StatusNoContent)
		rw.Write([]byte(`no orders found`))
	} else {
		body := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(body)
		encoder.Encode(&orders)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(body.Bytes())
	}
}
