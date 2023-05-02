package main

import "github.com/go-chi/chi/v5"

type MockDb struct {
	DbInterface
	storage  map[string]string
	sessions map[string]Session
	orders   []*Order
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
