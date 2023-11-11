package db

import (
	"context"
)

const addAccountBalance = `-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + $1
WHERE id = $2
RETURNING id, owner_name, balance, currency, created_at`

type AddAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID int64 `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, addAccountBalance, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
	owner_name,
	balance,
	currency
) VALUES (
	$1, $2, $3
) RETURNING id, owner_name, balance, currency, created_at`

type CreateAccountParams struct {
	OwnerName string `json:"owner_name"`
	Balance int64 `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.OwnerName, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, owner_name, balance, currency, created_at FROM accounts
WHERE id = $1 LIMIT 1`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

// FOR UPDATE - 1 transaction update at 1 time
// NO KEY - informs unique ID not being updated (safe), avoid deadlock
const getAccountForUpdate = `--name: GetAccountForUpdate :one
SELECT id, owner_name, balance, currency, created_at FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, owner_name, balance, currency, created_at FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2`

type ListAccountsParams struct {
	Limit int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.OwnerName,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING id, owner_name, balance, currency, created_at`

type UpdateAccountParams struct {
	ID int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
			&i.OwnerName,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
	)
	return i, err
}