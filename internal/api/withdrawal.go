package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nmramorov/gophemart/internal/models"
)

func (h *BalanceRouter) WithdrawMoney(rw http.ResponseWriter, r *http.Request) {
	withrawal := &models.WithdrawalPost{}
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
	h.Cursor.SaveWithdrawal(&models.Withdrawal{
		User:        username,
		Order:       withrawal.Order,
		Sum:         withrawal.Sum,
		ProcessedAt: time.Now(),
	})
	_, err = h.Cursor.UpdateUserBalance(username, &models.Balance{
		User:      username,
		Current:   resultedAccrual,
		Withdrawn: resultedWithdrawn,
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`success`))
}

func (h *BalanceRouter) GetWithdrawals(rw http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_token")
	sessionToken := cookie.Value
	username, err := h.Cursor.GetUsernameByToken(sessionToken)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	withdrawals, err := h.Cursor.GetWithdrawals(username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	encoder.Encode(withdrawals)
	rw.Header().Set("Content-Type", "application/json")

	if withdrawals == nil {
		rw.WriteHeader(http.StatusNoContent)

	}
	rw.Write(buff.Bytes())
}
