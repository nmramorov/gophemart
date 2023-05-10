package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/nmramorov/gophemart/internal/models"
)

func (h *Handler) RegisterUser(rw http.ResponseWriter, r *http.Request) {
	userInput := &models.UserInfo{}
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ValidateUserInfo(userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if ok := h.Cursor.SaveUserInfo(userInput); !ok {
		http.Error(rw, "user already exists", http.StatusConflict)
		return
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(600 * time.Second)

	h.Cursor.SaveSession(sessionToken, &models.Session{
		Username:  userInput.Username,
		ExpiresAt: expiresAt,
		Token:     sessionToken,
	})
	h.Cursor.SaveUserBalance(userInput.Username, &models.Balance{
		User:      userInput.Username,
		Current:   0.0,
		Withdrawn: 0.0,
	})

	http.SetCookie(rw, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`user created successfully`))
}
