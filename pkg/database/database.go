package database

import (
	"database/sql"
	"fmt"

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
