package todo_dto

import (
	"time"
)

// "dto" contains Data Transfert Objects which are used to transfer data between layers;
// the dto package should not depend on any other packages in the application (to prevent circular dependencies)

type Todo struct {
	Id          uint32 // mapping "int32(go) <-> int(mysql)" and "int64(go) <-> bigint(mysql)" ("int" in go is platform dependent)
	Title       string
	Completed   bool
	CreatedAt   time.Time  `db:"created_at"`
	CompletedAt *time.Time `db:"completed_at"` // nullable in the database, so we use a pointer
	Version     uint32     // used for optimistic locking
}
