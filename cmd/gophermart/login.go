package main

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Login(rw http.ResponseWriter, r *http.Request) {
	userInput := &UserInfo{}
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ValidateUserInfo(userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	dbData, err := h.Cursor.Get(userInput)
	if err != nil {
		http.Error(rw, "wrong password/username", http.StatusUnauthorized)
		return
	}
	if err := ValidateLogin(userInput, dbData.(*UserInfo)); err != nil {
		http.Error(rw, "wrong password/username", http.StatusUnauthorized)
		return
	}
	rw.WriteHeader(http.StatusOK)

	rw.Write([]byte(`success`))
}
