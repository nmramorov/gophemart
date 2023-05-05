package main

type DbInterface interface {
	Connect()
	SaveUserInfo(*UserInfo) bool
	GetUserInfo(*UserInfo) (*UserInfo, error)
	Update()
	SaveSession(string, *Session)
	GetSession(string) (*Session, bool)
	GetOrder(string) (*Order, error)
	SaveOrder(*Order)
	GetOrders() ([]*Order, error)
	GetUsernameByToken(string) (string, error)
	GetUserBalance(string) (*Balance, error)
	UpdateUserBalance(string, *Balance) *Balance
	GetWithdrawals(string) ([]*Withdrawal, error)
	SaveWithdrawal(*Withdrawal)
	UpdateOrder(*AccrualResponse)
}

type Cursor struct {
	DbInterface
}

func GetCursor(url string) *Cursor {
	return &Cursor{NewMock()}
}
