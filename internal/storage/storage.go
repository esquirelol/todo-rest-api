package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Storage struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func ConnectionStorage(ctx context.Context, storagePath string, logger *zap.Logger) *Storage {
	conn, err := pgx.Connect(ctx, storagePath)
	if err != nil {
		logger.Error("storage: failed to connection")
	}

	sqlQuery := `
	CREATE TABLE IF NOT EXISTS tasks(
	    id SERIAL PRIMARY KEY,
	    author TEXT NOT NULL,
	    title TEXT NOT NULL,
	    description TEXT,
	    status BOOLEAN DEFAULT false,
	    created_at TIMESTAMP DEFAULT now(),
	    completed_at TIMESTAMP
	);
`

	if _, err := conn.Exec(ctx, sqlQuery); err != nil {
		logger.Error("storage: failed to create table")
	}

	storage := &Storage{conn, logger}
	return storage
}
