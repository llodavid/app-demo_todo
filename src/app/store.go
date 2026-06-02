package app

import (
	"RobertTC32/example-demo_hello/src/commons"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	dsn string
	db  *sqlx.DB
}

func NewStore() (*Store, error) {
	slog.Debug("store::NewStore() - Executing")
	//dsnFormat := "%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local"
	//dsn := fmt.Sprintf(dsnFormat, Username, Password, Host, Port, DbName)
	dbdriver := os.Getenv("DB_DRIVER")
	dsn := os.Getenv("DB_DSN")
	if commons.IsRunningInDockerContainer() {
		// change database localhost to host.docker.internal
		dsn = strings.Replace(dsn, "localhost", "host.docker.internal", 1)
	}
	db, err := sqlx.Open(dbdriver, dsn)
	if err != nil {
		slog.Error("store::NewStore() - Failed to open SQL connection", "error", err)
		return &Store{}, fmt.Errorf("Failed to open SQL connection: %w", err)
	}
	// Pool configuration
	poolFlag, err := strconv.Atoi(os.Getenv("DB_POOL_MAX_OPEN_CONNS"))
	if err == nil {
		slog.Debug("store::NewStore() - Set", "maxOpenConns", poolFlag)
		db.SetMaxOpenConns(poolFlag) // max simultaneous connections
	}
	poolFlag, _ = strconv.Atoi(os.Getenv("DB_POOL_MAX_IDLE_CONNS"))
	if err == nil {
		slog.Debug("store::NewStore() - Set", "maxIdleConns", poolFlag)
		db.SetMaxIdleConns(poolFlag) // connections kept idle in pool
	}
	poolFlag, _ = strconv.Atoi(os.Getenv("DB_POOL_CONN_MAX_LIFETIME"))
	if err == nil {
		slog.Debug("store::NewStore() - Set", "connMaxLifetime", poolFlag)
		db.SetConnMaxLifetime(time.Duration(poolFlag) * time.Minute) // max connection age
	}
	poolFlag, _ = strconv.Atoi(os.Getenv("DB_POOL_CONN_MAX_IDLE_TIME"))
	if err == nil {
		slog.Debug("store::NewStore() - Set", "connMaxIdleTime", poolFlag)
		db.SetConnMaxIdleTime(time.Duration(poolFlag) * time.Minute) // evict idle connections after this
	}
	err = db.Ping()
	if err != nil {
		slog.Error("store::NewStore() - Failed to ping db", "error", err)
		return &Store{}, fmt.Errorf("Failed to ping db: %w", err)
	}
	slog.Debug("store::NewStore() - Database opened", "dsn", dsn)
	return &Store{
		dsn: dsn,
		db:  db,
	}, nil
}

func (s *Store) Destroy() error {
	slog.Debug("store::Destroy() - Executing")
	return s.db.Close()
}

func (s *Store) GetTodos() ([]Todo, error) {
	slog.Debug("store::GetTodos() - Executing")
	var tt []Todo
	sql := `SELECT id, title, completed, created_at, completed_at, version FROM todos ORDER BY id`
	if err := s.db.Select(&tt, sql); err != nil {
		slog.Error("store::GetTodos() - Failed to find todos", "error", err)
		return []Todo{}, fmt.Errorf("Failed to find todos: %w", err)
	}
	// no mapping, no null handling, no rows closing and no error handling needed with sqlx
	return tt, nil
}

func (s *Store) GetTodoById(id int32) (Todo, error) {
	slog.Debug("store::GetTodoById() - Executing")
	var t Todo
	sql := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = $1`
	if err := s.db.Get(&t, sql, id); err != nil {
		slog.Error("store::GetTodoById() - Failed to find todo", "error", err)
		return Todo{}, fmt.Errorf("Failed to find todo: %w", err)
	}
	return t, nil
}

func (s *Store) CreateTodo(t *Todo) error {
	slog.Debug("store::CreateTodo() - Executing")
	sql := `INSERT INTO todos (title, completed, created_at, completed_at) 
		VALUES ($1, $2, $3, $4) RETURNING *`
	// id and version are calculated by the database, so we need to get them back from database
	if err := s.db.Get(&t, sql, t.Title, t.Completed, t.CreatedAt, t.CompletedAt); err != nil {
		slog.Error("store::CreateTodo() - Failed to create todo", "error", err)
		return fmt.Errorf("Failed to create todo: %w", err)
	}
	return nil
}

func (s *Store) UpdateTodo(t *Todo) error {
	slog.Debug("store::UpdateTodo() - Executing")
	sql := `UPDATE todos SET title = $1, completed = $2, created_at = $3, completed_at = $4, version = version + 1 
		WHERE id = $5 RETURNING *`
	// optimistic locking
	if t.Version > 0 {
		sql = `UPDATE todos SET title = $1, completed = $2, created_at = $3, completed_at = $4, version = version + 1 
			WHERE id = $5 and version = $6 RETURNING *`
	}
	// id is calculated by the database, and version must be changed by client
	if err := s.db.Get(&t, sql, t.Title, t.Completed, t.CreatedAt, t.CompletedAt, t.Id, t.Version); err != nil {
		slog.Error("store::UpdateTodo() - Failed to update todo", "error", err)
		return fmt.Errorf("Failed to update todo: %w", err)
	}
	return nil
}

func (s *Store) DeleteTodo(id int32, version int32) error {
	slog.Debug("store::DeleteTodo() - Executing")
	sql := `DELETE todos WHERE id = $1`
	// optimistic locking
	if version > 0 {
		sql = `DELETE todos WHERE id = $1 AND version = $2`
	}
	if _, err := s.db.Exec(sql, id, version); err != nil {
		slog.Error("store::DeleteTodo() - Failed to delete todo", "error", err)
		return fmt.Errorf("Failed to delete todo: %w", err)
	}
	return nil
}

func (s *Store) SwitchTodoTitles(id1 int32, id2 int32) error {
	// example using transactions with optimistic locking
	slog.Debug("store::SwitchTodoTitles() - Executing")
	// we want to switch the titles of two todos
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to begin transaction", "error", err)
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	//
	var t1 Todo
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = $1`
	if err := s.db.Get(&t1, sqlSelect, id1); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to find todo", "error", err)
		return fmt.Errorf("Failed to find todo: %w", err)
	}
	var t2 Todo
	if err := s.db.Get(&t2, sqlSelect, id2); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to find todo", "error", err)
		return fmt.Errorf("Failed to find todo: %w", err)
	}
	sqlUpdate := `UPDATE todos SET title = $1, completed = $2, created_at = $3, completed_at = $4, version = version + 1 
			WHERE id = $5 and version = $6`
	if _, err := s.db.Exec(sqlUpdate, t2.Title, t2.Completed, t2.CreatedAt, t2.CompletedAt, t1.Id, t1.Version); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to update todo", "error", err)
		return fmt.Errorf("Failed to update todo: %w", err)
	}
	if _, err := s.db.Exec(sqlUpdate, t1.Title, t1.Completed, t1.CreatedAt, t1.CompletedAt, t2.Id, t2.Version); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to update todo", "error", err)
		return fmt.Errorf("Failed to update todo: %w", err)
	}
	tx.Commit()
	return nil
}
