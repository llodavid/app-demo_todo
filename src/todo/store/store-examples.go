package todo_store

import (
	todo_dto "RobertTC32/example-demo_hello/src/todo/dto"
	"encoding/json"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
)

// "store" (aka repository) contains adapter for DB in the "Hexagonal Architecture"

func (s *Store) GetDbmapping1() ([]todo_dto.Dbmapping1, error) {
	slog.Debug("store::GetDbmapping1() - Executing")
	var tt []todo_dto.Dbmapping1
	sql := `SELECT id, anint, abigint, anintunsighed, abigintunsigned, adecimal, afloat, adouble, aboolean, avarchar, adatetime, ablob
		FROM dbmapping ORDER BY id`
	if err := s.Db.Select(&tt, sql); err != nil {
		slog.Error("store::GetDbmapping1() - Failed to find dbmapping", "error", err)
		return []todo_dto.Dbmapping1{}, fmt.Errorf("Failed to find dbmapping: %w", err)
	}
	for i := range tt {
		t := tt[i]
		j, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println("GetDbmapping1() -", "index", i, "value \n", string(j))
	}
	return tt, nil
}

func (s *Store) GetDbmapping2() ([]todo_dto.Dbmapping2, error) {
	slog.Debug("store::GetDbmapping2() - Executing")
	var tt []todo_dto.Dbmapping2
	sql := `SELECT id, anint, abigint, anintunsighed, abigintunsigned, adecimal, afloat, adouble, aboolean, avarchar, adatetime, ablob
		FROM dbmapping ORDER BY id`
	if err := s.Db.Select(&tt, sql); err != nil {
		slog.Error("store::GetDbmapping2() - Failed to find dbmapping", "error", err)
		return []todo_dto.Dbmapping2{}, fmt.Errorf("Failed to find dbmapping: %w", err)
	}
	for i := range tt {
		t := tt[i]
		j, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println("GetDbmapping2() -", "index", i, "value \n", string(j))
	}
	return tt, nil
}
