package main

type DBInterface interface {
	Connect()
	SaveUserInfo(*UserInfo) bool
	GetUserInfo(*UserInfo) (*UserInfo, error)
	Update()
	SaveSession(string, *Session)
	GetSession(string) (*Session, bool)
	GetOrder(string, string) (*Order, error)
	SaveOrder(*Order)
	GetOrders(string) ([]*Order, error)
	GetUsernameByToken(string) (string, error)
	GetUserBalance(string) (*Balance, error)
	UpdateUserBalance(string, *Balance) *Balance
	GetWithdrawals(string) ([]*Withdrawal, error)
	SaveWithdrawal(*Withdrawal)
	UpdateOrder(string, *AccrualResponse)
	GetAllOrders() []*Order
}

type Cursor struct {
	DBInterface
}

func GetCursor(url string) *Cursor {
	return &Cursor{NewMock()}
}
