package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBInterface interface {
	Connect()
	SaveUserInfo(*UserInfo) bool
	GetUserInfo(*UserInfo) (*UserInfo, error)
	Update()
	SaveSession(string, *Session) //
	GetSession(string) (*Session, bool)
	GetOrder(string, string) (*Order, error)
	SaveOrder(*Order)
	GetOrders(string) ([]*Order, error)
	GetUsernameByToken(string) (string, error)
	GetUserBalance(string) (*Balance, error)
	UpdateUserBalance(string, *Balance) *Balance
	GetWithdrawals(string) ([]*Withdrawal, error)
	SaveWithdrawal(*Withdrawal)
	SaveUserBalance(string, *Balance) *Balance
	UpdateOrder(string, *AccrualResponse)
	GetAllOrders() []*Order
}

type Cursor struct {
	DBInterface
}

func GetCursor(url string) *Cursor {
	return &Cursor{NewCursor(url)}
}

type DBCursor struct {
	DBInterface
	DB      *sql.DB
	Context context.Context
	IsValid bool
}

func RunMigrations(databaseURL string) {
	m, err := migrate.New(
		"file://./migrations",
		databaseURL)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	if err := m.Up(); err != nil {
		ErrorLog.Fatal(err)
	}
	InfoLog.Println("Migrations successfully executed")
}

func NewCursor(DBURL string) *DBCursor {
	db, err := sql.Open("pgx", DBURL)
	if err != nil {
		ErrorLog.Printf("Unable to connect to database: %v\n", err)
		return nil
	}
	new := &DBCursor{
		DB:      db,
		Context: context.Background(),
		IsValid: true,
	}
	valid := new.Ping()
	if !valid {
		new.IsValid = false
	}
	RunMigrations(DBURL)
	return new
}

func (c *DBCursor) Close() {
	c.DB.Close()
}

func (c *DBCursor) Ping() bool {
	ctx, cancel := context.WithTimeout(c.Context, 1*time.Second)
	defer cancel()
	if err := c.DB.PingContext(ctx); err != nil {
		ErrorLog.Printf("ping error, database unreachable?: %e", err)
		return false
	}
	return true
}

func (c *DBCursor) Connect() {}

func (c *DBCursor) Update() {}

func (c *DBCursor) SaveSession(id string, session *Session) {
	_, err := c.DB.ExecContext(c.Context, SaveSession, session.Username, session.Token, session.ExpiresAt)
	if err != nil {
		ErrorLog.Fatalf("error inserting row %s to db: %e", id, err)
	}
}

func (c *DBCursor) SaveUserInfo(info *UserInfo) bool {
	_, err := c.DB.ExecContext(c.Context, SaveUserInfo, info.Username, info.Password)
	if err != nil {
		ErrorLog.Fatalf("error inserting row into Userinfo: %e", err)
		return false
	}
	return true
}

func (c *DBCursor) GetUserInfo(info *UserInfo) (*UserInfo, error) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetUserInfo, info.Username); row.Err() != nil {
		ErrorLog.Fatalf("error during getting user info from db: %e", row.Err())
		return nil, row.Err()
	}
	foundInfo := &UserInfo{}
	err := row.Scan(&foundInfo.Username, &foundInfo.Password)
	if err != nil {
		ErrorLog.Fatalf("error scanning userinfo from db: %e", err)
		return nil, err
	}
	return foundInfo, nil
}

func (c *DBCursor) GetOrder(username string, number string) (*Order, error) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetOrder, username, number); row.Err() != nil {
		ErrorLog.Fatalf("error during getting order %s from db: %e", number, row.Err())
		return nil, row.Err()
	}
	foundOrder := &Order{}
	err := row.Scan(&foundOrder.Username, &foundOrder.Number, &foundOrder.Status, &foundOrder.Accrual, &foundOrder.UploadedAt)
	if err == sql.ErrNoRows {
		ErrorLog.Printf("No rows found for order %s and user %s", number, username)
		return nil, nil
	}
	if err != nil {
		ErrorLog.Fatalf("error scanning single order from db: %e", err)
		return nil, err
	}
	return foundOrder, nil
}

func (c *DBCursor) SaveOrder(order *Order) {
	_, err := c.DB.ExecContext(c.Context, SaveOrder, order.Username, order.Number, order.Status, order.Accrual, order.UploadedAt)
	if err != nil {
		ErrorLog.Fatalf("error during saving order %s to db: %e", order.Number, err)
	}
}

func (c *DBCursor) GetOrders(username string) ([]*Order, error) {
	rows, err := c.DB.QueryContext(c.Context, GetOrders, username)
	if err != nil {
		ErrorLog.Fatalf("error during getting orders from db: %e", err)
		return nil, err
	}
	if rows.Err() != nil {
		ErrorLog.Fatalf("error during getting orders from db: %e", rows.Err())
		return nil, rows.Err()
	}
	foundOrders := []*Order{}
	for rows.Next() {
		var o Order
		if err = rows.Scan(&o.Username, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			ErrorLog.Fatalf("error scanning order for %s from db: %e", username, err)
			return foundOrders, err
		}
		foundOrders = append(foundOrders, &o)
	}
	return foundOrders, nil
}

func (c *DBCursor) GetUsernameByToken(token string) (string, error) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetSessionUser, token); row.Err() != nil {
		ErrorLog.Fatalf("error during getting current session user from db: %e", row.Err())
		return "", row.Err()
	}
	foundSession := &Session{}
	err := row.Scan(&foundSession.Username)
	if err != nil {
		ErrorLog.Fatalf("error scanning session username from db: %e", err)
		return "", err
	}
	return foundSession.Username, nil
}

func (c *DBCursor) GetUserBalance(username string) (*Balance, error) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetBalance, username); row.Err() != nil {
		ErrorLog.Fatalf("error during getting user balance from db: %e", row.Err())
		return nil, row.Err()
	}
	InfoLog.Printf("Getting balance for user %s", username)
	foundBalance := &Balance{}
	err := row.Scan(&foundBalance.User, &foundBalance.Current, &foundBalance.Withdrawn)
	if err == sql.ErrNoRows {
		return foundBalance, nil
	}
	if err != nil {
		ErrorLog.Fatalf("error scanning balance from db: %e", err)
		return nil, err
	}
	return foundBalance, nil
}

func (c *DBCursor) SaveUserBalance(username string, newBalance *Balance) *Balance {
	_, err := c.DB.ExecContext(c.Context, SaveBalance, username, newBalance.Current, newBalance.Withdrawn)
	if err != nil {
		ErrorLog.Fatalf("error during saving balance for user %s: %e", username, err)
	}
	InfoLog.Printf("Saved balance for %s, accrual is %f", username, newBalance.Current)
	newBalance.User = username
	return newBalance
}

func (c *DBCursor) UpdateUserBalance(username string, newBalance *Balance) *Balance {
	_, err := c.DB.ExecContext(c.Context, UpdateBalance, newBalance.Current, newBalance.Withdrawn, username)
	// if err == sql.ErrNoRows {
	// 	c.SaveUserBalance(username, newBalance)
	// 	_, err := c.DB.ExecContext(c.Context, UpdateBalance, newBalance.Current, newBalance.Withdrawn, username)
	// 	if err != nil {
	// 		ErrorLog.Fatalf("error during updating balance after saving: %e", err)
	// 	}
	// }
	if err != nil {
		ErrorLog.Fatalf("error during updating balance: %e", err)
	}
	InfoLog.Printf("Balance updated, Current: %f, Withdrawn: %f for user %s", newBalance.Current, newBalance.Withdrawn, username)
	return newBalance
}

func (c *DBCursor) GetWithdrawals(username string) ([]*Withdrawal, error) {
	rows, err := c.DB.QueryContext(c.Context, GetWithdrawals, username)

	if err != nil {
		ErrorLog.Fatalf("error during getting withdrawals from db: %e", err)
		return nil, err
	}
	if rows.Err() != nil {
		ErrorLog.Fatalf("error during getting withdrawals from db: %e", rows.Err())
		return nil, rows.Err()
	}
	foundWithdrawals := []*Withdrawal{}
	for rows.Next() {
		var w Withdrawal
		if err := rows.Scan(&w.User, &w.Order, &w.Sum, &w.ProcessedAt); err != nil {
			ErrorLog.Fatalf("error scanning withdrawal from db: %e", err)
			return foundWithdrawals, err
		}
		foundWithdrawals = append(foundWithdrawals, &w)
	}
	if err = rows.Err(); err != nil {
		return foundWithdrawals, err
	}
	return foundWithdrawals, nil
}

func (c *DBCursor) SaveWithdrawal(withdrawal *Withdrawal) {
	_, err := c.DB.ExecContext(c.Context, SaveWithdrawal, withdrawal.User, withdrawal.Order, withdrawal.Sum, withdrawal.ProcessedAt)
	if err != nil {
		ErrorLog.Fatalf("error during saving withdrawal to db: %e", err)
	}
}

func (c *DBCursor) UpdateOrder(username string, from *AccrualResponse) {
	var status string
	if from.Status == "REGISTERED" {
		status = "PROCESSING"
	} else {
		status = from.Status
	}
	_, err := c.DB.ExecContext(c.Context, UpdateOrder, status, from.Accrual, username, from.Order)
	if err != nil {
		ErrorLog.Fatalf("error during updating order: %e", err)
	}
}

func (c *DBCursor) GetSession(token string) (*Session, bool) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetSession, token); row.Err() != nil {
		ErrorLog.Fatalf("error during getting user session from db: %e", row.Err())
		return nil, false
	}
	foundSession := &Session{}

	err := row.Scan(&foundSession.Username, &foundSession.Token, &foundSession.ExpiresAt)
	if err != nil {
		ErrorLog.Fatalf("error scanning session from db: %e", err)
		return nil, false
	}
	return foundSession, true
}

func (c *DBCursor) GetAllOrders() []*Order {
	rows, err := c.DB.QueryContext(c.Context, GetAllOrders)

	if err != nil {
		ErrorLog.Fatalf("error during getting all orders from db: %e", err)
		return nil
	}
	if rows.Err() != nil {
		ErrorLog.Fatalf("error during getting all orders from db: %e", rows.Err())
		return nil
	}
	foundOrders := []*Order{}
	for rows.Next() {
		var o Order
		if err = rows.Scan(&o.Username, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			ErrorLog.Printf("error scanning order among orders from db: %e", err)
			ErrorLog.Println(foundOrders)
			return foundOrders
		}
		foundOrders = append(foundOrders, &o)
	}
	return foundOrders
}
