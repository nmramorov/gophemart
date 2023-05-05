package main

import (
	"github.com/go-chi/chi/v5"
)

type MockDb struct {
	DbInterface
	storage     map[string]string
	sessions    map[string]Session
	orders      []*Order
	balance     map[string]*Balance
	withdrawals map[string][]*Withdrawal
}

type TestHandler struct {
	*chi.Mux
	DbApi *MockDb
}

func NewMock() *MockDb {
	return &MockDb{
		storage:     make(map[string]string),
		sessions:    make(map[string]Session),
		orders:      make([]*Order, 0),
		balance:     make(map[string]*Balance),
		withdrawals: make(map[string][]*Withdrawal),
	}
}

func (mock *MockDb) Connect() {}

func (mock *MockDb) Update() {}

func (mock *MockDb) SaveSession(id string, session *Session) {
	mock.sessions[id] = *session
}

func (mock *MockDb) SaveUserInfo(info *UserInfo) bool {

	for k := range mock.storage {
		if k == info.Username {
			return false
		}
	}

	mock.storage[info.Username] = info.Password
	return true
}

func (mock *MockDb) GetUserInfo(info *UserInfo) (*UserInfo, error) {
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

func (mock *MockDb) GetOrder(number string) (*Order, error) {
	for _, order := range mock.orders {
		if order.Number == number {
			return order, nil
		}
	}

	return nil, nil
}

func (mock *MockDb) SaveOrder(order *Order) {
	mock.orders = append(mock.orders, order)
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

func (mock *MockDb) GetWithdrawals(username string) ([]*Withdrawal, error) {
	return mock.withdrawals[username], nil
}

func (mock *MockDb) SaveWithdrawal(withdrawal *Withdrawal) {
	mock.withdrawals[withdrawal.User] = append(mock.withdrawals[withdrawal.User], withdrawal)
}

func (mock *MockDb) UpdateOrder(from *AccrualResponse) {
	for _, order := range mock.orders {
		if order.Number == from.Order {
			order.Accrual = from.Accrual
			order.Status = from.Status
			break
		}
	}
}

func (mock *MockDb) GetSession(token string) (*Session, bool) {
	session, ok := mock.sessions[token]
	if !ok {
		return nil, false
	}
	return &session, true
}
