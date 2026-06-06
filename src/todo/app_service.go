package todo

import (
	todo_dto "RobertTC32/example-demo_hello/src/todo/dto"
	todo_store "RobertTC32/example-demo_hello/src/todo/store"
	"log/slog"
)

// "todo" contains application- and domain-services (aka use cases and entities) in the "Hexagonal Architecture";
// this simple crud application contains no real business logic,
// and no ports (interfaces) are defined for the adapters (implementations)

type AppService struct {
	Store *todo_store.Store
}

func NewAppService(store *todo_store.Store) (*AppService, error) {
	slog.Debug("todo::NewAppService() - Executing")
	return &AppService{
		Store: store,
	}, nil
}

func (t *AppService) GetTodos() ([]todo_dto.Todo, error) {
	return t.Store.GetTodos()
}

func (t *AppService) CreateTodoTitle(title string) error {
	return t.Store.CreateTodoTitle(title)
}

func (t *AppService) UpdateTodoCompleted(id uint32, version uint32) error {
	return t.Store.UpdateTodoCompleted(id, version, true)
}

func (t *AppService) DeleteTodo(id uint32, version uint32) error {
	return t.Store.DeleteTodo(id, version)
}

func (t *AppService) GetDbmapping1() ([]todo_dto.Dbmapping1, error) {
	return t.Store.GetDbmapping1()
}

func (t *AppService) GetDbmapping2() ([]todo_dto.Dbmapping2, error) {
	return t.Store.GetDbmapping2()
}

func (t *AppService) Destroy() error {
	return nil
}
