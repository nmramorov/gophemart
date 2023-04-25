package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) RegisterUser(rw http.ResponseWriter, r *http.Request) {
	userInput := &UserInfo{}
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ValidateUserInfo(userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if ok := h.Cursor.Save(userInput); !ok {
		http.Error(rw, "user already exists", http.StatusConflict)
		return
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(5 * time.Second)

	h.Cursor.SaveSession(sessionToken, &Session{
		Username:  userInput.Username,
		ExpiresAt: expiresAt,
	})

	rw.WriteHeader(http.StatusOK)
	http.SetCookie(rw, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	rw.Write([]byte(`user created successfully`))
}
