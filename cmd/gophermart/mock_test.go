package main

import (
	"github.com/go-chi/chi/v5"
)

type MockDb struct {
	DbInterface
	storage  map[string]string
	sessions map[string]Session
	orders   []*Order
	balance  map[string]*Balance
}

type TestHandler struct {
	*chi.Mux
	DbApi *MockDb
}

func NewMock() *MockDb {
	return &MockDb{
		storage:  make(map[string]string),
		sessions: make(map[string]Session),
		orders:   make([]*Order, 0),
		balance:  make(map[string]*Balance),
	}
}

func (mock *MockDb) Connect() {}

func (mock *MockDb) Update() {}

func (mock *MockDb) SaveSession(id string, data interface{}) {
	convertedData := *data.(*Session)
	mock.sessions[id] = convertedData
}

func (mock *MockDb) Save(data interface{}) bool {
	convertedData := *data.(*UserInfo)

	for k := range mock.storage {
		if k == convertedData.Username {
			return false
		}
	}

	mock.storage[convertedData.Username] = convertedData.Password
	return true
}

func (mock *MockDb) Get(data interface{}) (interface{}, error) {
	convertedData := *data.(*UserInfo)
	for k, v := range mock.storage {
		if k == convertedData.Username {
			return &UserInfo{
				Username: k,
				Password: v,
			}, nil
		}
	}
	return nil, ErrValidation
}

func (mock *MockDb) GetOrder(number interface{}) (interface{}, error) {
	mockNumber := number.(string)
	for _, order := range mock.orders {
		if order.Number == mockNumber {
			return order, nil
		}
	}

	return nil, nil
}

func (mock *MockDb) SaveOrder(order interface{}) {
	mockOrder := order.(*Order)
	mock.orders = append(mock.orders, mockOrder)
}

func (mock *MockDb) GetOrders() ([]*Order, error) {
	if len(mock.orders) == 0 {
		return nil, nil
	}
	return mock.orders, nil
}

func (mock *MockDb) GetUsernameByToken(token string) (string, error) {
	session, ok := mock.sessions[token]
	if !ok {
		return "", ErrValidation
	}
	return session.Username, nil
}

func (mock *MockDb) GetUserBalance(username string) (*Balance, error) {
	balance, ok := mock.balance[username]
	if !ok {
		return nil, ErrValidation
	}
	return balance, nil
}

func (mock *MockDb) UpdateUserBalance(username string, newBalance *Balance) *Balance {
	mock.balance[username] = newBalance
	return newBalance
}

func (mock *MockDb) GetUserTotalAccrual(username string) float64 {
	return mock.balance[username].Current
}
