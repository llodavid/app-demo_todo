package todo_store

import (
	todo_dto "RobertTC32/example-demo_hello/src/todo/dto"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

// "store" (aka repository) contains adapter for DB in the "Hexagonal Architecture"

func (s *Store) GetTodos() ([]todo_dto.Todo, error) {
	slog.Debug("store::GetTodos() - Executing")
	var tt []todo_dto.Todo
	sql := `SELECT id, title, completed, created_at, completed_at, version FROM todos ORDER BY id`
	if err := s.Db.Select(&tt, sql); err != nil {
		slog.Error("store::GetTodos() - Failed to find todos", "error", err)
		return []todo_dto.Todo{}, fmt.Errorf("Failed to find todos: %w", err)
	}
	// no mapping, no null handling, no rows closing and no error handling needed with sqlx
	slog.Debug("store::GetTodos() - Retrieved number", "todos", strconv.Itoa(len(tt)))
	return tt, nil
}

func (s *Store) GetTodoById(id int32) (todo_dto.Todo, error) {
	slog.Debug("store::GetTodoById() - Executing")
	var t todo_dto.Todo
	sql := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
	// MySQL/MariaDB only accept "?" and no "$1" in sql syntax
	if err := s.Db.Get(&t, sql, id); err != nil {
		slog.Error("store::GetTodoById() - Failed to find todo", "error", err)
		return todo_dto.Todo{}, fmt.Errorf("Failed to find todo: %w", err)
	}
	return t, nil
}

const TODOS_MAX = 20

func (s *Store) CreateTodo(t *todo_dto.Todo) error {
	slog.Debug("store::CreateTodo() - Executing")
	//
	var idCount int
	if err := s.Db.Get(&idCount, "SELECT count(id) as id_count FROM todos"); err != nil {
		slog.Error("store::CreateTodo() - Failed to count todo", "error", err)
		return fmt.Errorf("Failed to count todo: %w", err)
	}
	if idCount >= TODOS_MAX {
		slog.Error("store::CreateTodo() - Max todo exceeded", "maximum", TODOS_MAX)
		return fmt.Errorf("Max todo exceeded: maximum is %v", TODOS_MAX)
	}
	//
	sql := `INSERT INTO todos (title, completed, created_at, completed_at) 
		VALUES (?, ?, ?, ?) RETURNING *`
	// id and version are calculated by the database, so we need to get them back from database;
	// INSERT with RETURNING is allowed since MariaDB 10.5
	if err := s.Db.Get(&t, sql, t.Title, t.Completed, t.CreatedAt, t.CompletedAt); err != nil {
		slog.Error("store::CreateTodo() - Failed to insert todo", "error", err)
		return fmt.Errorf("Failed to insert todo: %w", err)
	}
	return nil
}

func (s *Store) CreateTodoTitle(title string) error {
	slog.Debug("store::CreateTodoTitle() - Executing")
	//
	var idCount int
	if err := s.Db.Get(&idCount, "SELECT count(id) as id_count FROM todos"); err != nil {
		slog.Error("store::CreateTodoTitle() - Failed to count todo", "error", err)
		return fmt.Errorf("Failed to count todo: %w", err)
	}
	if idCount >= TODOS_MAX {
		slog.Error("store::CreateTodoTitle() - Max todo exceeded", "maximum", TODOS_MAX)
		return fmt.Errorf("Max todo exceeded: maximum is %v", TODOS_MAX)
	}
	//
	sql := `INSERT INTO todos (title, completed, created_at) 
		VALUES (?, FALSE, NOW()) RETURNING *`
	// id and version are calculated by the database, so we need to get them back from database;
	// INSERT with RETURNING is allowed since MariaDB 10.5
	var t todo_dto.Todo
	if err := s.Db.Get(&t, sql, title); err != nil {
		slog.Error("store::CreateTodoTitle() - Failed to insert todo", "error", err)
		return fmt.Errorf("Failed to insert todo: %w", err)
	}
	return nil
}

func (s *Store) UpdateTodo(t *todo_dto.Todo) error {
	slog.Debug("store::UpdateTodo() - Executing")
	//
	sql := `UPDATE todos SET title = ?, completed = ?, created_at = ?, completed_at = ?, version = version + 1 
		WHERE id = ?`
	// optimistic locking
	if t.Version > 0 {
		sql = `UPDATE todos SET title = ?, completed = ?, created_at = ?, completed_at = ?, version = version + 1 
			WHERE id = ? and version = ?`
	}
	// id is calculated by the database, and version must be changed by client;
	// UPDATE with RETURNING is not allowed for MariaDB;
	// REPLACE with RETURNING is allowed since MariaDB 10.5, but can not be used for optimistic locking
	result, err := s.Db.Exec(sql, t.Title, t.Completed, t.CreatedAt, t.CompletedAt, t.Id, t.Version)
	rowsAffected, _ := result.RowsAffected()
	//
	if err != nil {
		slog.Error("store::UpdateTodo() - Failed to update todo", "error", err)
		return fmt.Errorf("Failed to update todo: %w", err)
	}
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
	if err2 := s.Db.Get(&t, sqlSelect, t.Id); err2 != nil {
		slog.Error("store::UpdateTodo() - Failed to find todo by id", "error", err2)
		return fmt.Errorf("Failed to find todo by id: %w", err2)
	}
	if rowsAffected == 0 {
		slog.Error("store::UpdateTodo() - Failed to find todo by id and rowversion", "error", err)
		return fmt.Errorf("Failed to find todo by id and rowversion: %w", err)
	}
	return nil
}

func (s *Store) UpdateTodoCompleted(id uint32, version uint32, completed bool) error {
	slog.Debug("store::UpdateTodoCompleted() - Executing")
	sql := `UPDATE todos SET completed = ?, completed_at = ?, version = version + 1 
		WHERE id = ?`
	// optimistic locking
	if version > 0 {
		sql = `UPDATE todos SET completed = ?, completed_at = ?, version = version + 1 
			WHERE id = ? and version = ?`
	}
	// id is calculated by the database, and version must be changed by client;
	// UPDATE with RETURNING is not allowed for MariaDB;
	// REPLACE with RETURNING is allowed since MariaDB 10.5, but can not be used for optimistic locking
	var completedAt *time.Time
	completedAt = nil
	if completed {
		now := time.Now()
		completedAt = &now
	}
	result, err := s.Db.Exec(sql, completed, completedAt, id, version)
	rowsAffected, _ := result.RowsAffected()
	//
	if err != nil {
		slog.Error("store::UpdateTodoCompleted() - Failed to update todo", "error", err)
		return fmt.Errorf("Failed to update todo: %w", err)
	}
	var t todo_dto.Todo
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
	if err2 := s.Db.Get(&t, sqlSelect, id); err2 != nil {
		slog.Error("store::UpdateTodoCompleted() - Failed to find todo by id", "error", err2)
		return fmt.Errorf("Failed to find todo by id: %w", err2)
	}
	if rowsAffected == 0 {
		slog.Error("store::UpdateTodoCompleted() - Failed to find todo by id and rowversion", "error", err)
		return fmt.Errorf("Failed to find todo by id and rowversion: %w", err)
	}
	return nil
}

func (s *Store) DeleteTodo(id uint32, version uint32) error {
	slog.Debug("store::DeleteTodo() - Executing")
	sql := `DELETE FROM todos WHERE id = ? RETURNING *`
	// optimistic locking
	if version > 0 {
		sql = `DELETE FROM todos WHERE id = ? AND version = ? RETURNING *`
	}
	// DELETE with RETURNING is allowed since MariaDB 10.0
	var t todo_dto.Todo
	if err := s.Db.Get(&t, sql, id, version); err != nil {
		// not found or sql error
		sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
		if err2 := s.Db.Get(&t, sqlSelect, id); err2 != nil {
			slog.Error("store::DeleteTodo() - Failed to find todo by id", "error", err2)
			return fmt.Errorf("Failed to find todo by id: %w", err2)
		}
		slog.Error("store::DeleteTodo() - Failed to delete todo", "error", err)
		return fmt.Errorf("Failed to delete todo: %w", err)
	}
	return nil
}

func (s *Store) SwitchTodoTitles(id1 uint32, id2 uint32) error {
	// example using transactions with optimistic locking
	slog.Debug("store::SwitchTodoTitles() - Executing")
	// we want to switch the titles of two todos
	tx, err := s.Db.BeginTxx(context.Background(), nil)
	if err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to begin transaction", "error", err)
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	//
	var t1 todo_dto.Todo
	sqlSelect := `SELECT id, title, completed, created_at, completed_at, version FROM todos WHERE id = ?`
	if err := s.Db.Get(&t1, sqlSelect, id1); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to find todo", "error", err)
		return fmt.Errorf("Failed to find todo: %w", err)
	}
	var t2 todo_dto.Todo
	if err := s.Db.Get(&t2, sqlSelect, id2); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to find todo", "error", err)
		return fmt.Errorf("Failed to find todo: %w", err)
	}
	sqlUpdate := `UPDATE todos SET title = ?, completed = ?, created_at = ?, completed_at = ?, version = version + 1 
			WHERE id = ? and version = ?`
	if _, err := s.Db.Exec(sqlUpdate, t2.Title, t2.Completed, t2.CreatedAt, t2.CompletedAt, t1.Id, t1.Version); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to update todo", "error", err)
		return fmt.Errorf("Failed to update todo: %w", err)
	}
	if _, err := s.Db.Exec(sqlUpdate, t1.Title, t1.Completed, t1.CreatedAt, t1.CompletedAt, t2.Id, t2.Version); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to update todo", "error", err)
		return fmt.Errorf("Failed to update todo: %w", err)
	}
	if err := tx.Commit(); err != nil {
		slog.Error("store::SwitchTodoTitles() - Failed to commit", "error", err)
		return fmt.Errorf("Failed to commit: %w", err)
	}
	return nil
}
