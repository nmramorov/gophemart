package main

import "time"

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	Username  string
	ExpiresAt time.Time
	Token     string
}

type Order struct {
	Number     string    `json:"number"`
	Username   string    `json:"-"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Balance struct {
	User string `json:"-"`
	Current float64 `json:"current"`
	Withdrawn int `json:"withdrawn"`
}
