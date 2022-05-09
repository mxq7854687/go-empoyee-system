package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

// Create a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// Executes a function within a database transaction
func (store *Store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *Store) CreateEmployeeTx(ctx context.Context, arg CreateEmployeeParams) (Employee, error) {
	var result Employee

	err := store.execTX(ctx, func(q *Queries) error {
		var err error
		result, err = q.CreateEmployee(context.Background(), arg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
