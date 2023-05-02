package main

import (
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
			Token:      sessionToken,
			UploadedAt: time.Now(),
		})
		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte(`new order created`))
		return
	}

	if order.Token != sessionToken {
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
