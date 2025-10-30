package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/codepnw/mini-ecommerce/pkg/config"
	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("db open failed: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db connect failed: %w", err)
	}
	return db, nil
}

type DBExec interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type TxManager struct {
	db *sql.DB
}

func InitTransaction(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Fatalf("rollback: %w", rbErr)
			}
		} else {
			cmErr := tx.Commit()
			err = cmErr
		}
	}()

	err = fn(tx)
	return err
}
