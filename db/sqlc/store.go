package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to run db queries individually and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a transaction within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transactions
type TransferTxParams struct {
	FromAccount int64 `json:"from_account_id"`
	ToAccount   int64 `json:"to_account_id"`
	Amount      int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

/*
TransferTx performs a money transfer from one account to another
It creates a transfer record, add account entries and updates account balances within a single database transaction
*/
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccount,
			ToAccountID:   arg.ToAccount,
			Amount:        arg.Amount,
		})

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccount,
			Amount:    -arg.Amount,
		})

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccount,
			Amount:    arg.Amount,
		})

		//TODO: Update account balance

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
