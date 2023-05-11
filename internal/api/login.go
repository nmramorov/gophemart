package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/nmramorov/gophemart/internal/models"
)

func (h *UserRouter) Login(rw http.ResponseWriter, r *http.Request) {
	userInput := &models.UserInfo{}
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ValidateUserInfo(userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	dbData, err := h.Cursor.GetUserInfo(userInput)

	if err != nil {
		http.Error(rw, "wrong password/username", http.StatusUnauthorized)
		return
	}
	if err := ValidateLogin(userInput, dbData); err != nil {
		http.Error(rw, "wrong password/username", http.StatusUnauthorized)
		return
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(600 * time.Second)

	h.Cursor.SaveSession(sessionToken, &models.Session{
		Username:  userInput.Username,
		ExpiresAt: expiresAt,
		Token:     sessionToken,
	})

	http.SetCookie(rw, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	rw.WriteHeader(http.StatusOK)

	rw.Write([]byte(`success`))
}
