package main

type DbInterface interface {
	Connect()
	Save(interface{}) bool
	Get(interface{}) (interface{}, error)
	Update()
	SaveSession(string, interface{})
	GetOrder(interface{}) (interface{}, error)
	SaveOrder(interface{})
	GetOrders() ([]*Order, error)
	GetUsernameByToken(token string) (string, error)
	GetUserBalance(username string) (*Balance, error)
	UpdateUserBalance(username string, newBalance *Balance) *Balance
	GetUserTotalAccrual(username string) float64
}

type Cursor struct {
	DbInterface
}
