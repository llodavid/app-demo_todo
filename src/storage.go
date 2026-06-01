package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	config *StorageConfig
	db     *sql.DB
}

type StorageConfig struct {
	Username string
	Password string
	DbName   string
	Port     string
	Host     string
}

func NewStorage(conf StorageConfig) (Storage, error) {
	slog.Debug("storage::NewStorage() - Started")
	dsnFormat := "%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local"
	dsn := fmt.Sprintf(dsnFormat, conf.Username, conf.Password, conf.Host, conf.Port, conf.DbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return Storage{}, fmt.Errorf("Failed to open SQL connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return Storage{}, fmt.Errorf("Failed to ping db: %w", err)
	}
	slog.Debug("storage::NewStorage() - Database opened", "dbName", conf.DbName)
	return Storage{
		config: &conf,
		db:     db,
	}, nil
}

func (s Storage) Destroy() error {
	slog.Debug("storage::Destroy() - Started")
	s.config = nil
	return s.db.Close()
}

func (s Storage) FindAllTodos() (*Todos, error) {
	slog.Debug("storage::FindAllTodos() - Started")
	sql := `SELECT id, title, completed, created_at, completed_at FROM todos ORDER BY created_at`
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//
	var todos Todos
	for rows.Next() {
		slog.Debug("storage::FindAllTodos() - Mapping row to todo")
		todo := &Todo{}
		var completedAt *time.Time
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Completed, &todo.CreatedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		if completedAt != nil {
			todo.CompletedAt = completedAt
		}
		slog.Debug("storage::FindAllTodos() - Todo added", "todo", *todo)
		todos.add(*todo)
	}
	return &todos, nil
}
