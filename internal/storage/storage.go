package storage

import (
	"context"
	"errors"
	"strconv"

	"github.com/esquirelol/todo-rest-api/internal/dto"
	"github.com/esquirelol/todo-rest-api/internal/models"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func (st *Storage) CreateUser(ctx context.Context, userName, password string) (int, error) {

	sqlQuery := `
	INSERT INTO users(user_name,password_hash)
	VALUES($1,$2)
	RETURNING id;
`
	var idUser int
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		st.logger.Error("failed to generate hash", zap.String("user_name", userName))
		return idUser, err
	}

	if err = st.conn.QueryRow(ctx, sqlQuery, userName, pass).Scan(&idUser); err != nil {
		st.logger.Error("failed to create user", zap.String("user_name", userName))
		return idUser, err
	}

	st.logger.Info("create user success", zap.String("user_name", userName))
	return idUser, nil

}

func (st *Storage) Create(ctx context.Context, todo dto.Todo) error {
	sqlQuery := `
	INSERT INTO tasks (user_id, title, description, status)
	VALUES($1,$2,$3,$4)
`
	if _, err := st.conn.Exec(ctx, sqlQuery, todo.UserId, todo.Title, todo.Description, todo.Status); err != nil {
		st.logger.Error("failed to create task", zap.Error(err))
		return err
	}

	st.logger.Info("created task success", zap.String("author:", todo.Author))
	return nil
}

func (st *Storage) Get(ctx context.Context, author string, idAuthor int) ([]models.ModelTodo, error) {
	storageTask := make([]models.ModelTodo, 0)

	sqlQuery := `
	SELECT t.id,u.user_name,t.title,t.description,t.status,t.created_at,t.completed_at 
	FROM users u JOIN tasks t 
	ON u.id = t.user_id
	WHERE u.id = $1 AND u.user_name = $2
`
	rows, err := st.conn.Query(ctx, sqlQuery, idAuthor, author)
	if err != nil {
		st.logger.Error("failed to select task", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var outTask models.ModelTodo
		if err := rows.Scan(&outTask.Id,
			&outTask.Author,
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
		st.logger.Error("failed to conv id", zap.String("id-task", idTask), zap.Error(err))
		return outTask, err
	}
	rds, err := NewRedis()
	if err != nil {
		st.logger.Error("failed to connect redis")
		return models.ModelTodo{}, err
	}

	task, err := rds.HGetAll(idTask).Result()

	if err == nil && len(task) == 0 {

		sqlQuery := `
			SELECT t.id,u.user_name,t.title,t.description,t.status FROM tasks t JOIN users u 
			ON u.id = t.user_id
			WHERE u.id = $1
			`
		if err := st.conn.QueryRow(ctx, sqlQuery, idTaskInt).Scan(
			&outTask.Id,
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
		rds.HSet(idTask, "id_task", outTask.Id)
		rds.HSet(idTask, "author", outTask.Author)
		rds.HSet(idTask, "title", outTask.Title)
		rds.HSet(idTask, "description", outTask.Description)
		rds.HSet(idTask, "status", outTask.Status)
		return outTask, nil

	}
	if err != nil {
		st.logger.Info("internal error redis", zap.Error(err))
		return models.ModelTodo{}, err
	}
	statusBool, err := strconv.ParseBool(task["status"])
	if err != nil {
		st.logger.Error("failed to parse", zap.String("id-task", idTask))
		return models.ModelTodo{}, err
	}
	outTask = models.ModelTodo{
		Id:          idTaskInt,
		Author:      task["author"],
		Title:       task["title"],
		Description: task["description"],
		Status:      statusBool,
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
	UPDATE tasks SET title = $1, description = $2,status = $3
	WHERE id = $4;
`
	res, err := st.conn.Exec(ctx, sqlQuery, title, description, status, idTaskInt)
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
