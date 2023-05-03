package main

type DbInterface interface {
	Connect()
	SaveUserInfo(*UserInfo) bool
	GetUserInfo(*UserInfo) (*UserInfo, error)
	Update()
	SaveSession(string, *Session)
	GetOrder(string) (*Order, error)
	SaveOrder(*Order)
	GetOrders() ([]*Order, error)
	GetUsernameByToken(string) (string, error)
	GetUserBalance(string) (*Balance, error)
	UpdateUserBalance(string, *Balance) *Balance
}

type Cursor struct {
	DbInterface
}
