package main

import "time"

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	Username  string
	ExpiresAt time.Time
}
