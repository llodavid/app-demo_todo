package main

import (
	"time"
)

type Todo struct {
	Id          int64
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type Todos []Todo

func (todos *Todos) add(todo Todo) {
	*todos = append(*todos, todo)
}

func (todos *Todos) addNew(title string) {
	todo := Todo{
		Id:          0,
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
	*todos = append(*todos, todo)
}

func (todos *Todos) delete(index int) error {
	t := *todos
	*todos = append(t[:index], t[index+1:]...)
	return nil
}

func (todos *Todos) toggle(index int) error {
	t := *todos
	todo := &t[index]
	if !todo.Completed {
		completedTime := time.Now()
		todo.CompletedAt = &completedTime
	} else {
		todo.CompletedAt = nil
	}
	todo.Completed = !todo.Completed
	return nil
}

func (todos *Todos) edit(index int, title string) error {
	(*todos)[index].Title = title
	return nil
}
