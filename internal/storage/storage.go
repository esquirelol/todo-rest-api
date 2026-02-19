package storage

import (
	"context"
	"errors"

	"github.com/esquirelol/todo-rest-api/internal/todo"
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

func (st *Storage) Create(ctx context.Context, todo *todo.Todo) error {
	sqlQuery := `
		INSERT INTO tasks(author,title,description,completed_at)
		VALUES($1,$2,$3,$4);
`
	if _, err := st.conn.Exec(ctx, sqlQuery, todo.Author, todo.Title, todo.Description, todo.Completed_at); err != nil {
		st.logger.Error("storage: failed to create")
		return err
	}
	st.logger.Info("todo created")
	return nil
}

func (st *Storage) Get(ctx context.Context, task *todo.Todo) (*todo.Todo, error) {
	sqlQuery := `
	SELECT title,description,status,created_at FROM tasks
	WHERE author = $1 and title = $2
	LIMIT 1;
`
	var outTask todo.Todo
	err := st.conn.QueryRow(ctx, sqlQuery, task.Author, task.Title).Scan(
		&outTask.Title,
		&outTask.Description,
		&outTask.Status,
		&outTask.Created_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			st.logger.Info("storage: dont found this task")
			return nil, err
		}
		st.logger.Error("storage:", zap.Error(err))
		return nil, err
	}
	st.logger.Info("storage: get success")
	return &outTask, nil
}

func (st *Storage) Done(ctx context.Context, todo *todo.Todo) error {
	sqlQuery := `
	UPDATE tasks SET status = true
	WHERE author = $1 and title = $2;
`
	cmdTag, err := st.conn.Exec(ctx, sqlQuery, todo.Author, todo.Title)
	if err != nil {
		st.logger.Error("storage: failed to update")
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		st.logger.Info("dont found this task")
	}

	st.logger.Info("done task success")
	return nil
}

func (st *Storage) Delete(ctx context.Context, todo *todo.Todo) error {
	sqlQuery := `
	DELETE FROM tasks
	WHERE author = $1 and title = $2;
`
	cmdTag, err := st.conn.Exec(ctx, sqlQuery, todo.Author, todo.Title)
	if err != nil {
		st.logger.Error("storage: failed to delete")
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		st.logger.Info("dont delete")
	}

	st.logger.Info("delete task success")
	return nil
}
