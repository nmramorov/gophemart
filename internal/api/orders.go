package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/theplant/luhn"

	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/logger"
	"github.com/nmramorov/gophemart/internal/models"
)

func (h *OrderRouter) UploadOrder(rw http.ResponseWriter, r *http.Request) {
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

	n, err := strconv.Atoi(requestNumber)
	if err != nil {
		http.Error(rw, "wrong number format", http.StatusUnprocessableEntity)
		return
	}
	if !luhn.Valid(n) {
		http.Error(rw, "wrong number format", http.StatusUnprocessableEntity)
		return
	}

	order, err := GetOrderFromDB(h.Cursor, username, requestNumber)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if order == nil {
		logger.InfoLog.Printf("Adding new order %s for user %s", requestNumber, username)
		newOrder := &models.Order{
			Number:     requestNumber,
			Username:   username,
			UploadedAt: time.Now(),
			Status:     "NEW",
		}
		err := ValidateOrder(h.Cursor, newOrder)
		if err != nil {
			logger.ErrorLog.Printf("Validation error for new order %s, token %s", newOrder.Number, sessionToken)
			http.Error(rw, "order was uploaded already by another user", http.StatusConflict)
			return
		}
		h.Cursor.SaveOrder(newOrder)
		err = h.Manager.AddJob(requestNumber, username)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte(`new order created`))
		return
	}

	logger.InfoLog.Println(order.Username, username)
	if order.Username != username {
		logger.ErrorLog.Printf("Validation error for order %s, token %s", order.Number, sessionToken)
		http.Error(rw, "order was uploaded already by another user", http.StatusConflict)
		return
	}
	logger.InfoLog.Printf("request number: %s", requestNumber)

	if order.Number == requestNumber {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`order created already`))
	}
}

func GetOrderFromDB(cursor *db.Cursor, username string, requestOrder string) (*models.Order, error) {
	order, err := cursor.GetOrder(username, requestOrder)
	if order == nil {
		return nil, err
	}
	return order, nil
}

func (h *OrderRouter) GetOrders(rw http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_token")
	sessionToken := cookie.Value
	username, err := h.Cursor.GetUsernameByToken(sessionToken)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	orders, err := h.Cursor.GetOrders(username)
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
		rw.WriteHeader(http.StatusOK)
		rw.Write(body.Bytes())
	}
}
