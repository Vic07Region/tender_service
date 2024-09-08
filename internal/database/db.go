package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func New(connection string) (*sql.DB, error) {
	dbo, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Проверьте соединение
	if err = dbo.Ping(); err != nil {
		dbo.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return dbo, nil
}
