package main

import (
	"encoding/json"
	"net/http"
)

type UserRegistrationInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) RegisterUser(rw http.ResponseWriter, r *http.Request) {
	userInput := &UserRegistrationInfo{}
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if ok := h.Cursor.Save(userInput); !ok {
		http.Error(rw, "user already exists", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)

	rw.Write([]byte(`user created successfully`))
}
