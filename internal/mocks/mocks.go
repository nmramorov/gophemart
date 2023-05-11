package mocks

import (
	"github.com/go-chi/chi/v5"

	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/errors"
	"github.com/nmramorov/gophemart/internal/models"
)

type MockDB struct {
	db.DBInterface
	storage     map[string]string
	sessions    map[string]models.Session
	orders      map[string][]*models.Order
	balance     map[string]*models.Balance
	withdrawals map[string][]*models.Withdrawal
}

type TestHandler struct {
	*chi.Mux
	DBAPI  *db.Cursor
	Cursor *db.Cursor
}

func NewMock() *MockDB {
	return &MockDB{
		storage:     make(map[string]string),
		sessions:    make(map[string]models.Session),
		orders:      make(map[string][]*models.Order),
		balance:     make(map[string]*models.Balance),
		withdrawals: make(map[string][]*models.Withdrawal),
	}
}

func (mock *MockDB) SaveSession(id string, session *models.Session) {
	mock.sessions[id] = *session
}

func (mock *MockDB) SaveUserInfo(info *models.UserInfo) bool {

	for k := range mock.storage {
		if k == info.Username {
			return false
		}
	}

	mock.storage[info.Username] = info.Password
	return true
}

func (mock *MockDB) GetUserInfo(info *models.UserInfo) (*models.UserInfo, error) {
	for k, v := range mock.storage {
		if k == info.Username {
			return &models.UserInfo{
				Username: k,
				Password: v,
			}, nil
		}
	}
	return nil, errors.ErrValidation
}

func (mock *MockDB) GetOrder(username string, number string) (*models.Order, error) {
	for user, orders := range mock.orders {
		if user == username {
			for _, order := range orders {
				if order.Number == number {
					return order, nil
				}
			}
		}
	}

	return nil, nil
}

func (mock *MockDB) SaveOrder(order *models.Order) {
	mock.orders[order.Username] = append(mock.orders[order.Username], order)
}

func (mock *MockDB) GetOrders(username string) ([]*models.Order, error) {
	if len(mock.orders[username]) == 0 {
		return nil, nil
	}
	return mock.orders[username], nil
}

func (mock *MockDB) GetUsernameByToken(token string) (string, error) {
	session, ok := mock.sessions[token]
	if !ok {
		return "", errors.ErrValidation
	}
	return session.Username, nil
}

func (mock *MockDB) GetUserBalance(username string) (*models.Balance, error) {
	balance, ok := mock.balance[username]
	if !ok {
		return nil, errors.ErrValidation
	}
	return balance, nil
}

func (mock *MockDB) UpdateUserBalance(username string, newBalance *models.Balance) *models.Balance {
	mock.balance[username] = newBalance
	return newBalance
}

func (mock *MockDB) GetWithdrawals(username string) ([]*models.Withdrawal, error) {
	return mock.withdrawals[username], nil
}

func (mock *MockDB) SaveWithdrawal(withdrawal *models.Withdrawal) {
	mock.withdrawals[withdrawal.User] = append(mock.withdrawals[withdrawal.User], withdrawal)
}

func (mock *MockDB) UpdateOrder(username string, from *models.AccrualResponse) {
	orders := mock.orders[username]
	for _, order := range orders {
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

func (mock *MockDB) GetSession(token string) (*models.Session, bool) {
	session, ok := mock.sessions[token]
	if !ok {
		return nil, false
	}
	return &session, true
}

func (mock *MockDB) GetAllOrders() []*models.Order {
	result := make([]*models.Order, 0)
	for _, orders := range mock.orders {
		result = append(result, orders...)
	}
	return result
}
