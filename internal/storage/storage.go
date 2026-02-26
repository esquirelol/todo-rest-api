package storage

import (
	"context"
	"errors"
	"strconv"

	"github.com/esquirelol/todo-rest-api/internal/dto"
	"github.com/esquirelol/todo-rest-api/internal/models"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Storage struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func ConnectionStorage(ctx context.Context, pathStorage string, logger *zap.Logger) (Storage, error) {
	conn, err := pgx.Connect(ctx, pathStorage)
	if err != nil {
		logger.Error("failed to connect storage", zap.Error(err))
		return Storage{}, err
	}
	storage := Storage{conn, logger}
	return storage, nil
}

func (st *Storage) Create(ctx context.Context, todo dto.Todo) error {
	sqlQuery := `
	INSERT INTO tasks("author","title","description")
	VALUES ($1,$2,$3);
`
	if _, err := st.conn.Exec(ctx, sqlQuery, todo.Author, todo.Title, todo.Description); err != nil {
		st.logger.Error("failed to create task", zap.Error(err))
		return err
	}
	st.logger.Info("created task success", zap.String("author:", todo.Author))
	return nil
}

func (st *Storage) Get(ctx context.Context, author string) ([]models.ModelTodo, error) {

	storageTask := make([]models.ModelTodo, 0)
	sqlQuery := `
	SELECT author,title,description,status,created_at,completed_at FROM tasks
	WHERE author = $1;
`
	rows, err := st.conn.Query(ctx, sqlQuery, author)
	if err != nil {
		st.logger.Error("failed to select task", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var outTask models.ModelTodo
		if err := rows.Scan(&outTask.Author,
			&outTask.Title,
			&outTask.Description,
			&outTask.Status,
			&outTask.CreatedAt,
			&outTask.CompletedAt,
		); err != nil {
			st.logger.Error("failed to scan data", zap.Error(err))
			return nil, err
		}
		storageTask = append(storageTask, outTask)
	}

	if len(storageTask) == 0 {
		st.logger.Info("task not exists")
		return nil, ErrNotExists
	}

	st.logger.Info("get success")
	return storageTask, nil
}

func (st *Storage) GetId(ctx context.Context, idTask string) (models.ModelTodo, error) {
	outTask := models.ModelTodo{}
	idTaskInt, err := strconv.Atoi(idTask)
	if err != nil {
		st.logger.Error("failed to conv id", zap.Error(err))
		return outTask, err
	}
	sqlQuery := `
	SELECT author,title,description,status FROM tasks
	WHERE id = $1
`
	if err := st.conn.QueryRow(ctx, sqlQuery, idTaskInt).Scan(
		&outTask.Author,
		&outTask.Title,
		&outTask.Description,
		&outTask.Status,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return outTask, ErrNotExists
		}
		return outTask, err
	}
	return outTask, nil
}

func (st *Storage) Update(ctx context.Context, todo dto.TodoUpdate, idTask string) error {
	task, err := st.GetId(ctx, idTask)

	if err != nil {
		if errors.Is(err, ErrNotExists) {
			st.logger.Info("task is not exists", zap.String("id_task:", idTask))
			return ErrNotExists
		}
		st.logger.Error("failed to get task", zap.Error(err))
	}
	idTaskInt, err := strconv.Atoi(idTask)
	if err != nil {
		st.logger.Error("failed to conv id", zap.Error(err))
		return err
	}
	author := task.Author
	if todo.Author != nil {
		author = *todo.Author
	}
	title := task.Title
	if todo.Title != nil {
		title = *todo.Title
	}
	description := task.Description
	if todo.Description != nil {
		description = *todo.Description
	}
	status := task.Status
	if todo.Status != nil {
		status = *todo.Status
	}
	sqlQuery := `
	UPDATE tasks SET author = $1,title = $2, description = $3,status = $4
	WHERE id = $5;
`
	res, err := st.conn.Exec(ctx, sqlQuery, author, title, description, status, idTaskInt)
	if err != nil {
		st.logger.Error("failed to update task", zap.Error(err))
		return err
	}
	if res.RowsAffected() == 0 {
		st.logger.Info("task is not exists", zap.String("id", idTask))
		return ErrNotExists
	}
	st.logger.Info("task update success")
	return nil
}

func (st *Storage) Delete(ctx context.Context, idTask string) error {
	idTaskInt, err := strconv.Atoi(idTask)
	if err != nil {
		st.logger.Error("failed to conv id", zap.Error(err))
		return err
	}
	sqlQuery := `
	DELETE FROM tasks
	WHERE id = $1;
`
	res, err := st.conn.Exec(ctx, sqlQuery, idTaskInt)
	if err != nil {
		st.logger.Error("failed to delete task", zap.String("id", idTask), zap.Error(err))
		return err
	}
	if res.RowsAffected() == 0 {
		st.logger.Info("task is not exists", zap.String("id", idTask))
		return ErrNotExists
	}
	st.logger.Info("task is deleted")
	return nil
}
