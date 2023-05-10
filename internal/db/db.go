package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/nmramorov/gophemart/internal/logger"
	"github.com/nmramorov/gophemart/internal/models"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBInterface interface {
	Connect()
	SaveUserInfo(*models.UserInfo) bool
	GetUserInfo(*models.UserInfo) (*models.UserInfo, error)
	Update()
	SaveSession(string, *models.Session)
	GetSession(string) (*models.Session, bool)
	GetOrder(string, string) (*models.Order, error)
	SaveOrder(*models.Order)
	GetOrders(string) ([]*models.Order, error)
	GetUsernameByToken(string) (string, error)
	GetUserBalance(string) (*models.Balance, error)
	UpdateUserBalance(string, *models.Balance) *models.Balance
	GetWithdrawals(string) ([]*models.Withdrawal, error)
	SaveWithdrawal(*models.Withdrawal)
	SaveUserBalance(string, *models.Balance) *models.Balance
	UpdateOrder(string, *models.AccrualResponse)
	GetAllOrders() []*models.Order
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
}

func RunMigrations(databaseURL string) {
	m, err := migrate.New(
		"file://./migrations",
		databaseURL)
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}
	if err := m.Up(); err != nil {
		logger.ErrorLog.Fatal(err)
	}
	logger.InfoLog.Println("Migrations successfully executed")
}

func NewCursor(DBURL string) *DBCursor {
	db, err := sql.Open("pgx", DBURL)
	if err != nil {
		logger.ErrorLog.Printf("Unable to connect to database: %v\n", err)
		return nil
	}
	new := &DBCursor{
		DB:      db,
		Context: context.Background(),
	}
	// valid := new.Ping()
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
		logger.ErrorLog.Printf("ping error, database unreachable?: %e", err)
		return false
	}
	return true
}

func (c *DBCursor) Connect() {}

func (c *DBCursor) Update() {}

func (c *DBCursor) SaveSession(id string, session *models.Session) {
	_, err := c.DB.ExecContext(c.Context, SaveSession, session.Username, session.Token, session.ExpiresAt)
	if err != nil {
		logger.ErrorLog.Fatalf("error inserting row %s to db: %e", id, err)
	}
}

func (c *DBCursor) SaveUserInfo(info *models.UserInfo) bool {
	_, err := c.DB.ExecContext(c.Context, SaveUserInfo, info.Username, info.Password)
	if err != nil {
		logger.ErrorLog.Fatalf("error inserting row into Userinfo: %e", err)
		return false
	}
	return true
}

func (c *DBCursor) GetUserInfo(info *models.UserInfo) (*models.UserInfo, error) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetUserInfo, info.Username); row.Err() != nil {
		logger.ErrorLog.Fatalf("error during getting user info from db: %e", row.Err())
		return nil, row.Err()
	}
	foundInfo := &models.UserInfo{}
	err := row.Scan(&foundInfo.Username, &foundInfo.Password)
	if err != nil {
		logger.ErrorLog.Fatalf("error scanning userinfo from db: %e", err)
		return nil, err
	}
	return foundInfo, nil
}

func (c *DBCursor) GetOrder(username string, number string) (*models.Order, error) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetOrder, username, number); row.Err() != nil {
		logger.ErrorLog.Fatalf("error during getting order %s from db: %e", number, row.Err())
		return nil, row.Err()
	}
	foundOrder := &models.Order{}
	err := row.Scan(&foundOrder.Username, &foundOrder.Number, &foundOrder.Status, &foundOrder.Accrual, &foundOrder.UploadedAt)
	if err == sql.ErrNoRows {
		logger.ErrorLog.Printf("No rows found for order %s and user %s", number, username)
		return nil, nil
	}
	if err != nil {
		logger.ErrorLog.Fatalf("error scanning single order from db: %e", err)
		return nil, err
	}
	return foundOrder, nil
}

func (c *DBCursor) SaveOrder(order *models.Order) {
	_, err := c.DB.ExecContext(c.Context, SaveOrder, order.Username, order.Number, order.Status, order.Accrual, order.UploadedAt)
	if err != nil {
		logger.ErrorLog.Fatalf("error during saving order %s to db: %e", order.Number, err)
	}
}

func (c *DBCursor) GetOrders(username string) ([]*models.Order, error) {
	rows, err := c.DB.QueryContext(c.Context, GetOrders, username)
	if err != nil {
		logger.ErrorLog.Fatalf("error during getting orders from db: %e", err)
		return nil, err
	}
	if rows.Err() != nil {
		logger.ErrorLog.Fatalf("error during getting orders from db: %e", rows.Err())
		return nil, rows.Err()
	}
	foundOrders := []*models.Order{}
	for rows.Next() {
		var o models.Order
		if err = rows.Scan(&o.Username, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			logger.ErrorLog.Fatalf("error scanning order for %s from db: %e", username, err)
			return foundOrders, err
		}
		foundOrders = append(foundOrders, &o)
	}
	return foundOrders, nil
}

func (c *DBCursor) GetUsernameByToken(token string) (string, error) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetSessionUser, token); row.Err() != nil {
		logger.ErrorLog.Fatalf("error during getting current session user from db: %e", row.Err())
		return "", row.Err()
	}
	foundSession := &models.Session{}
	err := row.Scan(&foundSession.Username)
	if err != nil {
		logger.ErrorLog.Fatalf("error scanning session username from db: %e", err)
		return "", err
	}
	return foundSession.Username, nil
}

func (c *DBCursor) GetUserBalance(username string) (*models.Balance, error) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetBalance, username); row.Err() != nil {
		logger.ErrorLog.Fatalf("error during getting user balance from db: %e", row.Err())
		return nil, row.Err()
	}
	logger.InfoLog.Printf("Getting balance for user %s", username)
	foundBalance := &models.Balance{}
	err := row.Scan(&foundBalance.User, &foundBalance.Current, &foundBalance.Withdrawn)
	if err == sql.ErrNoRows {
		return foundBalance, nil
	}
	if err != nil {
		logger.ErrorLog.Fatalf("error scanning balance from db: %e", err)
		return nil, err
	}
	return foundBalance, nil
}

func (c *DBCursor) SaveUserBalance(username string, newBalance *models.Balance) *models.Balance {
	_, err := c.DB.ExecContext(c.Context, SaveBalance, username, newBalance.Current, newBalance.Withdrawn)
	if err != nil {
		logger.ErrorLog.Fatalf("error during saving balance for user %s: %e", username, err)
	}
	logger.InfoLog.Printf("Saved balance for %s, accrual is %f", username, newBalance.Current)
	newBalance.User = username
	return newBalance
}

func (c *DBCursor) UpdateUserBalance(username string, newBalance *models.Balance) *models.Balance {
	_, err := c.DB.ExecContext(c.Context, UpdateBalance, newBalance.Current, newBalance.Withdrawn, username)
	if err != nil {
		logger.ErrorLog.Fatalf("error during updating balance: %e", err)
	}
	logger.InfoLog.Printf("Balance updated, Current: %f, Withdrawn: %f for user %s", newBalance.Current, newBalance.Withdrawn, username)
	return newBalance
}

func (c *DBCursor) GetWithdrawals(username string) ([]*models.Withdrawal, error) {
	rows, err := c.DB.QueryContext(c.Context, GetWithdrawals, username)

	if err != nil {
		logger.ErrorLog.Fatalf("error during getting withdrawals from db: %e", err)
		return nil, err
	}
	if rows.Err() != nil {
		logger.ErrorLog.Fatalf("error during getting withdrawals from db: %e", rows.Err())
		return nil, rows.Err()
	}
	foundWithdrawals := []*models.Withdrawal{}
	for rows.Next() {
		var w models.Withdrawal
		if err := rows.Scan(&w.User, &w.Order, &w.Sum, &w.ProcessedAt); err != nil {
			logger.ErrorLog.Fatalf("error scanning withdrawal from db: %e", err)
			return foundWithdrawals, err
		}
		foundWithdrawals = append(foundWithdrawals, &w)
	}
	if err = rows.Err(); err != nil {
		return foundWithdrawals, err
	}
	return foundWithdrawals, nil
}

func (c *DBCursor) SaveWithdrawal(withdrawal *models.Withdrawal) {
	_, err := c.DB.ExecContext(c.Context, SaveWithdrawal, withdrawal.User, withdrawal.Order, withdrawal.Sum, withdrawal.ProcessedAt)
	if err != nil {
		logger.ErrorLog.Fatalf("error during saving withdrawal to db: %e", err)
	}
}

func (c *DBCursor) UpdateOrder(username string, from *models.AccrualResponse) {
	var status string
	if from.Status == "REGISTERED" {
		status = "PROCESSING"
	} else {
		status = from.Status
	}
	_, err := c.DB.ExecContext(c.Context, UpdateOrder, status, from.Accrual, username, from.Order)
	if err != nil {
		logger.ErrorLog.Fatalf("error during updating order: %e", err)
	}
}

func (c *DBCursor) GetSession(token string) (*models.Session, bool) {
	var row *sql.Row
	if row = c.DB.QueryRowContext(c.Context, GetSession, token); row.Err() != nil {
		logger.ErrorLog.Fatalf("error during getting user session from db: %e", row.Err())
		return nil, false
	}
	foundSession := &models.Session{}

	err := row.Scan(&foundSession.Username, &foundSession.Token, &foundSession.ExpiresAt)
	if err != nil {
		logger.ErrorLog.Fatalf("error scanning session from db: %e", err)
		return nil, false
	}
	return foundSession, true
}

func (c *DBCursor) GetAllOrders() []*models.Order {
	rows, err := c.DB.QueryContext(c.Context, GetAllOrders)

	if err != nil {
		logger.ErrorLog.Fatalf("error during getting all orders from db: %e", err)
		return nil
	}
	if rows.Err() != nil {
		logger.ErrorLog.Fatalf("error during getting all orders from db: %e", rows.Err())
		return nil
	}
	foundOrders := []*models.Order{}
	for rows.Next() {
		var o models.Order
		if err = rows.Scan(&o.Username, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			logger.ErrorLog.Printf("error scanning order among orders from db: %e", err)
			logger.ErrorLog.Println(foundOrders)
			return foundOrders
		}
		foundOrders = append(foundOrders, &o)
	}
	return foundOrders
}
