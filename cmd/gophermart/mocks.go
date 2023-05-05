package main

import (
	"github.com/go-chi/chi/v5"
)

type MockDB struct {
	DBInterface
	storage     map[string]string
	sessions    map[string]Session
	orders      []*Order
	balance     map[string]*Balance
	withdrawals map[string][]*Withdrawal
}

type TestHandler struct {
	*chi.Mux
	DBAPI *MockDB
}

func NewMock() *MockDB {
	return &MockDB{
		storage:     make(map[string]string),
		sessions:    make(map[string]Session),
		orders:      make([]*Order, 0),
		balance:     make(map[string]*Balance),
		withdrawals: make(map[string][]*Withdrawal),
	}
}

func (mock *MockDB) Connect() {}

func (mock *MockDB) Update() {}

func (mock *MockDB) SaveSession(id string, session *Session) {
	mock.sessions[id] = *session
}

func (mock *MockDB) SaveUserInfo(info *UserInfo) bool {

	for k := range mock.storage {
		if k == info.Username {
			return false
		}
	}

	mock.storage[info.Username] = info.Password
	return true
}

func (mock *MockDB) GetUserInfo(info *UserInfo) (*UserInfo, error) {
	for k, v := range mock.storage {
		if k == info.Username {
			return &UserInfo{
				Username: k,
				Password: v,
			}, nil
		}
	}
	return nil, ErrValidation
}

func (mock *MockDB) GetOrder(number string) (*Order, error) {
	for _, order := range mock.orders {
		if order.Number == number {
			return order, nil
		}
	}

	return nil, nil
}

func (mock *MockDB) SaveOrder(order *Order) {
	mock.orders = append(mock.orders, order)
}

func (mock *MockDB) GetOrders() ([]*Order, error) {
	if len(mock.orders) == 0 {
		return nil, nil
	}
	return mock.orders, nil
}

func (mock *MockDB) GetUsernameByToken(token string) (string, error) {
	session, ok := mock.sessions[token]
	if !ok {
		return "", ErrValidation
	}
	return session.Username, nil
}

func (mock *MockDB) GetUserBalance(username string) (*Balance, error) {
	balance, ok := mock.balance[username]
	if !ok {
		return nil, ErrValidation
	}
	return balance, nil
}

func (mock *MockDB) UpdateUserBalance(username string, newBalance *Balance) *Balance {
	mock.balance[username] = newBalance
	return newBalance
}

func (mock *MockDB) GetWithdrawals(username string) ([]*Withdrawal, error) {
	return mock.withdrawals[username], nil
}

func (mock *MockDB) SaveWithdrawal(withdrawal *Withdrawal) {
	mock.withdrawals[withdrawal.User] = append(mock.withdrawals[withdrawal.User], withdrawal)
}

func (mock *MockDB) UpdateOrder(from *AccrualResponse) {
	for _, order := range mock.orders {
		if order.Number == from.Order {
			if from.Status == "REGISTERED" {
				order.Status = "PROCESSING"
				break
			}
			order.Accrual = from.Accrual
			order.Status = from.Status
			break
		}
	}
}

func (mock *MockDB) GetSession(token string) (*Session, bool) {
	session, ok := mock.sessions[token]
	if !ok {
		return nil, false
	}
	return &session, true
}
