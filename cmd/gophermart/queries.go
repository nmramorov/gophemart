package main

const (
	SaveSession    string = `INSERT INTO _sessions VALUES ($1, $2, $3);`
	GetUserInfo           = `SELECT * FROM userinfo WHERE username=$1;`
	GetOrder              = `SELECT * FROM orders WHERE username=$1 AND _number=$2;`
	SaveOrder             = `INSERT INTO orders VALUES ($1, $2, $3, $4, $5);`
	GetOrders             = `SELECT * FROM orders WHERE username=$1;`
	GetSessionUser        = `SELECT username FROM _session WHERE token=$1;`
	GetBalance            = `SELECT * FROM balances WHERE user=$1;`
	UpdateBalance  string = `UPDATE balances SET (current=$1, withdrawn=$2) WHERE user=$3;`
	GetWithdrawals        = `SELECT * FROM withdrawals WHERE user=$1;`
	SaveWithdrawal        = `INSERT INTO withdrawal VALUES ($1, $2, $3, $4);`
	UpdateOrder           = `UPDATE orders SET (_number=$1, status=$2, accrual=$3) WHERE username=$4;`
	GetSession            = `SELECT * FROM _session WHERE token=$1;`
	GetAllOrders          = `SELECT * FROM orders;`
)
