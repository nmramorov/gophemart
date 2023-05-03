package main

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) WithdrawMoney(rw http.ResponseWriter, r *http.Request) {
	withrawal := &WithdrawalPost{}
	if err := json.NewDecoder(r.Body).Decode(&withrawal); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, _ := r.Cookie("session_token")
	sessionToken := cookie.Value
	username, err := h.Cursor.GetUsernameByToken(sessionToken)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ValidateOrder(h.Cursor, username, withrawal.Order)
	if err != nil {
		http.Error(rw, "invalid order number", http.StatusUnprocessableEntity)
		return
	}

	userBalance, err := h.Cursor.GetUserBalance(username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	resultedAccrual := userBalance.Current - withrawal.Sum
	if resultedAccrual < 0 {
		http.Error(rw, "not enough money", http.StatusPaymentRequired)
		return
	}
	resultedWithdrawn := userBalance.Withdrawn + withrawal.Sum

	_ = h.Cursor.UpdateUserBalance(username, &Balance{
			User: username,
			Current: resultedAccrual,
			Withdrawn: resultedWithdrawn,
		})

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`success`))
}
