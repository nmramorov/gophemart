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
	GetWithdrawals(string) ([]*Withdrawal, error)
	SaveWithdrawal(string, *Withdrawal)
	UpdateOrder(*AccrualResponse)
}

type Cursor struct {
	DbInterface
}
